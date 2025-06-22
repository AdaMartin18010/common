# 01-企业应用

 (Enterprise Applications)

## 目录

- [01-企业应用](#01-企业应用)
  - [目录](#目录)
  - [1. 企业工作流概述](#1-企业工作流概述)
    - [1.1 业务流程管理](#11-业务流程管理)
      - [1.1.1 业务流程定义](#111-业务流程定义)
      - [1.1.2 工作流模式](#112-工作流模式)
    - [1.2 企业集成](#12-企业集成)
      - [1.2.1 企业集成模型](#121-企业集成模型)
  - [2. 核心业务流程](#2-核心业务流程)
    - [2.1 订单处理流程](#21-订单处理流程)
    - [2.2 审批流程](#22-审批流程)
    - [2.3 采购流程](#23-采购流程)
  - [3. 工作流建模](#3-工作流建模)
    - [3.1 BPMN建模](#31-bpmn建模)
    - [3.2 状态机建模](#32-状态机建模)
    - [3.3 规则引擎集成](#33-规则引擎集成)
  - [4. 企业集成模式](#4-企业集成模式)
    - [4.1 消息队列集成](#41-消息队列集成)
    - [4.2 数据库集成](#42-数据库集成)
    - [4.3 API集成](#43-api集成)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 企业工作流引擎](#51-企业工作流引擎)
    - [5.2 业务流程实现](#52-业务流程实现)
    - [5.3 集成组件](#53-集成组件)
  - [总结](#总结)

---

## 1. 企业工作流概述

### 1.1 业务流程管理

#### 1.1.1 业务流程定义

业务流程 ```latex
$BP$
``` 是一个五元组：

```latex
$```latex
$\text{BusinessProcess} = (A, R, D, C, O)$
```$
```

其中：

- ```latex
$A$
``` 是活动集合
- ```latex
$R$
``` 是角色集合
- ```latex
$D$
``` 是数据集合
- ```latex
$C$
``` 是控制流
- ```latex
$O$
``` 是业务目标

**业务流程特征**:

1. **结构化**: 活动之间有明确的依赖关系
2. **可重复**: 流程可以多次执行
3. **可测量**: 流程性能可以量化
4. **可优化**: 流程可以持续改进

#### 1.1.2 工作流模式

工作流模式集合：

```latex
$```latex
$\text{WorkflowPattern} = \{\text{Sequential}, \text{Parallel}, \text{Choice}, \text{Iteration}, \text{Compensation}\}$
```$
```

**常见模式**:

1. **顺序模式**: 活动按固定顺序执行
2. **并行模式**: 多个活动同时执行
3. **选择模式**: 根据条件选择执行路径
4. **迭代模式**: 活动重复执行
5. **补偿模式**: 错误时执行回滚操作

### 1.2 企业集成

#### 1.2.1 企业集成模型

企业集成 ```latex
$EI$
``` 定义为：

```latex
$```latex
$\text{EnterpriseIntegration} = \{\text{DataIntegration}, \text{ProcessIntegration}, \text{ApplicationIntegration}\}$
```$
```

**集成类型**:

1. **数据集成**: 统一数据格式和存储
2. **流程集成**: 协调跨系统业务流程
3. **应用集成**: 连接不同应用程序

## 2. 核心业务流程

### 2.1 订单处理流程

**流程定义**:

```go
// OrderProcessingWorkflow 订单处理工作流
type OrderProcessingWorkflow struct {
    OrderID     string
    CustomerID  string
    Items       []OrderItem
    Status      OrderStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// OrderItem 订单项
type OrderItem struct {
    ProductID   string
    Quantity    int
    UnitPrice   decimal.Decimal
    TotalPrice  decimal.Decimal
}

// OrderStatus 订单状态
type OrderStatus int

const (
    OrderCreated OrderStatus = iota
    OrderValidated
    OrderPaymentProcessed
    OrderInventoryChecked
    OrderShipped
    OrderDelivered
    OrderCompleted
    OrderCancelled
)

// 订单处理流程实现
func (opw *OrderProcessingWorkflow) Process() error {
    // 1. 创建订单
    if err := opw.createOrder(); err != nil {
        return err
    }
    
    // 2. 验证订单
    if err := opw.validateOrder(); err != nil {
        return opw.cancelOrder("订单验证失败")
    }
    
    // 3. 处理支付
    if err := opw.processPayment(); err != nil {
        return opw.cancelOrder("支付处理失败")
    }
    
    // 4. 检查库存
    if err := opw.checkInventory(); err != nil {
        return opw.refundPayment("库存不足")
    }
    
    // 5. 发货
    if err := opw.shipOrder(); err != nil {
        return opw.refundPayment("发货失败")
    }
    
    // 6. 确认交付
    if err := opw.confirmDelivery(); err != nil {
        return err
    }
    
    // 7. 完成订单
    return opw.completeOrder()
}

func (opw *OrderProcessingWorkflow) createOrder() error {
    opw.Status = OrderCreated
    opw.CreatedAt = time.Now()
    return nil
}

func (opw *OrderProcessingWorkflow) validateOrder() error {
    // 验证客户信息
    if opw.CustomerID == "" {
        return fmt.Errorf("客户ID不能为空")
    }
    
    // 验证订单项
    if len(opw.Items) == 0 {
        return fmt.Errorf("订单项不能为空")
    }
    
    // 验证价格
    for _, item := range opw.Items {
        if item.UnitPrice.LessThanOrEqual(decimal.Zero) {
            return fmt.Errorf("商品价格必须大于0")
        }
    }
    
    opw.Status = OrderValidated
    return nil
}

func (opw *OrderProcessingWorkflow) processPayment() error {
    // 计算总金额
    total := decimal.Zero
    for _, item := range opw.Items {
        total = total.Add(item.TotalPrice)
    }
    
    // 调用支付服务
    paymentService := NewPaymentService()
    if err := paymentService.ProcessPayment(opw.CustomerID, total); err != nil {
        return err
    }
    
    opw.Status = OrderPaymentProcessed
    return nil
}

func (opw *OrderProcessingWorkflow) checkInventory() error {
    inventoryService := NewInventoryService()
    
    for _, item := range opw.Items {
        if err := inventoryService.ReserveInventory(item.ProductID, item.Quantity); err != nil {
            return err
        }
    }
    
    opw.Status = OrderInventoryChecked
    return nil
}

func (opw *OrderProcessingWorkflow) shipOrder() error {
    shippingService := NewShippingService()
    
    if err := shippingService.CreateShipment(opw.OrderID, opw.Items); err != nil {
        return err
    }
    
    opw.Status = OrderShipped
    return nil
}

func (opw *OrderProcessingWorkflow) confirmDelivery() error {
    opw.Status = OrderDelivered
    return nil
}

func (opw *OrderProcessingWorkflow) completeOrder() error {
    opw.Status = OrderCompleted
    opw.UpdatedAt = time.Now()
    return nil
}

func (opw *OrderProcessingWorkflow) cancelOrder(reason string) error {
    opw.Status = OrderCancelled
    opw.UpdatedAt = time.Now()
    return fmt.Errorf("订单已取消: %s", reason)
}

func (opw *OrderProcessingWorkflow) refundPayment(reason string) error {
    // 处理退款
    paymentService := NewPaymentService()
    total := decimal.Zero
    for _, item := range opw.Items {
        total = total.Add(item.TotalPrice)
    }
    
    if err := paymentService.RefundPayment(opw.CustomerID, total); err != nil {
        return err
    }
    
    return opw.cancelOrder(reason)
}
```

### 2.2 审批流程

**流程定义**:

```go
// ApprovalWorkflow 审批工作流
type ApprovalWorkflow struct {
    RequestID   string
    RequesterID string
    RequestType string
    Amount      decimal.Decimal
    Status      ApprovalStatus
    Approvers   []Approver
    CurrentStep int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Approver 审批人
type Approver struct {
    UserID      string
    Role        string
    Priority    int
    Status      ApproverStatus
    Comments    string
    ApprovedAt  time.Time
}

// ApprovalStatus 审批状态
type ApprovalStatus int

const (
    ApprovalPending ApprovalStatus = iota
    ApprovalInProgress
    ApprovalApproved
    ApprovalRejected
    ApprovalCancelled
)

// ApproverStatus 审批人状态
type ApproverStatus int

const (
    ApproverPending ApproverStatus = iota
    ApproverApproved
    ApproverRejected
)

// 审批流程实现
func (aw *ApprovalWorkflow) Process() error {
    // 1. 创建审批请求
    if err := aw.createRequest(); err != nil {
        return err
    }
    
    // 2. 分配审批人
    if err := aw.assignApprovers(); err != nil {
        return err
    }
    
    // 3. 执行审批流程
    for aw.CurrentStep < len(aw.Approvers) {
        if err := aw.processApprovalStep(); err != nil {
            return err
        }
        
        if aw.Status == ApprovalRejected || aw.Status == ApprovalCancelled {
            break
        }
    }
    
    // 4. 完成审批
    return aw.completeApproval()
}

func (aw *ApprovalWorkflow) createRequest() error {
    aw.Status = ApprovalPending
    aw.CreatedAt = time.Now()
    return nil
}

func (aw *ApprovalWorkflow) assignApprovers() error {
    // 根据请求类型和金额确定审批人
    approvalService := NewApprovalService()
    approvers, err := approvalService.GetApprovers(aw.RequestType, aw.Amount)
    if err != nil {
        return err
    }
    
    aw.Approvers = approvers
    aw.Status = ApprovalInProgress
    return nil
}

func (aw *ApprovalWorkflow) processApprovalStep() error {
    if aw.CurrentStep >= len(aw.Approvers) {
        return fmt.Errorf("审批步骤超出范围")
    }
    
    currentApprover := &aw.Approvers[aw.CurrentStep]
    
    // 检查审批人状态
    if currentApprover.Status == ApproverPending {
        // 等待审批
        return nil
    }
    
    if currentApprover.Status == ApproverRejected {
        aw.Status = ApprovalRejected
        return nil
    }
    
    // 审批通过，进入下一步
    aw.CurrentStep++
    
    // 检查是否所有审批人都已审批
    if aw.CurrentStep >= len(aw.Approvers) {
        aw.Status = ApprovalApproved
    }
    
    return nil
}

func (aw *ApprovalWorkflow) completeApproval() error {
    aw.UpdatedAt = time.Now()
    
    // 根据审批结果执行相应操作
    switch aw.Status {
    case ApprovalApproved:
        return aw.executeApprovedRequest()
    case ApprovalRejected:
        return aw.handleRejectedRequest()
    default:
        return fmt.Errorf("未知的审批状态: %v", aw.Status)
    }
}

func (aw *ApprovalWorkflow) executeApprovedRequest() error {
    // 执行已批准的请求
    executionService := NewExecutionService()
    return executionService.ExecuteRequest(aw.RequestID, aw.RequestType, aw.Amount)
}

func (aw *ApprovalWorkflow) handleRejectedRequest() error {
    // 处理被拒绝的请求
    notificationService := NewNotificationService()
    return notificationService.NotifyRejection(aw.RequesterID, aw.RequestID)
}

// Approve 审批操作
func (aw *ApprovalWorkflow) Approve(userID string, approved bool, comments string) error {
    if aw.CurrentStep >= len(aw.Approvers) {
        return fmt.Errorf("没有待审批的步骤")
    }
    
    currentApprover := &aw.Approvers[aw.CurrentStep]
    if currentApprover.UserID != userID {
        return fmt.Errorf("用户无权进行此审批")
    }
    
    if approved {
        currentApprover.Status = ApproverApproved
    } else {
        currentApprover.Status = ApproverRejected
    }
    
    currentApprover.Comments = comments
    currentApprover.ApprovedAt = time.Now()
    
    return aw.processApprovalStep()
}
```

### 2.3 采购流程

**流程定义**:

```go
// ProcurementWorkflow 采购工作流
type ProcurementWorkflow struct {
    PurchaseID    string
    RequesterID   string
    Department    string
    Items         []PurchaseItem
    TotalAmount   decimal.Decimal
    Status        ProcurementStatus
    Suppliers     []Supplier
    SelectedSupplier *Supplier
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// PurchaseItem 采购项
type PurchaseItem struct {
    ProductID    string
    Description  string
    Quantity     int
    UnitPrice    decimal.Decimal
    TotalPrice   decimal.Decimal
}

// Supplier 供应商
type Supplier struct {
    SupplierID   string
    Name         string
    Rating       float64
    Quote        decimal.Decimal
    DeliveryTime int // 天数
}

// ProcurementStatus 采购状态
type ProcurementStatus int

const (
    ProcurementRequested ProcurementStatus = iota
    ProcurementApproved
    ProcurementQuoted
    ProcurementSupplierSelected
    ProcurementOrdered
    ProcurementDelivered
    ProcurementCompleted
    ProcurementCancelled
)

// 采购流程实现
func (pw *ProcurementWorkflow) Process() error {
    // 1. 创建采购请求
    if err := pw.createRequest(); err != nil {
        return err
    }
    
    // 2. 审批采购请求
    if err := pw.approveRequest(); err != nil {
        return pw.cancelProcurement("采购请求被拒绝")
    }
    
    // 3. 获取供应商报价
    if err := pw.getQuotes(); err != nil {
        return pw.cancelProcurement("无法获取供应商报价")
    }
    
    // 4. 选择供应商
    if err := pw.selectSupplier(); err != nil {
        return pw.cancelProcurement("无法选择供应商")
    }
    
    // 5. 下订单
    if err := pw.placeOrder(); err != nil {
        return pw.cancelProcurement("下订单失败")
    }
    
    // 6. 接收货物
    if err := pw.receiveDelivery(); err != nil {
        return err
    }
    
    // 7. 完成采购
    return pw.completeProcurement()
}

func (pw *ProcurementWorkflow) createRequest() error {
    pw.Status = ProcurementRequested
    pw.CreatedAt = time.Now()
    
    // 计算总金额
    total := decimal.Zero
    for _, item := range pw.Items {
        total = total.Add(item.TotalPrice)
    }
    pw.TotalAmount = total
    
    return nil
}

func (pw *ProcurementWorkflow) approveRequest() error {
    // 创建审批工作流
    approvalWorkflow := &ApprovalWorkflow{
        RequestID:   pw.PurchaseID,
        RequesterID: pw.RequesterID,
        RequestType: "Procurement",
        Amount:      pw.TotalAmount,
    }
    
    if err := approvalWorkflow.Process(); err != nil {
        return err
    }
    
    if approvalWorkflow.Status != ApprovalApproved {
        return fmt.Errorf("采购请求审批未通过")
    }
    
    pw.Status = ProcurementApproved
    return nil
}

func (pw *ProcurementWorkflow) getQuotes() error {
    // 获取供应商列表
    supplierService := NewSupplierService()
    suppliers, err := supplierService.GetSuppliers(pw.Items)
    if err != nil {
        return err
    }
    
    // 获取报价
    for i := range suppliers {
        quote, err := supplierService.GetQuote(suppliers[i].SupplierID, pw.Items)
        if err != nil {
            continue
        }
        suppliers[i].Quote = quote
    }
    
    pw.Suppliers = suppliers
    pw.Status = ProcurementQuoted
    return nil
}

func (pw *ProcurementWorkflow) selectSupplier() error {
    if len(pw.Suppliers) == 0 {
        return fmt.Errorf("没有可用的供应商")
    }
    
    // 选择最佳供应商（基于价格、评级、交付时间）
    bestSupplier := pw.Suppliers[0]
    bestScore := pw.calculateSupplierScore(bestSupplier)
    
    for _, supplier := range pw.Suppliers[1:] {
        score := pw.calculateSupplierScore(supplier)
        if score > bestScore {
            bestScore = score
            bestSupplier = supplier
        }
    }
    
    pw.SelectedSupplier = &bestSupplier
    pw.Status = ProcurementSupplierSelected
    return nil
}

func (pw *ProcurementWorkflow) calculateSupplierScore(supplier Supplier) float64 {
    // 评分算法：价格权重40%，评级权重40%，交付时间权重20%
    priceScore := 1.0 / supplier.Quote.InexactFloat64()
    ratingScore := supplier.Rating
    deliveryScore := 1.0 / float64(supplier.DeliveryTime)
    
    return 0.4*priceScore + 0.4*ratingScore + 0.2*deliveryScore
}

func (pw *ProcurementWorkflow) placeOrder() error {
    if pw.SelectedSupplier == nil {
        return fmt.Errorf("未选择供应商")
    }
    
    // 创建采购订单
    orderService := NewOrderService()
    if err := orderService.CreateOrder(pw.PurchaseID, pw.SelectedSupplier.SupplierID, pw.Items); err != nil {
        return err
    }
    
    pw.Status = ProcurementOrdered
    return nil
}

func (pw *ProcurementWorkflow) receiveDelivery() error {
    // 检查交付状态
    deliveryService := NewDeliveryService()
    if err := deliveryService.CheckDelivery(pw.PurchaseID); err != nil {
        return err
    }
    
    pw.Status = ProcurementDelivered
    return nil
}

func (pw *ProcurementWorkflow) completeProcurement() error {
    pw.Status = ProcurementCompleted
    pw.UpdatedAt = time.Now()
    return nil
}

func (pw *ProcurementWorkflow) cancelProcurement(reason string) error {
    pw.Status = ProcurementCancelled
    pw.UpdatedAt = time.Now()
    return fmt.Errorf("采购已取消: %s", reason)
}
```

## 3. 工作流建模

### 3.1 BPMN建模

**BPMN元素映射**:

```go
// BPMNElement BPMN元素接口
type BPMNElement interface {
    GetID() string
    GetType() BPMNType
    GetName() string
}

// BPMNType BPMN类型
type BPMNType int

const (
    BPMNStartEvent BPMNType = iota
    BPMNEndEvent
    BPMNTask
    BPMNGateway
    BPMNSequenceFlow
)

// BPMNProcess BPMN流程
type BPMNProcess struct {
    ID          string
    Name        string
    Elements    map[string]BPMNElement
    Flows       []BPMNSequenceFlow
    StartEvent  string
    EndEvents   []string
}

// BPMNTask BPMN任务
type BPMNTask struct {
    ID          string
    Name        string
    Type        string
    Assignee    string
    CandidateGroups []string
    FormKey     string
}

// BPMNGateway BPMN网关
type BPMNGateway struct {
    ID          string
    Name        string
    Type        GatewayType
    Conditions  map[string]string
}

// GatewayType 网关类型
type GatewayType int

const (
    GatewayExclusive GatewayType = iota
    GatewayParallel
    GatewayInclusive
)

// BPMNSequenceFlow BPMN顺序流
type BPMNSequenceFlow struct {
    ID          string
    SourceRef   string
    TargetRef   string
    Condition   string
}
```

### 3.2 状态机建模

**状态机实现**:

```go
// StateMachine 状态机
type StateMachine struct {
    ID          string
    Name        string
    States      map[string]State
    Transitions []Transition
    CurrentState string
    InitialState string
    FinalStates []string
}

// State 状态
type State struct {
    ID          string
    Name        string
    Actions     []Action
    EntryAction Action
    ExitAction  Action
}

// Transition 转换
type Transition struct {
    ID          string
    FromState   string
    ToState     string
    Event       string
    Condition   func(interface{}) bool
    Action      Action
}

// Action 动作
type Action func(interface{}) error

// 状态机执行器
func (sm *StateMachine) Execute(event string, data interface{}) error {
    // 查找可用的转换
    for _, transition := range sm.Transitions {
        if transition.FromState == sm.CurrentState && transition.Event == event {
            // 检查条件
            if transition.Condition != nil && !transition.Condition(data) {
                continue
            }
            
            // 执行退出动作
            if exitAction := sm.States[sm.CurrentState].ExitAction; exitAction != nil {
                if err := exitAction(data); err != nil {
                    return err
                }
            }
            
            // 执行转换动作
            if transition.Action != nil {
                if err := transition.Action(data); err != nil {
                    return err
                }
            }
            
            // 更新状态
            sm.CurrentState = transition.ToState
            
            // 执行进入动作
            if entryAction := sm.States[sm.CurrentState].EntryAction; entryAction != nil {
                if err := entryAction(data); err != nil {
                    return err
                }
            }
            
            return nil
        }
    }
    
    return fmt.Errorf("没有找到从状态 %s 在事件 %s 下的有效转换", sm.CurrentState, event)
}
```

### 3.3 规则引擎集成

**规则引擎**:

```go
// RuleEngine 规则引擎
type RuleEngine struct {
    Rules       []Rule
    Facts       map[string]interface{}
}

// Rule 规则
type Rule struct {
    ID          string
    Name        string
    Conditions  []Condition
    Actions     []Action
    Priority    int
}

// Condition 条件
type Condition func(map[string]interface{}) bool

// 规则引擎执行
func (re *RuleEngine) Execute() error {
    // 按优先级排序规则
    sort.Slice(re.Rules, func(i, j int) bool {
        return re.Rules[i].Priority > re.Rules[j].Priority
    })
    
    // 执行规则
    for _, rule := range re.Rules {
        if re.evaluateConditions(rule.Conditions) {
            if err := re.executeActions(rule.Actions); err != nil {
                return err
            }
        }
    }
    
    return nil
}

func (re *RuleEngine) evaluateConditions(conditions []Condition) bool {
    for _, condition := range conditions {
        if !condition(re.Facts) {
            return false
        }
    }
    return true
}

func (re *RuleEngine) executeActions(actions []Action) error {
    for _, action := range actions {
        if err := action(re.Facts); err != nil {
            return err
        }
    }
    return nil
}
```

## 4. 企业集成模式

### 4.1 消息队列集成

**消息队列集成**:

```go
// MessageQueueIntegration 消息队列集成
type MessageQueueIntegration struct {
    Producer    MessageProducer
    Consumer    MessageConsumer
    QueueName   string
}

// MessageProducer 消息生产者
type MessageProducer interface {
    SendMessage(queueName string, message []byte) error
    SendMessageWithHeaders(queueName string, message []byte, headers map[string]string) error
}

// MessageConsumer 消息消费者
type MessageConsumer interface {
    ConsumeMessage(queueName string, handler func([]byte) error) error
    ConsumeMessageWithHeaders(queueName string, handler func([]byte, map[string]string) error) error
}

// RabbitMQIntegration RabbitMQ集成
type RabbitMQIntegration struct {
    Connection  *amqp.Connection
    Channel     *amqp.Channel
}

func (rmq *RabbitMQIntegration) SendMessage(queueName string, message []byte) error {
    return rmq.Channel.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )
}

func (rmq *RabbitMQIntegration) ConsumeMessage(queueName string, handler func([]byte) error) error {
    msgs, err := rmq.Channel.Consume(
        queueName, // queue
        "",        // consumer
        true,      // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return err
    }
    
    go func() {
        for msg := range msgs {
            handler(msg.Body)
        }
    }()
    
    return nil
}
```

### 4.2 数据库集成

**数据库集成**:

```go
// DatabaseIntegration 数据库集成
type DatabaseIntegration struct {
    DB          *sql.DB
    Driver      string
    DSN         string
}

// WorkflowRepository 工作流仓库
type WorkflowRepository struct {
    db *sql.DB
}

func (wr *WorkflowRepository) SaveWorkflowInstance(instance *WorkflowInstance) error {
    query := `
        INSERT INTO workflow_instances (id, workflow_id, status, data, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    data, err := json.Marshal(instance.Data)
    if err != nil {
        return err
    }
    
    _, err = wr.db.Exec(query,
        instance.ID,
        instance.WorkflowID,
        instance.Status,
        data,
        instance.CreatedAt,
        instance.UpdatedAt,
    )
    
    return err
}

func (wr *WorkflowRepository) GetWorkflowInstance(id string) (*WorkflowInstance, error) {
    query := `
        SELECT id, workflow_id, status, data, created_at, updated_at
        FROM workflow_instances
        WHERE id = ?
    `
    
    var instance WorkflowInstance
    var data []byte
    
    err := wr.db.QueryRow(query, id).Scan(
        &instance.ID,
        &instance.WorkflowID,
        &instance.Status,
        &data,
        &instance.CreatedAt,
        &instance.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    
    if err := json.Unmarshal(data, &instance.Data); err != nil {
        return nil, err
    }
    
    return &instance, nil
}

func (wr *WorkflowRepository) UpdateWorkflowInstance(instance *WorkflowInstance) error {
    query := `
        UPDATE workflow_instances
        SET status = ?, data = ?, updated_at = ?
        WHERE id = ?
    `
    
    data, err := json.Marshal(instance.Data)
    if err != nil {
        return err
    }
    
    instance.UpdatedAt = time.Now()
    
    _, err = wr.db.Exec(query,
        instance.Status,
        data,
        instance.UpdatedAt,
        instance.ID,
    )
    
    return err
}
```

### 4.3 API集成

**API集成**:

```go
// APIIntegration API集成
type APIIntegration struct {
    Client      *http.Client
    BaseURL     string
    Headers     map[string]string
}

// RESTClient REST客户端
type RESTClient struct {
    client  *http.Client
    baseURL string
    headers map[string]string
}

func (rc *RESTClient) Get(path string, result interface{}) error {
    req, err := http.NewRequest("GET", rc.baseURL+path, nil)
    if err != nil {
        return err
    }
    
    // 添加请求头
    for key, value := range rc.headers {
        req.Header.Set(key, value)
    }
    
    resp, err := rc.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP错误: %d", resp.StatusCode)
    }
    
    return json.NewDecoder(resp.Body).Decode(result)
}

func (rc *RESTClient) Post(path string, data interface{}, result interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("POST", rc.baseURL+path, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    
    // 添加请求头
    req.Header.Set("Content-Type", "application/json")
    for key, value := range rc.headers {
        req.Header.Set(key, value)
    }
    
    resp, err := rc.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("HTTP错误: %d", resp.StatusCode)
    }
    
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    
    return nil
}
```

## 5. Go语言实现

### 5.1 企业工作流引擎

```go
// EnterpriseWorkflowEngine 企业工作流引擎
type EnterpriseWorkflowEngine struct {
    workflows       map[string]*WorkflowDefinition
    instances       map[string]*WorkflowInstance
    taskExecutors   map[string]TaskExecutor
    eventBus        EventBus
    repository      WorkflowRepository
    mu              sync.RWMutex
}

// WorkflowDefinition 工作流定义
type WorkflowDefinition struct {
    ID          string
    Name        string
    Version     string
    Tasks       []TaskDefinition
    Transitions []TransitionDefinition
    StartEvent  string
    EndEvents   []string
}

// TaskDefinition 任务定义
type TaskDefinition struct {
    ID              string
    Name            string
    Type            string
    Assignee        string
    CandidateGroups []string
    FormKey         string
    RetryPolicy     *RetryPolicy
    Timeout         time.Duration
}

// TransitionDefinition 转换定义
type TransitionDefinition struct {
    ID          string
    FromTask    string
    ToTask      string
    Condition   string
    Priority    int
}

// WorkflowInstance 工作流实例
type WorkflowInstance struct {
    ID          string
    WorkflowID  string
    Status      WorkflowStatus
    Data        map[string]interface{}
    Tasks       map[string]*TaskInstance
    CurrentTask string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// TaskInstance 任务实例
type TaskInstance struct {
    ID          string
    TaskID      string
    Status      TaskStatus
    Assignee    string
    Data        map[string]interface{}
    CreatedAt   time.Time
    StartedAt   time.Time
    CompletedAt time.Time
    Error       error
}

// 工作流引擎方法
func (ewe *EnterpriseWorkflowEngine) StartWorkflow(workflowID string, data map[string]interface{}) (*WorkflowInstance, error) {
    ewe.mu.Lock()
    defer ewe.mu.Unlock()
    
    // 获取工作流定义
    definition, exists := ewe.workflows[workflowID]
    if !exists {
        return nil, fmt.Errorf("工作流定义不存在: %s", workflowID)
    }
    
    // 创建工作流实例
    instance := &WorkflowInstance{
        ID:         generateID(),
        WorkflowID: workflowID,
        Status:     WorkflowStatusRunning,
        Data:       data,
        Tasks:      make(map[string]*TaskInstance),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    // 初始化任务实例
    for _, taskDef := range definition.Tasks {
        taskInstance := &TaskInstance{
            ID:       generateID(),
            TaskID:   taskDef.ID,
            Status:   TaskStatusPending,
            Data:     make(map[string]interface{}),
            CreatedAt: time.Now(),
        }
        instance.Tasks[taskDef.ID] = taskInstance
    }
    
    // 设置当前任务为开始事件
    instance.CurrentTask = definition.StartEvent
    
    // 保存实例
    ewe.instances[instance.ID] = instance
    if err := ewe.repository.SaveWorkflowInstance(instance); err != nil {
        return nil, err
    }
    
    // 发布工作流开始事件
    ewe.eventBus.Publish("workflow.started", instance)
    
    // 开始执行第一个任务
    return instance, ewe.executeCurrentTask(instance)
}

func (ewe *EnterpriseWorkflowEngine) executeCurrentTask(instance *WorkflowInstance) error {
    definition := ewe.workflows[instance.WorkflowID]
    taskDef := definition.Tasks[instance.CurrentTask]
    taskInstance := instance.Tasks[instance.CurrentTask]
    
    // 更新任务状态
    taskInstance.Status = TaskStatusRunning
    taskInstance.StartedAt = time.Now()
    instance.UpdatedAt = time.Now()
    
    // 保存状态
    if err := ewe.repository.UpdateWorkflowInstance(instance); err != nil {
        return err
    }
    
    // 执行任务
    executor, exists := ewe.taskExecutors[taskDef.Type]
    if !exists {
        return fmt.Errorf("任务执行器不存在: %s", taskDef.Type)
    }
    
    go func() {
        result, err := executor.Execute(taskDef, taskInstance, instance.Data)
        
        ewe.mu.Lock()
        defer ewe.mu.Unlock()
        
        // 更新任务状态
        if err != nil {
            taskInstance.Status = TaskStatusFailed
            taskInstance.Error = err
        } else {
            taskInstance.Status = TaskStatusCompleted
            taskInstance.Data = result
        }
        taskInstance.CompletedAt = time.Now()
        
        // 更新工作流状态
        instance.UpdatedAt = time.Now()
        
        // 保存状态
        ewe.repository.UpdateWorkflowInstance(instance)
        
        // 发布任务完成事件
        ewe.eventBus.Publish("task.completed", taskInstance)
        
        // 继续执行下一个任务
        ewe.continueWorkflow(instance)
    }()
    
    return nil
}

func (ewe *EnterpriseWorkflowEngine) continueWorkflow(instance *WorkflowInstance) error {
    definition := ewe.workflows[instance.WorkflowID]
    
    // 查找下一个任务
    nextTask := ewe.findNextTask(definition, instance.CurrentTask, instance.Data)
    
    if nextTask == "" {
        // 工作流完成
        instance.Status = WorkflowStatusCompleted
        instance.UpdatedAt = time.Now()
        ewe.repository.UpdateWorkflowInstance(instance)
        ewe.eventBus.Publish("workflow.completed", instance)
        return nil
    }
    
    // 设置下一个任务
    instance.CurrentTask = nextTask
    instance.UpdatedAt = time.Now()
    ewe.repository.UpdateWorkflowInstance(instance)
    
    // 执行下一个任务
    return ewe.executeCurrentTask(instance)
}

func (ewe *EnterpriseWorkflowEngine) findNextTask(definition *WorkflowDefinition, currentTask string, data map[string]interface{}) string {
    var candidates []string
    
    for _, transition := range definition.Transitions {
        if transition.FromTask == currentTask {
            // 检查条件
            if transition.Condition == "" || ewe.evaluateCondition(transition.Condition, data) {
                candidates = append(candidates, transition.ToTask)
            }
        }
    }
    
    if len(candidates) == 0 {
        return ""
    }
    
    // 如果有多个候选，选择优先级最高的
    if len(candidates) > 1 {
        // 简化实现：选择第一个
        return candidates[0]
    }
    
    return candidates[0]
}

func (ewe *EnterpriseWorkflowEngine) evaluateCondition(condition string, data map[string]interface{}) bool {
    // 简单的条件评估实现
    // 实际应用中可以使用表达式引擎
    return true
}
```

### 5.2 业务流程实现

```go
// BusinessProcessManager 业务流程管理器
type BusinessProcessManager struct {
    engine *EnterpriseWorkflowEngine
}

// OrderProcessManager 订单流程管理器
type OrderProcessManager struct {
    *BusinessProcessManager
    orderService    OrderService
    paymentService  PaymentService
    inventoryService InventoryService
    shippingService ShippingService
}

func (opm *OrderProcessManager) ProcessOrder(order *Order) error {
    // 创建工作流数据
    data := map[string]interface{}{
        "orderID":     order.ID,
        "customerID":  order.CustomerID,
        "items":       order.Items,
        "totalAmount": order.TotalAmount,
    }
    
    // 启动订单处理工作流
    instance, err := opm.engine.StartWorkflow("order-processing", data)
    if err != nil {
        return err
    }
    
    // 监控工作流执行
    go opm.monitorWorkflow(instance.ID)
    
    return nil
}

func (opm *OrderProcessManager) monitorWorkflow(instanceID string) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            instance, err := opm.engine.GetWorkflowInstance(instanceID)
            if err != nil {
                log.Printf("获取工作流实例失败: %v", err)
                return
            }
            
            if instance.Status == WorkflowStatusCompleted {
                log.Printf("订单处理工作流完成: %s", instanceID)
                return
            } else if instance.Status == WorkflowStatusFailed {
                log.Printf("订单处理工作流失败: %s", instanceID)
                return
            }
        }
    }
}
```

### 5.3 集成组件

```go
// IntegrationManager 集成管理器
type IntegrationManager struct {
    messageQueue MessageQueueIntegration
    database     DatabaseIntegration
    api          APIIntegration
    eventBus     EventBus
}

// EventBus 事件总线
type EventBus interface {
    Publish(topic string, event interface{}) error
    Subscribe(topic string, handler func(interface{}) error) error
    Unsubscribe(topic string) error
}

// InMemoryEventBus 内存事件总线
type InMemoryEventBus struct {
    handlers map[string][]func(interface{}) error
    mu       sync.RWMutex
}

func (eb *InMemoryEventBus) Publish(topic string, event interface{}) error {
    eb.mu.RLock()
    handlers := eb.handlers[topic]
    eb.mu.RUnlock()
    
    for _, handler := range handlers {
        go func(h func(interface{}) error) {
            if err := h(event); err != nil {
                log.Printf("事件处理失败: %v", err)
            }
        }(handler)
    }
    
    return nil
}

func (eb *InMemoryEventBus) Subscribe(topic string, handler func(interface{}) error) error {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    if eb.handlers == nil {
        eb.handlers = make(map[string][]func(interface{}) error)
    }
    
    eb.handlers[topic] = append(eb.handlers[topic], handler)
    return nil
}

func (eb *InMemoryEventBus) Unsubscribe(topic string) error {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    delete(eb.handlers, topic)
    return nil
}

// 集成管理器方法
func (im *IntegrationManager) HandleWorkflowEvent(event interface{}) error {
    // 处理工作流事件
    switch e := event.(type) {
    case *WorkflowInstance:
        // 发送到消息队列
        message, _ := json.Marshal(e)
        im.messageQueue.SendMessage("workflow-events", message)
        
        // 保存到数据库
        im.database.SaveWorkflowInstance(e)
        
        // 调用外部API
        im.api.Post("/workflow/events", e, nil)
        
    case *TaskInstance:
        // 处理任务事件
        message, _ := json.Marshal(e)
        im.messageQueue.SendMessage("task-events", message)
    }
    
    return nil
}
```

## 总结

本文档详细介绍了企业应用中的工作流系统，包括：

1. **业务流程管理**: 订单处理、审批、采购等核心流程
2. **工作流建模**: BPMN、状态机、规则引擎
3. **企业集成**: 消息队列、数据库、API集成
4. **Go语言实现**: 企业工作流引擎、业务流程实现、集成组件

企业工作流系统是业务流程自动化的核心，通过标准化的流程定义和执行引擎，提高业务效率和一致性。
