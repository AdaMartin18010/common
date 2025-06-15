# 03-图论 (Graph Theory)

## 概述

图论是数学的一个分支，研究图（Graph）的数学性质。图由顶点（Vertex）和边（Edge）组成，是描述对象之间关系的重要数学工具。在计算机科学中，图论广泛应用于算法设计、网络分析、数据结构等领域。

## 1. 基本概念

### 1.1 图的定义

**定义 1.1** (图)
一个图 $G = (V, E)$ 由以下两部分组成：
- $V$ 是顶点的有限集合，称为顶点集（Vertex Set）
- $E$ 是边的集合，每条边连接两个顶点，称为边集（Edge Set）

**形式化定义**：
```latex
G = (V, E) \text{ where } E \subseteq V \times V
```

### 1.2 图的类型

#### 1.2.1 无向图 (Undirected Graph)

**定义 1.2** (无向图)
在无向图中，边没有方向，$(u,v) \in E$ 表示顶点 $u$ 和 $v$ 之间有一条边。

**性质**：
- 对称性：$(u,v) \in E \Leftrightarrow (v,u) \in E$
- 自环：$(v,v) \in E$ 称为自环

#### 1.2.2 有向图 (Directed Graph)

**定义 1.3** (有向图)
在有向图中，边有方向，$(u,v) \in E$ 表示从顶点 $u$ 到顶点 $v$ 有一条有向边。

**性质**：
- 非对称性：$(u,v) \in E$ 不一定意味着 $(v,u) \in E$
- 入度和出度：每个顶点有入度和出度

### 1.3 图的基本性质

#### 1.3.1 度数 (Degree)

**定义 1.4** (度数)
顶点 $v$ 的度数 $deg(v)$ 是与该顶点相连的边的数量。

**定理 1.1** (握手定理)
对于任意图 $G = (V, E)$，所有顶点的度数之和等于边数的两倍：

$$\sum_{v \in V} deg(v) = 2|E|$$

**证明**：
每条边贡献给两个顶点的度数，因此所有度数之和等于边数的两倍。

#### 1.3.2 路径和连通性

**定义 1.5** (路径)
图 $G$ 中的路径是顶点序列 $v_0, v_1, \ldots, v_k$，其中 $(v_i, v_{i+1}) \in E$ 对所有 $0 \leq i < k$ 成立。

**定义 1.6** (连通图)
无向图 $G$ 是连通的，当且仅当任意两个顶点之间都存在路径。

## 2. 图的表示

### 2.1 邻接矩阵 (Adjacency Matrix)

**定义 2.1** (邻接矩阵)
图 $G = (V, E)$ 的邻接矩阵 $A$ 是一个 $|V| \times |V|$ 的矩阵，其中：

$$A[i][j] = \begin{cases}
1 & \text{if } (i,j) \in E \\
0 & \text{otherwise}
\end{cases}$$

### 2.2 邻接表 (Adjacency List)

**定义 2.2** (邻接表)
图 $G = (V, E)$ 的邻接表是一个数组，其中每个元素是一个链表，存储与该顶点相邻的所有顶点。

## 3. Go语言实现

### 3.1 图的基本结构

```go
package graph

import (
    "fmt"
    "math"
)

// Vertex 表示图中的顶点
type Vertex struct {
    ID   int
    Data interface{}
}

// Edge 表示图中的边
type Edge struct {
    From   int
    To     int
    Weight float64
}

// Graph 表示图的基本结构
type Graph struct {
    vertices map[int]*Vertex
    edges    map[int][]Edge
    directed bool
}

// NewGraph 创建新的图
func NewGraph(directed bool) *Graph {
    return &Graph{
        vertices: make(map[int]*Vertex),
        edges:    make(map[int][]Edge),
        directed: directed,
    }
}

// AddVertex 添加顶点
func (g *Graph) AddVertex(id int, data interface{}) {
    g.vertices[id] = &Vertex{ID: id, Data: data}
    if g.edges[id] == nil {
        g.edges[id] = make([]Edge, 0)
    }
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to int, weight float64) {
    edge := Edge{From: from, To: to, Weight: weight}
    g.edges[from] = append(g.edges[from], edge)
    
    // 如果是无向图，添加反向边
    if !g.directed {
        reverseEdge := Edge{From: to, To: from, Weight: weight}
        g.edges[to] = append(g.edges[to], reverseEdge)
    }
}

// GetDegree 获取顶点的度数
func (g *Graph) GetDegree(vertexID int) int {
    if edges, exists := g.edges[vertexID]; exists {
        return len(edges)
    }
    return 0
}

// GetTotalDegree 获取所有顶点的度数之和
func (g *Graph) GetTotalDegree() int {
    total := 0
    for _, edges := range g.edges {
        total += len(edges)
    }
    return total
}
```

### 3.2 图的遍历算法

#### 3.2.1 深度优先搜索 (DFS)

**算法描述**：
深度优先搜索从起始顶点开始，沿着图的边尽可能深入地访问顶点，直到无法继续前进，然后回溯。

**时间复杂度**：$O(|V| + |E|)$

```go
// DFS 深度优先搜索
func (g *Graph) DFS(startID int) []int {
    visited := make(map[int]bool)
    result := make([]int, 0)
    
    var dfs func(int)
    dfs = func(vertexID int) {
        if visited[vertexID] {
            return
        }
        
        visited[vertexID] = true
        result = append(result, vertexID)
        
        for _, edge := range g.edges[vertexID] {
            dfs(edge.To)
        }
    }
    
    dfs(startID)
    return result
}

// DFSIterative 迭代式深度优先搜索
func (g *Graph) DFSIterative(startID int) []int {
    visited := make(map[int]bool)
    result := make([]int, 0)
    stack := []int{startID}
    
    for len(stack) > 0 {
        vertexID := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if visited[vertexID] {
            continue
        }
        
        visited[vertexID] = true
        result = append(result, vertexID)
        
        // 将邻接顶点压入栈中（逆序以保持正确的访问顺序）
        for i := len(g.edges[vertexID]) - 1; i >= 0; i-- {
            neighbor := g.edges[vertexID][i].To
            if !visited[neighbor] {
                stack = append(stack, neighbor)
            }
        }
    }
    
    return result
}
```

#### 3.2.2 广度优先搜索 (BFS)

**算法描述**：
广度优先搜索从起始顶点开始，先访问所有相邻顶点，然后逐层向外扩展。

**时间复杂度**：$O(|V| + |E|)$

```go
// BFS 广度优先搜索
func (g *Graph) BFS(startID int) []int {
    visited := make(map[int]bool)
    result := make([]int, 0)
    queue := []int{startID}
    visited[startID] = true
    
    for len(queue) > 0 {
        vertexID := queue[0]
        queue = queue[1:]
        result = append(result, vertexID)
        
        for _, edge := range g.edges[vertexID] {
            neighbor := edge.To
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    
    return result
}
```

### 3.3 最短路径算法

#### 3.3.1 Dijkstra算法

**算法描述**：
Dijkstra算法用于找到从起始顶点到所有其他顶点的最短路径。

**时间复杂度**：$O((|V| + |E|) \log |V|)$

```go
// Dijkstra 最短路径算法
func (g *Graph) Dijkstra(startID int) map[int]float64 {
    distances := make(map[int]float64)
    visited := make(map[int]bool)
    
    // 初始化距离
    for vertexID := range g.vertices {
        distances[vertexID] = math.Inf(1)
    }
    distances[startID] = 0
    
    for len(visited) < len(g.vertices) {
        // 找到未访问顶点中距离最小的
        minVertex := -1
        minDist := math.Inf(1)
        
        for vertexID, dist := range distances {
            if !visited[vertexID] && dist < minDist {
                minDist = dist
                minVertex = vertexID
            }
        }
        
        if minVertex == -1 {
            break
        }
        
        visited[minVertex] = true
        
        // 更新邻接顶点的距离
        for _, edge := range g.edges[minVertex] {
            neighbor := edge.To
            newDist := distances[minVertex] + edge.Weight
            if newDist < distances[neighbor] {
                distances[neighbor] = newDist
            }
        }
    }
    
    return distances
}
```

#### 3.3.2 Floyd-Warshall算法

**算法描述**：
Floyd-Warshall算法用于找到所有顶点对之间的最短路径。

**时间复杂度**：$O(|V|^3)$

```go
// FloydWarshall 全源最短路径算法
func (g *Graph) FloydWarshall() [][]float64 {
    n := len(g.vertices)
    dist := make([][]float64, n)
    
    // 初始化距离矩阵
    for i := 0; i < n; i++ {
        dist[i] = make([]float64, n)
        for j := 0; j < n; j++ {
            if i == j {
                dist[i][j] = 0
            } else {
                dist[i][j] = math.Inf(1)
            }
        }
    }
    
    // 设置直接边的距离
    for vertexID, edges := range g.edges {
        for _, edge := range edges {
            dist[vertexID][edge.To] = edge.Weight
        }
    }
    
    // Floyd-Warshall算法核心
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

### 3.4 最小生成树算法

#### 3.4.1 Kruskal算法

**算法描述**：
Kruskal算法用于找到无向连通图的最小生成树。

**时间复杂度**：$O(|E| \log |E|)$

```go
// Edge 用于Kruskal算法
type KruskalEdge struct {
    From   int
    To     int
    Weight float64
}

// Kruskal 最小生成树算法
func (g *Graph) Kruskal() []Edge {
    // 收集所有边
    edges := make([]KruskalEdge, 0)
    for vertexID, edgeList := range g.edges {
        for _, edge := range edgeList {
            if !g.directed || edge.From < edge.To { // 避免重复边
                edges = append(edges, KruskalEdge{
                    From:   edge.From,
                    To:     edge.To,
                    Weight: edge.Weight,
                })
            }
        }
    }
    
    // 按权重排序
    sort.Slice(edges, func(i, j int) bool {
        return edges[i].Weight < edges[j].Weight
    })
    
    // 并查集
    parent := make(map[int]int)
    rank := make(map[int]int)
    
    // 初始化并查集
    for vertexID := range g.vertices {
        parent[vertexID] = vertexID
        rank[vertexID] = 0
    }
    
    // 查找根节点
    find := func(x int) int {
        if parent[x] != x {
            parent[x] = find(parent[x])
        }
        return parent[x]
    }
    
    // 合并两个集合
    union := func(x, y int) bool {
        rootX := find(x)
        rootY := find(y)
        
        if rootX == rootY {
            return false
        }
        
        if rank[rootX] < rank[rootY] {
            parent[rootX] = rootY
        } else if rank[rootX] > rank[rootY] {
            parent[rootY] = rootX
        } else {
            parent[rootY] = rootX
            rank[rootX]++
        }
        return true
    }
    
    // Kruskal算法核心
    mst := make([]Edge, 0)
    for _, edge := range edges {
        if union(edge.From, edge.To) {
            mst = append(mst, Edge{
                From:   edge.From,
                To:     edge.To,
                Weight: edge.Weight,
            })
        }
    }
    
    return mst
}
```

## 4. 图论定理和证明

### 4.1 欧拉路径定理

**定理 4.1** (欧拉路径定理)
连通图 $G$ 存在欧拉路径当且仅当：
- 所有顶点的度数都是偶数，或者
- 恰好有两个顶点的度数是奇数

**证明**：
1. **必要性**：如果存在欧拉路径，那么除了起点和终点外，每个顶点被访问的次数必须是偶数。
2. **充分性**：通过构造性证明，可以找到欧拉路径。

### 4.2 哈密顿路径问题

**定义 4.1** (哈密顿路径)
哈密顿路径是经过图中每个顶点恰好一次的路径。

**定理 4.2** (Dirac定理)
如果图 $G$ 有 $n$ 个顶点，且每个顶点的度数至少为 $n/2$，则 $G$ 存在哈密顿回路。

## 5. 应用实例

### 5.1 社交网络分析

```go
// SocialNetwork 社交网络图
type SocialNetwork struct {
    *Graph
}

// NewSocialNetwork 创建社交网络
func NewSocialNetwork() *SocialNetwork {
    return &SocialNetwork{
        Graph: NewGraph(false), // 无向图
    }
}

// AddFriendship 添加朋友关系
func (sn *SocialNetwork) AddFriendship(user1, user2 int) {
    sn.AddEdge(user1, user2, 1.0)
}

// GetFriendsOfFriends 获取朋友的朋友
func (sn *SocialNetwork) GetFriendsOfFriends(userID int) []int {
    friends := make(map[int]bool)
    
    // 获取直接朋友
    for _, edge := range sn.edges[userID] {
        friends[edge.To] = true
    }
    
    // 获取朋友的朋友
    friendsOfFriends := make(map[int]bool)
    for friend := range friends {
        for _, edge := range sn.edges[friend] {
            if edge.To != userID && !friends[edge.To] {
                friendsOfFriends[edge.To] = true
            }
        }
    }
    
    result := make([]int, 0, len(friendsOfFriends))
    for friend := range friendsOfFriends {
        result = append(result, friend)
    }
    
    return result
}

// GetShortestPath 获取两个用户之间的最短路径
func (sn *SocialNetwork) GetShortestPath(from, to int) []int {
    distances := sn.Dijkstra(from)
    if distances[to] == math.Inf(1) {
        return nil
    }
    
    // 重建路径
    path := []int{to}
    current := to
    
    for current != from {
        for _, edge := range sn.edges[current] {
            if distances[edge.To]+edge.Weight == distances[current] {
                path = append([]int{edge.To}, path...)
                current = edge.To
                break
            }
        }
    }
    
    return path
}
```

### 5.2 网络路由

```go
// NetworkRouter 网络路由器
type NetworkRouter struct {
    *Graph
}

// NewNetworkRouter 创建网络路由器
func NewNetworkRouter() *NetworkRouter {
    return &NetworkRouter{
        Graph: NewGraph(true), // 有向图
    }
}

// AddConnection 添加网络连接
func (nr *NetworkRouter) AddConnection(from, to int, bandwidth float64) {
    nr.AddEdge(from, to, 1.0/bandwidth) // 权重为延迟
}

// FindOptimalRoute 找到最优路由
func (nr *NetworkRouter) FindOptimalRoute(from, to int) []int {
    return nr.GetShortestPath(from, to)
}

// GetNetworkTopology 获取网络拓扑
func (nr *NetworkRouter) GetNetworkTopology() map[int][]int {
    topology := make(map[int][]int)
    for vertexID, edges := range nr.edges {
        neighbors := make([]int, 0, len(edges))
        for _, edge := range edges {
            neighbors = append(neighbors, edge.To)
        }
        topology[vertexID] = neighbors
    }
    return topology
}
```

## 6. 性能分析

### 6.1 时间复杂度分析

| 算法 | 时间复杂度 | 空间复杂度 |
|------|------------|------------|
| DFS | $O(|V| + |E|)$ | $O(|V|)$ |
| BFS | $O(|V| + |E|)$ | $O(|V|)$ |
| Dijkstra | $O((|V| + |E|) \log |V|)$ | $O(|V|)$ |
| Floyd-Warshall | $O(|V|^3)$ | $O(|V|^2)$ |
| Kruskal | $O(|E| \log |E|)$ | $O(|V|)$ |

### 6.2 空间复杂度分析

- **邻接矩阵**：$O(|V|^2)$
- **邻接表**：$O(|V| + |E|)$
- **稀疏图**：邻接表更优
- **稠密图**：邻接矩阵更优

## 7. 总结

图论是计算机科学中重要的数学基础，提供了描述和分析复杂关系的强大工具。通过Go语言的实现，我们可以：

1. **高效表示**：使用邻接表或邻接矩阵表示图
2. **快速遍历**：实现DFS和BFS算法
3. **路径优化**：使用Dijkstra和Floyd-Warshall算法
4. **最小生成树**：使用Kruskal算法
5. **实际应用**：社交网络分析、网络路由等

图论的理论基础和算法实现为构建复杂的软件系统提供了重要的数学支撑。

## 参考文献

1. Cormen, T. H., Leiserson, C. E., Rivest, R. L., & Stein, C. (2009). Introduction to Algorithms (3rd ed.). MIT Press.
2. Bondy, J. A., & Murty, U. S. R. (2008). Graph Theory. Springer.
3. West, D. B. (2001). Introduction to Graph Theory (2nd ed.). Prentice Hall.
