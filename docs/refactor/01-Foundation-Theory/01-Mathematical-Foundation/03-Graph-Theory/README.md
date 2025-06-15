# 03-图论 (Graph Theory)

## 目录

- [03-图论 (Graph Theory)](#03-图论-graph-theory)
  - [目录](#目录)
  - [1. 基础概念](#1-基础概念)
    - [1.1 图的定义](#11-图的定义)
    - [1.2 图的类型](#12-图的类型)
      - [1.2.1 按边的性质分类](#121-按边的性质分类)
      - [1.2.2 按连接性分类](#122-按连接性分类)
    - [1.3 基本术语](#13-基本术语)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 图的形式化模型](#21-图的形式化模型)
    - [2.2 图的性质](#22-图的性质)
    - [2.3 图的操作](#23-图的操作)
  - [3. 图算法](#3-图算法)
    - [3.1 遍历算法](#31-遍历算法)
      - [3.1.1 深度优先搜索 (DFS)](#311-深度优先搜索-dfs)
      - [3.1.2 广度优先搜索 (BFS)](#312-广度优先搜索-bfs)
    - [3.2 最短路径算法](#32-最短路径算法)
      - [3.2.1 Dijkstra算法](#321-dijkstra算法)
    - [3.3 最小生成树算法](#33-最小生成树算法)
      - [3.3.1 Kruskal算法](#331-kruskal算法)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 图的数据结构](#41-图的数据结构)
    - [4.2 算法实现](#42-算法实现)
      - [4.2.1 DFS实现](#421-dfs实现)
      - [4.2.2 BFS实现](#422-bfs实现)
      - [4.2.3 Dijkstra算法实现](#423-dijkstra算法实现)
    - [4.3 性能优化](#43-性能优化)
      - [4.3.1 内存优化](#431-内存优化)
      - [4.3.2 并发优化](#432-并发优化)
  - [5. 应用场景](#5-应用场景)
    - [5.1 网络路由](#51-网络路由)
    - [5.2 社交网络分析](#52-社交网络分析)
    - [5.3 生物信息学](#53-生物信息学)
    - [5.4 编译器优化](#54-编译器优化)
  - [6. 总结](#6-总结)

## 1. 基础概念

### 1.1 图的定义

**定义 1.1** (图): 图 $G = (V, E)$ 是一个有序对，其中：

- $V$ 是顶点集 (vertex set)，$V \neq \emptyset$
- $E$ 是边集 (edge set)，$E \subseteq V \times V$

**定义 1.2** (有向图): 有向图 $G = (V, E)$ 中，边集 $E$ 的元素是有序对 $(u, v)$，表示从顶点 $u$ 到顶点 $v$ 的有向边。

**定义 1.3** (无向图): 无向图 $G = (V, E)$ 中，边集 $E$ 的元素是无序对 $\{u, v\}$，表示顶点 $u$ 和 $v$ 之间的无向边。

### 1.2 图的类型

#### 1.2.1 按边的性质分类

- **简单图**: 没有自环和重边的图
- **多重图**: 允许重边的图
- **伪图**: 允许自环和重边的图

#### 1.2.2 按连接性分类

- **连通图**: 任意两个顶点之间都存在路径
- **强连通图**: 有向图中任意两个顶点之间都存在有向路径
- **完全图**: 任意两个顶点之间都有边相连

### 1.3 基本术语

**定义 1.4** (度数): 顶点 $v$ 的度数 $deg(v)$ 是与 $v$ 相邻的边的数量。

**定义 1.5** (路径): 路径是顶点序列 $v_0, v_1, \ldots, v_k$，其中 $(v_i, v_{i+1}) \in E$ 对所有 $0 \leq i < k$ 成立。

**定义 1.6** (环): 环是长度至少为3的路径，其中 $v_0 = v_k$。

## 2. 形式化定义

### 2.1 图的形式化模型

**公理 2.1** (图的基本公理):

1. 顶点集非空: $V \neq \emptyset$
2. 边集是顶点对的子集: $E \subseteq V \times V$
3. 自环允许性: $(v, v) \in E$ 在某些图中允许

**定理 2.1** (握手定理): 对于任何图 $G = (V, E)$，
$$\sum_{v \in V} deg(v) = 2|E|$$

**证明**:

- 每条边贡献给两个顶点的度数
- 因此所有顶点的度数之和等于边数的两倍

### 2.2 图的性质

**定义 2.1** (图的同构): 两个图 $G_1 = (V_1, E_1)$ 和 $G_2 = (V_2, E_2)$ 是同构的，如果存在双射 $f: V_1 \rightarrow V_2$，使得 $(u, v) \in E_1$ 当且仅当 $(f(u), f(v)) \in E_2$。

**定理 2.2** (欧拉回路定理): 连通图 $G$ 存在欧拉回路的充要条件是所有顶点的度数都是偶数。

**证明**:

- 必要性: 欧拉回路中每个顶点进入和离开的次数相等
- 充分性: 通过构造算法证明

### 2.3 图的操作

**定义 2.2** (图的并): $G_1 \cup G_2 = (V_1 \cup V_2, E_1 \cup E_2)$

**定义 2.3** (图的交): $G_1 \cap G_2 = (V_1 \cap V_2, E_1 \cap E_2)$

**定义 2.4** (图的补): $\overline{G} = (V, \overline{E})$，其中 $\overline{E} = V \times V - E$

## 3. 图算法

### 3.1 遍历算法

#### 3.1.1 深度优先搜索 (DFS)

**算法 3.1** (DFS):

```text
DFS(G, v):
    visited[v] = true
    for each neighbor u of v:
        if not visited[u]:
            DFS(G, u)
```

**时间复杂度**: $O(|V| + |E|)$

#### 3.1.2 广度优先搜索 (BFS)

**算法 3.2** (BFS):

```text
BFS(G, s):
    queue = [s]
    visited[s] = true
    while queue not empty:
        v = queue.dequeue()
        for each neighbor u of v:
            if not visited[u]:
                visited[u] = true
                queue.enqueue(u)
```

**时间复杂度**: $O(|V| + |E|)$

### 3.2 最短路径算法

#### 3.2.1 Dijkstra算法

**算法 3.3** (Dijkstra):

```text
Dijkstra(G, s):
    dist[s] = 0
    for all v ≠ s: dist[v] = ∞
    Q = V
    while Q not empty:
        u = extract-min(Q)
        for each neighbor v of u:
            if dist[v] > dist[u] + weight(u,v):
                dist[v] = dist[u] + weight(u,v)
```

**时间复杂度**: $O(|V|^2)$ (使用数组) 或 $O((|V| + |E|) \log |V|)$ (使用堆)

### 3.3 最小生成树算法

#### 3.3.1 Kruskal算法

**算法 3.4** (Kruskal):

```text
Kruskal(G):
    sort edges by weight
    for each edge (u,v) in sorted order:
        if find(u) ≠ find(v):
            union(u,v)
            add (u,v) to MST
```

**时间复杂度**: $O(|E| \log |E|)$

## 4. Go语言实现

### 4.1 图的数据结构

```go
// 顶点类型
type Vertex int

// 边结构
type Edge struct {
    From   Vertex
    To     Vertex
    Weight float64
}

// 图接口
type Graph interface {
    AddVertex(v Vertex)
    RemoveVertex(v Vertex)
    AddEdge(e Edge)
    RemoveEdge(e Edge)
    GetNeighbors(v Vertex) []Vertex
    GetEdgeWeight(from, to Vertex) float64
    GetVertices() []Vertex
    GetEdges() []Edge
}

// 邻接表实现
type AdjacencyListGraph struct {
    vertices map[Vertex]bool
    edges    map[Vertex]map[Vertex]float64
}

// 构造函数
func NewAdjacencyListGraph() *AdjacencyListGraph {
    return &AdjacencyListGraph{
        vertices: make(map[Vertex]bool),
        edges:    make(map[Vertex]map[Vertex]float64),
    }
}

// 添加顶点
func (g *AdjacencyListGraph) AddVertex(v Vertex) {
    g.vertices[v] = true
    if g.edges[v] == nil {
        g.edges[v] = make(map[Vertex]float64)
    }
}

// 添加边
func (g *AdjacencyListGraph) AddEdge(e Edge) {
    g.AddVertex(e.From)
    g.AddVertex(e.To)
    g.edges[e.From][e.To] = e.Weight
}

// 获取邻居
func (g *AdjacencyListGraph) GetNeighbors(v Vertex) []Vertex {
    neighbors := make([]Vertex, 0)
    for neighbor := range g.edges[v] {
        neighbors = append(neighbors, neighbor)
    }
    return neighbors
}
```

### 4.2 算法实现

#### 4.2.1 DFS实现

```go
// DFS实现
func (g *AdjacencyListGraph) DFS(start Vertex) []Vertex {
    visited := make(map[Vertex]bool)
    result := make([]Vertex, 0)
    
    var dfs func(Vertex)
    dfs = func(v Vertex) {
        visited[v] = true
        result = append(result, v)
        
        for _, neighbor := range g.GetNeighbors(v) {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
    }
    
    dfs(start)
    return result
}

// 并发安全的DFS
func (g *AdjacencyListGraph) DFSConcurrent(start Vertex) []Vertex {
    visited := make(map[Vertex]bool)
    var mu sync.RWMutex
    result := make([]Vertex, 0)
    
    var dfs func(Vertex)
    dfs = func(v Vertex) {
        mu.Lock()
        if visited[v] {
            mu.Unlock()
            return
        }
        visited[v] = true
        result = append(result, v)
        mu.Unlock()
        
        neighbors := g.GetNeighbors(v)
        var wg sync.WaitGroup
        for _, neighbor := range neighbors {
            wg.Add(1)
            go func(n Vertex) {
                defer wg.Done()
                dfs(n)
            }(neighbor)
        }
        wg.Wait()
    }
    
    dfs(start)
    return result
}
```

#### 4.2.2 BFS实现

```go
// BFS实现
func (g *AdjacencyListGraph) BFS(start Vertex) []Vertex {
    visited := make(map[Vertex]bool)
    queue := []Vertex{start}
    result := make([]Vertex, 0)
    
    visited[start] = true
    
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)
        
        for _, neighbor := range g.GetNeighbors(v) {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    
    return result
}

// 带层级的BFS
func (g *AdjacencyListGraph) BFSWithLevels(start Vertex) map[int][]Vertex {
    visited := make(map[Vertex]bool)
    queue := []Vertex{start}
    levels := make(map[int][]Vertex)
    level := 0
    
    visited[start] = true
    levels[level] = []Vertex{start}
    
    for len(queue) > 0 {
        levelSize := len(queue)
        level++
        
        for i := 0; i < levelSize; i++ {
            v := queue[0]
            queue = queue[1:]
            
            for _, neighbor := range g.GetNeighbors(v) {
                if !visited[neighbor] {
                    visited[neighbor] = true
                    queue = append(queue, neighbor)
                    levels[level] = append(levels[level], neighbor)
                }
            }
        }
    }
    
    return levels
}
```

#### 4.2.3 Dijkstra算法实现

```go
// 距离结构
type Distance struct {
    Vertex Vertex
    Dist   float64
}

// 优先队列
type PriorityQueue []Distance

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].Dist < pq[j].Dist
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    *pq = append(*pq, x.(Distance))
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    x := old[n-1]
    *pq = old[0 : n-1]
    return x
}

// Dijkstra算法实现
func (g *AdjacencyListGraph) Dijkstra(start Vertex) map[Vertex]float64 {
    dist := make(map[Vertex]float64)
    pq := &PriorityQueue{}
    heap.Init(pq)
    
    // 初始化距离
    for v := range g.vertices {
        if v == start {
            dist[v] = 0
        } else {
            dist[v] = math.Inf(1)
        }
    }
    
    heap.Push(pq, Distance{Vertex: start, Dist: 0})
    
    for pq.Len() > 0 {
        current := heap.Pop(pq).(Distance)
        
        if current.Dist > dist[current.Vertex] {
            continue
        }
        
        for neighbor, weight := range g.edges[current.Vertex] {
            newDist := dist[current.Vertex] + weight
            
            if newDist < dist[neighbor] {
                dist[neighbor] = newDist
                heap.Push(pq, Distance{Vertex: neighbor, Dist: newDist})
            }
        }
    }
    
    return dist
}
```

### 4.3 性能优化

#### 4.3.1 内存优化

```go
// 紧凑的图表示
type CompactGraph struct {
    vertices []Vertex
    edges    []Edge
    index    map[Vertex]int
}

func NewCompactGraph() *CompactGraph {
    return &CompactGraph{
        vertices: make([]Vertex, 0),
        edges:    make([]Edge, 0),
        index:    make(map[Vertex]int),
    }
}

// 使用位图表示访问状态
type BitSet struct {
    data []uint64
    size int
}

func NewBitSet(size int) *BitSet {
    return &BitSet{
        data: make([]uint64, (size+63)/64),
        size: size,
    }
}

func (bs *BitSet) Set(index int) {
    if index >= 0 && index < bs.size {
        bs.data[index/64] |= 1 << (index % 64)
    }
}

func (bs *BitSet) Get(index int) bool {
    if index >= 0 && index < bs.size {
        return (bs.data[index/64] & (1 << (index % 64))) != 0
    }
    return false
}
```

#### 4.3.2 并发优化

```go
// 并发图遍历
func (g *AdjacencyListGraph) ParallelBFS(start Vertex, numWorkers int) []Vertex {
    visited := make(map[Vertex]bool)
    var mu sync.RWMutex
    queue := make(chan Vertex, len(g.vertices))
    result := make([]Vertex, 0)
    var wg sync.WaitGroup
    
    // 启动工作协程
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for v := range queue {
                mu.Lock()
                if !visited[v] {
                    visited[v] = true
                    result = append(result, v)
                    
                    // 添加邻居到队列
                    for _, neighbor := range g.GetNeighbors(v) {
                        if !visited[neighbor] {
                            queue <- neighbor
                        }
                    }
                }
                mu.Unlock()
            }
        }()
    }
    
    // 启动遍历
    queue <- start
    close(queue)
    wg.Wait()
    
    return result
}
```

## 5. 应用场景

### 5.1 网络路由

图论在网络路由中的应用：

- 最短路径算法用于路由选择
- 最小生成树用于网络设计
- 连通性分析用于故障检测

### 5.2 社交网络分析

- 社区检测算法
- 影响力传播模型
- 推荐系统

### 5.3 生物信息学

- 蛋白质相互作用网络
- 基因调控网络
- 代谢网络分析

### 5.4 编译器优化

- 控制流图分析
- 数据流分析
- 寄存器分配

## 6. 总结

图论作为计算机科学的基础理论，在算法设计、数据结构、网络分析等领域有广泛应用。通过Go语言的实现，我们可以看到：

1. **理论到实践的转换**: 形式化定义可以转化为高效的代码实现
2. **性能优化**: 通过并发编程和内存优化提升算法效率
3. **实际应用**: 图论算法在多个领域都有重要应用

图论的研究不仅有助于理解复杂系统的结构，也为解决实际问题提供了强大的工具。

---

**相关链接**:

- [01-集合论](../01-Set-Theory/README.md)
- [02-逻辑学](../02-Logic/README.md)
- [04-概率论](../04-Probability-Theory/README.md)
- [03-设计模式](../../03-Design-Patterns/README.md)
- [02-软件架构](../../02-Software-Architecture/README.md)
