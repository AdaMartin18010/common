# 11.7.4 区块链共识算法

## 11.7.4.1 概述

区块链共识算法是分布式系统中确保节点间一致性的核心机制，通过数学和密码学原理实现去中心化的信任机制。

### 11.7.4.1.1 基本概念

**定义 11.7.4.1** (共识算法)
共识算法是分布式系统中多个节点就某个值达成一致的协议，满足：

- **安全性**: 诚实节点不会接受无效区块
- **活性**: 诚实节点最终会接受有效区块
- **一致性**: 所有诚实节点最终达成相同状态

### 11.7.4.1.2 共识算法分类

```go
// 共识算法类型枚举
type ConsensusType int

const (
    ProofOfWork ConsensusType = iota
    ProofOfStake
    DelegatedProofOfStake
    PracticalByzantineFaultTolerance
    Raft
    Paxos
)
```

## 11.7.4.2 工作量证明 (Proof of Work)

### 11.7.4.2.1 理论基础

**定理 11.7.4.1** (PoW安全性)
在PoW共识中，攻击者需要控制超过50%的算力才能成功进行双花攻击。

**证明**:
设攻击者算力为 $p$，诚实节点算力为 $1-p$。
攻击成功的概率为：
$$P_{attack} = \left(\frac{p}{1-p}\right)^z$$
其中 $z$ 为确认区块数。

当 $p < 0.5$ 时，$\lim_{z \to \infty} P_{attack} = 0$。

### 11.7.4.2.2 数学建模

**定义 11.7.4.2** (挖矿难度)
挖矿难度 $D$ 定义为：
$$D = \frac{2^{256}}{T}$$
其中 $T$ 为目标阈值。

**定义 11.7.4.3** (期望出块时间)
期望出块时间 $E[T]$ 为：
$$E[T] = \frac{D \cdot 2^{32}}{H}$$
其中 $H$ 为全网算力。

### 11.7.4.2.3 Go实现

```go
// PoW共识实现
type ProofOfWork struct {
    difficulty    uint64
    target        *big.Int
    blockTime     time.Duration
    miningReward  uint64
}

// 创建PoW共识
func NewProofOfWork(difficulty uint64, blockTime time.Duration) *ProofOfWork {
    target := new(big.Int)
    target.SetString("1", 16)
    target.Lsh(target, uint(256-difficulty))
    
    return &ProofOfWork{
        difficulty:   difficulty,
        target:       target,
        blockTime:    blockTime,
        miningReward: 50,
    }
}

// 挖矿过程
func (pow *ProofOfWork) Mine(block *Block) (uint64, []byte) {
    var nonce uint64
    target := pow.target
    
    for {
        block.Nonce = nonce
        hash := block.CalculateHash()
        
        hashInt := new(big.Int)
        hashInt.SetBytes(hash)
        
        if hashInt.Cmp(target) == -1 {
            return nonce, hash
        }
        
        nonce++
    }
}

// 验证区块
func (pow *ProofOfWork) ValidateBlock(block *Block) bool {
    hash := block.CalculateHash()
    hashInt := new(big.Int)
    hashInt.SetBytes(hash)
    
    return hashInt.Cmp(pow.target) == -1
}

// 难度调整
func (pow *ProofOfWork) AdjustDifficulty(actualTime time.Duration) {
    if actualTime < pow.blockTime/2 {
        pow.difficulty++
    } else if actualTime > pow.blockTime*2 {
        pow.difficulty--
    }
    
    // 重新计算目标
    pow.target = new(big.Int)
    pow.target.SetString("1", 16)
    pow.target.Lsh(pow.target, uint(256-pow.difficulty))
}
```

## 11.7.4.3 权益证明 (Proof of Stake)

### 11.7.4.3.1 理论基础

**定义 11.7.4.4** (权益证明)
权益证明是一种共识机制，其中验证者被选中的概率与其持有的代币数量成正比。

**定理 11.7.4.2** (PoS经济安全性)
在PoS中，攻击成本等于攻击者需要质押的代币价值。

**证明**:
设攻击者质押代币价值为 $V$，攻击成功概率为 $p$。
期望损失为：
$$E[Loss] = V \cdot p$$
当 $p$ 很小时，攻击成本接近 $V$。

### 11.7.4.3.2 数学建模

**定义 11.7.4.5** (验证者选择概率)
验证者 $i$ 被选中的概率为：
$$P_i = \frac{s_i}{\sum_{j=1}^n s_j}$$
其中 $s_i$ 为验证者 $i$ 的质押数量。

**定义 11.7.4.6** (惩罚机制)
恶意行为的惩罚函数为：
$$Penalty(s) = \alpha \cdot s + \beta \cdot s^2$$
其中 $\alpha, \beta$ 为惩罚参数。

### 11.7.4.3.3 Go实现

```go
// PoS共识实现
type ProofOfStake struct {
    validators    map[string]*Validator
    totalStake    uint64
    minStake      uint64
    blockTime     time.Duration
    rewardRate    float64
    penaltyRate   float64
}

// 验证者结构
type Validator struct {
    Address     string
    Stake       uint64
    Commission  float64
    Delegators  map[string]uint64
    LastBlock   time.Time
    Reputation  float64
}

// 创建PoS共识
func NewProofOfStake(minStake uint64, blockTime time.Duration) *ProofOfStake {
    return &ProofOfStake{
        validators:  make(map[string]*Validator),
        minStake:    minStake,
        blockTime:   blockTime,
        rewardRate:  0.05,  // 5%年化收益率
        penaltyRate: 0.1,   // 10%惩罚率
    }
}

// 添加验证者
func (pos *ProofOfStake) AddValidator(address string, stake uint64) error {
    if stake < pos.minStake {
        return fmt.Errorf("stake below minimum requirement")
    }
    
    pos.validators[address] = &Validator{
        Address:    address,
        Stake:      stake,
        Commission: 0.1,
        Delegators: make(map[string]uint64),
        Reputation: 1.0,
    }
    
    pos.totalStake += stake
    return nil
}

// 选择验证者
func (pos *ProofOfStake) SelectValidator() string {
    if len(pos.validators) == 0 {
        return ""
    }
    
    // 基于权益和声誉的加权随机选择
    totalWeight := 0.0
    weights := make(map[string]float64)
    
    for addr, validator := range pos.validators {
        weight := float64(validator.Stake) * validator.Reputation
        weights[addr] = weight
        totalWeight += weight
    }
    
    // 随机选择
    rand.Seed(time.Now().UnixNano())
    r := rand.Float64() * totalWeight
    
    currentWeight := 0.0
    for addr, weight := range weights {
        currentWeight += weight
        if r <= currentWeight {
            return addr
        }
    }
    
    // 返回第一个验证者（理论上不会到达这里）
    for addr := range pos.validators {
        return addr
    }
    return ""
}

// 验证区块
func (pos *ProofOfStake) ValidateBlock(block *Block, validator string) bool {
    v, exists := pos.validators[validator]
    if !exists {
        return false
    }
    
    // 检查验证者是否有资格
    if time.Since(v.LastBlock) < pos.blockTime {
        return false
    }
    
    // 验证区块内容
    if !block.IsValid() {
        // 惩罚恶意验证者
        pos.penalizeValidator(validator)
        return false
    }
    
    // 更新验证者状态
    v.LastBlock = time.Now()
    pos.rewardValidator(validator)
    
    return true
}

// 奖励验证者
func (pos *ProofOfStake) rewardValidator(validator string) {
    v := pos.validators[validator]
    reward := uint64(float64(v.Stake) * pos.rewardRate / 365 / 24 / 60) // 每分钟奖励
    
    v.Stake += reward
    pos.totalStake += reward
}

// 惩罚验证者
func (pos *ProofOfStake) penalizeValidator(validator string) {
    v := pos.validators[validator]
    penalty := uint64(float64(v.Stake) * pos.penaltyRate)
    
    if penalty > v.Stake {
        penalty = v.Stake
    }
    
    v.Stake -= penalty
    pos.totalStake -= penalty
    v.Reputation *= 0.9 // 降低声誉
}

// 委托代币
func (pos *ProofOfStake) Delegate(delegator, validator string, amount uint64) error {
    v, exists := pos.validators[validator]
    if !exists {
        return fmt.Errorf("validator not found")
    }
    
    v.Delegators[delegator] += amount
    v.Stake += amount
    pos.totalStake += amount
    
    return nil
}
```

## 11.7.4.4 委托权益证明 (Delegated Proof of Stake)

### 11.7.4.4.1 理论基础

**定义 11.7.4.7** (DPoS)
委托权益证明是PoS的变种，其中代币持有者通过投票选举出有限数量的验证者。

**定理 11.7.4.3** (DPoS效率)
DPoS通过减少验证者数量，将共识复杂度从 $O(n^2)$ 降低到 $O(n)$。

### 11.7.4.4.2 数学建模

**定义 11.7.4.8** (投票权重)
投票权重函数为：
$$W(v_i) = \sum_{j=1}^m w_j \cdot vote_{i,j}$$
其中 $w_j$ 为选民 $j$ 的权重，$vote_{i,j}$ 为选民 $j$ 对验证者 $i$ 的投票。

### 11.7.4.4.3 Go实现

```go
// DPoS共识实现
type DelegatedProofOfStake struct {
    validators    []*DPoSValidator
    voters        map[string]*Voter
    maxValidators int
    blockTime     time.Duration
    votingPeriod  time.Duration
}

// DPoS验证者
type DPoSValidator struct {
    Address     string
    Votes       uint64
    Rank        int
    IsActive    bool
    LastBlock   time.Time
    Performance float64
}

// 选民
type Voter struct {
    Address string
    Balance uint64
    Votes   map[string]uint64 // 验证者地址 -> 投票数量
}

// 创建DPoS共识
func NewDPoS(maxValidators int, blockTime time.Duration) *DelegatedProofOfStake {
    return &DelegatedProofOfStake{
        validators:    make([]*DPoSValidator, 0),
        voters:        make(map[string]*Voter),
        maxValidators: maxValidators,
        blockTime:     blockTime,
        votingPeriod:  24 * time.Hour,
    }
}

// 注册验证者
func (dpos *DelegatedProofOfStake) RegisterValidator(address string) {
    validator := &DPoSValidator{
        Address:     address,
        Votes:       0,
        Rank:        0,
        IsActive:    false,
        Performance: 1.0,
    }
    
    dpos.validators = append(dpos.validators, validator)
}

// 投票
func (dpos *DelegatedProofOfStake) Vote(voterAddr, validatorAddr string, amount uint64) error {
    voter, exists := dpos.voters[voterAddr]
    if !exists {
        voter = &Voter{
            Address: voterAddr,
            Balance: 0,
            Votes:   make(map[string]uint64),
        }
        dpos.voters[voterAddr] = voter
    }
    
    // 检查余额
    if voter.Balance < amount {
        return fmt.Errorf("insufficient balance")
    }
    
    // 更新投票
    voter.Balance -= amount
    voter.Votes[validatorAddr] += amount
    
    // 更新验证者票数
    for _, validator := range dpos.validators {
        if validator.Address == validatorAddr {
            validator.Votes += amount
            break
        }
    }
    
    return nil
}

// 选择活跃验证者
func (dpos *DelegatedProofOfStake) SelectActiveValidators() {
    // 按票数排序
    sort.Slice(dpos.validators, func(i, j int) bool {
        return dpos.validators[i].Votes > dpos.validators[j].Votes
    })
    
    // 选择前N个验证者
    for i, validator := range dpos.validators {
        if i < dpos.maxValidators {
            validator.IsActive = true
            validator.Rank = i + 1
        } else {
            validator.IsActive = false
            validator.Rank = 0
        }
    }
}

// 轮询验证者
func (dpos *DelegatedProofOfStake) GetNextValidator() string {
    activeValidators := make([]*DPoSValidator, 0)
    for _, validator := range dpos.validators {
        if validator.IsActive {
            activeValidators = append(activeValidators, validator)
        }
    }
    
    if len(activeValidators) == 0 {
        return ""
    }
    
    // 基于轮询和性能的选择
    currentTime := time.Now()
    index := int(currentTime.Unix() / int64(dpos.blockTime.Seconds())) % len(activeValidators)
    
    return activeValidators[index].Address
}

// 更新验证者性能
func (dpos *DelegatedProofOfStake) UpdatePerformance(validator string, success bool) {
    for _, v := range dpos.validators {
        if v.Address == validator {
            if success {
                v.Performance = math.Min(1.0, v.Performance+0.01)
            } else {
                v.Performance = math.Max(0.0, v.Performance-0.05)
            }
            break
        }
    }
}
```

## 11.7.4.5 实用拜占庭容错 (PBFT)

### 11.7.4.5.1 理论基础

**定义 11.7.4.9** (拜占庭故障)
拜占庭故障是指节点可能发送任意错误消息的故障类型。

**定理 11.7.4.4** (PBFT容错性)
PBFT可以在 $f$ 个拜占庭节点存在的情况下正常工作，其中 $f < \frac{n}{3}$，$n$ 为总节点数。

**证明**:
设诚实节点数为 $h$，拜占庭节点数为 $f$。
需要满足：

1. $h > f$ (诚实节点多于拜占庭节点)
2. $h + f = n$ (总节点数)

因此 $h > \frac{n}{2}$，即 $f < \frac{n}{3}$。

### 11.7.4.5.2 数学建模

**定义 11.7.4.10** (PBFT阶段)
PBFT包含三个阶段：

1. **预准备阶段**: $2f+1$ 个节点确认
2. **准备阶段**: $2f+1$ 个节点准备
3. **提交阶段**: $2f+1$ 个节点提交

### 11.7.4.5.3 Go实现

```go
// PBFT共识实现
type PBFT struct {
    nodes        map[string]*PBFTNode
    primary      string
    viewNumber   uint64
    sequenceNum  uint64
    f            int // 最大拜占庭节点数
    checkpoint   uint64
}

// PBFT节点
type PBFTNode struct {
    ID           string
    IsPrimary    bool
    ViewNumber   uint64
    SequenceNum  uint64
    State        NodeState
    Messages     map[string]*PBFTMessage
    Checkpoints  map[uint64]int // sequence -> 确认数
}

// 节点状态
type NodeState int

const (
    Normal NodeState = iota
    ViewChange
    Recovering
)

// PBFT消息类型
type MessageType int

const (
    PrePrepare MessageType = iota
    Prepare
    Commit
    ViewChange
    NewView
)

// PBFT消息
type PBFTMessage struct {
    Type        MessageType
    ViewNumber  uint64
    SequenceNum uint64
    Digest      []byte
    NodeID      string
    Signature   []byte
    Timestamp   time.Time
}

// 创建PBFT共识
func NewPBFT(nodeIDs []string) *PBFT {
    nodes := make(map[string]*PBFTNode)
    f := (len(nodeIDs) - 1) / 3
    
    for _, id := range nodeIDs {
        nodes[id] = &PBFTNode{
            ID:          id,
            IsPrimary:   false,
            ViewNumber:  0,
            SequenceNum: 0,
            State:       Normal,
            Messages:    make(map[string]*PBFTMessage),
            Checkpoints: make(map[uint64]int),
        }
    }
    
    // 设置主节点
    primary := nodeIDs[0]
    nodes[primary].IsPrimary = true
    
    return &PBFT{
        nodes:       nodes,
        primary:     primary,
        viewNumber:  0,
        sequenceNum: 0,
        f:           f,
        checkpoint:  100,
    }
}

// 预准备阶段
func (pbft *PBFT) PrePrepare(request []byte) error {
    if !pbft.nodes[pbft.primary].IsPrimary {
        return fmt.Errorf("not primary node")
    }
    
    pbft.sequenceNum++
    digest := pbft.calculateDigest(request)
    
    message := &PBFTMessage{
        Type:        PrePrepare,
        ViewNumber:  pbft.viewNumber,
        SequenceNum: pbft.sequenceNum,
        Digest:      digest,
        NodeID:      pbft.primary,
        Timestamp:   time.Now(),
    }
    
    // 广播预准备消息
    pbft.broadcast(message)
    
    return nil
}

// 准备阶段
func (pbft *PBFT) Prepare(message *PBFTMessage) error {
    // 验证消息
    if !pbft.validateMessage(message) {
        return fmt.Errorf("invalid message")
    }
    
    // 检查是否已收到足够的预准备消息
    prePrepareCount := pbft.countMessages(PrePrepare, message.ViewNumber, message.SequenceNum)
    if prePrepareCount < pbft.f+1 {
        return fmt.Errorf("insufficient pre-prepare messages")
    }
    
    // 发送准备消息
    prepareMsg := &PBFTMessage{
        Type:        Prepare,
        ViewNumber:  message.ViewNumber,
        SequenceNum: message.SequenceNum,
        Digest:      message.Digest,
        NodeID:      pbft.getCurrentNodeID(),
        Timestamp:   time.Now(),
    }
    
    pbft.broadcast(prepareMsg)
    
    return nil
}

// 提交阶段
func (pbft *PBFT) Commit(message *PBFTMessage) error {
    // 检查是否已收到足够的准备消息
    prepareCount := pbft.countMessages(Prepare, message.ViewNumber, message.SequenceNum)
    if prepareCount < 2*pbft.f+1 {
        return fmt.Errorf("insufficient prepare messages")
    }
    
    // 发送提交消息
    commitMsg := &PBFTMessage{
        Type:        Commit,
        ViewNumber:  message.ViewNumber,
        SequenceNum: message.SequenceNum,
        Digest:      message.Digest,
        NodeID:      pbft.getCurrentNodeID(),
        Timestamp:   time.Now(),
    }
    
    pbft.broadcast(commitMsg)
    
    return nil
}

// 执行阶段
func (pbft *PBFT) Execute(message *PBFTMessage) error {
    // 检查是否已收到足够的提交消息
    commitCount := pbft.countMessages(Commit, message.ViewNumber, message.SequenceNum)
    if commitCount < 2*pbft.f+1 {
        return fmt.Errorf("insufficient commit messages")
    }
    
    // 执行请求
    pbft.executeRequest(message.Digest)
    
    // 检查点
    if message.SequenceNum%pbft.checkpoint == 0 {
        pbft.createCheckpoint(message.SequenceNum)
    }
    
    return nil
}

// 视图变更
func (pbft *PBFT) ViewChange() error {
    pbft.viewNumber++
    
    // 选择新的主节点
    nodeIDs := make([]string, 0)
    for id := range pbft.nodes {
        nodeIDs = append(nodeIDs, id)
    }
    sort.Strings(nodeIDs)
    
    newPrimaryIndex := int(pbft.viewNumber) % len(nodeIDs)
    newPrimary := nodeIDs[newPrimaryIndex]
    
    // 更新主节点
    for _, node := range pbft.nodes {
        node.IsPrimary = false
    }
    pbft.nodes[newPrimary].IsPrimary = true
    pbft.primary = newPrimary
    
    // 发送视图变更消息
    viewChangeMsg := &PBFTMessage{
        Type:       ViewChange,
        ViewNumber: pbft.viewNumber,
        NodeID:     pbft.getCurrentNodeID(),
        Timestamp:  time.Now(),
    }
    
    pbft.broadcast(viewChangeMsg)
    
    return nil
}

// 辅助方法
func (pbft *PBFT) validateMessage(message *PBFTMessage) bool {
    // 验证视图号
    if message.ViewNumber < pbft.viewNumber {
        return false
    }
    
    // 验证序列号
    if message.SequenceNum <= pbft.sequenceNum {
        return false
    }
    
    // 验证签名（简化实现）
    return true
}

func (pbft *PBFT) countMessages(msgType MessageType, viewNumber, sequenceNum uint64) int {
    count := 0
    for _, node := range pbft.nodes {
        for _, msg := range node.Messages {
            if msg.Type == msgType && 
               msg.ViewNumber == viewNumber && 
               msg.SequenceNum == sequenceNum {
                count++
            }
        }
    }
    return count
}

func (pbft *PBFT) broadcast(message *PBFTMessage) {
    // 简化实现：直接存储消息
    for _, node := range pbft.nodes {
        key := fmt.Sprintf("%s_%d_%d", node.ID, message.ViewNumber, message.SequenceNum)
        node.Messages[key] = message
    }
}

func (pbft *PBFT) calculateDigest(data []byte) []byte {
    hash := sha256.Sum256(data)
    return hash[:]
}

func (pbft *PBFT) getCurrentNodeID() string {
    // 简化实现：返回第一个节点ID
    for id := range pbft.nodes {
        return id
    }
    return ""
}

func (pbft *PBFT) executeRequest(digest []byte) {
    // 执行请求的逻辑
    fmt.Printf("Executing request with digest: %x\n", digest)
}

func (pbft *PBFT) createCheckpoint(sequenceNum uint64) {
    for _, node := range pbft.nodes {
        node.Checkpoints[sequenceNum] = len(pbft.nodes)
    }
}
```

## 11.7.4.6 共识算法比较

### 11.7.4.6.1 性能对比

| 算法 | 吞吐量 | 延迟 | 能源效率 | 去中心化程度 |
|------|--------|------|----------|--------------|
| PoW | 低 | 高 | 低 | 高 |
| PoS | 中 | 中 | 高 | 中 |
| DPoS | 高 | 低 | 高 | 低 |
| PBFT | 高 | 低 | 高 | 中 |

### 11.7.4.6.2 安全性分析

**定理 11.7.4.5** (共识算法安全性比较)
在相同网络条件下，PBFT > PoS > DPoS > PoW 的安全性排序。

**证明**:

1. **PBFT**: 需要 $2f+1$ 个诚实节点，容错率最高
2. **PoS**: 攻击成本等于质押价值，经济安全性强
3. **DPoS**: 验证者数量有限，攻击成本相对较低
4. **PoW**: 51%攻击成本相对较低

### 11.7.4.6.3 Go实现比较框架

```go
// 共识算法接口
type ConsensusAlgorithm interface {
    ProposeBlock(data []byte) (*Block, error)
    ValidateBlock(block *Block) bool
    GetConsensus() ([]*Block, error)
    GetPerformance() ConsensusPerformance
}

// 共识性能指标
type ConsensusPerformance struct {
    Throughput    float64 // TPS
    Latency       time.Duration
    EnergyUsage   float64 // kWh
    SecurityLevel float64 // 0-1
}

// 共识算法比较器
type ConsensusComparator struct {
    algorithms map[string]ConsensusAlgorithm
}

// 创建比较器
func NewConsensusComparator() *ConsensusComparator {
    return &ConsensusComparator{
        algorithms: make(map[string]ConsensusAlgorithm),
    }
}

// 添加算法
func (cc *ConsensusComparator) AddAlgorithm(name string, algo ConsensusAlgorithm) {
    cc.algorithms[name] = algo
}

// 性能测试
func (cc *ConsensusComparator) Benchmark(iterations int) map[string]ConsensusPerformance {
    results := make(map[string]ConsensusPerformance)
    
    for name, algo := range cc.algorithms {
        start := time.Now()
        var totalLatency time.Duration
        var successCount int
        
        for i := 0; i < iterations; i++ {
            blockStart := time.Now()
            
            // 模拟区块提议
            data := []byte(fmt.Sprintf("test_data_%d", i))
            block, err := algo.ProposeBlock(data)
            
            if err == nil && algo.ValidateBlock(block) {
                successCount++
                totalLatency += time.Since(blockStart)
            }
        }
        
        duration := time.Since(start)
        throughput := float64(successCount) / duration.Seconds()
        avgLatency := totalLatency / time.Duration(successCount)
        
        performance := algo.GetPerformance()
        performance.Throughput = throughput
        performance.Latency = avgLatency
        
        results[name] = performance
    }
    
    return results
}

// 生成比较报告
func (cc *ConsensusComparator) GenerateReport() string {
    results := cc.Benchmark(1000)
    
    report := "共识算法性能比较报告\n"
    report += "==================\n\n"
    
    for name, perf := range results {
        report += fmt.Sprintf("算法: %s\n", name)
        report += fmt.Sprintf("  吞吐量: %.2f TPS\n", perf.Throughput)
        report += fmt.Sprintf("  延迟: %v\n", perf.Latency)
        report += fmt.Sprintf("  能源使用: %.2f kWh\n", perf.EnergyUsage)
        report += fmt.Sprintf("  安全等级: %.2f\n", perf.SecurityLevel)
        report += "\n"
    }
    
    return report
}
```

## 11.7.4.7 总结

本章详细介绍了区块链共识算法的理论基础、数学建模和Go语言实现，包括：

1. **工作量证明 (PoW)**: 基于计算难度的共识机制
2. **权益证明 (PoS)**: 基于代币质押的共识机制
3. **委托权益证明 (DPoS)**: 基于投票选举的共识机制
4. **实用拜占庭容错 (PBFT)**: 基于消息传递的共识机制

每种算法都有其独特的优势和适用场景，选择时需要综合考虑性能、安全性和去中心化程度等因素。

---

**相关链接**:

- [11.7.1 区块链基础理论](../01-Blockchain-Fundamentals/README.md)
- [11.7.2 智能合约](../02-Smart-Contracts/README.md)
- [11.7.3 DeFi协议](../03-DeFi-Protocols/README.md)
- [11.8 其他高级主题](../README.md)
