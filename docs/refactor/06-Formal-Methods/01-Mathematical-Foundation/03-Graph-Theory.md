# 03-图论 (Graph Theory)

## 概述

图论是离散数学的重要分支，在计算机科学中有广泛应用。本文档介绍图论的基本概念、算法以及在Go语言中的实现。

## 目录

1. [图的基本概念 (Basic Concepts)](#1-图的基本概念-basic-concepts)
2. [图的表示 (Graph Representation)](#2-图的表示-graph-representation)
3. [图遍历算法 (Graph Traversal)](#3-图遍历算法-graph-traversal)
4. [最短路径算法 (Shortest Path)](#4-最短路径算法-shortest-path)
5. [最小生成树 (Minimum Spanning Tree)](#5-最小生成树-minimum-spanning-tree)
6. [网络流 (Network Flow)](#6-网络流-network-flow)

---

## 1. 图的基本概念 (Basic Concepts)

### 1.1 定义

**定义 1.1.1** (图)
图 $G = (V, E)$ 由顶点集 $V$ 和边集 $E$ 组成，其中 $E \subseteq V \times V$。

**定义 1.1.2** (有向图)
有向图 $G = (V, E)$ 中，边是有序对 $(u, v) \in E$。

**定义 1.1.3** (无向图)
无向图中，边是无序对 $\{u, v\} \in E$。

### 1.2 基本性质

**定义 1.2.1** (度数)
- 入度: $deg^-(v) = |\{u : (u, v) \in E\}|$
- 出度: $deg^+(v) = |\{u : (v, u) \in E\}|$
- 度数: $deg(v) = deg^-(v) + deg^+(v)$

**定理 1.2.1** (握手定理)
对于无向图 $G = (V, E)$，$\sum_{v \in V} deg(v) = 2|E|$

**证明**:
每条边贡献给两个顶点的度数，因此总度数等于边数的两倍。

### 1.3 Go语言实现

```go
package graph

import (
    "fmt"
    "math"
)

// Vertex 顶点
type Vertex struct {
    ID   int
    Data interface{}
}

// Edge 边
type Edge struct {
    From   int
    To     int
    Weight float64
}

// Graph 图接口
type Graph interface {
    AddVertex(id int, data interface{})
    AddEdge(from, to int, weight float64)
    RemoveVertex(id int)
    RemoveEdge(from, to int)
    GetVertex(id int) (*Vertex, bool)
    GetEdge(from, to int) (*Edge, bool)
    Vertices() []*Vertex
    Edges() []*Edge
    Adjacent(vertexID int) []*Vertex
    IsDirected() bool
}

// AdjacencyListGraph 邻接表图
type AdjacencyListGraph struct {
    vertices map[int]*Vertex
    edges    map[int]map[int]*Edge
    directed bool
}

// NewAdjacencyListGraph 创建邻接表图
func NewAdjacencyListGraph(directed bool) *AdjacencyListGraph {
    return &AdjacencyListGraph{
        vertices: make(map[int]*Vertex),
        edges:    make(map[int]map[int]*Edge),
        directed: directed,
    }
}

func (g *AdjacencyListGraph) AddVertex(id int, data interface{}) {
    g.vertices[id] = &Vertex{ID: id, Data: data}
    if g.edges[id] == nil {
        g.edges[id] = make(map[int]*Edge)
    }
}

func (g *AdjacencyListGraph) AddEdge(from, to int, weight float64) {
    if g.vertices[from] == nil || g.vertices[to] == nil {
        return
    }
    
    edge := &Edge{From: from, To: to, Weight: weight}
    g.edges[from][to] = edge
    
    if !g.directed {
        g.edges[to][from] = &Edge{From: to, To: from, Weight: weight}
    }
}

func (g *AdjacencyListGraph) RemoveVertex(id int) {
    delete(g.vertices, id)
    delete(g.edges, id)
    
    // 删除指向该顶点的边
    for _, adjList := range g.edges {
        delete(adjList, id)
    }
}

func (g *AdjacencyListGraph) RemoveEdge(from, to int) {
    if g.edges[from] != nil {
        delete(g.edges[from], to)
    }
    if !g.directed && g.edges[to] != nil {
        delete(g.edges[to], from)
    }
}

func (g *AdjacencyListGraph) GetVertex(id int) (*Vertex, bool) {
    vertex, exists := g.vertices[id]
    return vertex, exists
}

func (g *AdjacencyListGraph) GetEdge(from, to int) (*Edge, bool) {
    if g.edges[from] != nil {
        edge, exists := g.edges[from][to]
        return edge, exists
    }
    return nil, false
}

func (g *AdjacencyListGraph) Vertices() []*Vertex {
    vertices := make([]*Vertex, 0, len(g.vertices))
    for _, vertex := range g.vertices {
        vertices = append(vertices, vertex)
    }
    return vertices
}

func (g *AdjacencyListGraph) Edges() []*Edge {
    edges := make([]*Edge, 0)
    for _, adjList := range g.edges {
        for _, edge := range adjList {
            edges = append(edges, edge)
        }
    }
    return edges
}

func (g *AdjacencyListGraph) Adjacent(vertexID int) []*Vertex {
    var adjacent []*Vertex
    if adjList, exists := g.edges[vertexID]; exists {
        adjacent = make([]*Vertex, 0, len(adjList))
        for toID := range adjList {
            if vertex, exists := g.vertices[toID]; exists {
                adjacent = append(adjacent, vertex)
            }
        }
    }
    return adjacent
}

func (g *AdjacencyListGraph) IsDirected() bool {
    return g.directed
}

// Degree 计算度数
func (g *AdjacencyListGraph) Degree(vertexID int) (inDegree, outDegree int) {
    if g.directed {
        // 入度
        for _, adjList := range g.edges {
            if _, exists := adjList[vertexID]; exists {
                inDegree++
            }
        }
        // 出度
        if adjList, exists := g.edges[vertexID]; exists {
            outDegree = len(adjList)
        }
    } else {
        // 无向图度数相等
        if adjList, exists := g.edges[vertexID]; exists {
            outDegree = len(adjList)
        }
        inDegree = outDegree
    }
    return
}
```

---

## 2. 图的表示 (Graph Representation)

### 2.1 邻接矩阵

**定义 2.1.1** (邻接矩阵)
对于图 $G = (V, E)$，邻接矩阵 $A$ 定义为：
$A[i][j] = \begin{cases} 
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$

### 2.2 Go语言实现

```go
package representation

// AdjacencyMatrixGraph 邻接矩阵图
type AdjacencyMatrixGraph struct {
    matrix   [][]float64
    vertices map[int]*Vertex
    directed bool
}

// NewAdjacencyMatrixGraph 创建邻接矩阵图
func NewAdjacencyMatrixGraph(size int, directed bool) *AdjacencyMatrixGraph {
    matrix := make([][]float64, size)
    for i := range matrix {
        matrix[i] = make([]float64, size)
        for j := range matrix[i] {
            matrix[i][j] = math.Inf(1) // 表示无边
        }
    }
    
    return &AdjacencyMatrixGraph{
        matrix:   matrix,
        vertices: make(map[int]*Vertex),
        directed: directed,
    }
}

func (g *AdjacencyMatrixGraph) AddVertex(id int, data interface{}) {
    g.vertices[id] = &Vertex{ID: id, Data: data}
}

func (g *AdjacencyMatrixGraph) AddEdge(from, to int, weight float64) {
    if from < len(g.matrix) && to < len(g.matrix) {
        g.matrix[from][to] = weight
        if !g.directed {
            g.matrix[to][from] = weight
        }
    }
}

func (g *AdjacencyMatrixGraph) GetEdge(from, to int) (float64, bool) {
    if from < len(g.matrix) && to < len(g.matrix) {
        weight := g.matrix[from][to]
        if !math.IsInf(weight, 1) {
            return weight, true
        }
    }
    return 0, false
}

// IncidenceMatrix 关联矩阵
type IncidenceMatrix struct {
    matrix   [][]int
    vertices []*Vertex
    edges    []*Edge
}

// NewIncidenceMatrix 创建关联矩阵
func NewIncidenceMatrix(vertices []*Vertex, edges []*Edge) *IncidenceMatrix {
    matrix := make([][]int, len(vertices))
    for i := range matrix {
        matrix[i] = make([]int, len(edges))
    }
    
    // 填充关联矩阵
    for j, edge := range edges {
        matrix[edge.From][j] = 1
        matrix[edge.To][j] = -1
    }
    
    return &IncidenceMatrix{
        matrix:   matrix,
        vertices: vertices,
        edges:    edges,
    }
}
```

---

## 3. 图遍历算法 (Graph Traversal)

### 3.1 深度优先搜索 (DFS)

**算法 3.1.1** (DFS)
```python
def DFS(G, start):
    visited = set()
    stack = [start]
    
    while stack:
        vertex = stack.pop()
        if vertex not in visited:
            visited.add(vertex)
            process(vertex)
            for neighbor in G.adjacent(vertex):
                if neighbor not in visited:
                    stack.append(neighbor)
```

**时间复杂度**: $O(|V| + |E|)$

### 3.2 广度优先搜索 (BFS)

**算法 3.2.1** (BFS)
```python
def BFS(G, start):
    visited = set()
    queue = [start]
    visited.add(start)
    
    while queue:
        vertex = queue.pop(0)
        process(vertex)
        for neighbor in G.adjacent(vertex):
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)
```

**时间复杂度**: $O(|V| + |E|)$

### 3.3 Go语言实现

```go
package traversal

import (
    "container/list"
    "fmt"
)

// DFS 深度优先搜索
func DFS(g Graph, startID int, process func(*Vertex)) {
    visited := make(map[int]bool)
    var dfs func(int)
    
    dfs = func(vertexID int) {
        if visited[vertexID] {
            return
        }
        
        visited[vertexID] = true
        if vertex, exists := g.GetVertex(vertexID); exists {
            process(vertex)
        }
        
        for _, neighbor := range g.Adjacent(vertexID) {
            dfs(neighbor.ID)
        }
    }
    
    dfs(startID)
}

// BFS 广度优先搜索
func BFS(g Graph, startID int, process func(*Vertex)) {
    visited := make(map[int]bool)
    queue := list.New()
    
    visited[startID] = true
    queue.PushBack(startID)
    
    for queue.Len() > 0 {
        element := queue.Front()
        queue.Remove(element)
        vertexID := element.Value.(int)
        
        if vertex, exists := g.GetVertex(vertexID); exists {
            process(vertex)
        }
        
        for _, neighbor := range g.Adjacent(vertexID) {
            if !visited[neighbor.ID] {
                visited[neighbor.ID] = true
                queue.PushBack(neighbor.ID)
            }
        }
    }
}

// TopologicalSort 拓扑排序
func TopologicalSort(g Graph) ([]int, error) {
    if !g.IsDirected() {
        return nil, fmt.Errorf("graph must be directed for topological sort")
    }
    
    inDegree := make(map[int]int)
    for _, vertex := range g.Vertices() {
        inDegree[vertex.ID] = 0
    }
    
    // 计算入度
    for _, edge := range g.Edges() {
        inDegree[edge.To]++
    }
    
    // 使用队列进行拓扑排序
    queue := list.New()
    for vertexID, degree := range inDegree {
        if degree == 0 {
            queue.PushBack(vertexID)
        }
    }
    
    result := make([]int, 0)
    for queue.Len() > 0 {
        element := queue.Front()
        queue.Remove(element)
        vertexID := element.Value.(int)
        
        result = append(result, vertexID)
        
        // 减少邻居的入度
        for _, neighbor := range g.Adjacent(vertexID) {
            inDegree[neighbor.ID]--
            if inDegree[neighbor.ID] == 0 {
                queue.PushBack(neighbor.ID)
            }
        }
    }
    
    if len(result) != len(g.Vertices()) {
        return nil, fmt.Errorf("graph contains cycle")
    }
    
    return result, nil
}

// ConnectedComponents 连通分量
func ConnectedComponents(g Graph) [][]int {
    visited := make(map[int]bool)
    components := make([][]int, 0)
    
    for _, vertex := range g.Vertices() {
        if !visited[vertex.ID] {
            component := make([]int, 0)
            
            // 使用DFS找到连通分量
            var dfs func(int)
            dfs = func(vertexID int) {
                if visited[vertexID] {
                    return
                }
                
                visited[vertexID] = true
                component = append(component, vertexID)
                
                for _, neighbor := range g.Adjacent(vertexID) {
                    dfs(neighbor.ID)
                }
            }
            
            dfs(vertex.ID)
            components = append(components, component)
        }
    }
    
    return components
}
```

---

## 4. 最短路径算法 (Shortest Path)

### 4.1 Dijkstra算法

**算法 4.1.1** (Dijkstra)
```python
def Dijkstra(G, start):
    distances = {v: float('inf') for v in G.vertices()}
    distances[start] = 0
    pq = PriorityQueue([(0, start)])
    
    while pq:
        dist, vertex = pq.pop()
        if dist > distances[vertex]:
            continue
            
        for neighbor, weight in G.adjacent(vertex):
            new_dist = dist + weight
            if new_dist < distances[neighbor]:
                distances[neighbor] = new_dist
                pq.push((new_dist, neighbor))
    
    return distances
```

**时间复杂度**: $O((|V| + |E|) \log |V|)$

### 4.2 Bellman-Ford算法

**算法 4.2.1** (Bellman-Ford)
```python
def BellmanFord(G, start):
    distances = {v: float('inf') for v in G.vertices()}
    distances[start] = 0
    
    for i in range(|V| - 1):
        for edge in G.edges():
            if distances[edge.from] + edge.weight < distances[edge.to]:
                distances[edge.to] = distances[edge.from] + edge.weight
    
    # 检查负环
    for edge in G.edges():
        if distances[edge.from] + edge.weight < distances[edge.to]:
            return None  # 存在负环
    
    return distances
```

**时间复杂度**: $O(|V| \cdot |E|)$

### 4.3 Go语言实现

```go
package shortest_path

import (
    "container/heap"
    "fmt"
    "math"
)

// PriorityQueue 优先队列
type PriorityQueue []*Item

type Item struct {
    vertexID int
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

// Dijkstra Dijkstra算法
func Dijkstra(g Graph, startID int) map[int]float64 {
    distances := make(map[int]float64)
    for _, vertex := range g.Vertices() {
        distances[vertex.ID] = math.Inf(1)
    }
    distances[startID] = 0
    
    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &Item{vertexID: startID, distance: 0})
    
    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        vertexID := item.vertexID
        
        if item.distance > distances[vertexID] {
            continue
        }
        
        for _, neighbor := range g.Adjacent(vertexID) {
            if edge, exists := g.GetEdge(vertexID, neighbor.ID); exists {
                newDist := distances[vertexID] + edge.Weight
                if newDist < distances[neighbor.ID] {
                    distances[neighbor.ID] = newDist
                    heap.Push(pq, &Item{
                        vertexID: neighbor.ID,
                        distance: newDist,
                    })
                }
            }
        }
    }
    
    return distances
}

// BellmanFord Bellman-Ford算法
func BellmanFord(g Graph, startID int) (map[int]float64, error) {
    distances := make(map[int]float64)
    for _, vertex := range g.Vertices() {
        distances[vertex.ID] = math.Inf(1)
    }
    distances[startID] = 0
    
    // 松弛操作
    for i := 0; i < len(g.Vertices())-1; i++ {
        for _, edge := range g.Edges() {
            if distances[edge.From] != math.Inf(1) {
                newDist := distances[edge.From] + edge.Weight
                if newDist < distances[edge.To] {
                    distances[edge.To] = newDist
                }
            }
        }
    }
    
    // 检查负环
    for _, edge := range g.Edges() {
        if distances[edge.From] != math.Inf(1) {
            if distances[edge.From]+edge.Weight < distances[edge.To] {
                return nil, fmt.Errorf("negative cycle detected")
            }
        }
    }
    
    return distances, nil
}

// FloydWarshall Floyd-Warshall算法
func FloydWarshall(g Graph) [][]float64 {
    vertices := g.Vertices()
    n := len(vertices)
    
    // 初始化距离矩阵
    dist := make([][]float64, n)
    for i := range dist {
        dist[i] = make([]float64, n)
        for j := range dist[i] {
            if i == j {
                dist[i][j] = 0
            } else {
                dist[i][j] = math.Inf(1)
            }
        }
    }
    
    // 填充初始边
    for _, edge := range g.Edges() {
        dist[edge.From][edge.To] = edge.Weight
    }
    
    // Floyd-Warshall算法
    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
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

---

## 5. 最小生成树 (Minimum Spanning Tree)

### 5.1 Kruskal算法

**算法 5.1.1** (Kruskal)
```python
def Kruskal(G):
    edges = sorted(G.edges(), key=lambda e: e.weight)
    uf = UnionFind(G.vertices())
    mst = []
    
    for edge in edges:
        if uf.find(edge.from) != uf.find(edge.to):
            uf.union(edge.from, edge.to)
            mst.append(edge)
    
    return mst
```

**时间复杂度**: $O(|E| \log |E|)$

### 5.2 Prim算法

**算法 5.2.1** (Prim)
```python
def Prim(G, start):
    mst = []
    pq = PriorityQueue([(0, start, None)])
    visited = set()
    
    while pq and len(visited) < len(G.vertices()):
        weight, vertex, parent = pq.pop()
        if vertex in visited:
            continue
            
        visited.add(vertex)
        if parent is not None:
            mst.append((parent, vertex, weight))
            
        for neighbor, weight in G.adjacent(vertex):
            if neighbor not in visited:
                pq.push((weight, neighbor, vertex))
    
    return mst
```

**时间复杂度**: $O((|V| + |E|) \log |V|)$

### 5.3 Go语言实现

```go
package mst

import (
    "container/heap"
    "sort"
)

// UnionFind 并查集
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(size int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, size),
        rank:   make([]int, size),
    }
    for i := range uf.parent {
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

// Kruskal Kruskal算法
func Kruskal(g Graph) []*Edge {
    edges := g.Edges()
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })
    
    uf := NewUnionFind(len(g.Vertices()))
    mst := make([]*Edge, 0)
    
    for _, edge := range edges {
        if uf.Find(edge.From) != uf.Find(edge.To) {
            uf.Union(edge.From, edge.To)
            mst = append(mst, edge)
        }
    }
    
    return mst
}

// PrimItem Prim算法优先队列项
type PrimItem struct {
    vertexID int
    weight   float64
    parent   int
    index    int
}

// PrimPriorityQueue Prim算法优先队列
type PrimPriorityQueue []*PrimItem

func (pq PrimPriorityQueue) Len() int { return len(pq) }

func (pq PrimPriorityQueue) Less(i, j int) bool {
    return pq[i].weight < pq[j].weight
}

func (pq PrimPriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PrimPriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*PrimItem)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PrimPriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

// Prim Prim算法
func Prim(g Graph, startID int) []*Edge {
    visited := make(map[int]bool)
    mst := make([]*Edge, 0)
    
    pq := &PrimPriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &PrimItem{vertexID: startID, weight: 0, parent: -1})
    
    for pq.Len() > 0 && len(visited) < len(g.Vertices()) {
        item := heap.Pop(pq).(*PrimItem)
        
        if visited[item.vertexID] {
            continue
        }
        
        visited[item.vertexID] = true
        
        if item.parent != -1 {
            mst = append(mst, &Edge{
                From:   item.parent,
                To:     item.vertexID,
                Weight: item.weight,
            })
        }
        
        for _, neighbor := range g.Adjacent(item.vertexID) {
            if !visited[neighbor.ID] {
                if edge, exists := g.GetEdge(item.vertexID, neighbor.ID); exists {
                    heap.Push(pq, &PrimItem{
                        vertexID: neighbor.ID,
                        weight:   edge.Weight,
                        parent:   item.vertexID,
                    })
                }
            }
        }
    }
    
    return mst
}
```

---

## 6. 网络流 (Network Flow)

### 6.1 最大流问题

**定义 6.1.1** (流网络)
流网络 $G = (V, E)$ 是一个有向图，每条边 $(u, v)$ 有容量 $c(u, v) \geq 0$。

**定义 6.1.2** (流)
流 $f: V \times V \rightarrow \mathbb{R}$ 满足：
1. 容量约束: $0 \leq f(u, v) \leq c(u, v)$
2. 流量守恒: $\sum_{v \in V} f(v, u) = \sum_{v \in V} f(u, v)$

### 6.2 Ford-Fulkerson算法

**算法 6.2.1** (Ford-Fulkerson)
```python
def FordFulkerson(G, source, sink):
    flow = 0
    residual = G.copy()
    
    while path = find_augmenting_path(residual, source, sink):
        bottleneck = min(path.weights())
        flow += bottleneck
        
        for edge in path:
            residual[edge.from][edge.to] -= bottleneck
            residual[edge.to][edge.from] += bottleneck
    
    return flow
```

### 6.3 Go语言实现

```go
package network_flow

import (
    "container/list"
    "math"
)

// FlowNetwork 流网络
type FlowNetwork struct {
    vertices map[int]*Vertex
    edges    map[int]map[int]*FlowEdge
    source   int
    sink     int
}

// FlowEdge 流边
type FlowEdge struct {
    From     int
    To       int
    Capacity float64
    Flow     float64
}

func NewFlowNetwork() *FlowNetwork {
    return &FlowNetwork{
        vertices: make(map[int]*Vertex),
        edges:    make(map[int]map[int]*FlowEdge),
    }
}

func (fn *FlowNetwork) AddVertex(id int, data interface{}) {
    fn.vertices[id] = &Vertex{ID: id, Data: data}
    if fn.edges[id] == nil {
        fn.edges[id] = make(map[int]*FlowEdge)
    }
}

func (fn *FlowNetwork) AddEdge(from, to int, capacity float64) {
    if fn.edges[from] == nil {
        fn.edges[from] = make(map[int]*FlowEdge)
    }
    
    fn.edges[from][to] = &FlowEdge{
        From:     from,
        To:       to,
        Capacity: capacity,
        Flow:     0,
    }
}

// FordFulkerson Ford-Fulkerson算法
func (fn *FlowNetwork) FordFulkerson(source, sink int) float64 {
    fn.source = source
    fn.sink = sink
    
    maxFlow := 0.0
    
    for {
        path, bottleneck := fn.findAugmentingPath()
        if bottleneck == 0 {
            break
        }
        
        maxFlow += bottleneck
        
        // 更新残余网络
        for i := 0; i < len(path)-1; i++ {
            u := path[i]
            v := path[i+1]
            
            // 前向边
            if edge, exists := fn.edges[u][v]; exists {
                edge.Flow += bottleneck
            }
            
            // 反向边
            if fn.edges[v] == nil {
                fn.edges[v] = make(map[int]*FlowEdge)
            }
            if edge, exists := fn.edges[v][u]; exists {
                edge.Flow -= bottleneck
            } else {
                fn.edges[v][u] = &FlowEdge{
                    From:     v,
                    To:       u,
                    Capacity: bottleneck,
                    Flow:     -bottleneck,
                }
            }
        }
    }
    
    return maxFlow
}

// findAugmentingPath 寻找增广路径
func (fn *FlowNetwork) findAugmentingPath() ([]int, float64) {
    visited := make(map[int]bool)
    parent := make(map[int]int)
    bottleneck := make(map[int]float64)
    
    queue := list.New()
    queue.PushBack(fn.source)
    visited[fn.source] = true
    bottleneck[fn.source] = math.Inf(1)
    
    for queue.Len() > 0 {
        element := queue.Front()
        queue.Remove(element)
        u := element.Value.(int)
        
        for v, edge := range fn.edges[u] {
            if !visited[v] && edge.Flow < edge.Capacity {
                visited[v] = true
                parent[v] = u
                bottleneck[v] = math.Min(bottleneck[u], edge.Capacity-edge.Flow)
                
                if v == fn.sink {
                    // 重建路径
                    path := make([]int, 0)
                    current := v
                    for current != fn.source {
                        path = append([]int{current}, path...)
                        current = parent[current]
                    }
                    path = append([]int{fn.source}, path...)
                    
                    return path, bottleneck[v]
                }
                
                queue.PushBack(v)
            }
        }
        
        // 检查反向边
        for v, edge := range fn.edges[u] {
            if !visited[v] && edge.Flow > 0 {
                visited[v] = true
                parent[v] = u
                bottleneck[v] = math.Min(bottleneck[u], edge.Flow)
                
                if v == fn.sink {
                    // 重建路径
                    path := make([]int, 0)
                    current := v
                    for current != fn.source {
                        path = append([]int{current}, path...)
                        current = parent[current]
                    }
                    path = append([]int{fn.source}, path...)
                    
                    return path, bottleneck[v]
                }
                
                queue.PushBack(v)
            }
        }
    }
    
    return nil, 0
}
```

---

## 总结

本文档介绍了图论的基本概念和算法：

1. **基本概念** - 图的定义、度数、握手定理
2. **图的表示** - 邻接表、邻接矩阵、关联矩阵
3. **图遍历** - DFS、BFS、拓扑排序、连通分量
4. **最短路径** - Dijkstra、Bellman-Ford、Floyd-Warshall
5. **最小生成树** - Kruskal、Prim算法
6. **网络流** - 最大流、Ford-Fulkerson算法

这些算法在计算机科学中有广泛应用，如路由算法、网络优化、任务调度等。

---

**相关链接**:
- [01-集合论 (Set Theory)](01-Set-Theory.md)
- [02-逻辑学 (Logic)](02-Logic.md)
- [04-概率论 (Probability Theory)](04-Probability-Theory.md)
- [02-形式化验证 (Formal Verification)](../02-Formal-Verification/README.md) 