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
  - [3. 图算法](#3-图算法)
    - [3.1 遍历算法](#31-遍历算法)
    - [3.2 最短路径算法](#32-最短路径算法)
    - [3.3 最小生成树算法](#33-最小生成树算法)
  - [4. 形式化定义](#4-形式化定义)
    - [4.1 图的数学定义](#41-图的数学定义)
    - [4.2 图的性质定理](#42-图的性质定理)
    - [4.3 算法正确性证明](#43-算法正确性证明)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 图的数据结构](#51-图的数据结构)
    - [5.2 基础算法实现](#52-基础算法实现)
    - [5.3 高级算法实现](#53-高级算法实现)
  - [6. 应用实例](#6-应用实例)
    - [6.1 网络拓扑分析](#61-网络拓扑分析)
    - [6.2 社交网络分析](#62-社交网络分析)
    - [6.3 路由算法](#63-路由算法)
  - [总结](#总结)

## 1. 基础定义

### 1.1 图的基本概念

**定义 1.1 (图)**
图 $G = (V, E)$ 是一个有序对，其中：

- $V$ 是顶点集（vertex set），$V \neq \emptyset$
- $E$ 是边集（edge set），$E \subseteq V \times V$

**定义 1.2 (有向图)**
有向图 $G = (V, E)$ 中，边集 $E$ 的元素是有序对 $(u, v)$，表示从顶点 $u$ 到顶点 $v$ 的有向边。

**定义 1.3 (无向图)**
无向图 $G = (V, E)$ 中，边集 $E$ 的元素是无序对 $\{u, v\}$，表示顶点 $u$ 和 $v$ 之间的无向边。

### 1.2 图的类型

**定义 1.4 (完全图)**
完全图 $K_n$ 是包含 $n$ 个顶点的图，其中任意两个顶点之间都有一条边。

**定义 1.5 (二分图)**
二分图 $G = (V_1 \cup V_2, E)$ 中，顶点集可以划分为两个不相交的子集 $V_1$ 和 $V_2$，使得每条边的两个端点分别属于不同的子集。

**定义 1.6 (加权图)**
加权图 $G = (V, E, w)$ 中，每条边 $e \in E$ 都有一个权重 $w(e) \in \mathbb{R}$。

### 1.3 图的表示

**邻接矩阵表示**：

对于图 $G = (V, E)$，邻接矩阵 $A$ 是一个 $|V| \times |V|$ 的矩阵，其中：

```latex
$$A[i][j] = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$$
```

**邻接表表示**：
每个顶点维护一个包含其邻居顶点的列表。

## 2. 图的性质

### 2.1 连通性

**定义 2.1 (连通图)**
无向图 $G$ 是连通的，当且仅当对于任意两个顶点 $u, v \in V$，存在从 $u$ 到 $v$ 的路径。

**定义 2.2 (强连通图)**
有向图 $G$ 是强连通的，当且仅当对于任意两个顶点 $u, v \in V$，存在从 $u$ 到 $v$ 的有向路径。

**定理 2.1 (连通性判定)**
图 $G = (V, E)$ 是连通的，当且仅当从任意顶点开始的深度优先搜索或广度优先搜索能够访问所有顶点。

**证明**：

- **必要性**：如果图是连通的，那么从任意顶点都存在到其他所有顶点的路径，因此DFS或BFS能够访问所有顶点。
- **充分性**：如果从某个顶点开始的遍历能够访问所有顶点，那么该顶点到所有其他顶点都存在路径，因此图是连通的。

### 2.2 路径和回路

**定义 2.3 (路径)**
路径是顶点序列 $v_0, v_1, \ldots, v_k$，其中 $(v_i, v_{i+1}) \in E$ 对所有 $0 \leq i < k$ 成立。

**定义 2.4 (简单路径)**
简单路径是不包含重复顶点的路径。

**定义 2.5 (回路)**
回路是起点和终点相同的路径。

**定理 2.2 (路径存在性)**
在连通图中，任意两个顶点之间都存在简单路径。

**证明**：
使用构造性证明。从起点开始，每次选择未访问的邻居顶点，直到到达终点。由于图是连通的，这样的路径总是存在。

### 2.3 树和森林

**定义 2.6 (树)**
树是连通的无环无向图。

**定义 2.7 (森林)**
森林是多个树的并集。

**定理 2.3 (树的性质)**
对于包含 $n$ 个顶点的树 $T$：

1. $T$ 有 $n-1$ 条边
2. $T$ 是连通的
3. $T$ 是无环的
4. 任意两个顶点之间有唯一路径

**证明**：
使用数学归纳法。对于 $n=1$，结论显然成立。假设对于 $n=k$ 成立，考虑 $n=k+1$ 的树。删除任意一条边，得到两个连通分量，每个都是树，根据归纳假设，总边数为 $(k_1-1) + (k_2-1) + 1 = k$，其中 $k_1 + k_2 = k+1$。

## 3. 图算法

### 3.1 遍历算法

**深度优先搜索 (DFS)**：

```go
func DFS(graph map[int][]int, start int) []int {
    visited := make(map[int]bool)
    result := []int{}

    var dfs func(int)
    dfs = func(node int) {
        if visited[node] {
            return
        }
        visited[node] = true
        result = append(result, node)

        for _, neighbor := range graph[node] {
            dfs(neighbor)
        }
    }

    dfs(start)
    return result
}
```

**广度优先搜索 (BFS)**：

```go
func BFS(graph map[int][]int, start int) []int {
    visited := make(map[int]bool)
    result := []int{}
    queue := []int{start}
    visited[start] = true

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)

        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }

    return result
}
```

### 3.2 最短路径算法

**Dijkstra算法**：

```go
type Edge struct {
    To     int
    Weight int
}

func Dijkstra(graph map[int][]Edge, start int) map[int]int {
    distances := make(map[int]int)
    for node := range graph {
        distances[node] = math.MaxInt32
    }
    distances[start] = 0

    pq := &PriorityQueue{}
    heap.Push(pq, &Item{node: start, distance: 0})

    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        node := item.node
        dist := item.distance

        if dist > distances[node] {
            continue
        }

        for _, edge := range graph[node] {
            newDist := dist + edge.Weight
            if newDist < distances[edge.To] {
                distances[edge.To] = newDist
                heap.Push(pq, &Item{node: edge.To, distance: newDist})
            }
        }
    }

    return distances
}
```

**定理 3.1 (Dijkstra算法正确性)**
Dijkstra算法能够正确计算从起点到所有其他顶点的最短路径。

**证明**：
使用反证法。假设存在某个顶点 $v$，算法计算的距离 $d[v]$ 不是最短距离。设 $d^*[v]$ 是真正的最短距离，则 $d[v] > d^*[v]$。考虑最短路径上的最后一个顶点 $u$，根据算法的贪心性质，$d[u] = d^*[u]$，这与 $d[v] > d^*[v]$ 矛盾。

### 3.3 最小生成树算法

**Kruskal算法**：

```go
type Edge struct {
    From   int
    To     int
    Weight int
}

func Kruskal(edges []Edge, n int) []Edge {
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })

    uf := NewUnionFind(n)
    result := []Edge{}

    for _, edge := range edges {
        if uf.Find(edge.From) != uf.Find(edge.To) {
            uf.Union(edge.From, edge.To)
            result = append(result, edge)
        }
    }

    return result
}
```

**定理 3.2 (Kruskal算法正确性)**
Kruskal算法能够找到图的最小生成树。

**证明**：
使用归纳法。假设算法已经选择了 $k$ 条边，形成森林 $F$。考虑下一条边 $e$，如果 $e$ 连接 $F$ 中的两个不同连通分量，那么 $e$ 属于某个最小生成树（根据割性质），因此算法选择 $e$ 是正确的。

## 4. 形式化定义

### 4.1 图的数学定义

**定义 4.1 (图的同构)**
两个图 $G_1 = (V_1, E_1)$ 和 $G_2 = (V_2, E_2)$ 是同构的，当且仅当存在双射 $f: V_1 \rightarrow V_2$，使得 $(u, v) \in E_1$ 当且仅当 $(f(u), f(v)) \in E_2$。

**定义 4.2 (图的补图)**
图 $G = (V, E)$ 的补图是 $G^c = (V, E^c)$，其中 $E^c = V \times V \setminus E$。

**定义 4.3 (图的度)**
顶点 $v$ 的度 $deg(v)$ 是与 $v$ 相邻的顶点数量。

**定理 4.1 (握手定理)**
对于图 $G = (V, E)$，$\sum_{v \in V} deg(v) = 2|E|$。

**证明**：
每条边贡献给两个顶点的度，因此所有顶点的度之和等于边数的两倍。

### 4.2 图的性质定理

**定理 4.2 (欧拉回路存在性)**
连通图 $G$ 存在欧拉回路，当且仅当所有顶点的度都是偶数。

**证明**：

- **必要性**：欧拉回路经过每条边恰好一次，因此每个顶点被进入和离开的次数相等，度为偶数。
- **充分性**：使用构造性证明，从任意顶点开始，每次选择未使用的边，直到无法继续，然后回溯并寻找其他路径。

**定理 4.3 (哈密顿回路存在性)**
对于 $n \geq 3$ 的完全图 $K_n$，存在哈密顿回路。

**证明**：
使用构造性证明。对于顶点 $v_1, v_2, \ldots, v_n$，路径 $v_1 \rightarrow v_2 \rightarrow \cdots \rightarrow v_n \rightarrow v_1$ 是一个哈密顿回路。

### 4.3 算法正确性证明

**定理 4.4 (DFS正确性)**
DFS能够访问从起点可达的所有顶点。

**证明**：
使用归纳法。对于距离起点为 $k$ 的顶点，DFS会在访问距离为 $k-1$ 的顶点后访问它们。

**定理 4.5 (BFS正确性)**
BFS能够按照距离起点的递增顺序访问所有可达顶点。

**证明**：
BFS使用队列，确保距离较小的顶点先被访问。

## 5. Go语言实现

### 5.1 图的数据结构

```go
// 图的基本接口
type Graph interface {
    AddVertex(v int)
    AddEdge(from, to int)
    RemoveVertex(v int)
    RemoveEdge(from, to int)
    GetNeighbors(v int) []int
    GetVertices() []int
    GetEdges() [][2]int
    IsConnected() bool
    HasCycle() bool
}

// 邻接表实现的图
type AdjacencyListGraph struct {
    vertices map[int][]int
    directed bool
}

func NewAdjacencyListGraph(directed bool) *AdjacencyListGraph {
    return &AdjacencyListGraph{
        vertices: make(map[int][]int),
        directed: directed,
    }
}

func (g *AdjacencyListGraph) AddVertex(v int) {
    if _, exists := g.vertices[v]; !exists {
        g.vertices[v] = []int{}
    }
}

func (g *AdjacencyListGraph) AddEdge(from, to int) {
    g.AddVertex(from)
    g.AddVertex(to)
    g.vertices[from] = append(g.vertices[from], to)
    if !g.directed {
        g.vertices[to] = append(g.vertices[to], from)
    }
}

func (g *AdjacencyListGraph) GetNeighbors(v int) []int {
    if neighbors, exists := g.vertices[v]; exists {
        return neighbors
    }
    return []int{}
}

func (g *AdjacencyListGraph) GetVertices() []int {
    vertices := make([]int, 0, len(g.vertices))
    for v := range g.vertices {
        vertices = append(vertices, v)
    }
    return vertices
}

func (g *AdjacencyListGraph) IsConnected() bool {
    if len(g.vertices) == 0 {
        return true
    }

    visited := make(map[int]bool)
    var dfs func(int)
    dfs = func(v int) {
        visited[v] = true
        for _, neighbor := range g.GetNeighbors(v) {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
    }

    // 从任意顶点开始DFS
    for v := range g.vertices {
        dfs(v)
        break
    }

    return len(visited) == len(g.vertices)
}

func (g *AdjacencyListGraph) HasCycle() bool {
    if g.directed {
        return g.hasDirectedCycle()
    }
    return g.hasUndirectedCycle()
}

func (g *AdjacencyListGraph) hasDirectedCycle() bool {
    visited := make(map[int]bool)
    recStack := make(map[int]bool)

    var dfs func(int) bool
    dfs = func(v int) bool {
        visited[v] = true
        recStack[v] = true

        for _, neighbor := range g.GetNeighbors(v) {
            if !visited[neighbor] {
                if dfs(neighbor) {
                    return true
                }
            } else if recStack[neighbor] {
                return true
            }
        }

        recStack[v] = false
        return false
    }

    for v := range g.vertices {
        if !visited[v] {
            if dfs(v) {
                return true
            }
        }
    }
    return false
}

func (g *AdjacencyListGraph) hasUndirectedCycle() bool {
    visited := make(map[int]bool)

    var dfs func(int, int) bool
    dfs = func(v, parent int) bool {
        visited[v] = true

        for _, neighbor := range g.GetNeighbors(v) {
            if !visited[neighbor] {
                if dfs(neighbor, v) {
                    return true
                }
            } else if neighbor != parent {
                return true
            }
        }
        return false
    }

    for v := range g.vertices {
        if !visited[v] {
            if dfs(v, -1) {
                return true
            }
        }
    }
    return false
}
```

### 5.2 基础算法实现

```go
// 拓扑排序
func TopologicalSort(graph map[int][]int) ([]int, error) {
    inDegree := make(map[int]int)
    for v := range graph {
        inDegree[v] = 0
    }

    // 计算入度
    for _, neighbors := range graph {
        for _, neighbor := range neighbors {
            inDegree[neighbor]++
        }
    }

    // 使用队列进行拓扑排序
    queue := []int{}
    for v, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, v)
        }
    }

    result := []int{}
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)

        for _, neighbor := range graph[v] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != len(graph) {
        return nil, fmt.Errorf("graph contains cycle")
    }

    return result, nil
}

// 强连通分量 (Tarjan算法)
func StronglyConnectedComponents(graph map[int][]int) [][]int {
    index := 0
    indices := make(map[int]int)
    lowLinks := make(map[int]int)
    onStack := make(map[int]bool)
    stack := []int{}
    result := [][]int{}

    var strongConnect func(int)
    strongConnect = func(v int) {
        indices[v] = index
        lowLinks[v] = index
        index++
        stack = append(stack, v)
        onStack[v] = true

        for _, neighbor := range graph[v] {
            if _, exists := indices[neighbor]; !exists {
                strongConnect(neighbor)
                lowLinks[v] = min(lowLinks[v], lowLinks[neighbor])
            } else if onStack[neighbor] {
                lowLinks[v] = min(lowLinks[v], indices[neighbor])
            }
        }

        if lowLinks[v] == indices[v] {
            component := []int{}
            for {
                w := stack[len(stack)-1]
                stack = stack[:len(stack)-1]
                onStack[w] = false
                component = append(component, w)
                if w == v {
                    break
                }
            }
            result = append(result, component)
        }
    }

    for v := range graph {
        if _, exists := indices[v]; !exists {
            strongConnect(v)
        }
    }

    return result
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

### 5.3 高级算法实现

```go
// Floyd-Warshall算法 (所有顶点对最短路径)
func FloydWarshall(graph [][]int) [][]int {
    n := len(graph)
    dist := make([][]int, n)
    for i := range dist {
        dist[i] = make([]int, n)
        copy(dist[i], graph[i])
    }

    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
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

// 最大流算法 (Ford-Fulkerson)
type FlowEdge struct {
    from   int
    to     int
    flow   int
    capacity int
}

func MaxFlow(graph map[int][]FlowEdge, source, sink int) int {
    // 构建残量网络
    residual := make(map[int][]FlowEdge)
    for from, edges := range graph {
        for _, edge := range edges {
            // 正向边
            residual[from] = append(residual[from], FlowEdge{
                from:     from,
                to:       edge.to,
                flow:     0,
                capacity: edge.capacity,
            })
            // 反向边
            residual[edge.to] = append(residual[edge.to], FlowEdge{
                from:     edge.to,
                to:       from,
                flow:     0,
                capacity: 0,
            })
        }
    }

    maxFlow := 0
    for {
        // 使用BFS寻找增广路径
        path := findAugmentingPath(residual, source, sink)
        if len(path) == 0 {
            break
        }

        // 计算路径上的最小容量
        minCapacity := math.MaxInt32
        for i := 0; i < len(path)-1; i++ {
            edge := findEdge(residual, path[i], path[i+1])
            if edge.capacity-edge.flow < minCapacity {
                minCapacity = edge.capacity - edge.flow
            }
        }

        // 更新残量网络
        for i := 0; i < len(path)-1; i++ {
            updateFlow(residual, path[i], path[i+1], minCapacity)
        }

        maxFlow += minCapacity
    }

    return maxFlow
}

func findAugmentingPath(residual map[int][]FlowEdge, source, sink int) []int {
    parent := make(map[int]int)
    visited := make(map[int]bool)
    queue := []int{source}
    visited[source] = true

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        if current == sink {
            // 重建路径
            path := []int{}
            for current != source {
                path = append([]int{current}, path...)
                current = parent[current]
            }
            path = append([]int{source}, path...)
            return path
        }

        for _, edge := range residual[current] {
            if !visited[edge.to] && edge.capacity > edge.flow {
                visited[edge.to] = true
                parent[edge.to] = current
                queue = append(queue, edge.to)
            }
        }
    }

    return []int{}
}
```

## 6. 应用实例

### 6.1 网络拓扑分析

```go
// 网络拓扑分析器
type NetworkTopologyAnalyzer struct {
    graph *AdjacencyListGraph
}

func NewNetworkTopologyAnalyzer() *NetworkTopologyAnalyzer {
    return &NetworkTopologyAnalyzer{
        graph: NewAdjacencyListGraph(false),
    }
}

func (nta *NetworkTopologyAnalyzer) AddConnection(node1, node2 int) {
    nta.graph.AddEdge(node1, node2)
}

func (nta *NetworkTopologyAnalyzer) AnalyzeConnectivity() map[string]interface{} {
    result := make(map[string]interface{})

    // 检查连通性
    result["is_connected"] = nta.graph.IsConnected()

    // 计算连通分量
    result["connected_components"] = nta.findConnectedComponents()

    // 计算网络直径
    result["diameter"] = nta.calculateDiameter()

    // 计算平均度
    result["average_degree"] = nta.calculateAverageDegree()

    return result
}

func (nta *NetworkTopologyAnalyzer) findConnectedComponents() [][]int {
    visited := make(map[int]bool)
    components := [][]int{}

    var dfs func(int, *[]int)
    dfs = func(v int, component *[]int) {
        visited[v] = true
        *component = append(*component, v)

        for _, neighbor := range nta.graph.GetNeighbors(v) {
            if !visited[neighbor] {
                dfs(neighbor, component)
            }
        }
    }

    for v := range nta.graph.vertices {
        if !visited[v] {
            component := []int{}
            dfs(v, &component)
            components = append(components, component)
        }
    }

    return components
}

func (nta *NetworkTopologyAnalyzer) calculateDiameter() int {
    vertices := nta.graph.GetVertices()
    maxDiameter := 0

    for _, start := range vertices {
        distances := nta.shortestPaths(start)
        for _, dist := range distances {
            if dist > maxDiameter && dist != math.MaxInt32 {
                maxDiameter = dist
            }
        }
    }

    return maxDiameter
}

func (nta *NetworkTopologyAnalyzer) calculateAverageDegree() float64 {
    totalDegree := 0
    vertexCount := len(nta.graph.vertices)

    for v := range nta.graph.vertices {
        totalDegree += len(nta.graph.GetNeighbors(v))
    }

    return float64(totalDegree) / float64(vertexCount)
}
```

### 6.2 社交网络分析

```go
// 社交网络分析器
type SocialNetworkAnalyzer struct {
    graph *AdjacencyListGraph
}

func NewSocialNetworkAnalyzer() *SocialNetworkAnalyzer {
    return &SocialNetworkAnalyzer{
        graph: NewAdjacencyListGraph(false),
    }
}

func (sna *SocialNetworkAnalyzer) AddFriendship(user1, user2 int) {
    sna.graph.AddEdge(user1, user2)
}

func (sna *SocialNetworkAnalyzer) FindInfluentialUsers() []int {
    // 使用度中心性找出有影响力的用户
    maxDegree := 0
    influentialUsers := []int{}

    for v := range sna.graph.vertices {
        degree := len(sna.graph.GetNeighbors(v))
        if degree > maxDegree {
            maxDegree = degree
            influentialUsers = []int{v}
        } else if degree == maxDegree {
            influentialUsers = append(influentialUsers, v)
        }
    }

    return influentialUsers
}

func (sna *SocialNetworkAnalyzer) FindCommunities() [][]int {
    // 使用标签传播算法发现社区
    labels := make(map[int]int)
    for i, v := range sna.graph.GetVertices() {
        labels[v] = i
    }

    changed := true
    for changed {
        changed = false
        for v := range sna.graph.vertices {
            // 计算邻居中最常见的标签
            labelCount := make(map[int]int)
            for _, neighbor := range sna.graph.GetNeighbors(v) {
                labelCount[labels[neighbor]]++
            }

            // 选择最常见的标签
            maxCount := 0
            newLabel := labels[v]
            for label, count := range labelCount {
                if count > maxCount {
                    maxCount = count
                    newLabel = label
                }
            }

            if newLabel != labels[v] {
                labels[v] = newLabel
                changed = true
            }
        }
    }

    // 按标签分组
    communities := make(map[int][]int)
    for v, label := range labels {
        communities[label] = append(communities[label], v)
    }

    result := [][]int{}
    for _, community := range communities {
        result = append(result, community)
    }

    return result
}
```

### 6.3 路由算法

```go
// 路由表项
type RoutingEntry struct {
    Destination int
    NextHop     int
    Cost        int
    Path        []int
}

// 路由算法
type RoutingAlgorithm struct {
    graph map[int][]Edge
}

func NewRoutingAlgorithm() *RoutingAlgorithm {
    return &RoutingAlgorithm{
        graph: make(map[int][]Edge),
    }
}

func (ra *RoutingAlgorithm) AddLink(from, to, cost int) {
    ra.graph[from] = append(ra.graph[from], Edge{To: to, Weight: cost})
    ra.graph[to] = append(ra.graph[to], Edge{To: from, Weight: cost})
}

func (ra *RoutingAlgorithm) ComputeRoutingTable(source int) map[int]RoutingEntry {
    distances := Dijkstra(ra.graph, source)
    routingTable := make(map[int]RoutingEntry)

    for dest := range ra.graph {
        if dest != source {
            path := ra.findPath(source, dest)
            nextHop := path[1] if len(path) > 1 else dest

            routingTable[dest] = RoutingEntry{
                Destination: dest,
                NextHop:     nextHop,
                Cost:        distances[dest],
                Path:        path,
            }
        }
    }

    return routingTable
}

func (ra *RoutingAlgorithm) findPath(source, dest int) []int {
    // 使用Dijkstra算法重建路径
    distances := make(map[int]int)
    previous := make(map[int]int)

    for node := range ra.graph {
        distances[node] = math.MaxInt32
    }
    distances[source] = 0

    pq := &PriorityQueue{}
    heap.Push(pq, &Item{node: source, distance: 0})

    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        node := item.node
        dist := item.distance

        if dist > distances[node] {
            continue
        }

        for _, edge := range ra.graph[node] {
            newDist := dist + edge.Weight
            if newDist < distances[edge.To] {
                distances[edge.To] = newDist
                previous[edge.To] = node
                heap.Push(pq, &Item{node: edge.To, distance: newDist})
            }
        }
    }

    // 重建路径
    path := []int{}
    current := dest
    for current != source {
        path = append([]int{current}, path...)
        current = previous[current]
    }
    path = append([]int{source}, path...)

    return path
}
```

## 总结

图论作为计算机科学的基础理论，在软件工程中有着广泛的应用。本章从数学定义出发，通过形式化证明建立了图论的理论基础，并提供了完整的Go语言实现。

主要内容包括：

1. **理论基础**：图的定义、性质、定理和证明
2. **算法实现**：遍历、最短路径、最小生成树等经典算法
3. **实际应用**：网络分析、社交网络、路由算法等

这些内容为后续的软件架构设计、算法分析和系统优化提供了重要的理论基础和实践指导。

---

**参考文献**：

1. Cormen, T. H., et al. "Introduction to Algorithms." MIT Press, 2009.
2. Bondy, J. A., & Murty, U. S. R. "Graph Theory." Springer, 2008.
3. West, D. B. "Introduction to Graph Theory." Prentice Hall, 2001.
