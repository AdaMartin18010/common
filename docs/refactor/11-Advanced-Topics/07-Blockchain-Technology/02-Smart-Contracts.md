# 2. 智能合约

## 概述

智能合约（Smart Contracts）是运行在区块链上的自动执行程序，能够在满足预设条件时自动执行合约条款，实现去中心化的自动化交易和业务逻辑。

## 2.1 智能合约定义

### 2.1.1 智能合约模型

智能合约 $SC$ 是一个五元组 $(S, F, T, E, G)$，其中：

```latex
$$SC = (S, F, T, E, G)$$
```

其中：

- $S$: 状态集合
- $F$: 函数集合
- $T$: 触发条件
- $E$: 执行环境
- $G$: 气体机制

### 2.1.2 合约状态

合约状态 $s \in S$ 包含：

```latex
$$s = (Storage, Balance, Code, Nonce)$$
```

其中：

- $Storage$: 存储数据
- $Balance$: 合约余额
- $Code$: 合约代码
- $Nonce$: 交易计数器

## 2.2 智能合约执行模型

### 2.2.1 虚拟机模型

智能合约虚拟机 $VM$ 执行指令：

```latex
$$VM = (PC, Stack, Memory, Storage, Gas)$$
```

其中：

- $PC$: 程序计数器
- $Stack$: 操作数栈
- $Memory$: 内存
- $Storage$: 存储
- $Gas$: 气体消耗

### 2.2.2 气体机制

气体消耗计算：

```latex
$$Gas_{total} = \sum_{i=1}^n Gas_{op_i} \times Count_{op_i}$$
```

其中 $Gas_{op_i}$ 是操作 $i$ 的气体消耗。

### 2.2.3 执行环境

执行环境 $E$ 包含：

```latex
$$E = (Address, Origin, Caller, Value, Data, GasPrice, Block)$$
```

## 2.3 智能合约语言

### 2.3.1 Solidity语法

Solidity合约结构：

```solidity
contract ContractName {
    // 状态变量
    mapping(address => uint) balances;
    
    // 事件
    event Transfer(address from, address to, uint amount);
    
    // 函数
    function transfer(address to, uint amount) public {
        require(balances[msg.sender] >= amount);
        balances[msg.sender] -= amount;
        balances[to] += amount;
        emit Transfer(msg.sender, to, amount);
    }
}
```

### 2.3.2 合约生命周期

合约生命周期：

```latex
$$Contract_{lifecycle} = (Deploy, Initialize, Execute, Terminate)$$
```

## 2.4 智能合约安全

### 2.4.1 常见漏洞

#### 重入攻击

重入攻击模式：

```latex
$$Attack_{reentrant} = (Call, Fallback, Reenter)$$
```

#### 整数溢出

整数溢出检查：

```latex
$$Safe_{add}(a, b) = \begin{cases} a + b & \text{if } a + b \geq 0 \\ \text{revert} & \text{otherwise} \end{cases}$$
```

#### 访问控制

访问控制检查：

```latex
$$Access_{control}(user, role) = \begin{cases} \text{allow} & \text{if } user \in role \\ \text{deny} & \text{otherwise} \end{cases}$$
```

### 2.4.2 安全模式

#### 检查-效果-交互模式

CEI模式：

```latex
$$CEI_{pattern} = (Check, Effect, Interaction)$$
```

#### 拉取支付模式

拉取支付：

```latex
$$Pull_{payment} = (Withdraw, Balance, Transfer)$$
```

## 2.5 Go语言实现

### 2.5.1 智能合约虚拟机

```go
package smartcontracts

import (
    "fmt"
    "math/big"
)

// VirtualMachine 智能合约虚拟机
type VirtualMachine struct {
    PC      int
    Stack   []*big.Int
    Memory  map[int]*big.Int
    Storage map[string]*big.Int
    Gas     int64
    GasUsed int64
}

// NewVirtualMachine 创建虚拟机
func NewVirtualMachine(gasLimit int64) *VirtualMachine {
    return &VirtualMachine{
        PC:      0,
        Stack:   make([]*big.Int, 0),
        Memory:  make(map[int]*big.Int),
        Storage: make(map[string]*big.Int),
        Gas:     gasLimit,
        GasUsed: 0,
    }
}

// Push 压栈操作
func (vm *VirtualMachine) Push(value *big.Int) error {
    if vm.GasUsed >= vm.Gas {
        return fmt.Errorf("out of gas")
    }
    
    vm.Stack = append(vm.Stack, value)
    vm.GasUsed += 3 // 基础气体消耗
    return nil
}

// Pop 出栈操作
func (vm *VirtualMachine) Pop() (*big.Int, error) {
    if vm.GasUsed >= vm.Gas {
        return nil, fmt.Errorf("out of gas")
    }
    
    if len(vm.Stack) == 0 {
        return nil, fmt.Errorf("stack underflow")
    }
    
    value := vm.Stack[len(vm.Stack)-1]
    vm.Stack = vm.Stack[:len(vm.Stack)-1]
    vm.GasUsed += 2
    
    return value, nil
}

// Add 加法操作
func (vm *VirtualMachine) Add() error {
    if vm.GasUsed >= vm.Gas {
        return fmt.Errorf("out of gas")
    }
    
    if len(vm.Stack) < 2 {
        return fmt.Errorf("stack underflow")
    }
    
    a := vm.Stack[len(vm.Stack)-2]
    b := vm.Stack[len(vm.Stack)-1]
    
    result := new(big.Int).Add(a, b)
    vm.Stack[len(vm.Stack)-2] = result
    vm.Stack = vm.Stack[:len(vm.Stack)-1]
    
    vm.GasUsed += 3
    return nil
}

// Store 存储操作
func (vm *VirtualMachine) Store(key string) error {
    if vm.GasUsed >= vm.Gas {
        return fmt.Errorf("out of gas")
    }
    
    if len(vm.Stack) == 0 {
        return fmt.Errorf("stack underflow")
    }
    
    value := vm.Stack[len(vm.Stack)-1]
    vm.Storage[key] = value
    vm.Stack = vm.Stack[:len(vm.Stack)-1]
    
    vm.GasUsed += 20000 // SSTORE气体消耗
    return nil
}

// Load 加载操作
func (vm *VirtualMachine) Load(key string) error {
    if vm.GasUsed >= vm.Gas {
        return fmt.Errorf("out of gas")
    }
    
    value, exists := vm.Storage[key]
    if !exists {
        value = big.NewInt(0)
    }
    
    vm.Stack = append(vm.Stack, value)
    vm.GasUsed += 200 // SLOAD气体消耗
    return nil
}
```

### 2.5.2 智能合约

```go
// SmartContract 智能合约
type SmartContract struct {
    Address     string
    Code        []byte
    Storage     map[string]*big.Int
    Balance     *big.Int
    Nonce       int64
    Owner       string
}

// NewSmartContract 创建智能合约
func NewSmartContract(address string, code []byte, owner string) *SmartContract {
    return &SmartContract{
        Address: address,
        Code:    code,
        Storage: make(map[string]*big.Int),
        Balance: big.NewInt(0),
        Nonce:   0,
        Owner:   owner,
    }
}

// Execute 执行合约
func (sc *SmartContract) Execute(vm *VirtualMachine, input []byte) ([]byte, error) {
    // 解析输入
    instruction, err := sc.parseInstruction(input)
    if err != nil {
        return nil, err
    }
    
    // 执行指令
    switch instruction.Opcode {
    case "PUSH":
        return sc.executePush(vm, instruction)
    case "ADD":
        return sc.executeAdd(vm, instruction)
    case "STORE":
        return sc.executeStore(vm, instruction)
    case "LOAD":
        return sc.executeLoad(vm, instruction)
    default:
        return nil, fmt.Errorf("unknown opcode: %s", instruction.Opcode)
    }
}

// parseInstruction 解析指令
func (sc *SmartContract) parseInstruction(input []byte) (*Instruction, error) {
    // 简化实现：实际应解析字节码
    return &Instruction{
        Opcode: "PUSH",
        Operand: big.NewInt(0),
    }, nil
}

// Instruction 指令结构
type Instruction struct {
    Opcode  string
    Operand *big.Int
}

// executePush 执行PUSH指令
func (sc *SmartContract) executePush(vm *VirtualMachine, inst *Instruction) ([]byte, error) {
    err := vm.Push(inst.Operand)
    if err != nil {
        return nil, err
    }
    return []byte("success"), nil
}

// executeAdd 执行ADD指令
func (sc *SmartContract) executeAdd(vm *VirtualMachine, inst *Instruction) ([]byte, error) {
    err := vm.Add()
    if err != nil {
        return nil, err
    }
    return []byte("success"), nil
}

// executeStore 执行STORE指令
func (sc *SmartContract) executeStore(vm *VirtualMachine, inst *Instruction) ([]byte, error) {
    key := fmt.Sprintf("key_%d", inst.Operand.Int64())
    err := vm.Store(key)
    if err != nil {
        return nil, err
    }
    return []byte("success"), nil
}

// executeLoad 执行LOAD指令
func (sc *SmartContract) executeLoad(vm *VirtualMachine, inst *Instruction) ([]byte, error) {
    key := fmt.Sprintf("key_%d", inst.Operand.Int64())
    err := vm.Load(key)
    if err != nil {
        return nil, err
    }
    return []byte("success"), nil
}
```

### 2.5.3 代币合约

```go
// TokenContract 代币合约
type TokenContract struct {
    *SmartContract
    Name     string
    Symbol   string
    Decimals int
    TotalSupply *big.Int
}

// NewTokenContract 创建代币合约
func NewTokenContract(address, name, symbol string, decimals int, totalSupply *big.Int) *TokenContract {
    sc := NewSmartContract(address, []byte{}, "owner")
    
    return &TokenContract{
        SmartContract: sc,
        Name:          name,
        Symbol:        symbol,
        Decimals:      decimals,
        TotalSupply:   totalSupply,
    }
}

// Transfer 转账
func (tc *TokenContract) Transfer(from, to string, amount *big.Int) error {
    // 检查余额
    fromKey := fmt.Sprintf("balance_%s", from)
    fromBalance, exists := tc.Storage[fromKey]
    if !exists {
        fromBalance = big.NewInt(0)
    }
    
    if fromBalance.Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance")
    }
    
    // 更新余额
    tc.Storage[fromKey] = new(big.Int).Sub(fromBalance, amount)
    
    toKey := fmt.Sprintf("balance_%s", to)
    toBalance, exists := tc.Storage[toKey]
    if !exists {
        toBalance = big.NewInt(0)
    }
    tc.Storage[toKey] = new(big.Int).Add(toBalance, amount)
    
    return nil
}

// BalanceOf 查询余额
func (tc *TokenContract) BalanceOf(address string) *big.Int {
    key := fmt.Sprintf("balance_%s", address)
    balance, exists := tc.Storage[key]
    if !exists {
        return big.NewInt(0)
    }
    return balance
}

// Mint 铸造代币
func (tc *TokenContract) Mint(to string, amount *big.Int) error {
    // 检查权限
    if to != tc.Owner {
        return fmt.Errorf("unauthorized")
    }
    
    // 增加总供应量
    tc.TotalSupply = new(big.Int).Add(tc.TotalSupply, amount)
    
    // 增加接收者余额
    toKey := fmt.Sprintf("balance_%s", to)
    toBalance, exists := tc.Storage[toKey]
    if !exists {
        toBalance = big.NewInt(0)
    }
    tc.Storage[toKey] = new(big.Int).Add(toBalance, amount)
    
    return nil
}

// Burn 销毁代币
func (tc *TokenContract) Burn(from string, amount *big.Int) error {
    // 检查余额
    fromKey := fmt.Sprintf("balance_%s", from)
    fromBalance, exists := tc.Storage[fromKey]
    if !exists {
        fromBalance = big.NewInt(0)
    }
    
    if fromBalance.Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance")
    }
    
    // 减少余额和总供应量
    tc.Storage[fromKey] = new(big.Int).Sub(fromBalance, amount)
    tc.TotalSupply = new(big.Int).Sub(tc.TotalSupply, amount)
    
    return nil
}
```

### 2.5.4 去中心化应用

```go
// DApp 去中心化应用
type DApp struct {
    Name        string
    Contracts   map[string]*SmartContract
    Users       map[string]*User
    Transactions []*Transaction
}

// User 用户
type User struct {
    Address string
    Balance *big.Int
    Tokens  map[string]*big.Int
}

// NewDApp 创建DApp
func NewDApp(name string) *DApp {
    return &DApp{
        Name:        name,
        Contracts:   make(map[string]*SmartContract),
        Users:       make(map[string]*User),
        Transactions: make([]*Transaction, 0),
    }
}

// DeployContract 部署合约
func (dapp *DApp) DeployContract(name string, contract *SmartContract) {
    dapp.Contracts[name] = contract
}

// AddUser 添加用户
func (dapp *DApp) AddUser(address string, initialBalance *big.Int) {
    dapp.Users[address] = &User{
        Address: address,
        Balance: initialBalance,
        Tokens:  make(map[string]*big.Int),
    }
}

// ExecuteTransaction 执行交易
func (dapp *DApp) ExecuteTransaction(tx *Transaction) error {
    // 验证交易
    if !tx.Verify("public_key") {
        return fmt.Errorf("invalid transaction")
    }
    
    // 检查余额
    user, exists := dapp.Users[tx.From]
    if !exists {
        return fmt.Errorf("user not found")
    }
    
    if user.Balance.Cmp(big.NewInt(int64(tx.Amount))) < 0 {
        return fmt.Errorf("insufficient balance")
    }
    
    // 执行交易
    user.Balance.Sub(user.Balance, big.NewInt(int64(tx.Amount)))
    
    recipient, exists := dapp.Users[tx.To]
    if exists {
        recipient.Balance.Add(recipient.Balance, big.NewInt(int64(tx.Amount)))
    }
    
    // 记录交易
    dapp.Transactions = append(dapp.Transactions, tx)
    
    return nil
}

// GetUserBalance 获取用户余额
func (dapp *DApp) GetUserBalance(address string) *big.Int {
    user, exists := dapp.Users[address]
    if !exists {
        return big.NewInt(0)
    }
    return user.Balance
}

// GetTransactionHistory 获取交易历史
func (dapp *DApp) GetTransactionHistory(address string) []*Transaction {
    var history []*Transaction
    
    for _, tx := range dapp.Transactions {
        if tx.From == address || tx.To == address {
            history = append(history, tx)
        }
    }
    
    return history
}
```

## 2.6 应用示例

### 2.6.1 代币合约示例

```go
// TokenContractExample 代币合约示例
func TokenContractExample() {
    // 创建代币合约
    totalSupply := big.NewInt(1000000)
    token := NewTokenContract("0x123", "MyToken", "MTK", 18, totalSupply)
    
    // 初始化用户
    alice := "0xAlice"
    bob := "0xBob"
    
    // 铸造代币给Alice
    token.Mint(alice, big.NewInt(1000))
    
    // Alice转账给Bob
    token.Transfer(alice, bob, big.NewInt(500))
    
    // 查询余额
    fmt.Printf("Alice balance: %s\n", token.BalanceOf(alice))
    fmt.Printf("Bob balance: %s\n", token.BalanceOf(bob))
    
    // 销毁代币
    token.Burn(alice, big.NewInt(100))
    fmt.Printf("Alice balance after burn: %s\n", token.BalanceOf(alice))
}
```

## 2.7 理论证明

### 2.7.1 智能合约安全性

**定理 2.1** (智能合约安全性)
在正确的访问控制和输入验证下，智能合约是安全的。

**证明**：
通过形式化验证和静态分析，可以证明智能合约的安全性。

### 2.7.2 气体机制正确性

**定理 2.2** (气体机制正确性)
气体机制能够防止无限循环和资源耗尽攻击。

**证明**：
通过分析气体消耗的单调性和有限性，可以证明气体机制的正确性。

## 2.8 总结

智能合约通过自动执行的程序代码，实现了去中心化的自动化交易和业务逻辑。Go语言实现展示了智能合约的核心概念和机制。

---

**参考文献**：

1. Wood, G. (2014). Ethereum: A secure decentralised generalised transaction ledger.
2. Szabo, N. (1996). Smart contracts: Building blocks for digital markets.
3. Buterin, V. (2014). Ethereum: A next-generation smart contract and decentralized application platform.
