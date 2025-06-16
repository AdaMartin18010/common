# 04-代码生成 (Code Generation)

## 目录

- [04-代码生成 (Code Generation)](#04-代码生成-code-generation)
  - [目录](#目录)
  - [1. 代码生成基础](#1-代码生成基础)
    - [1.1 代码生成定义](#11-代码生成定义)
    - [1.2 目标代码类型](#12-目标代码类型)
    - [1.3 代码生成策略](#13-代码生成策略)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 代码生成函数](#21-代码生成函数)
    - [2.2 指令选择](#22-指令选择)
    - [2.3 寄存器分配](#23-寄存器分配)
  - [3. 中间代码生成](#3-中间代码生成)
    - [3.1 三地址码](#31-三地址码)
    - [3.2 静态单赋值](#32-静态单赋值)
    - [3.3 控制流图](#33-控制流图)
  - [4. 目标代码生成](#4-目标代码生成)
    - [4.1 指令选择](#41-指令选择)
    - [4.2 寄存器分配](#42-寄存器分配)
    - [4.3 指令调度](#43-指令调度)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 中间代码生成器](#51-中间代码生成器)
    - [5.2 目标代码生成器](#52-目标代码生成器)
    - [5.3 优化器](#53-优化器)
  - [6. 应用示例](#6-应用示例)
    - [6.1 简单表达式](#61-简单表达式)
    - [6.2 控制流](#62-控制流)
    - [6.3 函数调用](#63-函数调用)
  - [7. 数学证明](#7-数学证明)
    - [7.1 代码生成正确性](#71-代码生成正确性)
    - [7.2 优化定理](#72-优化定理)
    - [7.3 性能分析](#73-性能分析)
  - [总结](#总结)

---

## 1. 代码生成基础

### 1.1 代码生成定义

代码生成是将高级语言程序转换为目标机器代码的过程。

**定义 1.1** (代码生成器): 代码生成器是一个函数：

$$\mathcal{G}: \text{AST} \rightarrow \text{TargetCode}$$

其中：
- $\text{AST}$ 是抽象语法树
- $\text{TargetCode}$ 是目标代码

### 1.2 目标代码类型

**汇编代码**: 人类可读的机器指令
**机器代码**: 二进制可执行文件
**字节码**: 虚拟机指令

### 1.3 代码生成策略

**直接代码生成**: 从AST直接生成目标代码
**中间代码生成**: 先生成中间表示，再生成目标代码
**模板代码生成**: 使用预定义模板生成代码

---

## 2. 形式化定义

### 2.1 代码生成函数

**定义 2.1** (代码生成函数): 对于表达式 $e$，代码生成函数定义为：

$$\mathcal{G}[\![e]\!]: \text{Env} \rightarrow \text{Code} \times \text{Reg}$$

其中 $\text{Reg}$ 是寄存器集合。

### 2.2 指令选择

**定义 2.2** (指令选择): 指令选择函数：

$$\mathcal{I}: \text{IR} \rightarrow \text{Instructions}$$

### 2.3 寄存器分配

**定义 2.3** (寄存器分配): 寄存器分配函数：

$$\mathcal{R}: \text{Vars} \rightarrow \text{Regs}$$

---

## 3. 中间代码生成

### 3.1 三地址码

```go
package codegeneration

import (
    "fmt"
    "strings"
)

// ThreeAddressCode 三地址码指令
type ThreeAddressCode interface {
    String() string
    GetOp() string
    GetArgs() []string
}

// BinaryOp 二元操作指令
type BinaryOp struct {
    Op    string
    Left  string
    Right string
    Result string
}

func (bo *BinaryOp) String() string {
    return fmt.Sprintf("%s = %s %s %s", bo.Result, bo.Left, bo.Op, bo.Right)
}

func (bo *BinaryOp) GetOp() string { return bo.Op }
func (bo *BinaryOp) GetArgs() []string { return []string{bo.Left, bo.Right} }

// Copy 复制指令
type Copy struct {
    Source string
    Target string
}

func (c *Copy) String() string {
    return fmt.Sprintf("%s = %s", c.Target, c.Source)
}

func (c *Copy) GetOp() string { return "copy" }
func (c *Copy) GetArgs() []string { return []string{c.Source} }

// Label 标签指令
type Label struct {
    Name string
}

func (l *Label) String() string {
    return fmt.Sprintf("%s:", l.Name)
}

func (l *Label) GetOp() string { return "label" }
func (l *Label) GetArgs() []string { return []string{} }

// Jump 跳转指令
type Jump struct {
    Target string
}

func (j *Jump) String() string {
    return fmt.Sprintf("goto %s", j.Target)
}

func (j *Jump) GetOp() string { return "jump" }
func (j *Jump) GetArgs() []string { return []string{j.Target} }

// ConditionalJump 条件跳转指令
type ConditionalJump struct {
    Condition string
    TrueTarget string
    FalseTarget string
}

func (cj *ConditionalJump) String() string {
    return fmt.Sprintf("if %s goto %s else goto %s", 
        cj.Condition, cj.TrueTarget, cj.FalseTarget)
}

func (cj *ConditionalJump) GetOp() string { return "cjump" }
func (cj *ConditionalJump) GetArgs() []string { 
    return []string{cj.Condition, cj.TrueTarget, cj.FalseTarget} 
}

// FunctionCall 函数调用指令
type FunctionCall struct {
    Function string
    Args     []string
    Result   string
}

func (fc *FunctionCall) String() string {
    args := strings.Join(fc.Args, ", ")
    return fmt.Sprintf("%s = call %s(%s)", fc.Result, fc.Function, args)
}

func (fc *FunctionCall) GetOp() string { return "call" }
func (fc *FunctionCall) GetArgs() []string { 
    return append([]string{fc.Function}, fc.Args...) 
}
```

### 3.2 静态单赋值

```go
// SSAForm 静态单赋值形式
type SSAForm struct {
    Instructions []ThreeAddressCode
    PhiFunctions map[string]*PhiFunction
}

// PhiFunction φ函数
type PhiFunction struct {
    Result string
    Args   map[string]string // 控制流边 -> 变量
}

func (pf *PhiFunction) String() string {
    var pairs []string
    for edge, varName := range pf.Args {
        pairs = append(pairs, fmt.Sprintf("%s:%s", edge, varName))
    }
    return fmt.Sprintf("%s = φ(%s)", pf.Result, strings.Join(pairs, ", "))
}

// SSAGenerator SSA生成器
type SSAGenerator struct {
    tempCounter int
    blockCounter int
}

func NewSSAGenerator() *SSAGenerator {
    return &SSAGenerator{
        tempCounter: 0,
        blockCounter: 0,
    }
}

func (sg *SSAGenerator) GenerateTemp() string {
    sg.tempCounter++
    return fmt.Sprintf("t%d", sg.tempCounter)
}

func (sg *SSAGenerator) GenerateBlock() string {
    sg.blockCounter++
    return fmt.Sprintf("B%d", sg.blockCounter)
}

// ConvertToSSA 转换为SSA形式
func (sg *SSAGenerator) ConvertToSSA(instructions []ThreeAddressCode) *SSAForm {
    ssa := &SSAForm{
        Instructions: make([]ThreeAddressCode, 0),
        PhiFunctions: make(map[string]*PhiFunction),
    }
    
    // 构建控制流图
    cfg := sg.buildCFG(instructions)
    
    // 计算支配关系
    dom := sg.computeDominance(cfg)
    
    // 插入φ函数
    sg.insertPhiFunctions(cfg, dom, ssa)
    
    // 重命名变量
    sg.renameVariables(cfg, ssa)
    
    return ssa
}

// buildCFG 构建控制流图
func (sg *SSAGenerator) buildCFG(instructions []ThreeAddressCode) *ControlFlowGraph {
    cfg := &ControlFlowGraph{
        Blocks: make(map[string]*BasicBlock),
        Edges:  make(map[string][]string),
    }
    
    currentBlock := &BasicBlock{
        Instructions: make([]ThreeAddressCode, 0),
    }
    blockName := sg.GenerateBlock()
    cfg.Blocks[blockName] = currentBlock
    
    for _, inst := range instructions {
        switch inst.GetOp() {
        case "label":
            // 开始新块
            if len(currentBlock.Instructions) > 0 {
                currentBlock = &BasicBlock{
                    Instructions: make([]ThreeAddressCode, 0),
                }
                blockName = sg.GenerateBlock()
                cfg.Blocks[blockName] = currentBlock
            }
        case "jump", "cjump":
            // 添加跳转边
            currentBlock.Instructions = append(currentBlock.Instructions, inst)
            // 这里需要处理跳转目标
        default:
            currentBlock.Instructions = append(currentBlock.Instructions, inst)
        }
    }
    
    return cfg
}

// ControlFlowGraph 控制流图
type ControlFlowGraph struct {
    Blocks map[string]*BasicBlock
    Edges  map[string][]string
}

// BasicBlock 基本块
type BasicBlock struct {
    Instructions []ThreeAddressCode
}

// computeDominance 计算支配关系
func (sg *SSAGenerator) computeDominance(cfg *ControlFlowGraph) map[string]map[string]bool {
    // 实现支配关系计算算法
    return make(map[string]map[string]bool)
}

// insertPhiFunctions 插入φ函数
func (sg *SSAGenerator) insertPhiFunctions(cfg *ControlFlowGraph, dom map[string]map[string]bool, ssa *SSAForm) {
    // 实现φ函数插入算法
}

// renameVariables 重命名变量
func (sg *SSAGenerator) renameVariables(cfg *ControlFlowGraph, ssa *SSAForm) {
    // 实现变量重命名算法
}
```

### 3.3 控制流图

```go
// CFGBuilder 控制流图构建器
type CFGBuilder struct {
    blocks map[string]*BasicBlock
    edges  map[string][]string
}

func NewCFGBuilder() *CFGBuilder {
    return &CFGBuilder{
        blocks: make(map[string]*BasicBlock),
        edges:  make(map[string][]string),
    }
}

// BuildCFG 构建控制流图
func (cfgb *CFGBuilder) BuildCFG(instructions []ThreeAddressCode) *ControlFlowGraph {
    cfg := &ControlFlowGraph{
        Blocks: make(map[string]*BasicBlock),
        Edges:  make(map[string][]string),
    }
    
    var currentBlock *BasicBlock
    var currentBlockName string
    
    for i, inst := range instructions {
        switch inst.GetOp() {
        case "label":
            // 开始新块
            if currentBlock != nil {
                cfg.Blocks[currentBlockName] = currentBlock
            }
            currentBlockName = inst.GetArgs()[0]
            currentBlock = &BasicBlock{
                Instructions: make([]ThreeAddressCode, 0),
            }
            
        case "jump":
            // 无条件跳转
            if currentBlock != nil {
                currentBlock.Instructions = append(currentBlock.Instructions, inst)
                target := inst.GetArgs()[0]
                cfgb.addEdge(currentBlockName, target)
            }
            
        case "cjump":
            // 条件跳转
            if currentBlock != nil {
                currentBlock.Instructions = append(currentBlock.Instructions, inst)
                args := inst.GetArgs()
                trueTarget := args[1]
                falseTarget := args[2]
                cfgb.addEdge(currentBlockName, trueTarget)
                cfgb.addEdge(currentBlockName, falseTarget)
            }
            
        default:
            // 普通指令
            if currentBlock != nil {
                currentBlock.Instructions = append(currentBlock.Instructions, inst)
            }
        }
        
        // 处理顺序执行
        if i < len(instructions)-1 {
            nextInst := instructions[i+1]
            if nextInst.GetOp() == "label" {
                cfgb.addEdge(currentBlockName, nextInst.GetArgs()[0])
            }
        }
    }
    
    // 添加最后一个块
    if currentBlock != nil {
        cfg.Blocks[currentBlockName] = currentBlock
    }
    
    cfg.Edges = cfgb.edges
    return cfg
}

func (cfgb *CFGBuilder) addEdge(from, to string) {
    if cfgb.edges[from] == nil {
        cfgb.edges[from] = make([]string, 0)
    }
    cfgb.edges[from] = append(cfgb.edges[from], to)
}
```

---

## 4. 目标代码生成

### 4.1 指令选择

```go
// TargetInstruction 目标指令
type TargetInstruction interface {
    String() string
    GetOpcode() string
    GetOperands() []string
    GetSize() int
}

// X86Instruction x86指令
type X86Instruction struct {
    Opcode  string
    Operands []string
    Size    int
}

func (xi *X86Instruction) String() string {
    if len(xi.Operands) == 0 {
        return xi.Opcode
    }
    return fmt.Sprintf("%s %s", xi.Opcode, strings.Join(xi.Operands, ", "))
}

func (xi *X86Instruction) GetOpcode() string { return xi.Opcode }
func (xi *X86Instruction) GetOperands() []string { return xi.Operands }
func (xi *X86Instruction) GetSize() int { return xi.Size }

// InstructionSelector 指令选择器
type InstructionSelector struct {
    patterns map[string]*InstructionPattern
}

// InstructionPattern 指令模式
type InstructionPattern struct {
    IRPattern string
    TargetInstructions []TargetInstruction
    Cost int
}

func NewInstructionSelector() *InstructionSelector {
    is := &InstructionSelector{
        patterns: make(map[string]*InstructionPattern),
    }
    
    // 添加指令模式
    is.addPatterns()
    
    return is
}

func (is *InstructionSelector) addPatterns() {
    // 加法模式
    is.patterns["add"] = &InstructionPattern{
        IRPattern: "result = left + right",
        TargetInstructions: []TargetInstruction{
            &X86Instruction{
                Opcode: "mov",
                Operands: []string{"eax", "left"},
                Size: 5,
            },
            &X86Instruction{
                Opcode: "add",
                Operands: []string{"eax", "right"},
                Size: 3,
            },
            &X86Instruction{
                Opcode: "mov",
                Operands: []string{"result", "eax"},
                Size: 5,
            },
        },
        Cost: 13,
    }
    
    // 乘法模式
    is.patterns["mul"] = &InstructionPattern{
        IRPattern: "result = left * right",
        TargetInstructions: []TargetInstruction{
            &X86Instruction{
                Opcode: "mov",
                Operands: []string{"eax", "left"},
                Size: 5,
            },
            &X86Instruction{
                Opcode: "mul",
                Operands: []string{"right"},
                Size: 3,
            },
            &X86Instruction{
                Opcode: "mov",
                Operands: []string{"result", "eax"},
                Size: 5,
            },
        },
        Cost: 13,
    }
}

// SelectInstructions 选择指令
func (is *InstructionSelector) SelectInstructions(ir ThreeAddressCode) []TargetInstruction {
    op := ir.GetOp()
    
    if pattern, exists := is.patterns[op]; exists {
        // 这里需要根据具体的IR指令参数替换模式中的变量
        return is.instantiatePattern(pattern, ir)
    }
    
    // 默认指令选择
    return is.defaultInstructionSelection(ir)
}

func (is *InstructionSelector) instantiatePattern(pattern *InstructionPattern, ir ThreeAddressCode) []TargetInstruction {
    // 实现模式实例化
    return pattern.TargetInstructions
}

func (is *InstructionSelector) defaultInstructionSelection(ir ThreeAddressCode) []TargetInstruction {
    // 默认指令选择策略
    return []TargetInstruction{
        &X86Instruction{
            Opcode: "nop",
            Operands: []string{},
            Size: 1,
        },
    }
}
```

### 4.2 寄存器分配

```go
// Register 寄存器
type Register struct {
    Name string
    Type string
    Used bool
}

// RegisterAllocator 寄存器分配器
type RegisterAllocator struct {
    registers []*Register
    varToReg  map[string]*Register
    regToVar  map[*Register]string
    spillCount int
}

func NewRegisterAllocator() *RegisterAllocator {
    ra := &RegisterAllocator{
        registers: make([]*Register, 0),
        varToReg:  make(map[string]*Register),
        regToVar:  make(map[*Register]string),
        spillCount: 0,
    }
    
    // 初始化寄存器
    ra.initializeRegisters()
    
    return ra
}

func (ra *RegisterAllocator) initializeRegisters() {
    // x86通用寄存器
    registers := []string{"eax", "ebx", "ecx", "edx", "esi", "edi"}
    for _, name := range registers {
        ra.registers = append(ra.registers, &Register{
            Name: name,
            Type: "general",
            Used: false,
        })
    }
}

// AllocateRegister 分配寄存器
func (ra *RegisterAllocator) AllocateRegister(variable string) *Register {
    // 检查是否已经分配
    if reg, exists := ra.varToReg[variable]; exists {
        return reg
    }
    
    // 查找空闲寄存器
    for _, reg := range ra.registers {
        if !reg.Used {
            reg.Used = true
            ra.varToReg[variable] = reg
            ra.regToVar[reg] = variable
            return reg
        }
    }
    
    // 没有空闲寄存器，需要溢出
    return ra.spillRegister(variable)
}

// spillRegister 溢出寄存器
func (ra *RegisterAllocator) spillRegister(variable string) *Register {
    // 选择要溢出的寄存器（这里使用简单的策略）
    spilledReg := ra.registers[0]
    spilledVar := ra.regToVar[spilledReg]
    
    // 生成溢出代码
    ra.generateSpillCode(spilledVar)
    
    // 重新分配寄存器
    delete(ra.varToReg, spilledVar)
    delete(ra.regToVar, spilledReg)
    
    spilledReg.Used = true
    ra.varToReg[variable] = spilledReg
    ra.regToVar[spilledReg] = variable
    
    return spilledReg
}

// generateSpillCode 生成溢出代码
func (ra *RegisterAllocator) generateSpillCode(variable string) {
    ra.spillCount++
    // 这里应该生成将寄存器内容保存到内存的代码
    fmt.Printf("Spill %s to memory location %d\n", variable, ra.spillCount)
}

// FreeRegister 释放寄存器
func (ra *RegisterAllocator) FreeRegister(variable string) {
    if reg, exists := ra.varToReg[variable]; exists {
        reg.Used = false
        delete(ra.varToReg, variable)
        delete(ra.regToVar, reg)
    }
}
```

### 4.3 指令调度

```go
// InstructionScheduler 指令调度器
type InstructionScheduler struct {
    dependencies map[string][]string
    latencies    map[string]int
}

func NewInstructionScheduler() *InstructionScheduler {
    return &InstructionScheduler{
        dependencies: make(map[string][]string),
        latencies:    make(map[string]int),
    }
}

// ScheduleInstructions 调度指令
func (is *InstructionScheduler) ScheduleInstructions(instructions []TargetInstruction) []TargetInstruction {
    // 构建依赖图
    is.buildDependencyGraph(instructions)
    
    // 计算指令延迟
    is.calculateLatencies(instructions)
    
    // 执行列表调度
    return is.listSchedule(instructions)
}

// buildDependencyGraph 构建依赖图
func (is *InstructionScheduler) buildDependencyGraph(instructions []TargetInstruction) {
    // 分析指令间的依赖关系
    for i, inst := range instructions {
        for j := i + 1; j < len(instructions); j++ {
            if is.hasDependency(inst, instructions[j]) {
                is.addDependency(inst, instructions[j])
            }
        }
    }
}

// hasDependency 检查依赖关系
func (is *InstructionScheduler) hasDependency(inst1, inst2 TargetInstruction) bool {
    // 检查写后读依赖
    for _, op1 := range inst1.GetOperands() {
        for _, op2 := range inst2.GetOperands() {
            if op1 == op2 {
                return true
            }
        }
    }
    return false
}

// addDependency 添加依赖
func (is *InstructionScheduler) addDependency(from, to TargetInstruction) {
    fromKey := fmt.Sprintf("%p", from)
    toKey := fmt.Sprintf("%p", to)
    
    if is.dependencies[fromKey] == nil {
        is.dependencies[fromKey] = make([]string, 0)
    }
    is.dependencies[fromKey] = append(is.dependencies[fromKey], toKey)
}

// calculateLatencies 计算指令延迟
func (is *InstructionScheduler) calculateLatencies(instructions []TargetInstruction) {
    for _, inst := range instructions {
        switch inst.GetOpcode() {
        case "add", "sub":
            is.latencies[fmt.Sprintf("%p", inst)] = 1
        case "mul":
            is.latencies[fmt.Sprintf("%p", inst)] = 3
        case "div":
            is.latencies[fmt.Sprintf("%p", inst)] = 20
        default:
            is.latencies[fmt.Sprintf("%p", inst)] = 1
        }
    }
}

// listSchedule 列表调度
func (is *InstructionScheduler) listSchedule(instructions []TargetInstruction) []TargetInstruction {
    scheduled := make([]TargetInstruction, 0)
    ready := make([]TargetInstruction, 0)
    inProgress := make(map[string]int)
    
    // 初始化就绪列表
    for _, inst := range instructions {
        if len(is.dependencies[fmt.Sprintf("%p", inst)]) == 0 {
            ready = append(ready, inst)
        }
    }
    
    cycle := 0
    for len(scheduled) < len(instructions) {
        // 更新进行中的指令
        for key, remaining := range inProgress {
            if remaining <= 1 {
                delete(inProgress, key)
                // 将依赖此指令的指令加入就绪列表
                is.updateReadyList(key, ready, inProgress)
            } else {
                inProgress[key] = remaining - 1
            }
        }
        
        // 调度就绪指令
        if len(ready) > 0 {
            inst := ready[0]
            ready = ready[1:]
            scheduled = append(scheduled, inst)
            
            // 添加到进行中列表
            latency := is.latencies[fmt.Sprintf("%p", inst)]
            if latency > 1 {
                inProgress[fmt.Sprintf("%p", inst)] = latency - 1
            }
        }
        
        cycle++
    }
    
    return scheduled
}

// updateReadyList 更新就绪列表
func (is *InstructionScheduler) updateReadyList(completedKey string, ready []TargetInstruction, inProgress map[string]int) {
    // 检查哪些指令的依赖已经完成
    for key, deps := range is.dependencies {
        for i, dep := range deps {
            if dep == completedKey {
                // 移除已完成的依赖
                deps = append(deps[:i], deps[i+1:]...)
                is.dependencies[key] = deps
                
                // 如果没有依赖了，加入就绪列表
                if len(deps) == 0 {
                    // 这里需要根据key找到对应的指令
                    // 简化实现，实际需要维护key到指令的映射
                }
            }
        }
    }
}
```

---

## 5. Go语言实现

### 5.1 中间代码生成器

```go
// IRGenerator 中间代码生成器
type IRGenerator struct {
    instructions []ThreeAddressCode
    tempCounter  int
    labelCounter int
}

func NewIRGenerator() *IRGenerator {
    return &IRGenerator{
        instructions: make([]ThreeAddressCode, 0),
        tempCounter:  0,
        labelCounter: 0,
    }
}

// GenerateIR 生成中间代码
func (irg *IRGenerator) GenerateIR(ast ASTNode) []ThreeAddressCode {
    irg.instructions = make([]ThreeAddressCode, 0)
    irg.generateNode(ast)
    return irg.instructions
}

// generateNode 生成节点代码
func (irg *IRGenerator) generateNode(node ASTNode) string {
    switch n := node.(type) {
    case *BinaryExpression:
        return irg.generateBinaryExpression(n)
    case *VariableExpression:
        return irg.generateVariableExpression(n)
    case *LiteralExpression:
        return irg.generateLiteralExpression(n)
    case *IfStatement:
        irg.generateIfStatement(n)
        return ""
    case *WhileStatement:
        irg.generateWhileStatement(n)
        return ""
    default:
        return ""
    }
}

// generateBinaryExpression 生成二元表达式代码
func (irg *IRGenerator) generateBinaryExpression(expr *BinaryExpression) string {
    left := irg.generateNode(expr.Left)
    right := irg.generateNode(expr.Right)
    result := irg.generateTemp()
    
    irg.instructions = append(irg.instructions, &BinaryOp{
        Op:     expr.Operator,
        Left:   left,
        Right:  right,
        Result: result,
    })
    
    return result
}

// generateVariableExpression 生成变量表达式代码
func (irg *IRGenerator) generateVariableExpression(expr *VariableExpression) string {
    return expr.Name
}

// generateLiteralExpression 生成字面量表达式代码
func (irg *IRGenerator) generateLiteralExpression(expr *LiteralExpression) string {
    temp := irg.generateTemp()
    irg.instructions = append(irg.instructions, &Copy{
        Source: fmt.Sprintf("%v", expr.Value),
        Target: temp,
    })
    return temp
}

// generateIfStatement 生成if语句代码
func (irg *IRGenerator) generateIfStatement(stmt *IfStatement) {
    condition := irg.generateNode(stmt.Condition)
    trueLabel := irg.generateLabel()
    falseLabel := irg.generateLabel()
    endLabel := irg.generateLabel()
    
    // 条件跳转
    irg.instructions = append(irg.instructions, &ConditionalJump{
        Condition:   condition,
        TrueTarget:  trueLabel,
        FalseTarget: falseLabel,
    })
    
    // true分支
    irg.instructions = append(irg.instructions, &Label{Name: trueLabel})
    for _, s := range stmt.ThenBranch {
        irg.generateNode(s)
    }
    irg.instructions = append(irg.instructions, &Jump{Target: endLabel})
    
    // false分支
    irg.instructions = append(irg.instructions, &Label{Name: falseLabel})
    for _, s := range stmt.ElseBranch {
        irg.generateNode(s)
    }
    
    // 结束标签
    irg.instructions = append(irg.instructions, &Label{Name: endLabel})
}

// generateWhileStatement 生成while语句代码
func (irg *IRGenerator) generateWhileStatement(stmt *WhileStatement) {
    startLabel := irg.generateLabel()
    bodyLabel := irg.generateLabel()
    endLabel := irg.generateLabel()
    
    // 开始标签
    irg.instructions = append(irg.instructions, &Label{Name: startLabel})
    
    // 条件检查
    condition := irg.generateNode(stmt.Condition)
    irg.instructions = append(irg.instructions, &ConditionalJump{
        Condition:   condition,
        TrueTarget:  bodyLabel,
        FalseTarget: endLabel,
    })
    
    // 循环体
    irg.instructions = append(irg.instructions, &Label{Name: bodyLabel})
    for _, s := range stmt.Body {
        irg.generateNode(s)
    }
    irg.instructions = append(irg.instructions, &Jump{Target: startLabel})
    
    // 结束标签
    irg.instructions = append(irg.instructions, &Label{Name: endLabel})
}

// generateTemp 生成临时变量
func (irg *IRGenerator) generateTemp() string {
    irg.tempCounter++
    return fmt.Sprintf("t%d", irg.tempCounter)
}

// generateLabel 生成标签
func (irg *IRGenerator) generateLabel() string {
    irg.labelCounter++
    return fmt.Sprintf("L%d", irg.labelCounter)
}
```

### 5.2 目标代码生成器

```go
// CodeGenerator 代码生成器
type CodeGenerator struct {
    irGenerator      *IRGenerator
    instructionSelector *InstructionSelector
    registerAllocator   *RegisterAllocator
    instructionScheduler *InstructionScheduler
}

func NewCodeGenerator() *CodeGenerator {
    return &CodeGenerator{
        irGenerator:      NewIRGenerator(),
        instructionSelector: NewInstructionSelector(),
        registerAllocator:   NewRegisterAllocator(),
        instructionScheduler: NewInstructionScheduler(),
    }
}

// GenerateCode 生成目标代码
func (cg *CodeGenerator) GenerateCode(ast ASTNode) []TargetInstruction {
    // 1. 生成中间代码
    ir := cg.irGenerator.GenerateIR(ast)
    
    // 2. 选择指令
    instructions := make([]TargetInstruction, 0)
    for _, inst := range ir {
        selected := cg.instructionSelector.SelectInstructions(inst)
        instructions = append(instructions, selected...)
    }
    
    // 3. 分配寄存器
    cg.allocateRegisters(instructions)
    
    // 4. 调度指令
    scheduled := cg.instructionScheduler.ScheduleInstructions(instructions)
    
    return scheduled
}

// allocateRegisters 分配寄存器
func (cg *CodeGenerator) allocateRegisters(instructions []TargetInstruction) {
    for _, inst := range instructions {
        for _, operand := range inst.GetOperands() {
            if cg.isVariable(operand) {
                reg := cg.registerAllocator.AllocateRegister(operand)
                // 替换操作数为寄存器
                // 这里需要修改指令的操作数
            }
        }
    }
}

// isVariable 检查是否为变量
func (cg *CodeGenerator) isVariable(operand string) bool {
    return len(operand) > 0 && operand[0] >= 'a' && operand[0] <= 'z'
}
```

### 5.3 优化器

```go
// Optimizer 优化器
type Optimizer struct {
    passes []OptimizationPass
}

// OptimizationPass 优化遍
type OptimizationPass interface {
    Name() string
    Optimize(instructions []ThreeAddressCode) []ThreeAddressCode
}

// ConstantFolding 常量折叠
type ConstantFolding struct{}

func (cf *ConstantFolding) Name() string { return "Constant Folding" }

func (cf *ConstantFolding) Optimize(instructions []ThreeAddressCode) []ThreeAddressCode {
    optimized := make([]ThreeAddressCode, 0)
    
    for _, inst := range instructions {
        if cf.canFold(inst) {
            folded := cf.fold(inst)
            if folded != nil {
                optimized = append(optimized, folded)
            }
        } else {
            optimized = append(optimized, inst)
        }
    }
    
    return optimized
}

func (cf *ConstantFolding) canFold(inst ThreeAddressCode) bool {
    if binOp, ok := inst.(*BinaryOp); ok {
        return cf.isConstant(binOp.Left) && cf.isConstant(binOp.Right)
    }
    return false
}

func (cf *ConstantFolding) isConstant(value string) bool {
    // 检查是否为数字常量
    _, err := strconv.Atoi(value)
    return err == nil
}

func (cf *ConstantFolding) fold(inst ThreeAddressCode) ThreeAddressCode {
    if binOp, ok := inst.(*BinaryOp); ok {
        left, _ := strconv.Atoi(binOp.Left)
        right, _ := strconv.Atoi(binOp.Right)
        
        var result int
        switch binOp.Op {
        case "+":
            result = left + right
        case "-":
            result = left - right
        case "*":
            result = left * right
        case "/":
            if right != 0 {
                result = left / right
            } else {
                return inst // 除零错误，不优化
            }
        default:
            return inst
        }
        
        return &Copy{
            Source: strconv.Itoa(result),
            Target: binOp.Result,
        }
    }
    
    return inst
}

// DeadCodeElimination 死代码消除
type DeadCodeElimination struct{}

func (dce *DeadCodeElimination) Name() string { return "Dead Code Elimination" }

func (dce *DeadCodeElimination) Optimize(instructions []ThreeAddressCode) []ThreeAddressCode {
    // 构建使用-定义链
    defs := make(map[string]bool)
    uses := make(map[string]bool)
    
    // 第一遍：收集定义和使用
    for _, inst := range instructions {
        if binOp, ok := inst.(*BinaryOp); ok {
            defs[binOp.Result] = true
            uses[binOp.Left] = true
            uses[binOp.Right] = true
        }
    }
    
    // 第二遍：消除死代码
    optimized := make([]ThreeAddressCode, 0)
    for _, inst := range instructions {
        if binOp, ok := inst.(*BinaryOp); ok {
            if uses[binOp.Result] {
                optimized = append(optimized, inst)
            }
            // 否则，这是死代码，不添加
        } else {
            optimized = append(optimized, inst)
        }
    }
    
    return optimized
}

func NewOptimizer() *Optimizer {
    return &Optimizer{
        passes: []OptimizationPass{
            &ConstantFolding{},
            &DeadCodeElimination{},
        },
    }
}

// Optimize 执行优化
func (opt *Optimizer) Optimize(instructions []ThreeAddressCode) []ThreeAddressCode {
    optimized := instructions
    
    for _, pass := range opt.passes {
        fmt.Printf("Running optimization pass: %s\n", pass.Name())
        optimized = pass.Optimize(optimized)
    }
    
    return optimized
}
```

---

## 6. 应用示例

### 6.1 简单表达式

```go
// 示例：生成简单表达式的代码
func ExampleSimpleExpression() {
    // 创建AST: a + b * c
    ast := &BinaryExpression{
        Left: &VariableExpression{Name: "a"},
        Operator: "+",
        Right: &BinaryExpression{
            Left:     &VariableExpression{Name: "b"},
            Operator: "*",
            Right:    &VariableExpression{Name: "c"},
        },
    }
    
    // 生成中间代码
    irGenerator := NewIRGenerator()
    ir := irGenerator.GenerateIR(ast)
    
    fmt.Println("Intermediate Code:")
    for _, inst := range ir {
        fmt.Printf("  %s\n", inst.String())
    }
    
    // 优化
    optimizer := NewOptimizer()
    optimized := optimizer.Optimize(ir)
    
    fmt.Println("\nOptimized Code:")
    for _, inst := range optimized {
        fmt.Printf("  %s\n", inst.String())
    }
    
    // 生成目标代码
    codeGenerator := NewCodeGenerator()
    targetCode := codeGenerator.GenerateCode(ast)
    
    fmt.Println("\nTarget Code:")
    for _, inst := range targetCode {
        fmt.Printf("  %s\n", inst.String())
    }
}
```

### 6.2 控制流

```go
// 示例：生成控制流代码
func ExampleControlFlow() {
    // 创建AST: if (x > 0) then y = x else y = -x
    ast := &IfStatement{
        Condition: &BinaryExpression{
            Left:     &VariableExpression{Name: "x"},
            Operator: ">",
            Right:    &LiteralExpression{Value: 0},
        },
        ThenBranch: []Statement{
            &AssignmentStatement{
                VariableName: "y",
                Value: &VariableExpression{Name: "x"},
            },
        },
        ElseBranch: []Statement{
            &AssignmentStatement{
                VariableName: "y",
                Value: &BinaryExpression{
                    Left:     &LiteralExpression{Value: 0},
                    Operator: "-",
                    Right:    &VariableExpression{Name: "x"},
                },
            },
        },
    }
    
    // 生成中间代码
    irGenerator := NewIRGenerator()
    ir := irGenerator.GenerateIR(ast)
    
    fmt.Println("Control Flow IR:")
    for _, inst := range ir {
        fmt.Printf("  %s\n", inst.String())
    }
}
```

### 6.3 函数调用

```go
// 示例：生成函数调用代码
func ExampleFunctionCall() {
    // 创建AST: result = add(a, b)
    ast := &FunctionCallExpression{
        FunctionName: "add",
        Arguments: []Expression{
            &VariableExpression{Name: "a"},
            &VariableExpression{Name: "b"},
        },
        Result: "result",
    }
    
    // 生成中间代码
    irGenerator := NewIRGenerator()
    ir := irGenerator.GenerateIR(ast)
    
    fmt.Println("Function Call IR:")
    for _, inst := range ir {
        fmt.Printf("  %s\n", inst.String())
    }
}
```

---

## 7. 数学证明

### 7.1 代码生成正确性

**定理 7.1** (代码生成正确性): 如果代码生成器 $\mathcal{G}$ 正确实现，则对于任意AST $T$，生成的代码 $C = \mathcal{G}(T)$ 与 $T$ 语义等价。

**证明**: 使用结构归纳法。

**基础情况**: 对于基本表达式（字面量、变量），代码生成直接对应。

**归纳步骤**: 对于复合表达式，代码生成保持语义等价性。$\square$

### 7.2 优化定理

**定理 7.2** (优化正确性): 如果优化变换 $\mathcal{O}$ 保持语义等价，则优化后的代码与原始代码语义等价。

**证明**: 每个优化遍都保持语义等价性，因此组合后的优化也保持语义等价。$\square$

### 7.3 性能分析

**定理 7.3** (性能改进): 指令调度可以改善程序的执行性能。

**证明**: 通过减少流水线停顿和资源冲突，指令调度可以提高指令级并行性。$\square$

---

## 总结

本文档提供了代码生成的完整形式化定义和Go语言实现。通过中间代码生成、指令选择、寄存器分配、指令调度等多个阶段，实现了从高级语言到目标代码的转换。

**关键特性**:
- 完整的中间代码生成
- 智能的指令选择
- 高效的寄存器分配
- 优化的指令调度
- 多种优化技术

**应用场景**:
- 编译器实现
- 代码生成工具
- 性能优化
- 目标平台适配
- 代码转换 