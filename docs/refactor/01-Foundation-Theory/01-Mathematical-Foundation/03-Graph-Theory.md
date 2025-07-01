# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [1. 基础定义](#1-基础定义)
    - [1.1 图的基本概念](#11-图的基本概念)
    - [1.2 图的类型](#12-图的类型)
    - [1.3 图的表示](#13-图的表示)
  - [2. 图的性质](#2-图的性质)
    - [2.1 连通性](#21-连通性)
    - [2.2 路径和回路](#22-路径和回路)
    - [2.3 树和森林](#23-树和森林)
  - [3. 图的算法](#3-图的算法)
    - [3.1 遍历算法](#31-遍历算法)
    - [3.2 最短路径算法](#32-最短路径算法)
    - [3.3 最小生成树](#33-最小生成树)
  - [4. 特殊图类](#4-特殊图类)
    - [4.1 平面图](#41-平面图)
    - [4.2 欧拉图](#42-欧拉图)
    - [4.3 哈密顿图](#43-哈密顿图)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 数据结构](#51-数据结构)
    - [5.2 基础算法实现](#52-基础算法实现)
    - [5.3 高级算法实现](#53-高级算法实现)
  - [6. 应用实例](#6-应用实例)
    - [6.1 网络拓扑分析](#61-网络拓扑分析)
    - [6.2 社交网络分析](#62-社交网络分析)
    - [6.3 路由算法](#63-路由算法)
  - [总结](#总结)

## 1. 基础定义

### 1.1 图的基本概念

**定义 1.1** (图)
图 $G = (V, E)$ 是一个有序对，其中：

- $V$ 是顶点集（vertex set），$V \neq \emptyset$
- $E$ 是边集（edge set），$E \subseteq V \times V$

**定义 1.2** (有向图)
有向图 $G = (V, E)$ 中，边集 $E$ 的元素是有序对 $(u, v)$，表示从顶点 $u$ 到顶点 $v$ 的有向边。

**定义 1.3** (无向图)
无向图 $G = (V, E)$ 中，边集 $E$ 的元素是无序对 $\{u, v\}$，表示顶点 $u$ 和 $v$ 之间的无向边。

### 1.2 图的类型

**定义 1.4** (完全图)
完全图 $K_n$ 是包含 $n$ 个顶点的图，其中任意两个顶点之间都有一条边。

**定义 1.5** (二分图)
二分图 $G = (V_1 \cup V_2, E)$ 中，顶点集可以划分为两个不相交的子集 $V_1$ 和 $V_2$，使得每条边的两个端点分别属于不同的子集。

**定义 1.6** (加权图)
加权图 $G = (V, E, w)$ 中，每条边 $e \in E$ 都有一个权重 $w(e) \in \mathbb{R}$。

### 1.3 图的表示

**邻接矩阵表示**：

对于图 $G = (V, E)$，邻接矩阵 $A$ 是一个 $|V| \times |V|$ 的矩阵，其中：

$$
A[i][j] = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}
$$

**邻接表表示**：
每个顶点维护一个包含其邻居顶点的列表。

## 2. 图的性质

### 2.1 连通性

**定义 2.1** (连通图)
无向图 $G$ 是连通的，当且仅当对于任意两个顶点 $u, v \in V$，存在从 $u$ 到 $v$ 的路径。

**定义 2.2** (强连通图)
有向图 $G$ 是强连通的，当且仅当对于任意两个顶点 $u, v \in V$，存在从 $u$ 到 $v$ 的有向路径。

**定理 2.1** (连通性判定)
图 $G = (V, E)$ 是连通的，当且仅当从任意顶点开始的深度优先搜索或广度优先搜索能够访问所有顶点。

**证明**：

- **必要性**：如果图是连通的，那么从任意顶点都存在到其他所有顶点的路径，因此DFS或BFS能够访问所有顶点。
- **充分性**：如果从某个顶点开始的遍历能够访问所有顶点，那么该顶点到所有其他顶点都存在路径，因此图是连通的。

### 2.2 路径和回路

**定义 2.3** (路径)
路径是顶点序列 $v_0, v_1, \ldots, v_k$，其中 $(v_i, v_{i+1}) \in E$ 对所有 $0 \leq i < k$ 成立。

**定义 2.4** (简单路径)
简单路径是不包含重复顶点的路径。

**定义 2.5** (回路)
回路是起点和终点相同的路径。

**定理 2.2** (路径存在性)
在连通图中，任意两个顶点之间都存在简单路径。

**证明**：
使用构造性证明。从起点开始，每次选择未访问的邻居顶点，直到到达终点。由于图是连通的，这样的路径总是存在。

### 2.3 树和森林

**定义 2.6** (树)
树是连通的无环无向图。

**定义 2.7** (森林)
森林是多个树的并集。

**定理 2.3** (树的性质)
对于包含 $n$ 个顶点的树 $T$：

1. $T$ 有 $n-1$ 条边
2. $T$ 是连通的
3. $T$ 中任意两个顶点之间有唯一的简单路径
4. 删除任意一条边会使图不连通
5. 添加任意一条边会形成一个环

## 3. 图的算法

### 3.1 遍历算法

```go
// 深度优先搜索
func DFS(g *Graph, start Vertex, visited map[Vertex]bool) {
    if visited[start] {
        return
    }
    
    visited[start] = true
    fmt.Printf("Visit vertex: %v\n", start)
    
    for _, neighbor := range g.Neighbors(start) {
        if !visited[neighbor] {
            DFS(g, neighbor, visited)
        }
    }
}

// 广度优先搜索
func BFS(g *Graph, start Vertex) {
    visited := make(map[Vertex]bool)
    queue := []Vertex{start}
    visited[start] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        fmt.Printf("Visit vertex: %v\n", current)
        
        for _, neighbor := range g.Neighbors(current) {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

### 3.2 最短路径算法

```go
// Dijkstra算法
func Dijkstra(g *WeightedGraph, start Vertex) map[Vertex]float64 {
    dist := make(map[Vertex]float64)
    visited := make(map[Vertex]bool)
    
    // 初始化距离
    for v := range g.Vertices() {
        dist[v] = math.Inf(1)
    }
    dist[start] = 0
    
    for len(visited) < len(g.Vertices()) {
        // 找到未访问的最小距离顶点
        u := getMinDistVertex(dist, visited)
        visited[u] = true
        
        // 更新邻居距离
        for _, v := range g.Neighbors(u) {
            if !visited[v] {
                newDist := dist[u] + g.Weight(u, v)
                if newDist < dist[v] {
                    dist[v] = newDist
                }
            }
        }
    }
    
    return dist
}
```

### 3.3 最小生成树

```go
// Kruskal算法
func Kruskal(g *WeightedGraph) []Edge {
    var mst []Edge
    edges := g.SortedEdges()
    uf := NewUnionFind(g.Vertices())
    
    for _, e := range edges {
        if !uf.Connected(e.From, e.To) {
            uf.Union(e.From, e.To)
            mst = append(mst, e)
        }
    }
    
    return mst
}

// Prim算法
func Prim(g *WeightedGraph, start Vertex) []Edge {
    var mst []Edge
    visited := make(map[Vertex]bool)
    pq := NewPriorityQueue()
    
    visited[start] = true
    for _, e := range g.EdgesFrom(start) {
        pq.Push(e)
    }
    
    for !pq.Empty() && len(mst) < len(g.Vertices())-1 {
        e := pq.Pop()
        if visited[e.To] {
            continue
        }
        
        visited[e.To] = true
        mst = append(mst, e)
        
        for _, next := range g.EdgesFrom(e.To) {
            if !visited[next.To] {
                pq.Push(next)
            }
        }
    }
    
    return mst
}
```

## 4. 特殊图类

### 4.1 平面图

**定义 4.1** (平面图)
平面图是可以在平面上画出的图，使得边只在顶点处相交。

**定理 4.1** (欧拉公式)
对于连通平面图 $G$，设 $v$ 是顶点数，$e$ 是边数，$f$ 是面数，则：
$$v - e + f = 2$$

### 4.2 欧拉图

**定义 4.2** (欧拉路径)
欧拉路径是遍历图中每条边恰好一次的路径。

**定义 4.3** (欧拉回路)
欧拉回路是遍历图中每条边恰好一次的回路。

**定理 4.2** (欧拉图判定)
无向图存在欧拉回路的充要条件是：图是连通的，且每个顶点的度数都是偶数。

### 4.3 哈密顿图

**定义 4.4** (哈密顿路径)
哈密顿路径是经过图中每个顶点恰好一次的路径。

**定义 4.5** (哈密顿回路)
哈密顿回路是经过图中每个顶点恰好一次的回路。

## 5. Go语言实现

### 5.1 数据结构

```go
// 图接口
type Graph interface {
    AddVertex(v Vertex)
    AddEdge(from, to Vertex)
    RemoveVertex(v Vertex)
    RemoveEdge(from, to Vertex)
    Vertices() []Vertex
    Edges() []Edge
    Neighbors(v Vertex) []Vertex
    HasEdge(from, to Vertex) bool
}

// 加权图接口
type WeightedGraph interface {
    Graph
    SetWeight(from, to Vertex, weight float64)
    Weight(from, to Vertex) float64
}

// 邻接表实现
type AdjListGraph struct {
    vertices map[Vertex]map[Vertex]float64
}

func NewAdjListGraph() *AdjListGraph {
    return &AdjListGraph{
        vertices: make(map[Vertex]map[Vertex]float64),
    }
}

func (g *AdjListGraph) AddVertex(v Vertex) {
    if _, exists := g.vertices[v]; !exists {
        g.vertices[v] = make(map[Vertex]float64)
    }
}

func (g *AdjListGraph) AddEdge(from, to Vertex, weight float64) {
    g.AddVertex(from)
    g.AddVertex(to)
    g.vertices[from][to] = weight
}
```

### 5.2 基础算法实现

```go
// 拓扑排序
func TopologicalSort(g *Graph) []Vertex {
    visited := make(map[Vertex]bool)
    stack := make([]Vertex, 0)
    
    var visit func(v Vertex)
    visit = func(v Vertex) {
        if visited[v] {
            return
        }
        visited[v] = true
        
        for _, neighbor := range g.Neighbors(v) {
            visit(neighbor)
        }
        stack = append([]Vertex{v}, stack...)
    }
    
    for _, v := range g.Vertices() {
        if !visited[v] {
            visit(v)
        }
    }
    
    return stack
}

// 强连通分量
func StronglyConnectedComponents(g *Graph) [][]Vertex {
    var components [][]Vertex
    visited := make(map[Vertex]bool)
    stack := make([]Vertex, 0)
    
    // 第一次DFS，填充栈
    for _, v := range g.Vertices() {
        if !visited[v] {
            fillOrder(g, v, visited, &stack)
        }
    }
    
    // 创建转置图
    gt := g.Transpose()
    
    // 清空访问标记
    for v := range visited {
        visited[v] = false
    }
    
    // 第二次DFS，找出强连通分量
    for len(stack) > 0 {
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if !visited[v] {
            var component []Vertex
            DFSUtil(gt, v, visited, &component)
            components = append(components, component)
        }
    }
    
    return components
}
```

### 5.3 高级算法实现

```go
// A*算法
func AStar(g *WeightedGraph, start, goal Vertex, heuristic func(Vertex) float64) []Vertex {
    openSet := NewPriorityQueue()
    openSet.Push(start, 0)
    
    cameFrom := make(map[Vertex]Vertex)
    gScore := make(map[Vertex]float64)
    fScore := make(map[Vertex]float64)
    
    for v := range g.Vertices() {
        gScore[v] = math.Inf(1)
        fScore[v] = math.Inf(1)
    }
    gScore[start] = 0
    fScore[start] = heuristic(start)
    
    for !openSet.Empty() {
        current := openSet.Pop()
        if current == goal {
            return reconstructPath(cameFrom, current)
        }
        
        for _, neighbor := range g.Neighbors(current) {
            tentativeGScore := gScore[current] + g.Weight(current, neighbor)
            
            if tentativeGScore < gScore[neighbor] {
                cameFrom[neighbor] = current
                gScore[neighbor] = tentativeGScore
                fScore[neighbor] = gScore[neighbor] + heuristic(neighbor)
                openSet.Push(neighbor, fScore[neighbor])
            }
        }
    }
    
    return nil
}
```

## 6. 应用实例

### 6.1 网络拓扑分析

```go
// 网络拓扑分析器
type NetworkAnalyzer struct {
    graph *WeightedGraph
}

func (na *NetworkAnalyzer) FindCriticalPaths() []Edge {
    var criticalPaths []Edge
    
    // 使用最小生成树找出关键路径
    mst := Kruskal(na.graph)
    
    // 分析每条边的重要性
    for _, edge := range mst {
        if na.isEdgeCritical(edge) {
            criticalPaths = append(criticalPaths, edge)
        }
    }
    
    return criticalPaths
}
```

### 6.2 社交网络分析

```go
// 社交网络分析器
type SocialNetworkAnalyzer struct {
    graph *Graph
}

func (sna *SocialNetworkAnalyzer) CalculateCentrality() map[Vertex]float64 {
    centrality := make(map[Vertex]float64)
    
    // 计算每个节点的中心性
    for v := range sna.graph.Vertices() {
        // 计算度中心性
        degreeCentrality := float64(len(sna.graph.Neighbors(v)))
        
        // 计算接近中心性
        closenessCentrality := sna.calculateCloseness(v)
        
        // 计算介数中心性
        betweennessCentrality := sna.calculateBetweenness(v)
        
        // 综合评分
        centrality[v] = (degreeCentrality + closenessCentrality + betweennessCentrality) / 3
    }
    
    return centrality
}
```

### 6.3 路由算法

```go
// 路由器
type Router struct {
    graph *WeightedGraph
    routingTable map[Vertex]map[Vertex]Vertex
}

func (r *Router) UpdateRoutingTable() {
    // 使用Floyd-Warshall算法更新路由表
    dist := make(map[Vertex]map[Vertex]float64)
    next := make(map[Vertex]map[Vertex]Vertex)
    
    // 初始化
    for u := range r.graph.Vertices() {
        dist[u] = make(map[Vertex]float64)
        next[u] = make(map[Vertex]Vertex)
        for v := range r.graph.Vertices() {
            dist[u][v] = math.Inf(1)
            if r.graph.HasEdge(u, v) {
                dist[u][v] = r.graph.Weight(u, v)
                next[u][v] = v
            }
        }
        dist[u][u] = 0
        next[u][u] = u
    }
    
    // Floyd-Warshall算法
    for k := range r.graph.Vertices() {
        for i := range r.graph.Vertices() {
            for j := range r.graph.Vertices() {
                if dist[i][k] + dist[k][j] < dist[i][j] {
                    dist[i][j] = dist[i][k] + dist[k][j]
                    next[i][j] = next[i][k]
                }
            }
        }
    }
    
    r.routingTable = next
}
```

## 总结

图论为计算机科学提供了强大的理论基础和实用工具，特别是在以下方面：

1. **网络分析**：
   - 社交网络分析
   - 通信网络优化
   - 路由算法设计

2. **算法设计**：
   - 最短路径算法
   - 最小生成树算法
   - 网络流算法

3. **应用领域**：
   - 交通规划
   - 资源调度
   - 电路设计
   - 生物信息学

通过Go语言的实现，我们可以将这些理论概念转化为实用的工程解决方案，为实际问题提供高效的解决方法。
