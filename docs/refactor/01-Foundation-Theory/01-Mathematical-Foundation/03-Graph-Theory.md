# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 图的定义](#11-图的定义)
    - [1.2 图的类型](#12-图的类型)
    - [1.3 图的基本性质](#13-图的基本性质)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 图的数学表示](#21-图的数学表示)
    - [2.2 图论定理](#22-图论定理)
  - [3. 图算法](#3-图算法)
    - [3.1 遍历算法](#31-遍历算法)
    - [3.2 最短路径算法](#32-最短路径算法)
    - [3.3 最小生成树算法](#33-最小生成树算法)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 图的数据结构](#41-图的数据结构)
    - [4.2 基础算法实现](#42-基础算法实现)
    - [4.3 高级算法实现](#43-高级算法实现)
  - [5. 应用示例](#5-应用示例)
    - [5.1 社交网络分析](#51-社交网络分析)
    - [5.2 路由算法](#52-路由算法)
    - [5.3 任务调度](#53-任务调度)
  - [总结](#总结)

## 1. 基本概念

### 1.1 图的定义

**定义 1.1**: 图 $G = (V, E)$ 由顶点集 $V$ 和边集 $E$ 组成，其中 $E \subseteq V \times V$。

**形式化表达**:
- 顶点集：$V = \{v_1, v_2, \ldots, v_n\}$
- 边集：$E = \{(u,v) \mid u,v \in V\}$
- 图的阶：$|V|$ (顶点数)
- 图的大小：$|E|$ (边数)

### 1.2 图的类型

**定义 1.2**: 图的分类

1. **无向图**: 边没有方向，$(u,v) = (v,u)$
2. **有向图**: 边有方向，$(u,v) \neq (v,u)$
3. **加权图**: 边有权重，$E \subseteq V \times V \times \mathbb{R}$
4. **多重图**: 允许重边和自环
5. **简单图**: 无重边无自环

### 1.3 图的基本性质

**定义 1.3**: 图的基本概念

- **度**: 与顶点相连的边数
- **路径**: 顶点序列 $v_1, v_2, \ldots, v_k$，其中 $(v_i, v_{i+1}) \in E$
- **连通性**: 任意两顶点间存在路径
- **完全图**: 任意两顶点间都有边相连

## 2. 形式化定义

### 2.1 图的数学表示

**定义 2.1**: 邻接矩阵表示

对于图 $G = (V, E)$，邻接矩阵 $A$ 定义为：

$$A_{ij} = \begin{cases}
1 & \text{if } (v_i, v_j) \in E \\
0 & \text{otherwise}
\end{cases}$$

**定义 2.2**: 邻接表表示

每个顶点维护一个邻接顶点列表：

$$\text{Adj}[v] = \{u \mid (v, u) \in E\}$$

### 2.2 图论定理

**定理 2.1** (握手定理): 对于任意图 $G = (V, E)$：

$$\sum_{v \in V} \deg(v) = 2|E|$$

**证明**:
每条边贡献给两个顶点的度数，因此所有顶点的度数之和等于边数的两倍。

**定理 2.2** (欧拉公式): 对于连通平面图：

$$|V| - |E| + |F| = 2$$

其中 $|F|$ 是面数。

## 3. 图算法

### 3.1 遍历算法

**算法 3.1** (深度优先搜索 - DFS):

```pseudocode
DFS(G, v):
    visited[v] = true
    for each neighbor u of v:
        if not visited[u]:
            DFS(G, u)
```

**算法 3.2** (广度优先搜索 - BFS):

```pseudocode
BFS(G, s):
    queue = [s]
    visited[s] = true
    while queue is not empty:
        v = queue.dequeue()
        for each neighbor u of v:
            if not visited[u]:
                visited[u] = true
                queue.enqueue(u)
```

### 3.2 最短路径算法

**算法 3.3** (Dijkstra算法):

```pseudocode
Dijkstra(G, s):
    dist[s] = 0
    for all v ≠ s: dist[v] = ∞
    Q = V
    while Q is not empty:
        u = extract-min(Q)
        for each neighbor v of u:
            if dist[v] > dist[u] + weight(u,v):
                dist[v] = dist[u] + weight(u,v)
```

### 3.3 最小生成树算法

**算法 3.4** (Kruskal算法):

```pseudocode
Kruskal(G):
    T = ∅
    sort edges by weight
    for each edge (u,v) in sorted order:
        if adding (u,v) doesn't create cycle:
            T = T ∪ {(u,v)}
    return T
```

## 4. Go语言实现

### 4.1 图的数据结构

```go
// Vertex 顶点
type Vertex[T comparable] struct {
    ID       T
    Data     interface{}
    Weight   float64
}

// Edge 边
type Edge[T comparable] struct {
    From     T
    To       T
    Weight   float64
}

// Graph 图接口
type Graph[T comparable] interface {
    AddVertex(vertex T, data interface{})
    RemoveVertex(vertex T)
    AddEdge(from, to T, weight float64)
    RemoveEdge(from, to T)
    GetNeighbors(vertex T) []T
    GetEdgeWeight(from, to T) (float64, bool)
    GetAllVertices() []T
    GetAllEdges() []Edge[T]
    IsDirected() bool
}

// AdjacencyListGraph 邻接表图
type AdjacencyListGraph[T comparable] struct {
    vertices map[T]*Vertex[T]
    edges    map[T]map[T]float64
    directed bool
}

// NewAdjacencyListGraph 创建邻接表图
func NewAdjacencyListGraph[T comparable](directed bool) *AdjacencyListGraph[T] {
    return &AdjacencyListGraph[T]{
        vertices: make(map[T]*Vertex[T]),
        edges:    make(map[T]map[T]float64),
        directed: directed,
    }
}

// AddVertex 添加顶点
func (g *AdjacencyListGraph[T]) AddVertex(vertex T, data interface{}) {
    g.vertices[vertex] = &Vertex[T]{
        ID:     vertex,
        Data:   data,
        Weight: 0,
    }
    if g.edges[vertex] == nil {
        g.edges[vertex] = make(map[T]float64)
    }
}

// RemoveVertex 移除顶点
func (g *AdjacencyListGraph[T]) RemoveVertex(vertex T) {
    delete(g.vertices, vertex)
    delete(g.edges, vertex)
    
    // 移除指向该顶点的边
    for _, neighbors := range g.edges {
        delete(neighbors, vertex)
    }
}

// AddEdge 添加边
func (g *AdjacencyListGraph[T]) AddEdge(from, to T, weight float64) {
    if g.edges[from] == nil {
        g.edges[from] = make(map[T]float64)
    }
    g.edges[from][to] = weight
    
    if !g.directed {
        if g.edges[to] == nil {
            g.edges[to] = make(map[T]float64)
        }
        g.edges[to][from] = weight
    }
}

// RemoveEdge 移除边
func (g *AdjacencyListGraph[T]) RemoveEdge(from, to T) {
    if neighbors, exists := g.edges[from]; exists {
        delete(neighbors, to)
    }
    
    if !g.directed {
        if neighbors, exists := g.edges[to]; exists {
            delete(neighbors, from)
        }
    }
}

// GetNeighbors 获取邻居
func (g *AdjacencyListGraph[T]) GetNeighbors(vertex T) []T {
    neighbors := make([]T, 0)
    if vertexEdges, exists := g.edges[vertex]; exists {
        for neighbor := range vertexEdges {
            neighbors = append(neighbors, neighbor)
        }
    }
    return neighbors
}

// GetEdgeWeight 获取边权重
func (g *AdjacencyListGraph[T]) GetEdgeWeight(from, to T) (float64, bool) {
    if neighbors, exists := g.edges[from]; exists {
        if weight, exists := neighbors[to]; exists {
            return weight, true
        }
    }
    return 0, false
}

// GetAllVertices 获取所有顶点
func (g *AdjacencyListGraph[T]) GetAllVertices() []T {
    vertices := make([]T, 0, len(g.vertices))
    for vertex := range g.vertices {
        vertices = append(vertices, vertex)
    }
    return vertices
}

// GetAllEdges 获取所有边
func (g *AdjacencyListGraph[T]) GetAllEdges() []Edge[T] {
    edges := make([]Edge[T], 0)
    for from, neighbors := range g.edges {
        for to, weight := range neighbors {
            edges = append(edges, Edge[T]{
                From:   from,
                To:     to,
                Weight: weight,
            })
        }
    }
    return edges
}

// IsDirected 是否为有向图
func (g *AdjacencyListGraph[T]) IsDirected() bool {
    return g.directed
}
```

### 4.2 基础算法实现

```go
// GraphTraversal 图遍历
type GraphTraversal[T comparable] struct {
    graph Graph[T]
}

// NewGraphTraversal 创建图遍历器
func NewGraphTraversal[T comparable](graph Graph[T]) *GraphTraversal[T] {
    return &GraphTraversal[T]{
        graph: graph,
    }
}

// DFS 深度优先搜索
func (gt *GraphTraversal[T]) DFS(start T) []T {
    visited := make(map[T]bool)
    result := make([]T, 0)
    
    gt.dfsHelper(start, visited, &result)
    return result
}

// dfsHelper DFS辅助函数
func (gt *GraphTraversal[T]) dfsHelper(vertex T, visited map[T]bool, result *[]T) {
    visited[vertex] = true
    *result = append(*result, vertex)
    
    neighbors := gt.graph.GetNeighbors(vertex)
    for _, neighbor := range neighbors {
        if !visited[neighbor] {
            gt.dfsHelper(neighbor, visited, result)
        }
    }
}

// BFS 广度优先搜索
func (gt *GraphTraversal[T]) BFS(start T) []T {
    visited := make(map[T]bool)
    queue := []T{start}
    result := make([]T, 0)
    
    visited[start] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        result = append(result, current)
        
        neighbors := gt.graph.GetNeighbors(current)
        for _, neighbor := range neighbors {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    
    return result
}

// TopologicalSort 拓扑排序
func (gt *GraphTraversal[T]) TopologicalSort() ([]T, error) {
    if !gt.graph.IsDirected() {
        return nil, fmt.Errorf("topological sort requires directed graph")
    }
    
    visited := make(map[T]bool)
    temp := make(map[T]bool)
    result := make([]T, 0)
    
    vertices := gt.graph.GetAllVertices()
    for _, vertex := range vertices {
        if !visited[vertex] {
            if err := gt.topologicalSortHelper(vertex, visited, temp, &result); err != nil {
                return nil, err
            }
        }
    }
    
    // 反转结果
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    
    return result, nil
}

// topologicalSortHelper 拓扑排序辅助函数
func (gt *GraphTraversal[T]) topologicalSortHelper(vertex T, visited, temp map[T]bool, result *[]T) error {
    if temp[vertex] {
        return fmt.Errorf("cycle detected")
    }
    
    if visited[vertex] {
        return nil
    }
    
    temp[vertex] = true
    
    neighbors := gt.graph.GetNeighbors(vertex)
    for _, neighbor := range neighbors {
        if err := gt.topologicalSortHelper(neighbor, visited, temp, result); err != nil {
            return err
        }
    }
    
    temp[vertex] = false
    visited[vertex] = true
    *result = append(*result, vertex)
    
    return nil
}
```

### 4.3 高级算法实现

```go
// ShortestPath 最短路径算法
type ShortestPath[T comparable] struct {
    graph Graph[T]
}

// NewShortestPath 创建最短路径算法
func NewShortestPath[T comparable](graph Graph[T]) *ShortestPath[T] {
    return &ShortestPath[T]{
        graph: graph,
    }
}

// Dijkstra Dijkstra算法
func (sp *ShortestPath[T]) Dijkstra(start T) (map[T]float64, map[T]T) {
    vertices := sp.graph.GetAllVertices()
    dist := make(map[T]float64)
    prev := make(map[T]T)
    unvisited := make(map[T]bool)
    
    // 初始化
    for _, vertex := range vertices {
        dist[vertex] = math.Inf(1)
        unvisited[vertex] = true
    }
    dist[start] = 0
    
    for len(unvisited) > 0 {
        // 找到未访问顶点中距离最小的
        var current T
        minDist := math.Inf(1)
        for vertex := range unvisited {
            if dist[vertex] < minDist {
                minDist = dist[vertex]
                current = vertex
            }
        }
        
        if minDist == math.Inf(1) {
            break // 无法到达的顶点
        }
        
        delete(unvisited, current)
        
        // 更新邻居距离
        neighbors := sp.graph.GetNeighbors(current)
        for _, neighbor := range neighbors {
            if !unvisited[neighbor] {
                continue
            }
            
            weight, exists := sp.graph.GetEdgeWeight(current, neighbor)
            if !exists {
                continue
            }
            
            alt := dist[current] + weight
            if alt < dist[neighbor] {
                dist[neighbor] = alt
                prev[neighbor] = current
            }
        }
    }
    
    return dist, prev
}

// BellmanFord Bellman-Ford算法
func (sp *ShortestPath[T]) BellmanFord(start T) (map[T]float64, map[T]T, error) {
    vertices := sp.graph.GetAllVertices()
    dist := make(map[T]float64)
    prev := make(map[T]T)
    
    // 初始化
    for _, vertex := range vertices {
        dist[vertex] = math.Inf(1)
    }
    dist[start] = 0
    
    // 松弛操作
    for i := 0; i < len(vertices)-1; i++ {
        edges := sp.graph.GetAllEdges()
        for _, edge := range edges {
            if dist[edge.From] != math.Inf(1) {
                alt := dist[edge.From] + edge.Weight
                if alt < dist[edge.To] {
                    dist[edge.To] = alt
                    prev[edge.To] = edge.From
                }
            }
        }
    }
    
    // 检查负环
    edges := sp.graph.GetAllEdges()
    for _, edge := range edges {
        if dist[edge.From] != math.Inf(1) {
            alt := dist[edge.From] + edge.Weight
            if alt < dist[edge.To] {
                return nil, nil, fmt.Errorf("negative cycle detected")
            }
        }
    }
    
    return dist, prev, nil
}

// MinimumSpanningTree 最小生成树算法
type MinimumSpanningTree[T comparable] struct {
    graph Graph[T]
}

// NewMinimumSpanningTree 创建最小生成树算法
func NewMinimumSpanningTree[T comparable](graph Graph[T]) *MinimumSpanningTree[T] {
    return &MinimumSpanningTree[T]{
        graph: graph,
    }
}

// Kruskal Kruskal算法
func (mst *MinimumSpanningTree[T]) Kruskal() ([]Edge[T], error) {
    if mst.graph.IsDirected() {
        return nil, fmt.Errorf("Kruskal algorithm requires undirected graph")
    }
    
    edges := mst.graph.GetAllEdges()
    vertices := mst.graph.GetAllVertices()
    
    // 排序边
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })
    
    // 并查集
    uf := NewUnionFind[T]()
    for _, vertex := range vertices {
        uf.MakeSet(vertex)
    }
    
    result := make([]Edge[T], 0)
    
    for _, edge := range edges {
        if uf.Find(edge.From) != uf.Find(edge.To) {
            result = append(result, edge)
            uf.Union(edge.From, edge.To)
        }
    }
    
    return result, nil
}

// UnionFind 并查集
type UnionFind[T comparable] struct {
    parent map[T]T
    rank   map[T]int
}

// NewUnionFind 创建并查集
func NewUnionFind[T comparable]() *UnionFind[T] {
    return &UnionFind[T]{
        parent: make(map[T]T),
        rank:   make(map[T]int),
    }
}

// MakeSet 创建集合
func (uf *UnionFind[T]) MakeSet(x T) {
    uf.parent[x] = x
    uf.rank[x] = 0
}

// Find 查找根节点
func (uf *UnionFind[T]) Find(x T) T {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // 路径压缩
    }
    return uf.parent[x]
}

// Union 合并集合
func (uf *UnionFind[T]) Union(x, y T) {
    rootX := uf.Find(x)
    rootY := uf.Find(y)
    
    if rootX == rootY {
        return
    }
    
    if uf.rank[rootX] < uf.rank[rootY] {
        uf.parent[rootX] = rootY
    } else if uf.rank[rootX] > uf.rank[rootY] {
        uf.parent[rootY] = rootX
    } else {
        uf.parent[rootY] = rootX
        uf.rank[rootX]++
    }
}
```

## 5. 应用示例

### 5.1 社交网络分析

```go
// SocialNetwork 社交网络
type SocialNetwork struct {
    graph *AdjacencyListGraph[string]
}

// NewSocialNetwork 创建社交网络
func NewSocialNetwork() *SocialNetwork {
    return &SocialNetwork{
        graph: NewAdjacencyListGraph[string](false),
    }
}

// AddUser 添加用户
func (sn *SocialNetwork) AddUser(userID string, name string) {
    sn.graph.AddVertex(userID, name)
}

// AddFriendship 添加好友关系
func (sn *SocialNetwork) AddFriendship(user1, user2 string) {
    sn.graph.AddEdge(user1, user2, 1.0)
}

// GetFriends 获取好友列表
func (sn *SocialNetwork) GetFriends(userID string) []string {
    return sn.graph.GetNeighbors(userID)
}

// GetMutualFriends 获取共同好友
func (sn *SocialNetwork) GetMutualFriends(user1, user2 string) []string {
    friends1 := NewSet[string]()
    friends2 := NewSet[string]()
    
    for _, friend := range sn.GetFriends(user1) {
        friends1.Add(friend)
    }
    
    for _, friend := range sn.GetFriends(user2) {
        friends2.Add(friend)
    }
    
    mutual := friends1.Intersection(friends2)
    return mutual.ToSlice()
}

// GetShortestPath 获取最短路径（最少中间人）
func (sn *SocialNetwork) GetShortestPath(from, to string) ([]string, error) {
    sp := NewShortestPath(sn.graph)
    dist, prev := sp.Dijkstra(from)
    
    if dist[to] == math.Inf(1) {
        return nil, fmt.Errorf("no path exists")
    }
    
    // 重建路径
    path := make([]string, 0)
    current := to
    for current != from {
        path = append([]string{current}, path...)
        current = prev[current]
    }
    path = append([]string{from}, path...)
    
    return path, nil
}
```

### 5.2 路由算法

```go
// NetworkNode 网络节点
type NetworkNode struct {
    ID       string
    IP       string
    Latency  float64
}

// NetworkGraph 网络图
type NetworkGraph struct {
    graph *AdjacencyListGraph[string]
    nodes map[string]*NetworkNode
}

// NewNetworkGraph 创建网络图
func NewNetworkGraph() *NetworkGraph {
    return &NetworkGraph{
        graph: NewAdjacencyListGraph[string](true),
        nodes: make(map[string]*NetworkNode),
    }
}

// AddNode 添加节点
func (ng *NetworkGraph) AddNode(id, ip string, latency float64) {
    node := &NetworkNode{
        ID:      id,
        IP:      ip,
        Latency: latency,
    }
    ng.nodes[id] = node
    ng.graph.AddVertex(id, node)
}

// AddLink 添加链路
func (ng *NetworkGraph) AddLink(from, to string, bandwidth float64) {
    // 使用带宽的倒数作为权重（带宽越大，权重越小）
    weight := 1.0 / bandwidth
    ng.graph.AddEdge(from, to, weight)
}

// FindOptimalRoute 寻找最优路由
func (ng *NetworkGraph) FindOptimalRoute(source, destination string) ([]string, float64, error) {
    sp := NewShortestPath(ng.graph)
    dist, prev := sp.Dijkstra(source)
    
    if dist[destination] == math.Inf(1) {
        return nil, 0, fmt.Errorf("no route exists")
    }
    
    // 重建路径
    path := make([]string, 0)
    current := destination
    for current != source {
        path = append([]string{current}, path...)
        current = prev[current]
    }
    path = append([]string{source}, path...)
    
    return path, dist[destination], nil
}

// GetNetworkTopology 获取网络拓扑
func (ng *NetworkGraph) GetNetworkTopology() map[string][]string {
    topology := make(map[string][]string)
    vertices := ng.graph.GetAllVertices()
    
    for _, vertex := range vertices {
        topology[vertex] = ng.graph.GetNeighbors(vertex)
    }
    
    return topology
}
```

### 5.3 任务调度

```go
// Task 任务
type Task struct {
    ID       string
    Duration int
    Priority int
    Dependencies []string
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
    graph *AdjacencyListGraph[string]
    tasks map[string]*Task
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler() *TaskScheduler {
    return &TaskScheduler{
        graph: NewAdjacencyListGraph[string](true),
        tasks: make(map[string]*Task),
    }
}

// AddTask 添加任务
func (ts *TaskScheduler) AddTask(task *Task) {
    ts.tasks[task.ID] = task
    ts.graph.AddVertex(task.ID, task)
    
    // 添加依赖关系
    for _, dep := range task.Dependencies {
        ts.graph.AddEdge(dep, task.ID, float64(task.Duration))
    }
}

// GetSchedule 获取调度顺序
func (ts *TaskScheduler) GetSchedule() ([]string, error) {
    traversal := NewGraphTraversal(ts.graph)
    return traversal.TopologicalSort()
}

// CalculateCriticalPath 计算关键路径
func (ts *TaskScheduler) CalculateCriticalPath() ([]string, int) {
    vertices := ts.graph.GetAllVertices()
    
    // 计算最早开始时间
    earliest := make(map[string]int)
    for _, vertex := range vertices {
        earliest[vertex] = 0
    }
    
    // 拓扑排序
    schedule, err := ts.GetSchedule()
    if err != nil {
        return nil, 0
    }
    
    // 计算最早开始时间
    for _, taskID := range schedule {
        task := ts.tasks[taskID]
        maxEarliest := 0
        
        neighbors := ts.graph.GetNeighbors(taskID)
        for _, neighbor := range neighbors {
            if earliest[neighbor]+task.Duration > maxEarliest {
                maxEarliest = earliest[neighbor] + task.Duration
            }
        }
        
        earliest[taskID] = maxEarliest
    }
    
    // 找到最晚完成的任务
    maxTime := 0
    var lastTask string
    for taskID, time := range earliest {
        if time > maxTime {
            maxTime = time
            lastTask = taskID
        }
    }
    
    // 重建关键路径
    criticalPath := make([]string, 0)
    current := lastTask
    for current != "" {
        criticalPath = append([]string{current}, criticalPath...)
        
        // 找到前驱任务
        var prevTask string
        maxTime = 0
        for _, vertex := range vertices {
            if weight, exists := ts.graph.GetEdgeWeight(vertex, current); exists {
                if earliest[vertex]+int(weight) == earliest[current] {
                    if earliest[vertex] > maxTime {
                        maxTime = earliest[vertex]
                        prevTask = vertex
                    }
                }
            }
        }
        current = prevTask
    }
    
    return criticalPath, earliest[lastTask]
}
```

## 总结

图论为计算机科学提供了强大的建模和分析工具，通过Go语言的实现，我们可以构建高效的图算法库。这些算法在社交网络分析、网络路由、任务调度等领域有广泛应用。

**关键特性**:
- 支持有向图和无向图
- 完整的图遍历算法
- 最短路径和最小生成树算法
- 拓扑排序和关键路径分析
- 实际应用场景的示例

**应用领域**:
- 社交网络分析
- 网络路由和通信
- 任务调度和项目管理
- 生物信息学和化学
- 地理信息系统
- 编译器优化 