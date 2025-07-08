# 04-物理引擎

 (Physics Engine)

## 目录

- [04-物理引擎](#04-物理引擎)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 核心功能](#11-核心功能)
    - [1.2 设计原则](#12-设计原则)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 物理系统](#21-物理系统)
    - [4.2 碰撞检测](#42-碰撞检测)
    - [4.3 碰撞响应](#43-碰撞响应)
    - [4.4 约束系统](#44-约束系统)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 物理引擎主类](#51-物理引擎主类)
    - [5.2 空间分割](#52-空间分割)
    - [5.3 窄相检测](#53-窄相检测)
    - [5.4 软体物理](#54-软体物理)
  - [6. 性能优化](#6-性能优化)
    - [6.1 空间优化](#61-空间优化)
    - [6.2 时间优化](#62-时间优化)
    - [6.3 内存优化](#63-内存优化)
  - [7. 总结](#7-总结)
    - [7.1 关键特性](#71-关键特性)
    - [7.2 扩展方向](#72-扩展方向)

## 1. 概述

物理引擎是游戏引擎的核心组件，负责模拟现实世界的物理现象，包括重力、碰撞、约束、流体等。现代物理引擎需要提供准确、稳定的物理模拟，同时保证实时性能。

### 1.1 核心功能

- **刚体动力学**: 处理刚体的运动和旋转
- **碰撞检测**: 检测物体之间的碰撞
- **碰撞响应**: 处理碰撞后的物理反应
- **约束系统**: 处理关节、弹簧等约束
- **软体物理**: 模拟布料、绳索等软体
- **流体模拟**: 模拟液体和气体的流动

### 1.2 设计原则

- **准确性**: 提供准确的物理模拟
- **稳定性**: 保证数值计算的稳定性
- **高性能**: 支持实时物理计算
- **可扩展**: 支持自定义物理行为
- **易用性**: 提供简单的API接口

## 2. 形式化定义

### 2.1 物理系统

**定义 2.1.1** (物理系统)
物理系统是一个六元组 ```latex
P = (B, C, F, T, S, I)

```，其中：

- ```latex
B
``` 是刚体集合 (Body Set)
- ```latex
C
``` 是约束集合 (Constraint Set)
- ```latex
F
``` 是力集合 (Force Set)
- ```latex
T
``` 是时间系统 (Time System)
- ```latex
S
``` 是求解器 (Solver)
- ```latex
I
``` 是积分器 (Integrator)

**定义 2.1.2** (刚体)
刚体是一个五元组 ```latex
R = (m, I, p, v, \omega)
```，其中：

- ```latex
m
``` 是质量 (Mass)
- ```latex
I
``` 是惯性张量 (Inertia Tensor)
- ```latex
p
``` 是位置 (Position)
- ```latex
v
``` 是速度 (Velocity)
- ```latex
\omega
``` 是角速度 (Angular Velocity)

### 2.2 运动方程

**定义 2.2.1** (牛顿第二定律)
对于刚体 ```latex
R
```，运动方程为：
$```latex
F = m \cdot a
```$
$```latex
\tau = I \cdot \alpha
```$

其中 ```latex
F
``` 是合力，```latex
\tau
``` 是合力矩，```latex
a
``` 是加速度，```latex
\alpha
``` 是角加速度。

**定义 2.2.2** (运动积分)
位置和速度的积分方程为：
$```latex
p(t + \Delta t) = p(t) + v(t) \cdot \Delta t + \frac{1}{2} a(t) \cdot \Delta t^2
```$
$```latex
v(t + \Delta t) = v(t) + a(t) \cdot \Delta t
```$

## 3. 数学基础

### 3.1 线性代数

**定理 3.1.1** (惯性张量的对角化)
对于刚体的惯性张量 ```latex
I
```，存在正交矩阵 ```latex
Q
``` 使得：
$```latex
I = Q \cdot D \cdot Q^T
```$

其中 ```latex
D
``` 是对角矩阵，对角线元素是主惯性矩。

**证明**:
根据谱定理，实对称矩阵可以对角化。惯性张量是实对称矩阵，因此存在正交矩阵 ```latex
Q
``` 使得 ```latex
I = Q \cdot D \cdot Q^T
```。

### 3.2 碰撞检测

**定义 3.2.1** (分离轴定理)
对于两个凸多面体 ```latex
A
``` 和 ```latex
B
```，如果存在一个轴 ```latex
n
``` 使得 ```latex
A
``` 和 ```latex
B
``` 在该轴上的投影不重叠，则 ```latex
A
``` 和 ```latex
B
``` 不相交。

**定理 3.2.1** (GJK算法)
GJK算法可以在 ```latex
O(n)
``` 时间内检测两个凸多面体的碰撞，其中 ```latex
n
``` 是顶点数。

### 3.3 约束求解

**定义 3.3.1** (约束方程)
约束方程定义为：
$```latex
C(q) = 0
```$

其中 ```latex
q
``` 是广义坐标。

**定理 3.3.1** (拉格朗日乘数法)
约束力的计算为：
$```latex
\lambda = -J \cdot M^{-1} \cdot J^T \cdot C
```$

其中 ```latex
J
``` 是雅可比矩阵，```latex
M
``` 是质量矩阵。

## 4. 物理系统

### 4.1 刚体动力学

```go
// 刚体
type RigidBody struct {
    ID          uint32
    mass        float32
    inverseMass float32
    inertia     mgl32.Mat3
    inverseInertia mgl32.Mat3
    
    position    mgl32.Vec3
    rotation    mgl32.Quat
    velocity    mgl32.Vec3
    angularVel  mgl32.Vec3
    
    force       mgl32.Vec3
    torque      mgl32.Vec3
    
    collider    Collider
    material    *PhysicsMaterial
    
    // 状态标志
    isStatic    bool
    isAwake     bool
    sleepTimer  float32
}

// 物理材质
type PhysicsMaterial struct {
    friction    float32
    restitution float32
    density     float32
}

// 更新刚体
func (rb *RigidBody) Update(deltaTime float32) {
    if rb.isStatic {
        return
    }
    
    // 计算加速度
    acceleration := rb.force.Mul(rb.inverseMass)
    angularAccel := rb.inverseInertia.Mul3x1(rb.torque)
    
    // 积分速度和角速度
    rb.velocity = rb.velocity.Add(acceleration.Mul(deltaTime))
    rb.angularVel = rb.angularVel.Add(angularAccel.Mul(deltaTime))
    
    // 积分位置和旋转
    rb.position = rb.position.Add(rb.velocity.Mul(deltaTime))
    
    // 旋转积分
    angularVelQuat := mgl32.Quat{
        W: 0,
        V: rb.angularVel.Mul(deltaTime * 0.5),
    }
    rb.rotation = rb.rotation.Add(angularVelQuat.Mul(rb.rotation))
    rb.rotation = rb.rotation.Normalize()
    
    // 清除力
    rb.force = mgl32.Vec3{0, 0, 0}
    rb.torque = mgl32.Vec3{0, 0, 0}
    
    // 更新惯性张量
    rb.updateInertia()
}

// 更新惯性张量
func (rb *RigidBody) updateInertia() {
    if rb.collider == nil {
        return
    }
    
    // 计算局部惯性张量
    localInertia := rb.collider.GetInertia(rb.mass)
    
    // 变换到世界坐标系
    rotationMatrix := rb.rotation.Mat4().Mat3()
    rb.inertia = rotationMatrix.Mul3(localInertia).Mul3(rotationMatrix.Transpose())
    rb.inverseInertia = rb.inertia.Inverse()
}
```

### 4.2 碰撞检测

```go
// 碰撞器接口
type Collider interface {
    GetType() ColliderType
    GetBounds() AABB
    GetInertia(mass float32) mgl32.Mat3
    TestCollision(other Collider) *Contact
}

// 碰撞器类型
type ColliderType int

const (
    ColliderSphere ColliderType = iota
    ColliderBox
    ColliderCapsule
    ColliderMesh
)

// 球体碰撞器
type SphereCollider struct {
    radius float32
}

func (sc *SphereCollider) GetType() ColliderType {
    return ColliderSphere
}

func (sc *SphereCollider) GetBounds() AABB {
    return AABB{
        Min: mgl32.Vec3{-sc.radius, -sc.radius, -sc.radius},
        Max: mgl32.Vec3{sc.radius, sc.radius, sc.radius},
    }
}

func (sc *SphereCollider) GetInertia(mass float32) mgl32.Mat3 {
    inertia := 0.4 * mass * sc.radius * sc.radius
    return mgl32.Mat3{
        inertia, 0, 0,
        0, inertia, 0,
        0, 0, inertia,
    }
}

func (sc *SphereCollider) TestCollision(other Collider) *Contact {
    switch other.GetType() {
    case ColliderSphere:
        return sc.testSphereSphere(other.(*SphereCollider))
    case ColliderBox:
        return sc.testSphereBox(other.(*BoxCollider))
    default:
        return nil
    }
}

// 球体与球体碰撞
func (sc *SphereCollider) testSphereSphere(other *SphereCollider) *Contact {
    // 计算距离
    distance := sc.center.Sub(other.center).Len()
    radiusSum := sc.radius + other.radius
    
    if distance >= radiusSum {
        return nil
    }
    
    // 计算接触点
    normal := sc.center.Sub(other.center).Normalize()
    penetration := radiusSum - distance
    contactPoint := sc.center.Sub(normal.Mul(sc.radius))
    
    return &Contact{
        Normal:       normal,
        Penetration:  penetration,
        ContactPoint: contactPoint,
    }
}

// 盒体碰撞器
type BoxCollider struct {
    halfExtents mgl32.Vec3
}

func (bc *BoxCollider) GetType() ColliderType {
    return ColliderBox
}

func (bc *BoxCollider) GetBounds() AABB {
    return AABB{
        Min: bc.halfExtents.Mul(-1),
        Max: bc.halfExtents,
    }
}

func (bc *BoxCollider) GetInertia(mass float32) mgl32.Mat3 {
    x2 := bc.halfExtents.X * bc.halfExtents.X
    y2 := bc.halfExtents.Y * bc.halfExtents.Y
    z2 := bc.halfExtents.Z * bc.halfExtents.Z
    
    return mgl32.Mat3{
        mass * (y2 + z2) / 12, 0, 0,
        0, mass * (x2 + z2) / 12, 0,
        0, 0, mass * (x2 + y2) / 12,
    }
}
```

### 4.3 碰撞响应

```go
// 接触点
type Contact struct {
    Normal       mgl32.Vec3
    Penetration  float32
    ContactPoint mgl32.Vec3
    BodyA        *RigidBody
    BodyB        *RigidBody
    Friction     float32
    Restitution  float32
}

// 碰撞求解器
type ContactSolver struct {
    contacts     []*Contact
    iterations   int
    tolerance    float32
}

// 求解碰撞
func (cs *ContactSolver) Solve(deltaTime float32) {
    for iteration := 0; iteration < cs.iterations; iteration++ {
        maxPenetration := float32(0)
        
        for _, contact := range cs.contacts {
            penetration := cs.resolveContact(contact, deltaTime)
            if penetration > maxPenetration {
                maxPenetration = penetration
            }
        }
        
        // 收敛检查
        if maxPenetration < cs.tolerance {
            break
        }
    }
}

// 解析单个接触
func (cs *ContactSolver) resolveContact(contact *Contact, deltaTime float32) float32 {
    bodyA := contact.BodyA
    bodyB := contact.BodyB
    
    if bodyA.isStatic && bodyB.isStatic {
        return contact.Penetration
    }
    
    // 计算相对速度
    relativeVel := bodyB.velocity.Sub(bodyA.velocity)
    angularVelA := bodyA.angularVel.Cross(contact.ContactPoint.Sub(bodyA.position))
    angularVelB := bodyB.angularVel.Cross(contact.ContactPoint.Sub(bodyB.position))
    relativeVel = relativeVel.Add(angularVelB.Sub(angularVelA))
    
    // 计算冲量
    normalVel := mgl32.Dot(relativeVel, contact.Normal)
    
    if normalVel > 0 {
        return contact.Penetration
    }
    
    // 计算有效质量
    invMassA := bodyA.inverseMass
    invMassB := bodyB.inverseMass
    
    rA := contact.ContactPoint.Sub(bodyA.position)
    rB := contact.ContactPoint.Sub(bodyB.position)
    
    angularA := bodyA.inverseInertia.Mul3x1(rA.Cross(contact.Normal)).Cross(rA)
    angularB := bodyB.inverseInertia.Mul3x1(rB.Cross(contact.Normal)).Cross(rB)
    
    effectiveMass := invMassA + invMassB + mgl32.Dot(contact.Normal, angularA.Add(angularB))
    
    if effectiveMass <= 0 {
        return contact.Penetration
    }
    
    // 计算冲量
    restitution := contact.Restitution
    impulse := -(1 + restitution) * normalVel / effectiveMass
    
    // 应用冲量
    if !bodyA.isStatic {
        bodyA.velocity = bodyA.velocity.Sub(contact.Normal.Mul(impulse * invMassA))
        bodyA.angularVel = bodyA.angularVel.Add(bodyA.inverseInertia.Mul3x1(rA.Cross(contact.Normal.Mul(impulse))))
    }
    
    if !bodyB.isStatic {
        bodyB.velocity = bodyB.velocity.Add(contact.Normal.Mul(impulse * invMassB))
        bodyB.angularVel = bodyB.angularVel.Sub(bodyB.inverseInertia.Mul3x1(rB.Cross(contact.Normal.Mul(impulse))))
    }
    
    // 位置修正
    penetrationCorrection := contact.Penetration * 0.8 // 松弛因子
    correction := contact.Normal.Mul(penetrationCorrection / effectiveMass)
    
    if !bodyA.isStatic {
        bodyA.position = bodyA.position.Sub(correction.Mul(invMassA))
    }
    
    if !bodyB.isStatic {
        bodyB.position = bodyB.position.Add(correction.Mul(invMassB))
    }
    
    return contact.Penetration - penetrationCorrection
}
```

### 4.4 约束系统

```go
// 约束接口
type Constraint interface {
    GetType() ConstraintType
    Solve(deltaTime float32)
    GetBodies() []*RigidBody
}

// 约束类型
type ConstraintType int

const (
    ConstraintPoint ConstraintType = iota
    ConstraintHinge
    ConstraintPrismatic
    ConstraintDistance
)

// 点约束
type PointConstraint struct {
    bodyA       *RigidBody
    bodyB       *RigidBody
    pointA      mgl32.Vec3
    pointB      mgl32.Vec3
}

func (pc *PointConstraint) GetType() ConstraintType {
    return ConstraintPoint
}

func (pc *PointConstraint) GetBodies() []*RigidBody {
    return []*RigidBody{pc.bodyA, pc.bodyB}
}

func (pc *PointConstraint) Solve(deltaTime float32) {
    // 计算约束误差
    worldPointA := pc.bodyA.position.Add(pc.bodyA.rotation.Rotate(pc.pointA))
    worldPointB := pc.bodyB.position.Add(pc.bodyB.rotation.Rotate(pc.pointB))
    
    error := worldPointB.Sub(worldPointA)
    
    // 计算雅可比矩阵
    jacobianA := mgl32.Mat3{
        -1, 0, 0,
        0, -1, 0,
        0, 0, -1,
    }
    
    jacobianB := mgl32.Mat3{
        1, 0, 0,
        0, 1, 0,
        0, 0, 1,
    }
    
    // 计算有效质量
    invMassA := pc.bodyA.inverseMass
    invMassB := pc.bodyB.inverseMass
    
    effectiveMass := invMassA + invMassB
    
    if effectiveMass <= 0 {
        return
    }
    
    // 计算拉格朗日乘数
    lambda := error.Mul(-1.0 / effectiveMass)
    
    // 应用约束力
    if !pc.bodyA.isStatic {
        pc.bodyA.velocity = pc.bodyA.velocity.Add(lambda.Mul(invMassA))
    }
    
    if !pc.bodyB.isStatic {
        pc.bodyB.velocity = pc.bodyB.velocity.Sub(lambda.Mul(invMassB))
    }
}

// 距离约束
type DistanceConstraint struct {
    bodyA       *RigidBody
    bodyB       *RigidBody
    pointA      mgl32.Vec3
    pointB      mgl32.Vec3
    distance    float32
}

func (dc *DistanceConstraint) GetType() ConstraintType {
    return ConstraintDistance
}

func (dc *DistanceConstraint) GetBodies() []*RigidBody {
    return []*RigidBody{dc.bodyA, dc.bodyB}
}

func (dc *DistanceConstraint) Solve(deltaTime float32) {
    // 计算当前距离
    worldPointA := dc.bodyA.position.Add(dc.bodyA.rotation.Rotate(dc.pointA))
    worldPointB := dc.bodyB.position.Add(dc.bodyB.rotation.Rotate(dc.pointB))
    
    currentDistance := worldPointB.Sub(worldPointA).Len()
    
    if currentDistance < 0.0001 {
        return
    }
    
    // 计算约束误差
    error := currentDistance - dc.distance
    
    // 计算约束方向
    direction := worldPointB.Sub(worldPointA).Normalize()
    
    // 计算有效质量
    invMassA := dc.bodyA.inverseMass
    invMassB := dc.bodyB.inverseMass
    
    effectiveMass := invMassA + invMassB
    
    if effectiveMass <= 0 {
        return
    }
    
    // 计算拉格朗日乘数
    lambda := -error / effectiveMass
    
    // 应用约束力
    if !dc.bodyA.isStatic {
        dc.bodyA.velocity = dc.bodyA.velocity.Sub(direction.Mul(lambda * invMassA))
    }
    
    if !dc.bodyB.isStatic {
        dc.bodyB.velocity = dc.bodyB.velocity.Add(direction.Mul(lambda * invMassB))
    }
}
```

## 5. Go语言实现

### 5.1 物理引擎主类

```go
// 物理引擎
type PhysicsEngine struct {
    bodies       []*RigidBody
    constraints  []Constraint
    gravity      mgl32.Vec3
    timeStep     float32
    iterations   int
    
    broadPhase   *BroadPhase
    narrowPhase  *NarrowPhase
    solver       *ContactSolver
    
    // 性能统计
    stats        *PhysicsStats
}

// 物理统计
type PhysicsStats struct {
    bodyCount    int
    contactCount int
    constraintCount int
    updateTime   time.Duration
    collisionTime time.Duration
    solverTime   time.Duration
}

// 更新物理引擎
func (pe *PhysicsEngine) Update(deltaTime float32) {
    startTime := time.Now()
    
    // 应用重力
    pe.applyGravity()
    
    // 更新刚体
    pe.updateBodies(deltaTime)
    
    // 碰撞检测
    collisionStart := time.Now()
    contacts := pe.detectCollisions()
    pe.stats.collisionTime = time.Since(collisionStart)
    
    // 约束求解
    constraintStart := time.Now()
    pe.solveConstraints(deltaTime)
    pe.stats.solverTime = time.Since(constraintStart)
    
    // 碰撞响应
    pe.solver.contacts = contacts
    pe.solver.Solve(deltaTime)
    
    // 更新统计
    pe.stats.bodyCount = len(pe.bodies)
    pe.stats.contactCount = len(contacts)
    pe.stats.constraintCount = len(pe.constraints)
    pe.stats.updateTime = time.Since(startTime)
}

// 应用重力
func (pe *PhysicsEngine) applyGravity() {
    for _, body := range pe.bodies {
        if !body.isStatic {
            body.force = body.force.Add(pe.gravity.Mul(body.mass))
        }
    }
}

// 更新刚体
func (pe *PhysicsEngine) updateBodies(deltaTime float32) {
    for _, body := range pe.bodies {
        body.Update(deltaTime)
    }
}

// 检测碰撞
func (pe *PhysicsEngine) detectCollisions() []*Contact {
    // 宽相检测
    pairs := pe.broadPhase.GetPotentialPairs()
    
    // 窄相检测
    contacts := make([]*Contact, 0)
    for _, pair := range pairs {
        contact := pe.narrowPhase.TestCollision(pair.bodyA, pair.bodyB)
        if contact != nil {
            contacts = append(contacts, contact)
        }
    }
    
    return contacts
}

// 求解约束
func (pe *PhysicsEngine) solveConstraints(deltaTime float32) {
    for iteration := 0; iteration < pe.iterations; iteration++ {
        for _, constraint := range pe.constraints {
            constraint.Solve(deltaTime)
        }
    }
}
```

### 5.2 空间分割

```go
// 宽相检测
type BroadPhase struct {
    grid        *SpatialGrid
    pairs       []BodyPair
}

// 空间网格
type SpatialGrid struct {
    cellSize    float32
    cells       map[GridKey][]*RigidBody
}

// 网格键
type GridKey struct {
    X, Y, Z int
}

// 获取潜在碰撞对
func (bp *BroadPhase) GetPotentialPairs() []BodyPair {
    bp.pairs = bp.pairs[:0]
    
    // 遍历所有网格单元
    for _, bodies := range bp.grid.cells {
        if len(bodies) < 2 {
            continue
        }
        
        // 检查同一单元内的所有物体对
        for i := 0; i < len(bodies); i++ {
            for j := i + 1; j < len(bodies); j++ {
                bodyA := bodies[i]
                bodyB := bodies[j]
                
                // 检查AABB重叠
                if bodyA.collider.GetBounds().Overlaps(bodyB.collider.GetBounds()) {
                    bp.pairs = append(bp.pairs, BodyPair{
                        bodyA: bodyA,
                        bodyB: bodyB,
                    })
                }
            }
        }
    }
    
    return bp.pairs
}

// 更新网格
func (bp *BroadPhase) UpdateGrid() {
    // 清空网格
    bp.grid.cells = make(map[GridKey][]*RigidBody)
    
    // 重新插入所有物体
    for _, body := range bp.bodies {
        bounds := body.collider.GetBounds()
        keys := bp.grid.GetKeys(bounds)
        
        for _, key := range keys {
            bp.grid.cells[key] = append(bp.grid.cells[key], body)
        }
    }
}

// 获取包围盒覆盖的网格键
func (sg *SpatialGrid) GetKeys(bounds AABB) []GridKey {
    minKey := GridKey{
        X: int(bounds.Min.X / sg.cellSize),
        Y: int(bounds.Min.Y / sg.cellSize),
        Z: int(bounds.Min.Z / sg.cellSize),
    }
    
    maxKey := GridKey{
        X: int(bounds.Max.X / sg.cellSize),
        Y: int(bounds.Max.Y / sg.cellSize),
        Z: int(bounds.Max.Z / sg.cellSize),
    }
    
    keys := make([]GridKey, 0)
    for x := minKey.X; x <= maxKey.X; x++ {
        for y := minKey.Y; y <= maxKey.Y; y++ {
            for z := minKey.Z; z <= maxKey.Z; z++ {
                keys = append(keys, GridKey{X: x, Y: y, Z: z})
            }
        }
    }
    
    return keys
}
```

### 5.3 窄相检测

```go
// 窄相检测
type NarrowPhase struct {
    algorithms  map[ColliderType]map[ColliderType]CollisionAlgorithm
}

// 碰撞算法接口
type CollisionAlgorithm interface {
    TestCollision(colliderA, colliderB Collider) *Contact
}

// 球体-球体碰撞算法
type SphereSphereAlgorithm struct{}

func (ssa *SphereSphereAlgorithm) TestCollision(colliderA, colliderB Collider) *Contact {
    sphereA := colliderA.(*SphereCollider)
    sphereB := colliderB.(*SphereCollider)
    
    return sphereA.TestCollision(sphereB)
}

// 球体-盒体碰撞算法
type SphereBoxAlgorithm struct{}

func (sba *SphereBoxAlgorithm) TestCollision(colliderA, colliderB Collider) *Contact {
    sphere := colliderA.(*SphereCollider)
    box := colliderB.(*BoxCollider)
    
    // 计算球心在盒体坐标系中的位置
    localCenter := box.transform.Inverse().Mul4x1(sphere.center)
    
    // 计算最近点
    closest := mgl32.Vec3{
        mgl32.Clamp(localCenter.X, -box.halfExtents.X, box.halfExtents.X),
        mgl32.Clamp(localCenter.Y, -box.halfExtents.Y, box.halfExtents.Y),
        mgl32.Clamp(localCenter.Z, -box.halfExtents.Z, box.halfExtents.Z),
    }
    
    // 计算距离
    distance := localCenter.Sub(closest).Len()
    
    if distance >= sphere.radius {
        return nil
    }
    
    // 计算接触信息
    worldClosest := box.transform.Mul4x1(closest)
    normal := sphere.center.Sub(worldClosest).Normalize()
    penetration := sphere.radius - distance
    
    return &Contact{
        Normal:       normal,
        Penetration:  penetration,
        ContactPoint: worldClosest,
    }
}
```

### 5.4 软体物理

```go
// 软体
type SoftBody struct {
    particles    []*Particle
    springs      []*Spring
    triangles    []*Triangle
    
    // 物理属性
    mass         float32
    stiffness    float32
    damping      float32
    gravity      mgl32.Vec3
}

// 粒子
type Particle struct {
    position     mgl32.Vec3
    velocity     mgl32.Vec3
    force        mgl32.Vec3
    mass         float32
    inverseMass  float32
    fixed        bool
}

// 弹簧
type Spring struct {
    particleA    *Particle
    particleB    *Particle
    restLength   float32
    stiffness    float32
    damping      float32
}

// 三角形
type Triangle struct {
    particles    [3]*Particle
    normal       mgl32.Vec3
    area         float32
}

// 更新软体
func (sb *SoftBody) Update(deltaTime float32) {
    // 应用重力
    for _, particle := range sb.particles {
        if !particle.fixed {
            particle.force = particle.force.Add(sb.gravity.Mul(particle.mass))
        }
    }
    
    // 计算弹簧力
    for _, spring := range sb.springs {
        sb.calculateSpringForce(spring)
    }
    
    // 计算空气阻力
    for _, triangle := range sb.triangles {
        sb.calculateAirResistance(triangle)
    }
    
    // 积分运动
    for _, particle := range sb.particles {
        if !particle.fixed {
            sb.integrateParticle(particle, deltaTime)
        }
    }
    
    // 更新三角形法线
    for _, triangle := range sb.triangles {
        sb.updateTriangleNormal(triangle)
    }
}

// 计算弹簧力
func (sb *SoftBody) calculateSpringForce(spring *Spring) {
    delta := spring.particleB.position.Sub(spring.particleA.position)
    distance := delta.Len()
    
    if distance < 0.0001 {
        return
    }
    
    // 弹簧力
    stretch := distance - spring.restLength
    springForce := delta.Normalize().Mul(stretch * spring.stiffness)
    
    // 阻尼力
    relativeVel := spring.particleB.velocity.Sub(spring.particleA.velocity)
    dampingForce := delta.Normalize().Mul(mgl32.Dot(relativeVel, delta.Normalize()) * spring.damping)
    
    // 应用力
    totalForce := springForce.Add(dampingForce)
    
    if !spring.particleA.fixed {
        spring.particleA.force = spring.particleA.force.Add(totalForce)
    }
    
    if !spring.particleB.fixed {
        spring.particleB.force = spring.particleB.force.Sub(totalForce)
    }
}

// 积分粒子运动
func (sb *SoftBody) integrateParticle(particle *Particle, deltaTime float32) {
    // 计算加速度
    acceleration := particle.force.Mul(particle.inverseMass)
    
    // 积分速度和位置
    particle.velocity = particle.velocity.Add(acceleration.Mul(deltaTime))
    particle.position = particle.position.Add(particle.velocity.Mul(deltaTime))
    
    // 清除力
    particle.force = mgl32.Vec3{0, 0, 0}
}
```

## 6. 性能优化

### 6.1 空间优化

**定理 6.1.1** (空间分割复杂度)
使用空间分割可以将碰撞检测复杂度从 ```latex
O(n^2)
``` 降低到 ```latex
O(n \log n)

```。

**实现**:

```go
// 八叉树
type Octree struct {
    bounds      AABB
    children    [8]*Octree
    bodies      []*RigidBody
    maxBodies   int
    maxDepth    int
}

// 插入物体
func (ot *Octree) Insert(body *RigidBody) {
    if !ot.bounds.Contains(body.collider.GetBounds()) {
        return
    }
    
    if len(ot.bodies) < ot.maxBodies || ot.maxDepth <= 0 {
        ot.bodies = append(ot.bodies, body)
        return
    }
    
    // 分割节点
    if ot.children[0] == nil {
        ot.split()
    }
    
    // 插入到子节点
    for _, child := range ot.children {
        child.Insert(body)
    }
}

// 分割节点
func (ot *Octree) split() {
    center := ot.bounds.GetCenter()
    halfSize := ot.bounds.GetSize().Mul(0.5)
    
    for i := 0; i < 8; i++ {
        childBounds := AABB{
            Min: center.Add(mgl32.Vec3{
                float32(i&1) * halfSize.X,
                float32((i>>1)&1) * halfSize.Y,
                float32((i>>2)&1) * halfSize.Z,
            }).Sub(halfSize),
            Max: center.Add(mgl32.Vec3{
                float32(i&1) * halfSize.X,
                float32((i>>1)&1) * halfSize.Y,
                float32((i>>2)&1) * halfSize.Z,
            }),
        }
        
        ot.children[i] = &Octree{
            bounds:    childBounds,
            maxBodies: ot.maxBodies,
            maxDepth:  ot.maxDepth - 1,
        }
    }
}
```

### 6.2 时间优化

```go
// 时间步长自适应
type AdaptiveTimeStep struct {
    minStep      float32
    maxStep      float32
    targetError  float32
    currentStep  float32
}

// 计算最优时间步长
func (ats *AdaptiveTimeStep) CalculateOptimalStep(error float32) float32 {
    if error < ats.targetError {
        // 误差较小，可以增大时间步长
        ats.currentStep = mgl32.Min(ats.currentStep*1.1, ats.maxStep)
    } else {
        // 误差较大，需要减小时间步长
        ats.currentStep = mgl32.Max(ats.currentStep*0.9, ats.minStep)
    }
    
    return ats.currentStep
}
```

### 6.3 内存优化

```go
// 对象池
type PhysicsObjectPool struct {
    bodyPool     *ObjectPool[RigidBody]
    contactPool  *ObjectPool[Contact]
    constraintPool *ObjectPool[Constraint]
}

// 创建对象池
func NewPhysicsObjectPool() *PhysicsObjectPool {
    return &PhysicsObjectPool{
        bodyPool: &ObjectPool[RigidBody]{
            factory: func() RigidBody { return RigidBody{} },
            reset: func(body RigidBody) RigidBody {
                // 重置刚体状态
                body.force = mgl32.Vec3{0, 0, 0}
                body.torque = mgl32.Vec3{0, 0, 0}
                return body
            },
        },
        contactPool: &ObjectPool[Contact]{
            factory: func() Contact { return Contact{} },
            reset: func(contact Contact) Contact {
                // 重置接触点
                contact.Penetration = 0
                return contact
            },
        },
    }
}
```

## 7. 总结

物理引擎是现代游戏引擎的核心组件，需要提供准确、稳定、高性能的物理模拟。通过合理的算法设计和优化策略，可以构建出满足游戏需求的物理引擎。

### 7.1 关键特性

- **准确的物理模拟**: 支持刚体动力学、碰撞检测、约束求解
- **高性能计算**: 使用空间分割、对象池等优化技术
- **稳定的数值计算**: 采用合适的积分方法和求解器
- **可扩展架构**: 支持自定义碰撞器和约束
- **易用性**: 提供简单的API接口

### 7.2 扩展方向

- **软体物理**: 支持布料、绳索、流体等软体模拟
- **粒子系统**: 支持大规模粒子效果
- **破坏系统**: 支持物体的破碎和变形
- **车辆物理**: 支持车辆动力学模拟
- **角色控制器**: 支持角色动画和物理交互

---

**参考文献**:

1. Physics for Game Developers, David M. Bourg
2. Game Physics Engine Development, Ian Millington
3. Real-Time Collision Detection, Christer Ericson
4. Computational Physics, Nicholas J. Giordano

**相关链接**:

- [01-游戏引擎架构](./01-Game-Engine-Architecture.md)
- [02-网络游戏服务器](./02-Network-Game-Server.md)
- [03-实时渲染系统](./03-Real-time-Rendering-System.md)
