# 负载均衡 (Load Balancing)

## 目录

- [负载均衡 (Load Balancing)](#负载均衡-load-balancing)
  - [目录](#目录)
  - [1. 理论基础](#1-理论基础)
    - [1.1 负载均衡定义](#11-负载均衡定义)
    - [1.2 负载均衡分类](#12-负载均衡分类)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 负载均衡系统模型](#21-负载均衡系统模型)
    - [2.2 负载均衡目标函数](#22-负载均衡目标函数)
    - [2.3 算法复杂度分析](#23-算法复杂度分析)
  - [3. 算法实现](#3-算法实现)
    - [3.1 轮询算法 (Round Robin)](#31-轮询算法-round-robin)
    - [3.2 加权轮询算法 (Weighted Round Robin)](#32-加权轮询算法-weighted-round-robin)
    - [3.3 最少连接算法 (Least Connections)](#33-最少连接算法-least-connections)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 基础接口定义](#41-基础接口定义)
    - [4.2 轮询算法实现](#42-轮询算法实现)
  - [5. 性能分析](#5-性能分析)
  - [6. 实际应用](#6-实际应用)

---

## 1. 理论基础

### 1.1 负载均衡定义

负载均衡是一种分布式系统技术，用于在多个服务器之间分配工作负载，以提高系统的整体性能、可靠性和可扩展性。

**形式化定义**:
设 $S = \{s_1, s_2, ..., s_n\}$ 为服务器集合，$R = \{r_1, r_2, ..., r_m\}$ 为请求集合。负载均衡函数 $L: R \to S$ 将每个请求映射到一个服务器。
目标是最小化所有服务器中的最大负载：
$$
\min \max_{s_i \in S} \sum_{r_j \in L^{-1}(s_i)} w(r_j)
$$
其中 $w(r_j)$ 为请求 $r_j$ 的权重或成本。

### 1.2 负载均衡分类

#### 1.2.1 按层次分类

1.  **应用层负载均衡 (Layer 7)**
    -   基于HTTP/HTTPS协议。
    -   支持内容感知路由（例如，基于URL路径或Header）。
    -   可进行SSL终止、HTTP压缩等高级功能。

2.  **传输层负载均衡 (Layer 4)**
    -   基于TCP/UDP协议。
    -   高性能，低延迟。
    -   不感知应用层内容，仅根据IP地址和端口进行转发。

#### 1.2.2 按算法分类

1.  **静态算法**
    -   轮询 (Round Robin)
    -   加权轮询 (Weighted Round Robin)
    -   IP哈希 (IP Hash)

2.  **动态算法**
    -   最少连接 (Least Connections)
    -   加权最少连接 (Weighted Least Connections)
    -   最快响应时间 (Fastest Response Time)

---

## 2. 形式化定义

### 2.1 负载均衡系统模型

一个负载均衡系统可以形式化地定义为一个五元组：
$LB = (S, R, L, M, C)$

其中：
- $S$: 服务器集合, $S = \{s_1, s_2, ..., s_n\}$
- $R$: 请求集合, $R = \{r_1, r_2, ..., r_m\}$
- $L$: 负载均衡函数, $L: R \to S$
- $M$: 监控函数, $M: S \to \mathbb{R}^+$ (例如，返回服务器的连接数、CPU利用率等指标)
- $C$: 约束条件集合 (例如，服务器容量限制)

### 2.2 负载均衡目标函数

目标是最小化最大负载：
$$
\min \max_{s_i \in S} \left( \sum_{r_j \in L^{-1}(s_i)} w(r_j) \right)
$$

**约束条件**:
1.  所有请求必须被分配: $\bigcup_{s_i \in S} L^{-1}(s_i) = R$
2.  每个服务器的负载不能超过其容量: $\forall s_i \in S, \sum_{r_j \in L^{-1}(s_i)} w(r_j) \le C_i$

### 2.3 算法复杂度分析

**定理 2.1**: 最优负载均衡问题（即最小化最大负载）是NP难问题。

**证明**:
可以将该问题规约到著名的**分区问题 (Partition Problem)**。
给定一个多重集 $A$ of 整数，判断是否存在一个子集 $A' \subseteq A$ 使得 $\sum_{a \in A'} a = \sum_{b \in A \setminus A'} b$。
这等价于一个有2台服务器的负载均衡问题，其中请求的权重对应集合中的整数。
由于分区问题是NP难的，因此更通用的负载均衡问题也是NP难的。

---

## 3. 算法实现

### 3.1 轮询算法 (Round Robin)

$$
L_{RR}(r_i) = s_{(i-1) \pmod n + 1}
$$
其中 $i$ 是请求的序号，$n = |S|$ 是服务器数量。

- **时间复杂度**: O(1)
- **空间复杂度**: O(1)

### 3.2 加权轮询算法 (Weighted Round Robin)

设 $W = \{w_1, w_2, ..., w_n\}$ 为服务器权重。一种常见的实现是平滑加权轮询（Smooth Weighted Round-Robin）。

### 3.3 最少连接算法 (Least Connections)

$$
L_{LC}(r_i) = \arg\min_{s_j \in S} \text{connections}(s_j)
$$
其中 `connections(s_j)` 是服务器 $s_j$ 的当前活动连接数。

---

## 4. Go语言实现

### 4.1 基础接口定义

```go
package loadbalancer

import (
    "errors"
    "sync"
    "time"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
    SelectServer(request *Request) (*Server, error)
    AddServer(server *Server) error
    RemoveServer(serverID string) error
    UpdateServerStatus(serverID string, status ServerStatus) error
}

// Server 服务器定义
type Server struct {
    ID           string
    Address      string
    Port         int
    Weight       int
    Status       ServerStatus
    Connections  int
    ResponseTime time.Duration
}

// ServerStatus 服务器状态
type ServerStatus int

const (
    ServerStatusHealthy ServerStatus = iota
    ServerStatusUnhealthy
    ServerStatusMaintenance
)

// Request 请求定义 (用于策略决策)
type Request struct {
    ID       string
    ClientIP string
    Headers  map[string]string
}
```

### 4.2 轮询算法实现

```go
// RoundRobinLoadBalancer 轮询负载均衡器
type RoundRobinLoadBalancer struct {
    servers []*Server
    current uint32
    mu      sync.RWMutex
}

// NewRoundRobinLoadBalancer 创建轮询负载均衡器
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
    return &RoundRobinLoadBalancer{
        servers: make([]*Server, 0),
        current: 0,
    }
}

// SelectServer 选择服务器
func (rr *RoundRobinLoadBalancer) SelectServer(request *Request) (*Server, error) {
    rr.mu.RLock()
    defer rr.mu.RUnlock()

    if len(rr.servers) == 0 {
        return nil, errors.New("no available servers")
    }

    // 使用原子操作保证并发安全
    next := atomic.AddUint32(&rr.current, 1)
    server := rr.servers[(next-1)%uint32(len(rr.servers))]
    
    return server, nil
}

// AddServer 添加服务器
func (rr *RoundRobinLoadBalancer) AddServer(server *Server) error {
    rr.mu.Lock()
    defer rr.mu.Unlock()

    rr.servers = append(rr.servers, server)
    return nil
}

// RemoveServer 移除服务器
func (rr *RoundRobinLoadBalancer) RemoveServer(serverID string) error {
    rr.mu.Lock()
    defer rr.mu.Unlock()

    for i, server := range rr.servers {
        if server.ID == serverID {
            rr.servers = append(rr.servers[:i], rr.servers[i+1:]...)
            return nil
        }
    }
    return errors.New("server not found")
}

// UpdateServerStatus 更新服务器状态 (简单实现，实际中可能需要将不健康的服务器临时移出轮询列表)
func (rr *RoundRobinLoadBalancer) UpdateServerStatus(serverID string, status ServerStatus) error {
    rr.mu.Lock()
    defer rr.mu.Unlock()

    for _, server := range rr.servers {
        if server.ID == serverID {
            server.Status = status
            return nil
        }
    }
    return errors.New("server not found")
}
```

---

## 5. 性能分析

| 算法 | 时间复杂度 | 空间复杂度 | 动态适应性 |
| :--- | :--- | :--- | :--- |
| 轮询 | O(1) | O(N) | 差 |
| 加权轮询 | O(1) (平滑) | O(N) | 中 |
| IP哈希 | O(1) | O(N) | 中 |
| 最少连接 | O(log N) | O(N) | 好 |

---

## 6. 实际应用

- **Nginx**: 同时支持Layer 4和Layer 7负载均衡，提供轮询、加权轮询、IP哈希、最少连接等多种算法。
- **HAProxy**: 高性能的TCP/HTTP负载均衡器。
- **Kubernetes Service**: 使用iptables或IPVS实现基本的Layer 4负载均衡。
- **云服务商 (AWS ELB, GCP Cloud Load Balancing)**: 提供托管的、高可用的负载均衡服务。 