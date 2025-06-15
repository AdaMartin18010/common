# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [概述](#概述)
  - [基本概念](#基本概念)
    - [图的定义](#图的定义)
    - [图的类型](#图的类型)
    - [图的基本性质](#图的基本性质)
  - [形式化理论](#形式化理论)
    - [图的形式化定义](#图的形式化定义)
    - [图的基本定理](#图的基本定理)
    - [图的代数结构](#图的代数结构)
  - [算法与复杂度](#算法与复杂度)
    - [图遍历算法](#图遍历算法)
      - [深度优先搜索 (DFS)](#深度优先搜索-dfs)
      - [广度优先搜索 (BFS)](#广度优先搜索-bfs)
    - [最短路径算法](#最短路径算法)
      - [Dijkstra算法](#dijkstra算法)
    - [最小生成树算法](#最小生成树算法)
      - [Kruskal算法](#kruskal算法)
    - [网络流算法](#网络流算法)
      - [Ford-Fulkerson算法](#ford-fulkerson算法)
  - [Go语言实现](#go语言实现)
    - [图的数据结构](#图的数据结构)
    - [图的基本操作](#图的基本操作)
    - [图算法实现](#图算法实现)
  - [应用领域](#应用领域)
    - [网络分析](#网络分析)
    - [社交网络](#社交网络)
    - [路由算法](#路由算法)
    - [依赖分析](#依赖分析)
  - [相关链接](#相关链接)

## 概述

图论是研究图（Graph）的数学分支，图是由顶点（Vertex）和边（Edge）组成的数学结构。图论在计算机科学中有着广泛的应用，包括网络分析、算法设计、数据结构等领域。

## 基本概念

### 图的定义

**定义 1 (图)**: 图 $G = (V, E)$ 是一个有序对，其中：

- $V$ 是顶点的有限集合，称为顶点集
- $E$ 是边的集合，每条边是顶点集 $V$ 中两个顶点的无序对

**定义 2 (有向图)**: 有向图 $D = (V, A)$ 是一个有序对，其中：

- $V$ 是顶点的有限集合
- $A$ 是弧的集合，每条弧是顶点集 $V$ 中两个顶点的有序对

### 图的类型

1. **无向图**: 边没有方向
2. **有向图**: 边有方向
3. **加权图**: 边有权重
4. **多重图**: 允许重边
5. **简单图**: 无重边无自环

### 图的基本性质

**定义 3 (度数)**: 顶点 $v$ 的度数 $deg(v)$ 是与 $v$ 相邻的边数

**定理 1 (握手定理)**: 对于任意图 $G = (V, E)$，
$$\sum_{v \in V} deg(v) = 2|E|$$

**证明**: 每条边贡献给两个顶点的度数，因此所有顶点的度数之和等于边数的两倍。

## 形式化理论

### 图的形式化定义

**定义 4 (邻接矩阵)**: 图 $G = (V, E)$ 的邻接矩阵 $A$ 是一个 $n \times n$ 矩阵，其中：
$$A_{ij} = \begin{cases}
1 & \text{if } (v_i, v_j) \in E \\
0 & \text{otherwise}
\end{cases}$$

**定义 5 (邻接表)**: 图 $G = (V, E)$ 的邻接表是一个映射 $adj: V \rightarrow 2^V$，其中 $adj(v)$ 是与 $v$ 相邻的顶点集合。

### 图的基本定理

**定理 2 (欧拉定理)**: 连通图 $G$ 存在欧拉回路的充要条件是所有顶点的度数都是偶数。

**定理 3 (哈密顿定理)**: 对于 $n \geq 3$ 的完全图 $K_n$，存在哈密顿回路。

**定理 4 (平面图欧拉公式)**: 对于连通平面图 $G = (V, E, F)$，
$$|V| - |E| + |F| = 2$$
其中 $F$ 是面的集合。

### 图的代数结构

**定义 6 (图同构)**: 两个图 $G_1 = (V_1, E_1)$ 和 $G_2 = (V_2, E_2)$ 同构，如果存在双射 $f: V_1 \rightarrow V_2$，使得：
$$(u, v) \in E_1 \Leftrightarrow (f(u), f(v)) \in E_2$$

## 算法与复杂度

### 图遍历算法

#### 深度优先搜索 (DFS)

**算法 1 (DFS)**:
```go
func DFS(g *Graph, start Vertex) {
    visited := make(map[Vertex]bool)
    dfsHelper(g, start, visited)
}

func dfsHelper(g *Graph, v Vertex, visited map[Vertex]bool) {
    visited[v] = true
    fmt.Printf("Visit: %v\n", v)

    for _, neighbor := range g.AdjacencyList[v] {
        if !visited[neighbor] {
            dfsHelper(g, neighbor, visited)
        }
    }
}
```

**时间复杂度**: $O(|V| + |E|)$
**空间复杂度**: $O(|V|)$

#### 广度优先搜索 (BFS)

**算法 2 (BFS)**:
```go
func BFS(g *Graph, start Vertex) {
    visited := make(map[Vertex]bool)
    queue := []Vertex{start}
    visited[start] = true

    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        fmt.Printf("Visit: %v\n", v)

        for _, neighbor := range g.AdjacencyList[v] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

**时间复杂度**: $O(|V| + |E|)$
**空间复杂度**: $O(|V|)$

### 最短路径算法

#### Dijkstra算法

**算法 3 (Dijkstra)**:
```go
func Dijkstra(g *Graph, start Vertex) map[Vertex]int {
    distances := make(map[Vertex]int)
    for v := range g.AdjacencyList {
        distances[v] = math.MaxInt32
    }
    distances[start] = 0

    pq := &PriorityQueue{}
    heap.Push(pq, &Item{vertex: start, distance: 0})

    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        u := item.vertex

        if item.distance > distances[u] {
            continue
        }

        for v, weight := range g.AdjacencyList[u] {
            if distances[u]+weight < distances[v] {
                distances[v] = distances[u] + weight
                heap.Push(pq, &Item{vertex: v, distance: distances[v]})
            }
        }
    }

    return distances
}
```

**时间复杂度**: $O((|V| + |E|) \log |V|)$
**空间复杂度**: $O(|V|)$

### 最小生成树算法

#### Kruskal算法

**算法 4 (Kruskal)**:
```go
func Kruskal(g *Graph) []Edge {
    var edges []Edge
    for u := range g.AdjacencyList {
        for v, weight := range g.AdjacencyList[u] {
            edges = append(edges, Edge{u, v, weight})
        }
    }

    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })

    uf := NewUnionFind(len(g.AdjacencyList))
    var mst []Edge

    for _, edge := range edges {
        if uf.Find(edge.From) != uf.Find(edge.To) {
            mst = append(mst, edge)
            uf.Union(edge.From, edge.To)
        }
    }

    return mst
}
```

**时间复杂度**: $O(|E| \log |E|)$
**空间复杂度**: $O(|V|)$

### 网络流算法

#### Ford-Fulkerson算法

**算法 5 (Ford-Fulkerson)**:
```go
func FordFulkerson(g *Graph, source, sink Vertex) int {
    residual := g.Copy()
    maxFlow := 0

    for {
        path := findAugmentingPath(residual, source, sink)
        if path == nil {
            break
        }

        minCapacity := findMinCapacity(path)
        maxFlow += minCapacity

        // 更新残量图
        for i := 0; i < len(path)-1; i++ {
            u, v := path[i], path[i+1]
            residual.AdjacencyList[u][v] -= minCapacity
            residual.AdjacencyList[v][u] += minCapacity
        }
    }

    return maxFlow
}
```

**时间复杂度**: $O(|E| \cdot f^*)$，其中 $f^*$ 是最大流值
**空间复杂度**: $O(|V| + |E|)$

## Go语言实现

### 图的数据结构

```go
// Vertex 表示图的顶点
type Vertex int

// Edge 表示图的边
type Edge struct {
    From   Vertex
    To     Vertex
    Weight int
}

// Graph 表示图
type Graph struct {
    AdjacencyList map[Vertex]map[Vertex]int
    Directed      bool
}

// NewGraph 创建新图
func NewGraph(directed bool) *Graph {
    return &Graph{
        AdjacencyList: make(map[Vertex]map[Vertex]int),
        Directed:      directed,
    }
}

// AddVertex 添加顶点
func (g *Graph) AddVertex(v Vertex) {
    if g.AdjacencyList[v] == nil {
        g.AdjacencyList[v] = make(map[Vertex]int)
    }
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to Vertex, weight int) {
    g.AddVertex(from)
    g.AddVertex(to)

    g.AdjacencyList[from][to] = weight
    if !g.Directed {
        g.AdjacencyList[to][from] = weight
    }
}

// RemoveEdge 删除边
func (g *Graph) RemoveEdge(from, to Vertex) {
    delete(g.AdjacencyList[from], to)
    if !g.Directed {
        delete(g.AdjacencyList[to], from)
    }
}

// GetDegree 获取顶点度数
func (g *Graph) GetDegree(v Vertex) int {
    return len(g.AdjacencyList[v])
}

// IsConnected 检查图是否连通
func (g *Graph) IsConnected() bool {
    if len(g.AdjacencyList) == 0 {
        return true
    }

    visited := make(map[Vertex]bool)
    var start Vertex
    for v := range g.AdjacencyList {
        start = v
        break
    }

    g.dfs(start, visited)

    return len(visited) == len(g.AdjacencyList)
}

// dfs 深度优先搜索辅助函数
func (g *Graph) dfs(v Vertex, visited map[Vertex]bool) {
    visited[v] = true
    for neighbor := range g.AdjacencyList[v] {
        if !visited[neighbor] {
            g.dfs(neighbor, visited)
        }
    }
}
```

### 图的基本操作

```go
// Copy 复制图
func (g *Graph) Copy() *Graph {
    newGraph := NewGraph(g.Directed)
    for v, neighbors := range g.AdjacencyList {
        newGraph.AdjacencyList[v] = make(map[Vertex]int)
        for neighbor, weight := range neighbors {
            newGraph.AdjacencyList[v][neighbor] = weight
        }
    }
    return newGraph
}

// Transpose 转置图（仅对有向图有效）
func (g *Graph) Transpose() *Graph {
    if !g.Directed {
        return g.Copy()
    }

    transposed := NewGraph(true)
    for v, neighbors := range g.AdjacencyList {
        for neighbor, weight := range neighbors {
            transposed.AddEdge(neighbor, v, weight)
        }
    }
    return transposed
}

// GetConnectedComponents 获取连通分量
func (g *Graph) GetConnectedComponents() [][]Vertex {
    visited := make(map[Vertex]bool)
    var components [][]Vertex

    for v := range g.AdjacencyList {
        if !visited[v] {
            var component []Vertex
            g.dfsComponent(v, visited, &component)
            components = append(components, component)
        }
    }

    return components
}

// dfsComponent 获取连通分量的DFS辅助函数
func (g *Graph) dfsComponent(v Vertex, visited map[Vertex]bool, component *[]Vertex) {
    visited[v] = true
    *component = append(*component, v)

    for neighbor := range g.AdjacencyList[v] {
        if !visited[neighbor] {
            g.dfsComponent(neighbor, visited, component)
        }
    }
}
```

### 图算法实现

```go
// TopologicalSort 拓扑排序（仅对有向无环图）
func (g *Graph) TopologicalSort() ([]Vertex, error) {
    if !g.Directed {
        return nil, fmt.Errorf("topological sort only works on directed graphs")
    }

    inDegree := make(map[Vertex]int)
    for v := range g.AdjacencyList {
        inDegree[v] = 0
    }

    // 计算入度
    for _, neighbors := range g.AdjacencyList {
        for neighbor := range neighbors {
            inDegree[neighbor]++
        }
    }

    // 使用队列进行拓扑排序
    var queue []Vertex
    for v, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, v)
        }
    }

    var result []Vertex
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)

        for neighbor := range g.AdjacencyList[v] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != len(g.AdjacencyList) {
        return nil, fmt.Errorf("graph contains cycle")
    }

    return result, nil
}

// FloydWarshall 全源最短路径算法
func (g *Graph) FloydWarshall() map[Vertex]map[Vertex]int {
    // 初始化距离矩阵
    dist := make(map[Vertex]map[Vertex]int)
    for u := range g.AdjacencyList {
        dist[u] = make(map[Vertex]int)
        for v := range g.AdjacencyList {
            if u == v {
                dist[u][v] = 0
            } else if weight, exists := g.AdjacencyList[u][v]; exists {
                dist[u][v] = weight
            } else {
                dist[u][v] = math.MaxInt32
            }
        }
    }

    // Floyd-Warshall算法核心
    for k := range g.AdjacencyList {
        for i := range g.AdjacencyList {
            for j := range g.AdjacencyList {
                if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 {
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

## 应用领域

### 网络分析

图论在网络分析中用于：
- 网络拓扑分析
- 路由算法设计
- 网络流量优化
- 故障检测和恢复

### 社交网络

图论在社交网络分析中用于：
- 社区发现
- 影响力分析
- 信息传播建模
- 推荐系统

### 路由算法

图论在路由算法中用于：
- 最短路径计算
- 负载均衡
- 网络规划
- 流量工程

### 依赖分析

图论在软件工程中用于：
- 模块依赖分析
- 编译顺序确定
- 循环依赖检测
- 软件架构分析

## 相关链接

- [01-集合论 (Set Theory)](../01-Set-Theory/README.md)
- [02-逻辑学 (Logic)](../02-Logic/README.md)
- [04-概率论 (Probability Theory)](../04-Probability-Theory/README.md)
- [08-软件工程形式化 (Software Engineering Formalization)](../../08-Software-Engineering-Formalization/README.md)
- [09-编程语言理论 (Programming Language Theory)](../../09-Programming-Language-Theory/README.md)
