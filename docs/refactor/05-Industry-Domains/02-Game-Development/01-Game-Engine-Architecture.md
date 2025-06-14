# 01-游戏引擎架构 (Game Engine Architecture)

## 目录

1. [理论基础](#1-理论基础)
2. [形式化定义](#2-形式化定义)
3. [核心系统](#3-核心系统)
4. [Go语言实现](#4-go语言实现)
5. [性能优化](#5-性能优化)

## 1. 理论基础

### 1.1 游戏引擎定义

游戏引擎是游戏开发的核心框架，提供渲染、物理、音频、输入、网络等基础功能。

**形式化定义**：
```math
游戏引擎定义为六元组：
GE = (R, P, A, I, N, S)

其中：
- R: 渲染系统
- P: 物理系统
- A: 音频系统
- I: 输入系统
- N: 网络系统
- S: 场景系统
```

### 1.2 核心特性

1. **实时渲染**: 60FPS以上的渲染性能
2. **物理模拟**: 精确的物理计算
3. **音频处理**: 3D音频和音效
4. **输入处理**: 多设备输入支持
5. **网络同步**: 实时多人游戏

## 2. 形式化定义

### 2.1 游戏循环

```math
游戏循环定义为：
Loop(t) = Update(t) \circ Render(t) \circ ProcessInput(t)

其中：
- t: 时间步长
- Update: 游戏逻辑更新
- Render: 画面渲染
- ProcessInput: 输入处理
```

### 2.2 实体组件系统

```math
ECS系统定义为：
ECS = (E, C, S)

其中：
- E: 实体集合
- C: 组件集合
- S: 系统集合

实体定义：
Entity = \{c_1, c_2, ..., c_n\} \subseteq C
```

## 3. 核心系统

### 3.1 渲染系统

```go
// RenderSystem 渲染系统
type RenderSystem struct {
    renderer    *Renderer
    shaderManager *ShaderManager
    textureManager *TextureManager
    camera      *Camera
    scene       *Scene
}

// Renderer 渲染器
type Renderer struct {
    window      *Window
    context     *Context
    pipeline    *Pipeline
}

// Render 渲染场景
func (rs *RenderSystem) Render() error {
    // 清除缓冲区
    rs.renderer.Clear()
    
    // 设置相机
    rs.renderer.SetCamera(rs.camera)
    
    // 渲染所有对象
    for _, object := range rs.scene.Objects {
        rs.renderObject(object)
    }
    
    // 交换缓冲区
    rs.renderer.SwapBuffers()
    
    return nil
}
```

### 3.2 物理系统

```go
// PhysicsSystem 物理系统
type PhysicsSystem struct {
    world       *PhysicsWorld
    bodies      map[string]*RigidBody
    constraints []*Constraint
}

// PhysicsWorld 物理世界
type PhysicsWorld struct {
    gravity     Vector3
    timeStep    float64
    iterations  int
}

// Update 更新物理
func (ps *PhysicsSystem) Update(deltaTime float64) {
    // 应用重力
    ps.applyGravity()
    
    // 碰撞检测
    ps.detectCollisions()
    
    // 求解约束
    ps.solveConstraints()
    
    // 更新位置
    ps.updatePositions(deltaTime)
}

// applyGravity 应用重力
func (ps *PhysicsSystem) applyGravity() {
    for _, body := range ps.bodies {
        if body.IsStatic {
            continue
        }
        body.Velocity = body.Velocity.Add(ps.world.gravity.Multiply(ps.world.timeStep))
    }
}
```

### 3.3 音频系统

```go
// AudioSystem 音频系统
type AudioSystem struct {
    device      *AudioDevice
    sources     map[string]*AudioSource
    listener    *AudioListener
}

// AudioSource 音频源
type AudioSource struct {
    ID          string
    Position    Vector3
    Velocity    Vector3
    Buffer      *AudioBuffer
    Volume      float64
    Pitch       float64
    Loop        bool
}

// PlaySound 播放音效
func (as *AudioSystem) PlaySound(sourceID string) error {
    source, exists := as.sources[sourceID]
    if !exists {
        return errors.New("audio source not found")
    }
    
    // 计算3D音频
    distance := as.listener.Position.Distance(source.Position)
    volume := as.calculate3DVolume(distance, source.Volume)
    
    // 播放音频
    return as.device.Play(source.Buffer, volume, source.Pitch)
}

// calculate3DVolume 计算3D音量
func (as *AudioSystem) calculate3DVolume(distance, baseVolume float64) float64 {
    // 距离衰减
    attenuation := 1.0 / (1.0 + distance*0.1)
    return baseVolume * attenuation
}
```

## 4. Go语言实现

### 4.1 游戏引擎核心

```go
// GameEngine 游戏引擎
type GameEngine struct {
    renderSystem   *RenderSystem
    physicsSystem  *PhysicsSystem
    audioSystem    *AudioSystem
    inputSystem    *InputSystem
    networkSystem  *NetworkSystem
    sceneSystem    *SceneSystem
    running        bool
    lastTime       time.Time
}

// NewGameEngine 创建游戏引擎
func NewGameEngine() *GameEngine {
    return &GameEngine{
        renderSystem:   NewRenderSystem(),
        physicsSystem:  NewPhysicsSystem(),
        audioSystem:    NewAudioSystem(),
        inputSystem:    NewInputSystem(),
        networkSystem:  NewNetworkSystem(),
        sceneSystem:    NewSceneSystem(),
        running:        false,
        lastTime:       time.Now(),
    }
}

// Start 启动引擎
func (ge *GameEngine) Start() error {
    ge.running = true
    ge.lastTime = time.Now()
    
    // 初始化所有系统
    if err := ge.initializeSystems(); err != nil {
        return err
    }
    
    // 开始游戏循环
    go ge.gameLoop()
    
    return nil
}

// Stop 停止引擎
func (ge *GameEngine) Stop() {
    ge.running = false
}

// gameLoop 游戏循环
func (ge *GameEngine) gameLoop() {
    for ge.running {
        currentTime := time.Now()
        deltaTime := currentTime.Sub(ge.lastTime).Seconds()
        ge.lastTime = currentTime
        
        // 处理输入
        ge.inputSystem.ProcessInput()
        
        // 更新物理
        ge.physicsSystem.Update(deltaTime)
        
        // 更新场景
        ge.sceneSystem.Update(deltaTime)
        
        // 渲染
        ge.renderSystem.Render()
        
        // 网络同步
        ge.networkSystem.Update()
        
        // 控制帧率
        ge.limitFrameRate()
    }
}

// limitFrameRate 限制帧率
func (ge *GameEngine) limitFrameRate() {
    targetFrameTime := 1.0 / 60.0 // 60 FPS
    elapsed := time.Since(ge.lastTime).Seconds()
    
    if elapsed < targetFrameTime {
        sleepTime := time.Duration((targetFrameTime - elapsed) * float64(time.Second))
        time.Sleep(sleepTime)
    }
}

// initializeSystems 初始化系统
func (ge *GameEngine) initializeSystems() error {
    if err := ge.renderSystem.Initialize(); err != nil {
        return err
    }
    
    if err := ge.physicsSystem.Initialize(); err != nil {
        return err
    }
    
    if err := ge.audioSystem.Initialize(); err != nil {
        return err
    }
    
    if err := ge.inputSystem.Initialize(); err != nil {
        return err
    }
    
    if err := ge.networkSystem.Initialize(); err != nil {
        return err
    }
    
    if err := ge.sceneSystem.Initialize(); err != nil {
        return err
    }
    
    return nil
}
```

### 4.2 实体组件系统

```go
// Entity 实体
type Entity struct {
    ID       string
    Components map[string]Component
    Active   bool
}

// Component 组件接口
type Component interface {
    GetType() string
    GetID() string
}

// TransformComponent 变换组件
type TransformComponent struct {
    ID       string
    Position Vector3
    Rotation Vector3
    Scale    Vector3
}

// GetType 获取组件类型
func (tc *TransformComponent) GetType() string {
    return "Transform"
}

// GetID 获取组件ID
func (tc *TransformComponent) GetID() string {
    return tc.ID
}

// RenderComponent 渲染组件
type RenderComponent struct {
    ID       string
    Mesh     *Mesh
    Material *Material
    Visible  bool
}

// GetType 获取组件类型
func (rc *RenderComponent) GetType() string {
    return "Render"
}

// GetID 获取组件ID
func (rc *RenderComponent) GetID() string {
    return rc.ID
}

// PhysicsComponent 物理组件
type PhysicsComponent struct {
    ID       string
    Body     *RigidBody
    Collider *Collider
}

// GetType 获取组件类型
func (pc *PhysicsComponent) GetType() string {
    return "Physics"
}

// GetID 获取组件ID
func (pc *PhysicsComponent) GetID() string {
    return pc.ID
}

// EntityManager 实体管理器
type EntityManager struct {
    entities map[string]*Entity
    systems  []System
    mu       sync.RWMutex
}

// NewEntityManager 创建实体管理器
func NewEntityManager() *EntityManager {
    return &EntityManager{
        entities: make(map[string]*Entity),
        systems:  make([]System, 0),
    }
}

// CreateEntity 创建实体
func (em *EntityManager) CreateEntity() *Entity {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    entity := &Entity{
        ID:        em.generateID(),
        Components: make(map[string]Component),
        Active:    true,
    }
    
    em.entities[entity.ID] = entity
    return entity
}

// AddComponent 添加组件
func (em *EntityManager) AddComponent(entityID string, component Component) error {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    entity, exists := em.entities[entityID]
    if !exists {
        return errors.New("entity not found")
    }
    
    entity.Components[component.GetType()] = component
    return nil
}

// GetComponent 获取组件
func (em *EntityManager) GetComponent(entityID, componentType string) (Component, error) {
    em.mu.RLock()
    defer em.mu.RUnlock()
    
    entity, exists := em.entities[entityID]
    if !exists {
        return nil, errors.New("entity not found")
    }
    
    component, exists := entity.Components[componentType]
    if !exists {
        return nil, errors.New("component not found")
    }
    
    return component, nil
}

// RemoveEntity 移除实体
func (em *EntityManager) RemoveEntity(entityID string) error {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    if _, exists := em.entities[entityID]; !exists {
        return errors.New("entity not found")
    }
    
    delete(em.entities, entityID)
    return nil
}

// AddSystem 添加系统
func (em *EntityManager) AddSystem(system System) {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    em.systems = append(em.systems, system)
}

// Update 更新所有系统
func (em *EntityManager) Update(deltaTime float64) {
    em.mu.RLock()
    defer em.mu.RUnlock()
    
    for _, system := range em.systems {
        system.Update(em.entities, deltaTime)
    }
}

// generateID 生成ID
func (em *EntityManager) generateID() string {
    return fmt.Sprintf("entity_%d", len(em.entities)+1)
}
```

### 4.3 系统接口

```go
// System 系统接口
type System interface {
    Update(entities map[string]*Entity, deltaTime float64)
    GetName() string
}

// RenderSystem 渲染系统
type RenderSystem struct {
    name string
}

// GetName 获取系统名称
func (rs *RenderSystem) GetName() string {
    return "RenderSystem"
}

// Update 更新渲染系统
func (rs *RenderSystem) Update(entities map[string]*Entity, deltaTime float64) {
    for _, entity := range entities {
        if !entity.Active {
            continue
        }
        
        // 获取变换组件
        transformComp, err := entity.Components["Transform"]
        if err != nil {
            continue
        }
        
        // 获取渲染组件
        renderComp, err := entity.Components["Render"]
        if err != nil {
            continue
        }
        
        transform := transformComp.(*TransformComponent)
        render := renderComp.(*RenderComponent)
        
        if render.Visible {
            // 渲染对象
            rs.renderObject(transform, render)
        }
    }
}

// renderObject 渲染对象
func (rs *RenderSystem) renderObject(transform *TransformComponent, render *RenderComponent) {
    // 设置变换矩阵
    matrix := rs.calculateTransformMatrix(transform)
    
    // 绑定材质
    rs.bindMaterial(render.Material)
    
    // 渲染网格
    rs.renderMesh(render.Mesh, matrix)
}

// PhysicsSystem 物理系统
type PhysicsSystem struct {
    name string
}

// GetName 获取系统名称
func (ps *PhysicsSystem) GetName() string {
    return "PhysicsSystem"
}

// Update 更新物理系统
func (ps *PhysicsSystem) Update(entities map[string]*Entity, deltaTime float64) {
    for _, entity := range entities {
        if !entity.Active {
            continue
        }
        
        // 获取物理组件
        physicsComp, err := entity.Components["Physics"]
        if err != nil {
            continue
        }
        
        physics := physicsComp.(*PhysicsComponent)
        
        // 更新物理
        ps.updatePhysics(physics, deltaTime)
    }
}

// updatePhysics 更新物理
func (ps *PhysicsSystem) updatePhysics(physics *PhysicsComponent, deltaTime float64) {
    // 更新刚体
    if physics.Body != nil {
        physics.Body.Update(deltaTime)
    }
    
    // 碰撞检测
    if physics.Collider != nil {
        ps.detectCollisions(physics.Collider)
    }
}
```

## 5. 性能优化

### 5.1 空间分区

```go
// SpatialPartition 空间分区
type SpatialPartition struct {
    gridSize    float64
    grid        map[Vector3][]*Entity
}

// NewSpatialPartition 创建空间分区
func NewSpatialPartition(gridSize float64) *SpatialPartition {
    return &SpatialPartition{
        gridSize: gridSize,
        grid:     make(map[Vector3][]*Entity),
    }
}

// Insert 插入实体
func (sp *SpatialPartition) Insert(entity *Entity) {
    transform, err := entity.Components["Transform"]
    if err != nil {
        return
    }
    
    pos := transform.(*TransformComponent).Position
    gridPos := sp.worldToGrid(pos)
    
    sp.grid[gridPos] = append(sp.grid[gridPos], entity)
}

// Query 查询区域内的实体
func (sp *SpatialPartition) Query(bounds Bounds) []*Entity {
    var entities []*Entity
    
    minGrid := sp.worldToGrid(bounds.Min)
    maxGrid := sp.worldToGrid(bounds.Max)
    
    for x := minGrid.X; x <= maxGrid.X; x++ {
        for y := minGrid.Y; y <= maxGrid.Y; y++ {
            for z := minGrid.Z; z <= maxGrid.Z; z++ {
                gridPos := Vector3{X: x, Y: y, Z: z}
                if gridEntities, exists := sp.grid[gridPos]; exists {
                    entities = append(entities, gridEntities...)
                }
            }
        }
    }
    
    return entities
}

// worldToGrid 世界坐标转网格坐标
func (sp *SpatialPartition) worldToGrid(pos Vector3) Vector3 {
    return Vector3{
        X: math.Floor(pos.X / sp.gridSize),
        Y: math.Floor(pos.Y / sp.gridSize),
        Z: math.Floor(pos.Z / sp.gridSize),
    }
}
```

### 5.2 对象池

```go
// ObjectPool 对象池
type ObjectPool struct {
    pool    []interface{}
    factory func() interface{}
    reset   func(interface{})
    mu      sync.Mutex
}

// NewObjectPool 创建对象池
func NewObjectPool(factory func() interface{}, reset func(interface{}), initialSize int) *ObjectPool {
    pool := &ObjectPool{
        pool:    make([]interface{}, 0, initialSize),
        factory: factory,
        reset:   reset,
    }
    
    // 预创建对象
    for i := 0; i < initialSize; i++ {
        pool.pool = append(pool.pool, factory())
    }
    
    return pool
}

// Get 获取对象
func (op *ObjectPool) Get() interface{} {
    op.mu.Lock()
    defer op.mu.Unlock()
    
    if len(op.pool) > 0 {
        obj := op.pool[len(op.pool)-1]
        op.pool = op.pool[:len(op.pool)-1]
        return obj
    }
    
    return op.factory()
}

// Return 返回对象
func (op *ObjectPool) Return(obj interface{}) {
    op.mu.Lock()
    defer op.mu.Unlock()
    
    if op.reset != nil {
        op.reset(obj)
    }
    
    op.pool = append(op.pool, obj)
}
```

## 总结

游戏引擎架构是游戏开发的核心，需要平衡性能、功能和易用性。本文档提供了完整的理论基础、形式化定义和Go语言实现。

### 关键要点

1. **模块化设计**: 清晰的系统分离
2. **性能优化**: 空间分区、对象池等技术
3. **实时性**: 60FPS的游戏循环
4. **可扩展性**: 支持插件和自定义组件
5. **跨平台**: 支持多种平台

### 扩展阅读

- [网络游戏服务器](./02-Network-Game-Server.md)
- [实时渲染系统](./03-Real-time-Rendering-System.md)
- [物理引擎](./04-Physics-Engine.md) 