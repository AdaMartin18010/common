# 03-图论 (Graph Theory)

## 概述

图论是数学的一个分支，研究图（Graph）的数学性质。图是由顶点（Vertex）和边（Edge）组成的数学结构，广泛应用于计算机科学、网络分析、算法设计等领域。

## 1. 基本概念

### 1.1 图的定义

**定义 1.1** (图)
一个图 $G = (V, E)$ 由以下两部分组成：
- $V$ 是顶点的有限集合，称为顶点集（Vertex Set）
- $E$ 是边的集合，每条边连接两个顶点，称为边集（Edge Set）

**形式化定义**：
```go
// 图的基本结构
type Graph[T comparable] struct {
    Vertices map[T]bool
    Edges    map[T]map[T]bool
}

// 创建新图
func NewGraph[T comparable]() *Graph[T] {
    return &Graph[T]{
        Vertices: make(map[T]bool),
        Edges:    make(map[T]map[T]bool),
    }
}
```

### 1.2 图的类型

#### 1.2.1 无向图 (Undirected Graph)

**定义 1.2** (无向图)
在无向图中，边没有方向，如果顶点 $u$ 和 $v$ 之间有边，则 $(u,v) = (v,u)$。

```go
// 无向图实现
type UndirectedGraph[T comparable] struct {
    *Graph[T]
}

func NewUndirectedGraph[T comparable]() *UndirectedGraph[T] {
    return &UndirectedGraph[T]{
        Graph: NewGraph[T](),
    }
}

// 添加边（无向）
func (g *UndirectedGraph[T]) AddEdge(u, v T) {
    g.addVertex(u)
    g.addVertex(v)
    
    if g.Edges[u] == nil {
        g.Edges[u] = make(map[T]bool)
    }
    if g.Edges[v] == nil {
        g.Edges[v] = make(map[T]bool)
    }
    
    g.Edges[u][v] = true
    g.Edges[v][u] = true
}
```

#### 1.2.2 有向图 (Directed Graph)

**定义 1.3** (有向图)
在有向图中，边有方向，从顶点 $u$ 指向顶点 $v$ 的边表示为 $(u,v)$。

```go
// 有向图实现
type DirectedGraph[T comparable] struct {
    *Graph[T]
}

func NewDirectedGraph[T comparable]() *DirectedGraph[T] {
    return &DirectedGraph[T]{
        Graph: NewGraph[T](),
    }
}

// 添加边（有向）
func (g *DirectedGraph[T]) AddEdge(from, to T) {
    g.addVertex(from)
    g.addVertex(to)
    
    if g.Edges[from] == nil {
        g.Edges[from] = make(map[T]bool)
    }
    
    g.Edges[from][to] = true
}
```

### 1.3 图的基本操作

```go
// 添加顶点
func (g *Graph[T]) addVertex(v T) {
    g.Vertices[v] = true
}

// 添加顶点（公开方法）
func (g *Graph[T]) AddVertex(v T) {
    g.addVertex(v)
}

// 删除顶点
func (g *Graph[T]) RemoveVertex(v T) {
    delete(g.Vertices, v)
    delete(g.Edges, v)
    
    // 删除所有指向该顶点的边
    for _, edges := range g.Edges {
        delete(edges, v)
    }
}

// 检查顶点是否存在
func (g *Graph[T]) HasVertex(v T) bool {
    return g.Vertices[v]
}

// 检查边是否存在
func (g *Graph[T]) HasEdge(u, v T) bool {
    if edges, exists := g.Edges[u]; exists {
        return edges[v]
    }
    return false
}

// 获取顶点的邻居
func (g *Graph[T]) GetNeighbors(v T) []T {
    var neighbors []T
    if edges, exists := g.Edges[v]; exists {
        for neighbor := range edges {
            neighbors = append(neighbors, neighbor)
        }
    }
    return neighbors
}

// 获取顶点数量
func (g *Graph[T]) VertexCount() int {
    return len(g.Vertices)
}

// 获取边数量
func (g *Graph[T]) EdgeCount() int {
    count := 0
    for _, edges := range g.Edges {
        count += len(edges)
    }
    return count
}
```

## 2. 图的表示方法

### 2.1 邻接矩阵 (Adjacency Matrix)

**定义 2.1** (邻接矩阵)
对于图 $G = (V, E)$，邻接矩阵 $A$ 是一个 $|V| \times |V|$ 的矩阵，其中：
$$A_{ij} = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$$

```go
// 邻接矩阵表示
type AdjacencyMatrix[T comparable] struct {
    vertices []T
    matrix   [][]bool
    vertexMap map[T]int
}

func NewAdjacencyMatrix[T comparable]() *AdjacencyMatrix[T] {
    return &AdjacencyMatrix[T]{
        vertices:  make([]T, 0),
        matrix:    make([][]bool, 0),
        vertexMap: make(map[T]int),
    }
}

// 添加顶点
func (am *AdjacencyMatrix[T]) AddVertex(v T) {
    if _, exists := am.vertexMap[v]; !exists {
        am.vertexMap[v] = len(am.vertices)
        am.vertices = append(am.vertices, v)
        
        // 扩展矩阵
        for i := range am.matrix {
            am.matrix[i] = append(am.matrix[i], false)
        }
        newRow := make([]bool, len(am.vertices))
        am.matrix = append(am.matrix, newRow)
    }
}

// 添加边
func (am *AdjacencyMatrix[T]) AddEdge(u, v T) {
    am.AddVertex(u)
    am.AddVertex(v)
    
    i := am.vertexMap[u]
    j := am.vertexMap[v]
    
    am.matrix[i][j] = true
    am.matrix[j][i] = true // 无向图
}

// 检查边是否存在
func (am *AdjacencyMatrix[T]) HasEdge(u, v T) bool {
    i, existsI := am.vertexMap[u]
    j, existsJ := am.vertexMap[v]
    
    if !existsI || !existsJ {
        return false
    }
    
    return am.matrix[i][j]
}
```

### 2.2 邻接表 (Adjacency List)

**定义 2.2** (邻接表)
邻接表是图的另一种表示方法，对于每个顶点维护一个包含其所有邻居的列表。

```go
// 邻接表表示
type AdjacencyList[T comparable] struct {
    vertices map[T][]T
}

func NewAdjacencyList[T comparable]() *AdjacencyList[T] {
    return &AdjacencyList[T]{
        vertices: make(map[T][]T),
    }
}

// 添加顶点
func (al *AdjacencyList[T]) AddVertex(v T) {
    if _, exists := al.vertices[v]; !exists {
        al.vertices[v] = make([]T, 0)
    }
}

// 添加边
func (al *AdjacencyList[T]) AddEdge(u, v T) {
    al.AddVertex(u)
    al.AddVertex(v)
    
    // 检查边是否已存在
    for _, neighbor := range al.vertices[u] {
        if neighbor == v {
            return
        }
    }
    
    al.vertices[u] = append(al.vertices[u], v)
    al.vertices[v] = append(al.vertices[v], u) // 无向图
}

// 获取邻居
func (al *AdjacencyList[T]) GetNeighbors(v T) []T {
    if neighbors, exists := al.vertices[v]; exists {
        return neighbors
    }
    return []T{}
}
```

## 3. 图的遍历算法

### 3.1 深度优先搜索 (DFS)

**算法 3.1** (深度优先搜索)
深度优先搜索是一种图遍历算法，沿着图的边尽可能深入地搜索。

**时间复杂度**: $O(|V| + |E|)$
**空间复杂度**: $O(|V|)$

```go
// DFS实现
func (g *Graph[T]) DFS(start T, visit func(T)) {
    visited := make(map[T]bool)
    g.dfsHelper(start, visited, visit)
}

func (g *Graph[T]) dfsHelper(v T, visited map[T]bool, visit func(T)) {
    if visited[v] {
        return
    }
    
    visited[v] = true
    visit(v)
    
    for neighbor := range g.Edges[v] {
        g.dfsHelper(neighbor, visited, visit)
    }
}

// 非递归DFS
func (g *Graph[T]) DFSIterative(start T, visit func(T)) {
    visited := make(map[T]bool)
    stack := []T{start}
    
    for len(stack) > 0 {
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if visited[v] {
            continue
        }
        
        visited[v] = true
        visit(v)
        
        // 将邻居按相反顺序压入栈中
        neighbors := g.GetNeighbors(v)
        for i := len(neighbors) - 1; i >= 0; i-- {
            if !visited[neighbors[i]] {
                stack = append(stack, neighbors[i])
            }
        }
    }
}
```

### 3.2 广度优先搜索 (BFS)

**算法 3.2** (广度优先搜索)
广度优先搜索是一种图遍历算法，先访问所有相邻顶点，再访问下一层顶点。

**时间复杂度**: $O(|V| + |E|)$
**空间复杂度**: $O(|V|)$

```go
// BFS实现
func (g *Graph[T]) BFS(start T, visit func(T)) {
    visited := make(map[T]bool)
    queue := []T{start}
    visited[start] = true
    
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        
        visit(v)
        
        for neighbor := range g.Edges[v] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}

// BFS with level information
func (g *Graph[T]) BFSWithLevels(start T) map[T]int {
    levels := make(map[T]int)
    visited := make(map[T]bool)
    queue := []T{start}
    
    levels[start] = 0
    visited[start] = true
    
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        
        for neighbor := range g.Edges[v] {
            if !visited[neighbor] {
                visited[neighbor] = true
                levels[neighbor] = levels[v] + 1
                queue = append(queue, neighbor)
            }
        }
    }
    
    return levels
}
```

## 4. 最短路径算法

### 4.1 Dijkstra算法

**算法 4.1** (Dijkstra算法)
Dijkstra算法用于在带权图中找到从源点到所有其他顶点的最短路径。

**时间复杂度**: $O((|V| + |E|) \log |V|)$ (使用优先队列)
**空间复杂度**: $O(|V|)$

```go
// 带权图
type WeightedGraph[T comparable] struct {
    vertices map[T]bool
    edges    map[T]map[T]float64
}

func NewWeightedGraph[T comparable]() *WeightedGraph[T] {
    return &WeightedGraph[T]{
        vertices: make(map[T]bool),
        edges:    make(map[T]map[T]float64),
    }
}

// 添加带权边
func (wg *WeightedGraph[T]) AddEdge(u, v T, weight float64) {
    wg.addVertex(u)
    wg.addVertex(v)
    
    if wg.edges[u] == nil {
        wg.edges[u] = make(map[T]float64)
    }
    if wg.edges[v] == nil {
        wg.edges[v] = make(map[T]float64)
    }
    
    wg.edges[u][v] = weight
    wg.edges[v][u] = weight // 无向图
}

// Dijkstra算法实现
func (wg *WeightedGraph[T]) Dijkstra(start T) map[T]float64 {
    distances := make(map[T]float64)
    visited := make(map[T]bool)
    
    // 初始化距离
    for vertex := range wg.vertices {
        distances[vertex] = math.Inf(1)
    }
    distances[start] = 0
    
    // 使用优先队列
    pq := make([]T, 0)
    heap.Push(&pq, start)
    
    for len(pq) > 0 {
        u := heap.Pop(&pq).(T)
        
        if visited[u] {
            continue
        }
        
        visited[u] = true
        
        for v, weight := range wg.edges[u] {
            if !visited[v] {
                newDist := distances[u] + weight
                if newDist < distances[v] {
                    distances[v] = newDist
                    heap.Push(&pq, v)
                }
            }
        }
    }
    
    return distances
}

// 优先队列实现
type PriorityQueue[T comparable] []T

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
    // 这里需要根据实际的距离值进行比较
    // 简化实现，实际应用中需要更复杂的比较逻辑
    return i < j
}

func (pq PriorityQueue[T]) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
    *pq = append(*pq, x.(T))
}

func (pq *PriorityQueue[T]) Pop() interface{} {
    old := *pq
    n := len(old)
    x := old[n-1]
    *pq = old[0 : n-1]
    return x
}
```

### 4.2 Floyd-Warshall算法

**算法 4.2** (Floyd-Warshall算法)
Floyd-Warshall算法用于找到图中所有顶点对之间的最短路径。

**时间复杂度**: $O(|V|^3)$
**空间复杂度**: $O(|V|^2)$

```go
// Floyd-Warshall算法实现
func (wg *WeightedGraph[T]) FloydWarshall() map[T]map[T]float64 {
    vertices := make([]T, 0, len(wg.vertices))
    for v := range wg.vertices {
        vertices = append(vertices, v)
    }
    
    // 创建距离矩阵
    dist := make(map[T]map[T]float64)
    for _, u := range vertices {
        dist[u] = make(map[T]float64)
        for _, v := range vertices {
            if u == v {
                dist[u][v] = 0
            } else if weight, exists := wg.edges[u][v]; exists {
                dist[u][v] = weight
            } else {
                dist[u][v] = math.Inf(1)
            }
        }
    }
    
    // Floyd-Warshall核心算法
    for _, k := range vertices {
        for _, i := range vertices {
            for _, j := range vertices {
                if dist[i][k] != math.Inf(1) && dist[k][j] != math.Inf(1) {
                    if dist[i][k]+dist[k][j] < dist[i][j] {
                        dist[i][j] = dist[i][k] + dist[k][j]
                    }
                }
            }
        }
    }
    
    return dist
}
```

## 5. 最小生成树算法

### 5.1 Kruskal算法

**算法 5.1** (Kruskal算法)
Kruskal算法用于在带权无向图中找到最小生成树。

**时间复杂度**: $O(|E| \log |E|)$
**空间复杂度**: $O(|V|)$

```go
// 边结构
type Edge[T comparable] struct {
    From   T
    To     T
    Weight float64
}

// Kruskal算法实现
func (wg *WeightedGraph[T]) Kruskal() []Edge[T] {
    var edges []Edge[T]
    
    // 收集所有边
    for u := range wg.edges {
        for v, weight := range wg.edges[u] {
            if u < v { // 避免重复边
                edges = append(edges, Edge[T]{From: u, To: v, Weight: weight})
            }
        }
    }
    
    // 按权重排序
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })
    
    // 并查集
    uf := NewUnionFind[T]()
    for vertex := range wg.vertices {
        uf.MakeSet(vertex)
    }
    
    var mst []Edge[T]
    for _, edge := range edges {
        if uf.Find(edge.From) != uf.Find(edge.To) {
            uf.Union(edge.From, edge.To)
            mst = append(mst, edge)
        }
    }
    
    return mst
}

// 并查集实现
type UnionFind[T comparable] struct {
    parent map[T]T
    rank   map[T]int
}

func NewUnionFind[T comparable]() *UnionFind[T] {
    return &UnionFind[T]{
        parent: make(map[T]T),
        rank:   make(map[T]int),
    }
}

func (uf *UnionFind[T]) MakeSet(x T) {
    uf.parent[x] = x
    uf.rank[x] = 0
}

func (uf *UnionFind[T]) Find(x T) T {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // 路径压缩
    }
    return uf.parent[x]
}

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

## 6. 图论定理与证明

### 6.1 握手定理 (Handshaking Lemma)

**定理 6.1** (握手定理)
在任何图中，所有顶点的度数之和等于边数的两倍。

**证明**：
设图 $G = (V, E)$，对于每条边 $(u,v) \in E$，它贡献给顶点 $u$ 和 $v$ 各一个度数。因此：
$$\sum_{v \in V} \deg(v) = 2|E|$$

```go
// 验证握手定理
func (g *Graph[T]) VerifyHandshakingLemma() bool {
    totalDegree := 0
    for v := range g.Vertices {
        totalDegree += len(g.Edges[v])
    }
    
    return totalDegree == 2*g.EdgeCount()
}
```

### 6.2 欧拉路径定理

**定理 6.2** (欧拉路径定理)
一个无向图存在欧拉路径当且仅当：
1. 图是连通的
2. 恰好有0个或2个顶点的度数为奇数

**证明**：
- 必要性：欧拉路径必须经过每条边恰好一次，因此除了起点和终点外，每个顶点的度数必须是偶数。
- 充分性：可以通过构造性证明，从奇数度数的顶点开始，按照特定规则遍历图。

```go
// 检查是否存在欧拉路径
func (g *UndirectedGraph[T]) HasEulerPath() bool {
    // 检查连通性
    if !g.IsConnected() {
        return false
    }
    
    // 计算奇数度数的顶点数量
    oddDegreeCount := 0
    for v := range g.Vertices {
        if len(g.Edges[v])%2 == 1 {
            oddDegreeCount++
        }
    }
    
    return oddDegreeCount == 0 || oddDegreeCount == 2
}

// 检查图的连通性
func (g *Graph[T]) IsConnected() bool {
    if len(g.Vertices) == 0 {
        return true
    }
    
    // 从任意顶点开始DFS
    var start T
    for v := range g.Vertices {
        start = v
        break
    }
    
    visited := make(map[T]bool)
    g.dfsHelper(start, visited, func(v T) {})
    
    // 检查是否所有顶点都被访问
    return len(visited) == len(g.Vertices)
}
```

## 7. 应用实例

### 7.1 社交网络分析

```go
// 社交网络图
type SocialNetwork struct {
    graph *UndirectedGraph[string]
}

func NewSocialNetwork() *SocialNetwork {
    return &SocialNetwork{
        graph: NewUndirectedGraph[string](),
    }
}

// 添加朋友关系
func (sn *SocialNetwork) AddFriendship(user1, user2 string) {
    sn.graph.AddEdge(user1, user2)
}

// 计算用户的朋友数量
func (sn *SocialNetwork) GetFriendCount(user string) int {
    return len(sn.graph.GetNeighbors(user))
}

// 找到共同朋友
func (sn *SocialNetwork) GetCommonFriends(user1, user2 string) []string {
    friends1 := make(map[string]bool)
    for _, friend := range sn.graph.GetNeighbors(user1) {
        friends1[friend] = true
    }
    
    var common []string
    for _, friend := range sn.graph.GetNeighbors(user2) {
        if friends1[friend] {
            common = append(common, friend)
        }
    }
    
    return common
}

// 计算社交网络中的连通分量
func (sn *SocialNetwork) GetConnectedComponents() [][]string {
    visited := make(map[string]bool)
    var components [][]string
    
    for user := range sn.graph.Vertices {
        if !visited[user] {
            var component []string
            sn.graph.DFS(user, func(v string) {
                if !visited[v] {
                    visited[v] = true
                    component = append(component, v)
                }
            })
            components = append(components, component)
        }
    }
    
    return components
}
```

### 7.2 网络路由

```go
// 网络拓扑
type NetworkTopology struct {
    graph *WeightedGraph[string]
}

func NewNetworkTopology() *NetworkTopology {
    return &NetworkTopology{
        graph: NewWeightedGraph[string](),
    }
}

// 添加网络连接
func (nt *NetworkTopology) AddConnection(node1, node2 string, latency float64) {
    nt.graph.AddEdge(node1, node2, latency)
}

// 找到最短路径
func (nt *NetworkTopology) FindShortestPath(from, to string) ([]string, float64) {
    distances := nt.graph.Dijkstra(from)
    
    if distances[to] == math.Inf(1) {
        return nil, math.Inf(1)
    }
    
    // 重建路径（简化实现）
    path := []string{to}
    current := to
    
    for current != from {
        minDist := math.Inf(1)
        var prev string
        
        for neighbor := range nt.graph.edges[current] {
            if distances[neighbor] < minDist {
                minDist = distances[neighbor]
                prev = neighbor
            }
        }
        
        if prev == "" {
            break
        }
        
        path = append([]string{prev}, path...)
        current = prev
    }
    
    return path, distances[to]
}
```

## 总结

图论为计算机科学提供了强大的数学工具，用于建模和分析各种复杂系统。通过Go语言的实现，我们可以：

1. **高效表示**: 使用邻接矩阵和邻接表高效表示图结构
2. **算法实现**: 实现经典的图算法如DFS、BFS、Dijkstra等
3. **实际应用**: 将图论应用于社交网络、网络路由等实际问题
4. **性能优化**: 利用Go的并发特性优化大规模图处理

图论的理论基础与Go语言的实践相结合，为构建复杂的网络应用和算法提供了坚实的基础。

