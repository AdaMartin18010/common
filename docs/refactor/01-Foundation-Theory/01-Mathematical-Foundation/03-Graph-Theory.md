# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 图的定义](#11-图的定义)
    - [1.2 图的类型](#12-图的类型)
    - [1.3 基本术语](#13-基本术语)
  - [2. 图的表示](#2-图的表示)
    - [2.1 邻接矩阵](#21-邻接矩阵)
    - [2.2 邻接表](#22-邻接表)
    - [2.3 边列表](#23-边列表)
  - [3. 图算法](#3-图算法)
    - [3.1 深度优先搜索](#31-深度优先搜索)
    - [3.2 广度优先搜索](#32-广度优先搜索)
    - [3.3 最短路径算法](#33-最短路径算法)
    - [3.4 最小生成树](#34-最小生成树)
  - [4. 图论定理](#4-图论定理)
    - [4.1 欧拉定理](#41-欧拉定理)
    - [4.2 哈密顿定理](#42-哈密顿定理)
    - [4.3 四色定理](#43-四色定理)
  - [5. 应用领域](#5-应用领域)
    - [5.1 网络分析](#51-网络分析)
    - [5.2 社交网络](#52-社交网络)
    - [5.3 路由算法](#53-路由算法)
  - [6. 总结](#6-总结)
  - [参考文献](#参考文献)

## 1. 基础概念

### 1.1 图的定义

**定义 1.1** (图): 图 $G = (V, E)$ 是一个有序对，其中：
- $V$ 是顶点集 (vertex set)，$V \neq \emptyset$
- $E$ 是边集 (edge set)，$E \subseteq V \times V$

**定义 1.2** (有向图): 有向图 $G = (V, E)$ 中，边是有序对 $(u, v) \in E$，表示从顶点 $u$ 到顶点 $v$ 的有向边。

**定义 1.3** (无向图): 无向图 $G = (V, E)$ 中，边是无序对 $\{u, v\} \in E$，表示顶点 $u$ 和 $v$ 之间的无向边。

### 1.2 图的类型

```go
// 图的基本类型定义
type GraphType int

const (
    UndirectedGraph GraphType = iota
    DirectedGraph
    WeightedGraph
    BipartiteGraph
    Tree
    DAG // Directed Acyclic Graph
)

// 顶点接口
type Vertex interface {
    ID() int
    String() string
}

// 边接口
type Edge interface {
    From() Vertex
    To() Vertex
    Weight() float64
    String() string
}

// 基础图结构
type Graph struct {
    vertices map[int]Vertex
    edges    map[int]map[int]Edge
    graphType GraphType
    directed  bool
    weighted  bool
}

// 创建新图
func NewGraph(graphType GraphType) *Graph {
    return &Graph{
        vertices: make(map[int]Vertex),
        edges:    make(map[int]map[int]Edge),
        graphType: graphType,
        directed:  graphType == DirectedGraph || graphType == DAG,
        weighted:  graphType == WeightedGraph,
    }
}
```

### 1.3 基本术语

**定义 1.4** (度数): 顶点 $v$ 的度数 $deg(v)$ 是与 $v$ 相邻的边数。

**定理 1.1** (握手定理): 对于任意图 $G = (V, E)$，
$$\sum_{v \in V} deg(v) = 2|E|$$

**证明**: 每条边贡献给两个顶点的度数，因此所有顶点的度数之和等于边数的两倍。

```go
// 度数计算
func (g *Graph) Degree(vertexID int) int {
    if g.directed {
        return g.InDegree(vertexID) + g.OutDegree(vertexID)
    }
    
    if neighbors, exists := g.edges[vertexID]; exists {
        return len(neighbors)
    }
    return 0
}

// 入度 (有向图)
func (g *Graph) InDegree(vertexID int) int {
    if !g.directed {
        return g.Degree(vertexID)
    }
    
    count := 0
    for _, neighbors := range g.edges {
        if _, exists := neighbors[vertexID]; exists {
            count++
        }
    }
    return count
}

// 出度 (有向图)
func (g *Graph) OutDegree(vertexID int) int {
    if !g.directed {
        return g.Degree(vertexID)
    }
    
    if neighbors, exists := g.edges[vertexID]; exists {
        return len(neighbors)
    }
    return 0
}
```

## 2. 图的表示

### 2.1 邻接矩阵

**定义 2.1** (邻接矩阵): 对于图 $G = (V, E)$，邻接矩阵 $A$ 是一个 $|V| \times |V|$ 的矩阵，其中：
$$A[i][j] = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$$

```go
// 邻接矩阵表示
type AdjacencyMatrix struct {
    matrix [][]int
    size   int
}

func NewAdjacencyMatrix(size int) *AdjacencyMatrix {
    matrix := make([][]int, size)
    for i := range matrix {
        matrix[i] = make([]int, size)
    }
    return &AdjacencyMatrix{matrix: matrix, size: size}
}

func (am *AdjacencyMatrix) AddEdge(from, to int) {
    if from >= 0 && from < am.size && to >= 0 && to < am.size {
        am.matrix[from][to] = 1
    }
}

func (am *AdjacencyMatrix) HasEdge(from, to int) bool {
    if from >= 0 && from < am.size && to >= 0 && to < am.size {
        return am.matrix[from][to] == 1
    }
    return false
}

// 时间复杂度分析
// 空间复杂度: O(V²)
// 添加边: O(1)
// 检查边: O(1)
// 遍历邻居: O(V)
```

### 2.2 邻接表

**定义 2.2** (邻接表): 邻接表是图的另一种表示方法，为每个顶点维护一个邻居列表。

```go
// 邻接表表示
type AdjacencyList struct {
    vertices map[int][]int
    size     int
}

func NewAdjacencyList() *AdjacencyList {
    return &AdjacencyList{
        vertices: make(map[int][]int),
        size:     0,
    }
}

func (al *AdjacencyList) AddVertex(vertexID int) {
    if _, exists := al.vertices[vertexID]; !exists {
        al.vertices[vertexID] = []int{}
        al.size++
    }
}

func (al *AdjacencyList) AddEdge(from, to int) {
    al.AddVertex(from)
    al.AddVertex(to)
    
    // 检查边是否已存在
    for _, neighbor := range al.vertices[from] {
        if neighbor == to {
            return
        }
    }
    
    al.vertices[from] = append(al.vertices[from], to)
}

func (al *AdjacencyList) GetNeighbors(vertexID int) []int {
    if neighbors, exists := al.vertices[vertexID]; exists {
        return neighbors
    }
    return []int{}
}

// 时间复杂度分析
// 空间复杂度: O(V + E)
// 添加边: O(1) 平均情况
// 检查边: O(deg(v))
// 遍历邻居: O(deg(v))
```

### 2.3 边列表

**定义 2.3** (边列表): 边列表是图的简单表示方法，直接存储所有边。

```go
// 边列表表示
type EdgeList struct {
    edges []Edge
    size  int
}

type SimpleEdge struct {
    from   int
    to     int
    weight float64
}

func (se SimpleEdge) From() Vertex {
    return &SimpleVertex{id: se.from}
}

func (se SimpleEdge) To() Vertex {
    return &SimpleVertex{id: se.to}
}

func (se SimpleEdge) Weight() float64 {
    return se.weight
}

func (se SimpleEdge) String() string {
    return fmt.Sprintf("(%d -> %d, weight: %.2f)", se.from, se.to, se.weight)
}

type SimpleVertex struct {
    id int
}

func (sv SimpleVertex) ID() int {
    return sv.id
}

func (sv SimpleVertex) String() string {
    return fmt.Sprintf("V%d", sv.id)
}

func NewEdgeList() *EdgeList {
    return &EdgeList{
        edges: []Edge{},
        size:  0,
    }
}

func (el *EdgeList) AddEdge(from, to int, weight float64) {
    edge := &SimpleEdge{from: from, to: to, weight: weight}
    el.edges = append(el.edges, edge)
    el.size++
}

// 时间复杂度分析
// 空间复杂度: O(E)
// 添加边: O(1)
// 检查边: O(E)
// 遍历邻居: O(E)
```

## 3. 图算法

### 3.1 深度优先搜索

**算法 3.1** (深度优先搜索): DFS是一种图遍历算法，沿着图的深度优先遍历。

```go
// 深度优先搜索
func (g *Graph) DFS(startID int) []int {
    visited := make(map[int]bool)
    result := []int{}
    
    var dfs func(vertexID int)
    dfs = func(vertexID int) {
        visited[vertexID] = true
        result = append(result, vertexID)
        
        if neighbors, exists := g.edges[vertexID]; exists {
            for neighborID := range neighbors {
                if !visited[neighborID] {
                    dfs(neighborID)
                }
            }
        }
    }
    
    dfs(startID)
    return result
}

// 非递归版本
func (g *Graph) DFSIterative(startID int) []int {
    visited := make(map[int]bool)
    result := []int{}
    stack := []int{startID}
    
    for len(stack) > 0 {
        vertexID := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if !visited[vertexID] {
            visited[vertexID] = true
            result = append(result, vertexID)
            
            if neighbors, exists := g.edges[vertexID]; exists {
                for neighborID := range neighbors {
                    if !visited[neighborID] {
                        stack = append(stack, neighborID)
                    }
                }
            }
        }
    }
    
    return result
}

// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
```

### 3.2 广度优先搜索

**算法 3.2** (广度优先搜索): BFS是一种图遍历算法，按照层次遍历图。

```go
// 广度优先搜索
func (g *Graph) BFS(startID int) []int {
    visited := make(map[int]bool)
    result := []int{}
    queue := []int{startID}
    visited[startID] = true
    
    for len(queue) > 0 {
        vertexID := queue[0]
        queue = queue[1:]
        result = append(result, vertexID)
        
        if neighbors, exists := g.edges[vertexID]; exists {
            for neighborID := range neighbors {
                if !visited[neighborID] {
                    visited[neighborID] = true
                    queue = append(queue, neighborID)
                }
            }
        }
    }
    
    return result
}

// 带距离的BFS
func (g *Graph) BFSWithDistance(startID int) map[int]int {
    visited := make(map[int]bool)
    distance := make(map[int]int)
    queue := []int{startID}
    visited[startID] = true
    distance[startID] = 0
    
    for len(queue) > 0 {
        vertexID := queue[0]
        queue = queue[1:]
        
        if neighbors, exists := g.edges[vertexID]; exists {
            for neighborID := range neighbors {
                if !visited[neighborID] {
                    visited[neighborID] = true
                    distance[neighborID] = distance[vertexID] + 1
                    queue = append(queue, neighborID)
                }
            }
        }
    }
    
    return distance
}

// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
```

### 3.3 最短路径算法

**算法 3.3** (Dijkstra算法): 用于计算带权图中单源最短路径。

```go
// Dijkstra算法
func (g *Graph) Dijkstra(startID int) map[int]float64 {
    distances := make(map[int]float64)
    visited := make(map[int]bool)
    
    // 初始化距离
    for vertexID := range g.vertices {
        distances[vertexID] = math.Inf(1)
    }
    distances[startID] = 0
    
    // 优先队列 (使用堆)
    pq := &PriorityQueue{}
    heap.Push(pq, &Item{vertexID: startID, distance: 0})
    
    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        vertexID := item.vertexID
        
        if visited[vertexID] {
            continue
        }
        visited[vertexID] = true
        
        if neighbors, exists := g.edges[vertexID]; exists {
            for neighborID, edge := range neighbors {
                newDistance := distances[vertexID] + edge.Weight()
                if newDistance < distances[neighborID] {
                    distances[neighborID] = newDistance
                    heap.Push(pq, &Item{vertexID: neighborID, distance: newDistance})
                }
            }
        }
    }
    
    return distances
}

// 优先队列实现
type Item struct {
    vertexID int
    distance float64
    index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

// 时间复杂度: O((V + E) log V) 使用二叉堆
// 空间复杂度: O(V)
```

### 3.4 最小生成树

**定义 3.1** (最小生成树): 对于连通无向图 $G = (V, E)$，最小生成树是包含所有顶点且边权重和最小的树。

**算法 3.4** (Kruskal算法): 用于计算最小生成树。

```go
// Kruskal算法
func (g *Graph) KruskalMST() []Edge {
    var edges []Edge
    
    // 收集所有边
    for _, neighbors := range g.edges {
        for _, edge := range neighbors {
            edges = append(edges, edge)
        }
    }
    
    // 按权重排序
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight() < edges[j].Weight()
    })
    
    // 并查集
    uf := NewUnionFind(len(g.vertices))
    mst := []Edge{}
    
    for _, edge := range edges {
        fromID := edge.From().ID()
        toID := edge.To().ID()
        
        if uf.Find(fromID) != uf.Find(toID) {
            mst = append(mst, edge)
            uf.Union(fromID, toID)
        }
    }
    
    return mst
}

// 并查集实现
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(size int) *UnionFind {
    parent := make([]int, size)
    rank := make([]int, size)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{parent: parent, rank: rank}
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
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

// 时间复杂度: O(E log E)
// 空间复杂度: O(V)
```

## 4. 图论定理

### 4.1 欧拉定理

**定理 4.1** (欧拉定理): 连通图 $G$ 存在欧拉回路的充要条件是所有顶点的度数都是偶数。

**证明**: 
- 必要性：欧拉回路经过每条边恰好一次，每个顶点被进入和离开的次数相等。
- 充分性：通过构造性证明，可以找到欧拉回路。

```go
// 欧拉回路检测
func (g *Graph) HasEulerCircuit() bool {
    // 检查连通性
    if !g.IsConnected() {
        return false
    }
    
    // 检查所有顶点度数是否为偶数
    for vertexID := range g.vertices {
        if g.Degree(vertexID)%2 != 0 {
            return false
        }
    }
    
    return true
}

// 连通性检测
func (g *Graph) IsConnected() bool {
    if len(g.vertices) == 0 {
        return true
    }
    
    // 从任意顶点开始DFS
    var startID int
    for vertexID := range g.vertices {
        startID = vertexID
        break
    }
    
    visited := g.DFS(startID)
    return len(visited) == len(g.vertices)
}
```

### 4.2 哈密顿定理

**定理 4.2** (哈密顿定理): 对于完全图 $K_n$，存在哈密顿回路。

**证明**: 完全图中任意两个顶点都相邻，可以通过多种方式构造哈密顿回路。

```go
// 哈密顿回路检测 (NP完全问题)
func (g *Graph) HasHamiltonianCycle() bool {
    if len(g.vertices) < 3 {
        return false
    }
    
    // 使用回溯法
    path := []int{}
    visited := make(map[int]bool)
    
    return g.hamiltonianBacktrack(path, visited)
}

func (g *Graph) hamiltonianBacktrack(path []int, visited map[int]bool) bool {
    if len(path) == len(g.vertices) {
        // 检查最后一个顶点是否与起始顶点相邻
        if len(path) > 0 {
            first := path[0]
            last := path[len(path)-1]
            if g.HasEdge(last, first) {
                return true
            }
        }
        return false
    }
    
    for vertexID := range g.vertices {
        if !visited[vertexID] {
            if len(path) == 0 || g.HasEdge(path[len(path)-1], vertexID) {
                visited[vertexID] = true
                path = append(path, vertexID)
                
                if g.hamiltonianBacktrack(path, visited) {
                    return true
                }
                
                visited[vertexID] = false
                path = path[:len(path)-1]
            }
        }
    }
    
    return false
}
```

### 4.3 四色定理

**定理 4.3** (四色定理): 任何平面图都可以用四种颜色着色，使得相邻顶点颜色不同。

```go
// 图着色算法
func (g *Graph) ColorGraph() map[int]int {
    colors := make(map[int]int)
    
    for vertexID := range g.vertices {
        // 找到可用的最小颜色
        usedColors := make(map[int]bool)
        
        if neighbors, exists := g.edges[vertexID]; exists {
            for neighborID := range neighbors {
                if color, hasColor := colors[neighborID]; hasColor {
                    usedColors[color] = true
                }
            }
        }
        
        // 分配最小可用颜色
        for color := 0; ; color++ {
            if !usedColors[color] {
                colors[vertexID] = color
                break
            }
        }
    }
    
    return colors
}

// 检查着色是否有效
func (g *Graph) IsValidColoring(colors map[int]int) bool {
    for vertexID, neighbors := range g.edges {
        vertexColor := colors[vertexID]
        for neighborID := range neighbors {
            if colors[neighborID] == vertexColor {
                return false
            }
        }
    }
    return true
}
```

## 5. 应用领域

### 5.1 网络分析

图论在网络分析中有广泛应用：

```go
// 网络中心性分析
type CentralityAnalysis struct {
    graph *Graph
}

// 度中心性
func (ca *CentralityAnalysis) DegreeCentrality(vertexID int) float64 {
    degree := ca.graph.Degree(vertexID)
    totalVertices := len(ca.graph.vertices)
    return float64(degree) / float64(totalVertices-1)
}

// 接近中心性
func (ca *CentralityAnalysis) ClosenessCentrality(vertexID int) float64 {
    distances := ca.graph.BFSWithDistance(vertexID)
    totalDistance := 0
    reachableVertices := 0
    
    for _, distance := range distances {
        if distance > 0 {
            totalDistance += distance
            reachableVertices++
        }
    }
    
    if reachableVertices == 0 {
        return 0
    }
    
    return float64(reachableVertices) / float64(totalDistance)
}

// 介数中心性
func (ca *CentralityAnalysis) BetweennessCentrality(vertexID int) float64 {
    // 简化实现
    totalPaths := 0
    pathsThroughVertex := 0
    
    for startID := range ca.graph.vertices {
        for endID := range ca.graph.vertices {
            if startID != endID && startID != vertexID && endID != vertexID {
                totalPaths++
                // 检查最短路径是否经过vertexID
                if ca.isOnShortestPath(startID, endID, vertexID) {
                    pathsThroughVertex++
                }
            }
        }
    }
    
    if totalPaths == 0 {
        return 0
    }
    
    return float64(pathsThroughVertex) / float64(totalPaths)
}

func (ca *CentralityAnalysis) isOnShortestPath(start, end, vertex int) bool {
    // 简化实现：检查vertex是否在start到end的最短路径上
    distances := ca.graph.BFSWithDistance(start)
    if distance, exists := distances[end]; exists && distance > 0 {
        vertexDistance := distances[vertex]
        return vertexDistance > 0 && vertexDistance < distance
    }
    return false
}
```

### 5.2 社交网络

```go
// 社交网络分析
type SocialNetwork struct {
    graph *Graph
    users map[int]*User
}

type User struct {
    ID       int
    Name     string
    Friends  []int
    Posts    []Post
}

type Post struct {
    ID      int
    Content string
    Likes   int
    Shares  int
}

// 影响力分析
func (sn *SocialNetwork) InfluenceScore(userID int) float64 {
    // 基于度中心性和接近中心性计算影响力
    centrality := &CentralityAnalysis{graph: sn.graph}
    
    degreeCentrality := centrality.DegreeCentrality(userID)
    closenessCentrality := centrality.ClosenessCentrality(userID)
    
    return (degreeCentrality + closenessCentrality) / 2.0
}

// 社区检测 (简化版)
func (sn *SocialNetwork) DetectCommunities() [][]int {
    visited := make(map[int]bool)
    communities := [][]int{}
    
    for userID := range sn.users {
        if !visited[userID] {
            community := sn.DFS(userID)
            communities = append(communities, community)
            
            for _, memberID := range community {
                visited[memberID] = true
            }
        }
    }
    
    return communities
}
```

### 5.3 路由算法

```go
// 网络路由
type NetworkRouter struct {
    graph *Graph
    routingTable map[int]map[int]int // [source][destination] -> nextHop
}

func NewNetworkRouter(graph *Graph) *NetworkRouter {
    return &NetworkRouter{
        graph: graph,
        routingTable: make(map[int]map[int]int),
    }
}

// 构建路由表
func (nr *NetworkRouter) BuildRoutingTable() {
    for sourceID := range nr.graph.vertices {
        nr.routingTable[sourceID] = make(map[int]int)
        distances := nr.graph.BFSWithDistance(sourceID)
        
        for destID, distance := range distances {
            if distance > 0 {
                // 找到下一跳
                nextHop := nr.findNextHop(sourceID, destID)
                nr.routingTable[sourceID][destID] = nextHop
            }
        }
    }
}

func (nr *NetworkRouter) findNextHop(source, dest int) int {
    // 简化实现：返回第一个邻居
    if neighbors, exists := nr.graph.edges[source]; exists {
        for neighborID := range neighbors {
            return neighborID
        }
    }
    return source
}

// 路由查找
func (nr *NetworkRouter) Route(source, dest int) []int {
    path := []int{source}
    current := source
    
    for current != dest {
        if nextHop, exists := nr.routingTable[current][dest]; exists {
            path = append(path, nextHop)
            current = nextHop
        } else {
            return nil // 无路径
        }
    }
    
    return path
}
```

## 6. 总结

图论是计算机科学和软件工程的重要理论基础，提供了：

1. **理论基础**: 为算法设计和分析提供数学基础
2. **数据结构**: 邻接矩阵、邻接表等高效的图表示方法
3. **算法设计**: DFS、BFS、最短路径、最小生成树等经典算法
4. **应用领域**: 网络分析、社交网络、路由算法等实际应用

通过Go语言的实现，我们可以：
- 高效地表示和操作图结构
- 实现各种图算法
- 解决实际工程问题
- 验证图论定理和性质

图论为软件工程提供了强大的工具，特别是在网络应用、社交平台、路由系统等领域有广泛应用。

## 参考文献

1. Cormen, T. H., Leiserson, C. E., Rivest, R. L., & Stein, C. (2009). Introduction to Algorithms (3rd ed.). MIT Press.
2. Bondy, J. A., & Murty, U. S. R. (2008). Graph Theory. Springer.
3. West, D. B. (2001). Introduction to Graph Theory (2nd ed.). Prentice Hall.
