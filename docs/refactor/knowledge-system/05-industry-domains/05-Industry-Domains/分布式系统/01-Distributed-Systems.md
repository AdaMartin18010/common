# 分布式系统 (Distributed Systems)

## 1. 基本概念

### 1.1 分布式系统定义

**分布式系统 (Distributed System)** 是由多个独立计算机组成的系统，这些计算机通过网络进行通信和协调，共同完成一个任务。分布式系统的目标是提供比单个计算机更高的性能、可靠性和可扩展性。

### 1.2 核心特征

- **并发性**: 多个组件可以同时执行
- **缺乏全局时钟**: 不同节点的时间可能不同步
- **故障独立性**: 单个节点故障不影响整个系统
- **消息传递**: 节点间通过消息进行通信
- **部分故障**: 系统可能部分可用

### 1.3 设计目标

- **透明性**: 用户感知不到系统的分布式特性
- **开放性**: 支持异构系统和标准接口
- **可扩展性**: 能够处理负载增长
- **容错性**: 在部分故障时仍能正常工作
- **一致性**: 保证数据的一致性

## 2. 分布式系统架构

### 2.1 客户端-服务器架构

```go
// 服务器端
type Server struct {
    listener net.Listener
    handlers map[string]HandlerFunc
}

type HandlerFunc func(conn net.Conn, request []byte) ([]byte, error)

func (s *Server) Start(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    s.listener = listener
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        return
    }
    
    request := buffer[:n]
    response, err := s.processRequest(request)
    if err != nil {
        return
    }
    
    conn.Write(response)
}

// 客户端
type Client struct {
    serverAddr string
}

func (c *Client) SendRequest(data []byte) ([]byte, error) {
    conn, err := net.Dial("tcp", c.serverAddr)
    if err != nil {
        return nil, err
    }
    defer conn.Close()
    
    _, err = conn.Write(data)
    if err != nil {
        return nil, err
    }
    
    response := make([]byte, 1024)
    n, err := conn.Read(response)
    if err != nil {
        return nil, err
    }
    
    return response[:n], nil
}
```

### 2.2 对等网络架构

```go
// 对等节点
type Peer struct {
    id       string
    address  string
    peers    map[string]*PeerConnection
    data     map[string]interface{}
    mu       sync.RWMutex
}

type PeerConnection struct {
    id     string
    conn   net.Conn
    addr   string
}

func (p *Peer) Start(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    
    p.address = addr
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go p.handleConnection(conn)
    }
}

func (p *Peer) ConnectToPeer(peerAddr string) error {
    conn, err := net.Dial("tcp", peerAddr)
    if err != nil {
        return err
    }
    
    peerID := generatePeerID()
    peerConn := &PeerConnection{
        id:   peerID,
        conn: conn,
        addr: peerAddr,
    }
    
    p.mu.Lock()
    p.peers[peerID] = peerConn
    p.mu.Unlock()
    
    go p.handlePeerConnection(peerConn)
    return nil
}

func (p *Peer) Broadcast(message []byte) {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    for _, peer := range p.peers {
        go func(conn net.Conn) {
            conn.Write(message)
        }(peer.conn)
    }
}
```

### 2.3 微服务架构

```go
// 服务注册中心
type ServiceRegistry struct {
    services map[string][]*ServiceInstance
    mu       sync.RWMutex
}

type ServiceInstance struct {
    ID       string            `json:"id"`
    Name     string            `json:"name"`
    Address  string            `json:"address"`
    Port     int               `json:"port"`
    Metadata map[string]string `json:"metadata"`
    Health   string            `json:"health"`
}

func (sr *ServiceRegistry) Register(service *ServiceInstance) error {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    
    if sr.services[service.Name] == nil {
        sr.services[service.Name] = make([]*ServiceInstance, 0)
    }
    
    sr.services[service.Name] = append(sr.services[service.Name], service)
    return nil
}

func (sr *ServiceRegistry) GetService(name string) ([]*ServiceInstance, error) {
    sr.mu.RLock()
    defer sr.mu.RUnlock()
    
    instances, exists := sr.services[name]
    if !exists {
        return nil, fmt.Errorf("service %s not found", name)
    }
    
    return instances, nil
}

// 负载均衡器
type LoadBalancer struct {
    registry *ServiceRegistry
    strategy LoadBalancingStrategy
}

type LoadBalancingStrategy interface {
    Select(instances []*ServiceInstance) *ServiceInstance
}

type RoundRobinStrategy struct {
    current int
    mu      sync.Mutex
}

func (rr *RoundRobinStrategy) Select(instances []*ServiceInstance) *ServiceInstance {
    if len(instances) == 0 {
        return nil
    }
    
    rr.mu.Lock()
    defer rr.mu.Unlock()
    
    instance := instances[rr.current%len(instances)]
    rr.current++
    
    return instance
}
```

## 3. 分布式一致性

### 3.1 CAP定理

**CAP定理**指出，在分布式系统中，不可能同时满足以下三个特性：

- **一致性 (Consistency)**: 所有节点看到的数据是一致的
- **可用性 (Availability)**: 每个请求都能得到响应
- **分区容错性 (Partition Tolerance)**: 系统在网络分区时仍能正常工作

### 3.2 一致性模型

#### 3.2.1 强一致性

```go
// 强一致性实现
type StrongConsistency struct {
    nodes    map[string]*Node
    quorum   int
    mu       sync.Mutex
}

type Node struct {
    id      string
    data    map[string]interface{}
    version map[string]int
}

func (sc *StrongConsistency) Write(key string, value interface{}) error {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    
    // 获取当前版本号
    maxVersion := 0
    for _, node := range sc.nodes {
        if version, exists := node.version[key]; exists && version > maxVersion {
            maxVersion = version
        }
    }
    
    newVersion := maxVersion + 1
    
    // 向所有节点写入
    successCount := 0
    for _, node := range sc.nodes {
        if sc.writeToNode(node, key, value, newVersion) {
            successCount++
        }
    }
    
    // 检查是否达到法定人数
    if successCount < sc.quorum {
        return fmt.Errorf("failed to reach quorum: %d/%d", successCount, sc.quorum)
    }
    
    return nil
}

func (sc *StrongConsistency) Read(key string) (interface{}, error) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    
    // 从所有节点读取
    responses := make(map[int][]interface{})
    for _, node := range sc.nodes {
        if value, version, err := sc.readFromNode(node, key); err == nil {
            responses[version] = append(responses[version], value)
        }
    }
    
    // 找到最高版本的数据
    maxVersion := 0
    for version := range responses {
        if version > maxVersion {
            maxVersion = version
        }
    }
    
    if maxVersion == 0 {
        return nil, fmt.Errorf("key %s not found", key)
    }
    
    // 检查是否达到法定人数
    if len(responses[maxVersion]) < sc.quorum {
        return nil, fmt.Errorf("failed to reach quorum for read")
    }
    
    return responses[maxVersion][0], nil
}
```

#### 3.2.2 最终一致性

```go
// 最终一致性实现
type EventualConsistency struct {
    nodes    map[string]*Node
    events   chan Event
    mu       sync.RWMutex
}

type Event struct {
    Type      string      `json:"type"`
    Key       string      `json:"key"`
    Value     interface{} `json:"value"`
    Version   int         `json:"version"`
    Timestamp time.Time   `json:"timestamp"`
    Source    string      `json:"source"`
}

func (ec *EventualConsistency) Write(key string, value interface{}) error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    
    // 获取当前版本号
    maxVersion := 0
    for _, node := range ec.nodes {
        if version, exists := node.version[key]; exists && version > maxVersion {
            maxVersion = version
        }
    }
    
    newVersion := maxVersion + 1
    
    // 创建事件
    event := Event{
        Type:      "WRITE",
        Key:       key,
        Value:     value,
        Version:   newVersion,
        Timestamp: time.Now(),
        Source:    "client",
    }
    
    // 异步传播事件
    go ec.propagateEvent(event)
    
    return nil
}

func (ec *EventualConsistency) propagateEvent(event Event) {
    ec.events <- event
    
    // 向其他节点传播
    for nodeID, node := range ec.nodes {
        if nodeID != event.Source {
            go ec.sendEventToNode(node, event)
        }
    }
}

func (ec *EventualConsistency) Read(key string) (interface{}, error) {
    ec.mu.RLock()
    defer ec.mu.RUnlock()
    
    // 从本地节点读取（可能不是最新数据）
    for _, node := range ec.nodes {
        if value, exists := node.data[key]; exists {
            return value, nil
        }
    }
    
    return nil, fmt.Errorf("key %s not found", key)
}
```

### 3.3 分布式算法

#### 3.3.1 Paxos算法

```go
// Paxos算法实现
type PaxosNode struct {
    id           string
    proposers    map[string]*Proposer
    acceptors    map[string]*Acceptor
    learners     map[string]*Learner
    mu           sync.Mutex
}

type Proposer struct {
    id        string
    proposal  int
    value     interface{}
    promises  map[string]bool
    accepted  map[string]interface{}
}

type Acceptor struct {
    id           string
    promisedID   int
    acceptedID   int
    acceptedValue interface{}
}

type Learner struct {
    id      string
    learned map[string]interface{}
}

func (pn *PaxosNode) Propose(value interface{}) error {
    pn.mu.Lock()
    defer pn.mu.Unlock()
    
    proposerID := generateID()
    proposer := &Proposer{
        id:       proposerID,
        proposal: pn.generateProposalID(),
        value:    value,
        promises: make(map[string]bool),
        accepted: make(map[string]interface{}),
    }
    
    pn.proposers[proposerID] = proposer
    
    // Phase 1: Prepare
    promises := pn.prepare(proposer)
    
    // Phase 2: Accept
    if len(promises) > len(pn.acceptors)/2 {
        return pn.accept(proposer, promises)
    }
    
    return fmt.Errorf("failed to get majority promises")
}

func (pn *PaxosNode) prepare(proposer *Proposer) map[string]bool {
    promises := make(map[string]bool)
    
    for acceptorID, acceptor := range pn.acceptors {
        if proposer.proposal > acceptor.promisedID {
            acceptor.promisedID = proposer.proposal
            promises[acceptorID] = true
            
            if acceptor.acceptedID > 0 {
                proposer.accepted[acceptorID] = acceptor.acceptedValue
            }
        }
    }
    
    return promises
}

func (pn *PaxosNode) accept(proposer *Proposer, promises map[string]bool) error {
    // 如果有已接受的值，使用最高提案号的值
    if len(proposer.accepted) > 0 {
        maxID := 0
        for _, value := range proposer.accepted {
            if id, ok := value.(int); ok && id > maxID {
                maxID = id
                proposer.value = value
            }
        }
    }
    
    accepts := 0
    for acceptorID, acceptor := range pn.acceptors {
        if promises[acceptorID] {
            acceptor.acceptedID = proposer.proposal
            acceptor.acceptedValue = proposer.value
            accepts++
        }
    }
    
    if accepts > len(pn.acceptors)/2 {
        // 学习阶段
        pn.learn(proposer.value)
        return nil
    }
    
    return fmt.Errorf("failed to get majority accepts")
}

func (pn *PaxosNode) learn(value interface{}) {
    for _, learner := range pn.learners {
        learner.learned["consensus"] = value
    }
}
```

#### 3.3.2 Raft算法

```go
// Raft算法实现
type RaftNode struct {
    id          string
    term        int
    votedFor    string
    state       NodeState
    leaderID    string
    
    // 日志相关
    log         []LogEntry
    commitIndex int
    lastApplied int
    
    // 领导者相关
    nextIndex   map[string]int
    matchIndex  map[string]int
    
    // 选举相关
    electionTimeout time.Duration
    lastHeartbeat   time.Time
    
    mu sync.Mutex
}

type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

type LogEntry struct {
    Term    int         `json:"term"`
    Index   int         `json:"index"`
    Command interface{} `json:"command"`
}

func (rn *RaftNode) Start() {
    rn.state = Follower
    rn.electionTimeout = time.Duration(rand.Intn(300)+300) * time.Millisecond
    
    go rn.runElectionTimer()
    go rn.runHeartbeat()
}

func (rn *RaftNode) runElectionTimer() {
    for {
        time.Sleep(rn.electionTimeout)
        
        rn.mu.Lock()
        if rn.state != Leader {
            rn.startElection()
        }
        rn.mu.Unlock()
    }
}

func (rn *RaftNode) startElection() {
    rn.state = Candidate
    rn.term++
    rn.votedFor = rn.id
    rn.leaderID = ""
    
    votes := 1 // 自己的一票
    
    // 向其他节点请求投票
    for _, nodeID := range rn.getOtherNodes() {
        go func(id string) {
            if rn.requestVote(id) {
                rn.mu.Lock()
                votes++
                if votes > len(rn.getOtherNodes())/2+1 {
                    rn.becomeLeader()
                }
                rn.mu.Unlock()
            }
        }(nodeID)
    }
}

func (rn *RaftNode) becomeLeader() {
    rn.state = Leader
    rn.leaderID = rn.id
    
    // 初始化领导者状态
    for _, nodeID := range rn.getOtherNodes() {
        rn.nextIndex[nodeID] = len(rn.log)
        rn.matchIndex[nodeID] = 0
    }
    
    // 发送心跳
    rn.sendHeartbeat()
}

func (rn *RaftNode) sendHeartbeat() {
    for _, nodeID := range rn.getOtherNodes() {
        go rn.appendEntries(nodeID)
    }
}

func (rn *RaftNode) appendEntries(nodeID string) bool {
    // 实现AppendEntries RPC
    // 这里简化实现
    return true
}
```

## 4. 分布式存储

### 4.1 数据分片

```go
// 一致性哈希
type ConsistentHash struct {
    ring     map[uint32]string
    nodes    []string
    replicas int
    mu       sync.RWMutex
}

func (ch *ConsistentHash) AddNode(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()
    
    ch.nodes = append(ch.nodes, node)
    
    for i := 0; i < ch.replicas; i++ {
        hash := ch.hash(fmt.Sprintf("%s:%d", node, i))
        ch.ring[hash] = node
    }
}

func (ch *ConsistentHash) GetNode(key string) string {
    ch.mu.RLock()
    defer ch.mu.RUnlock()
    
    if len(ch.ring) == 0 {
        return ""
    }
    
    hash := ch.hash(key)
    
    // 找到大于等于hash的第一个节点
    var targetHash uint32
    var targetNode string
    
    for ringHash, node := range ch.ring {
        if ringHash >= hash {
            if targetNode == "" || ringHash < targetHash {
                targetHash = ringHash
                targetNode = node
            }
        }
    }
    
    // 如果没有找到，返回第一个节点
    if targetNode == "" {
        for ringHash, node := range ch.ring {
            if targetNode == "" || ringHash < targetHash {
                targetHash = ringHash
                targetNode = node
            }
        }
    }
    
    return targetNode
}

func (ch *ConsistentHash) hash(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}
```

### 4.2 数据复制

```go
// 主从复制
type MasterSlaveReplication struct {
    master   *Node
    slaves   map[string]*Node
    mu       sync.RWMutex
}

func (msr *MasterSlaveReplication) Write(key string, value interface{}) error {
    msr.mu.Lock()
    defer msr.mu.Unlock()
    
    // 写入主节点
    if err := msr.master.Write(key, value); err != nil {
        return err
    }
    
    // 异步复制到从节点
    go msr.replicateToSlaves(key, value)
    
    return nil
}

func (msr *MasterSlaveReplication) Read(key string) (interface{}, error) {
    msr.mu.RLock()
    defer msr.mu.RUnlock()
    
    // 从主节点读取
    return msr.master.Read(key)
}

func (msr *MasterSlaveReplication) replicateToSlaves(key string, value interface{}) {
    for slaveID, slave := range msr.slaves {
        go func(id string, s *Node) {
            if err := s.Write(key, value); err != nil {
                log.Printf("Failed to replicate to slave %s: %v", id, err)
            }
        }(slaveID, slave)
    }
}
```

## 5. 分布式事务

### 5.1 两阶段提交

```go
// 两阶段提交
type TwoPhaseCommit struct {
    coordinator *Coordinator
    participants map[string]*Participant
}

type Coordinator struct {
    id           string
    participants map[string]*Participant
    state        CommitState
    mu           sync.Mutex
}

type Participant struct {
    id    string
    state CommitState
    data  map[string]interface{}
    mu    sync.Mutex
}

type CommitState int

const (
    Initial CommitState = iota
    Prepared
    Committed
    Aborted
)

func (tpc *TwoPhaseCommit) BeginTransaction() string {
    txnID := generateTransactionID()
    tpc.coordinator.state = Initial
    return txnID
}

func (tpc *TwoPhaseCommit) Prepare(txnID string) error {
    tpc.coordinator.mu.Lock()
    defer tpc.coordinator.mu.Unlock()
    
    // 阶段1：准备阶段
    prepared := 0
    for _, participant := range tpc.coordinator.participants {
        if participant.Prepare(txnID) {
            prepared++
        }
    }
    
    // 检查是否所有参与者都准备就绪
    if prepared == len(tpc.coordinator.participants) {
        tpc.coordinator.state = Prepared
        return tpc.commit(txnID)
    } else {
        tpc.coordinator.state = Aborted
        return tpc.abort(txnID)
    }
}

func (tpc *TwoPhaseCommit) commit(txnID string) error {
    // 阶段2：提交阶段
    for _, participant := range tpc.coordinator.participants {
        participant.Commit(txnID)
    }
    tpc.coordinator.state = Committed
    return nil
}

func (tpc *TwoPhaseCommit) abort(txnID string) error {
    // 阶段2：回滚阶段
    for _, participant := range tpc.coordinator.participants {
        participant.Abort(txnID)
    }
    tpc.coordinator.state = Aborted
    return nil
}

func (p *Participant) Prepare(txnID string) bool {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // 检查是否可以准备
    if p.canPrepare(txnID) {
        p.state = Prepared
        return true
    }
    
    p.state = Aborted
    return false
}

func (p *Participant) Commit(txnID string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.state == Prepared {
        p.applyTransaction(txnID)
        p.state = Committed
    }
}

func (p *Participant) Abort(txnID string) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.rollbackTransaction(txnID)
    p.state = Aborted
}
```

## 6. 故障检测与恢复

### 6.1 心跳机制

```go
// 心跳检测
type HeartbeatDetector struct {
    nodes       map[string]*NodeInfo
    timeout     time.Duration
    checkInterval time.Duration
    mu          sync.RWMutex
}

type NodeInfo struct {
    id            string
    address       string
    lastHeartbeat time.Time
    status        NodeStatus
}

type NodeStatus int

const (
    Online NodeStatus = iota
    Offline
    Suspicious
)

func (hd *HeartbeatDetector) Start() {
    go hd.checkHeartbeats()
}

func (hd *HeartbeatDetector) checkHeartbeats() {
    ticker := time.NewTicker(hd.checkInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        hd.mu.Lock()
        now := time.Now()
        
        for _, node := range hd.nodes {
            if now.Sub(node.lastHeartbeat) > hd.timeout {
                if node.status == Online {
                    node.status = Suspicious
                    log.Printf("Node %s is suspicious", node.id)
                } else if node.status == Suspicious {
                    node.status = Offline
                    log.Printf("Node %s is offline", node.id)
                    hd.handleNodeFailure(node)
                }
            } else if node.status != Online {
                node.status = Online
                log.Printf("Node %s is back online", node.id)
            }
        }
        hd.mu.Unlock()
    }
}

func (hd *HeartbeatDetector) handleNodeFailure(node *NodeInfo) {
    // 处理节点故障
    // 1. 通知其他节点
    // 2. 重新分配负载
    // 3. 启动故障恢复流程
}
```

### 6.2 故障恢复

```go
// 故障恢复
type FailureRecovery struct {
    nodes        map[string]*Node
    backupNodes  map[string]*Node
    recoveryPlan map[string]RecoveryStrategy
}

type RecoveryStrategy interface {
    Recover(node *Node) error
}

type SimpleRecovery struct{}

func (sr *SimpleRecovery) Recover(node *Node) error {
    // 简单的重启策略
    return node.Restart()
}

type StateTransferRecovery struct {
    sourceNode *Node
}

func (str *StateTransferRecovery) Recover(node *Node) error {
    // 状态转移策略
    return node.TransferState(str.sourceNode)
}

func (fr *FailureRecovery) HandleFailure(nodeID string) error {
    node, exists := fr.nodes[nodeID]
    if !exists {
        return fmt.Errorf("node %s not found", nodeID)
    }
    
    strategy, exists := fr.recoveryPlan[nodeID]
    if !exists {
        strategy = &SimpleRecovery{}
    }
    
    return strategy.Recover(node)
}
```

## 总结

分布式系统是现代软件架构的重要组成部分，它通过将系统分解为多个独立节点来提高性能、可靠性和可扩展性。然而，分布式系统也带来了复杂性，包括网络分区、节点故障、数据一致性等问题。

**关键设计原则**：

1. 容错性：系统应该能够处理节点故障
2. 可扩展性：系统应该能够处理负载增长
3. 一致性：在可能的情况下保证数据一致性
4. 性能：在满足其他要求的前提下优化性能

**常见挑战**：

1. 网络延迟和分区
2. 时钟同步问题
3. 数据一致性
4. 故障检测和恢复
5. 负载均衡

分布式系统的设计需要在一致性、可用性和性能之间找到平衡，选择合适的算法和架构模式来满足特定的需求。
