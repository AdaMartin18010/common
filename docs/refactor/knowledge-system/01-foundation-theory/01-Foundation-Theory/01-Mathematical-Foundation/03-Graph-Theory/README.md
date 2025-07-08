# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 图的定义](#11-图的定义)
    - [1.2 图的类型](#12-图的类型)
    - [1.3 基本性质](#13-基本性质)
  - [2. 图的表示](#2-图的表示)
    - [2.1 邻接矩阵](#21-邻接矩阵)
    - [2.2 邻接表](#22-邻接表)
    - [2.3 边列表](#23-边列表)
  - [3. 图算法](#3-图算法)
    - [3.1 遍历算法](#31-遍历算法)
      - [3.1.1 深度优先搜索 (DFS)](#311-深度优先搜索-dfs)
      - [3.1.2 广度优先搜索 (BFS)](#312-广度优先搜索-bfs)
    - [3.2 最短路径](#32-最短路径)
      - [3.2.1 Dijkstra算法](#321-dijkstra算法)
      - [3.2.2 Floyd-Warshall算法](#322-floyd-warshall算法)
    - [3.3 最小生成树](#33-最小生成树)
      - [3.3.1 Kruskal算法](#331-kruskal算法)
      - [3.3.2 Prim算法](#332-prim算法)
  - [4. 形式化定义](#4-形式化定义)
    - [4.1 图的同构](#41-图的同构)
    - [4.2 图的连通性](#42-图的连通性)
    - [4.3 图的着色](#43-图的着色)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 图的数据结构](#51-图的数据结构)
    - [5.2 Dijkstra算法实现](#52-dijkstra算法实现)
    - [5.3 最小生成树实现](#53-最小生成树实现)
    - [5.4 使用示例](#54-使用示例)
  - [总结](#总结)
    - [关键要点](#关键要点)
    - [进一步研究方向](#进一步研究方向)

## 1. 基础概念

### 1.1 图的定义

**定义 1.1**: 图 ```latex
G = (V, E)
``` 是一个有序对，其中：

- ```latex
V
``` 是顶点集（vertex set），```latex
V \neq \emptyset
```
- ```latex
E
``` 是边集（edge set），```latex
E \subseteq V \times V
```

**定义 1.2**: 对于边 ```latex
e = (u, v) \in E
```：

- ```latex
u
``` 和 ```latex
v
``` 是 ```latex
e
``` 的端点
- ```latex
u
``` 和 ```latex
v
``` 是相邻的（adjacent）
- ```latex
e
``` 与 ```latex
u
``` 和 ```latex
v
``` 相关联（incident）

### 1.2 图的类型

**定义 1.3**: 无向图（Undirected Graph）

- 边是无序对：```latex
(u, v) = (v, u)
```
- 邻接关系是对称的

**定义 1.4**: 有向图（Directed Graph）

- 边是有序对：```latex
(u, v) \neq (v, u)
```
- 邻接关系不对称

**定义 1.5**: 加权图（Weighted Graph）

- 每条边 ```latex
e
``` 关联一个权重 ```latex
w(e) \in \mathbb{R}
```

### 1.3 基本性质

**定义 1.6**: 顶点的度（Degree）

- 无向图：```latex
deg(v) = |\{e \in E : v \text{ 是 } e \text{ 的端点}\}|
```
- 有向图：入度 ```latex
deg^-(v)
``` 和出度 ```latex
deg^+(v)
```

**定理 1.1**: 握手定理
对于无向图 ```latex
G = (V, E)
```：
$```latex
\sum_{v \in V} deg(v) = 2|E|
```$

**证明**:
每条边贡献给两个顶点的度，因此总度数等于边数的两倍。

## 2. 图的表示

### 2.1 邻接矩阵

**定义 2.1**: 邻接矩阵 ```latex
A
``` 是一个 ```latex
|V| \times |V|
``` 的矩阵：

$
A[i][j] = \begin{cases}
1 & \text{if } (i, j) \in E \\
0 & \text{otherwise}
\end{cases}
$

**性质**:

- 无向图的邻接矩阵是对称的
- 对角线元素为0（无自环）
- 空间复杂度：```latex
O(|V|^2)
```

### 2.2 邻接表

**定义 2.2**: 邻接表为每个顶点维护一个邻接顶点列表：
$```latex
Adj[v] = \{u \in V : (v, u) \in E\}
```$

**性质**:

- 空间复杂度：```latex
O(|V| + |E|)
```
- 适合稀疏图
- 遍历邻接顶点效率高

### 2.3 边列表

**定义 2.3**: 边列表直接存储所有边：
$```latex
E = \{(u_1, v_1), (u_2, v_2), \ldots, (u_m, v_m)\}
```$

## 3. 图算法

### 3.1 遍历算法

#### 3.1.1 深度优先搜索 (DFS)

**算法描述**:

```go
func DFS(graph *Graph, start int, visited map[int]bool) {
    visited[start] = true
    fmt.Printf("Visit: %d\n", start)

    for _, neighbor := range graph.adjList[start] {
        if !visited[neighbor] {
            DFS(graph, neighbor, visited)
        }
    }
}
```

**时间复杂度**: ```latex
O(|V| + |E|)
```

#### 3.1.2 广度优先搜索 (BFS)

**算法描述**:

```go
func BFS(graph *Graph, start int) {
    queue := []int{start}
    visited := make(map[int]bool)
    visited[start] = true

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        fmt.Printf("Visit: %d\n", current)

        for _, neighbor := range graph.adjList[current] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

**时间复杂度**: ```latex
O(|V| + |E|)
```

### 3.2 最短路径

#### 3.2.1 Dijkstra算法

**定理 3.1**: Dijkstra算法正确性
对于非负权图，Dijkstra算法能找到从源点到所有其他顶点的最短路径。

**算法复杂度**: ```latex
O((|V| + |E|) \log |V|)
```

#### 3.2.2 Floyd-Warshall算法

**定理 3.2**: Floyd-Warshall算法
对于任意图，Floyd-Warshall算法能找到所有顶点对之间的最短路径。

**算法复杂度**: ```latex
O(|V|^3)
```

### 3.3 最小生成树

#### 3.3.1 Kruskal算法

**定理 3.3**: Kruskal算法正确性
Kruskal算法能找到图的最小生成树。

**算法复杂度**: ```latex
O(|E| \log |E|)
```

#### 3.3.2 Prim算法

**定理 3.4**: Prim算法正确性
Prim算法能找到图的最小生成树。

**算法复杂度**: ```latex
O((|V| + |E|) \log |V|)
```

## 4. 形式化定义

### 4.1 图的同构

**定义 4.1**: 图同构
两个图 ```latex
G_1 = (V_1, E_1)
``` 和 ```latex
G_2 = (V_2, E_2)
``` 是同构的，如果存在双射 ```latex
f: V_1 \rightarrow V_2
``` 使得：
$```latex
(u, v) \in E_1 \iff (f(u), f(v)) \in E_2
```$

### 4.2 图的连通性

**定义 4.2**: 连通图
无向图 ```latex
G
``` 是连通的，如果对于任意两个顶点 ```latex
u, v \in V
```，存在从 ```latex
u
``` 到 ```latex
v
``` 的路径。

**定义 4.3**: 强连通图
有向图 ```latex
G
``` 是强连通的，如果对于任意两个顶点 ```latex
u, v \in V
```，存在从 ```latex
u
``` 到 ```latex
v
``` 的有向路径。

### 4.3 图的着色

**定义 4.4**: 图着色
图 ```latex
G
``` 的 ```latex
k
```-着色是一个函数 ```latex
c: V \rightarrow \{1, 2, \ldots, k\}
```，使得相邻顶点有不同的颜色。

**定理 4.1**: 四色定理
任何平面图都可以用四种颜色着色。

## 5. Go语言实现

### 5.1 图的数据结构

```go
package graph

import (
    "container/heap"
    "fmt"
    "math"
)

// Graph 表示图的数据结构
type Graph struct {
    vertices int
    adjList  map[int][]Edge
}

// Edge 表示边
type Edge struct {
    to     int
    weight float64
}

// NewGraph 创建新图
func NewGraph(vertices int) *Graph {
    return &Graph{
        vertices: vertices,
        adjList:  make(map[int][]Edge),
    }
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to int, weight float64) {
    g.adjList[from] = append(g.adjList[from], Edge{to: to, weight: weight})
}

// AddUndirectedEdge 添加无向边
func (g *Graph) AddUndirectedEdge(u, v int, weight float64) {
    g.AddEdge(u, v, weight)
    g.AddEdge(v, u, weight)
}
```

### 5.2 Dijkstra算法实现

```go
// Dijkstra 实现Dijkstra最短路径算法
func (g *Graph) Dijkstra(start int) map[int]float64 {
    distances := make(map[int]float64)
    for i := 0; i < g.vertices; i++ {
        distances[i] = math.Inf(1)
    }
    distances[start] = 0

    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &Item{vertex: start, distance: 0})

    for pq.Len() > 0 {
        current := heap.Pop(pq).(*Item)

        if current.distance > distances[current.vertex] {
            continue
        }

        for _, edge := range g.adjList[current.vertex] {
            newDist := distances[current.vertex] + edge.weight
            if newDist < distances[edge.to] {
                distances[edge.to] = newDist
                heap.Push(pq, &Item{vertex: edge.to, distance: newDist})
            }
        }
    }

    return distances
}

// PriorityQueue 优先队列实现
type PriorityQueue []*Item

type Item struct {
    vertex   int
    distance float64
    index    int
}

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
```

### 5.3 最小生成树实现

```go
// Kruskal 实现Kruskal最小生成树算法
func (g *Graph) Kruskal() []Edge {
    var edges []Edge
    for from, edgeList := range g.adjList {
        for _, edge := range edgeList {
            if from < edge.to { // 避免重复边
                edges = append(edges, Edge{from: from, to: edge.to, weight: edge.weight})
            }
        }
    }

    // 按权重排序
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].weight < edges[j].weight
    })

    uf := NewUnionFind(g.vertices)
    var mst []Edge

    for _, edge := range edges {
        if uf.Find(edge.from) != uf.Find(edge.to) {
            uf.Union(edge.from, edge.to)
            mst = append(mst, edge)
        }
    }

    return mst
}

// UnionFind 并查集实现
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(n int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, n),
        rank:   make([]int, n),
    }
    for i := 0; i < n; i++ {
        uf.parent[i] = i
    }
    return uf
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
    px, py := uf.Find(x), uf.Find(y)
    if px == py {
        return
    }

    if uf.rank[px] < uf.rank[py] {
        uf.parent[px] = py
    } else if uf.rank[px] > uf.rank[py] {
        uf.parent[py] = px
    } else {
        uf.parent[py] = px
        uf.rank[px]++
    }
}
```

### 5.4 使用示例

```go
func main() {
    // 创建图
    g := NewGraph(5)

    // 添加边
    g.AddUndirectedEdge(0, 1, 4)
    g.AddUndirectedEdge(0, 2, 2)
    g.AddUndirectedEdge(1, 2, 1)
    g.AddUndirectedEdge(1, 3, 5)
    g.AddUndirectedEdge(2, 3, 8)
    g.AddUndirectedEdge(2, 4, 10)
    g.AddUndirectedEdge(3, 4, 2)

    // 计算最短路径
    distances := g.Dijkstra(0)
    fmt.Println("Shortest distances from vertex 0:")
    for vertex, distance := range distances {
        fmt.Printf("To %d: %.2f\n", vertex, distance)
    }

    // 计算最小生成树
    mst := g.Kruskal()
    fmt.Println("\nMinimum Spanning Tree:")
    totalWeight := 0.0
    for _, edge := range mst {
        fmt.Printf("Edge %d-%d: %.2f\n", edge.from, edge.to, edge.weight)
        totalWeight += edge.weight
    }
    fmt.Printf("Total weight: %.2f\n", totalWeight)
}
```

## 总结

图论是计算机科学中的基础理论，广泛应用于网络分析、路由算法、社交网络分析等领域。通过形式化定义和Go语言实现，我们建立了从理论到实践的完整框架。

### 关键要点

1. **理论基础**: 图的定义、性质和基本定理
2. **算法设计**: 遍历、最短路径、最小生成树等经典算法
3. **实现技术**: 邻接表、优先队列、并查集等数据结构
4. **应用场景**: 网络路由、社交网络、生物信息学等

### 进一步研究方向

1. **高级图算法**: 最大流、匹配、着色等
2. **图神经网络**: 深度学习在图上的应用
3. **动态图**: 随时间变化的图结构
4. **大规模图处理**: 分布式图计算框架

## 详细内容
- 背景与定义：
- 关键概念：
- 相关原理：
- 实践应用：
- 典型案例：
- 拓展阅读：

## 参考文献
- [示例参考文献1](#)
- [示例参考文献2](#)

## 标签
- #待补充 #知识点 #标签