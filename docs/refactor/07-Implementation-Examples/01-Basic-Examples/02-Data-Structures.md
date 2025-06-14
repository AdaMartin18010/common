# 02-数据结构 (Data Structures)

## 概述

数据结构是计算机科学的基础，用于组织和存储数据。本文档介绍基本数据结构的Go语言实现，包括线性结构和非线性结构。

## 目录

- [02-数据结构 (Data Structures)](#02-数据结构-data-structures)
  - [概述](#概述)
  - [目录](#目录)
  - [1. 线性数据结构 (Linear Data Structures)](#1-线性数据结构-linear-data-structures)
    - [1.1 链表 (Linked List)](#11-链表-linked-list)
    - [1.2 栈 (Stack)](#12-栈-stack)
    - [1.3 队列 (Queue)](#13-队列-queue)
  - [2. 树形数据结构 (Tree Data Structures)](#2-树形数据结构-tree-data-structures)
    - [2.1 二叉树 (Binary Tree)](#21-二叉树-binary-tree)
    - [2.2 二叉搜索树 (Binary Search Tree)](#22-二叉搜索树-binary-search-tree)
  - [3. 图数据结构 (Graph Data Structures)](#3-图数据结构-graph-data-structures)
    - [3.1 邻接表图](#31-邻接表图)
  - [4. 散列表 (Hash Tables)](#4-散列表-hash-tables)
    - [4.1 基本散列表](#41-基本散列表)
  - [5. 高级数据结构 (Advanced Data Structures)](#5-高级数据结构-advanced-data-structures)
    - [5.1 红黑树](#51-红黑树)
  - [总结](#总结)

---

## 1. 线性数据结构 (Linear Data Structures)

### 1.1 链表 (Linked List)

```go
package linkedlist

import "fmt"

// Node 链表节点
type Node struct {
    Data interface{}
    Next *Node
}

// LinkedList 链表
type LinkedList struct {
    Head *Node
    Size int
}

// NewLinkedList 创建链表
func NewLinkedList() *LinkedList {
    return &LinkedList{Head: nil, Size: 0}
}

// InsertAtBeginning 在开头插入
func (ll *LinkedList) InsertAtBeginning(data interface{}) {
    newNode := &Node{Data: data, Next: ll.Head}
    ll.Head = newNode
    ll.Size++
}

// InsertAtEnd 在末尾插入
func (ll *LinkedList) InsertAtEnd(data interface{}) {
    newNode := &Node{Data: data, Next: nil}
    
    if ll.Head == nil {
        ll.Head = newNode
    } else {
        current := ll.Head
        for current.Next != nil {
            current = current.Next
        }
        current.Next = newNode
    }
    ll.Size++
}

// InsertAtPosition 在指定位置插入
func (ll *LinkedList) InsertAtPosition(data interface{}, position int) error {
    if position < 0 || position > ll.Size {
        return fmt.Errorf("invalid position")
    }
    
    if position == 0 {
        ll.InsertAtBeginning(data)
        return nil
    }
    
    newNode := &Node{Data: data, Next: nil}
    current := ll.Head
    for i := 0; i < position-1; i++ {
        current = current.Next
    }
    
    newNode.Next = current.Next
    current.Next = newNode
    ll.Size++
    return nil
}

// DeleteAtBeginning 删除开头元素
func (ll *LinkedList) DeleteAtBeginning() (interface{}, error) {
    if ll.Head == nil {
        return nil, fmt.Errorf("list is empty")
    }
    
    data := ll.Head.Data
    ll.Head = ll.Head.Next
    ll.Size--
    return data, nil
}

// DeleteAtEnd 删除末尾元素
func (ll *LinkedList) DeleteAtEnd() (interface{}, error) {
    if ll.Head == nil {
        return nil, fmt.Errorf("list is empty")
    }
    
    if ll.Head.Next == nil {
        data := ll.Head.Data
        ll.Head = nil
        ll.Size--
        return data, nil
    }
    
    current := ll.Head
    for current.Next.Next != nil {
        current = current.Next
    }
    
    data := current.Next.Data
    current.Next = nil
    ll.Size--
    return data, nil
}

// Search 搜索元素
func (ll *LinkedList) Search(data interface{}) (int, bool) {
    current := ll.Head
    position := 0
    
    for current != nil {
        if current.Data == data {
            return position, true
        }
        current = current.Next
        position++
    }
    
    return -1, false
}

// Reverse 反转链表
func (ll *LinkedList) Reverse() {
    var prev *Node
    current := ll.Head
    
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    
    ll.Head = prev
}

// Print 打印链表
func (ll *LinkedList) Print() {
    current := ll.Head
    for current != nil {
        fmt.Printf("%v -> ", current.Data)
        current = current.Next
    }
    fmt.Println("nil")
}

// DoublyLinkedList 双向链表
type DoublyNode struct {
    Data interface{}
    Prev *DoublyNode
    Next *DoublyNode
}

type DoublyLinkedList struct {
    Head *DoublyNode
    Tail *DoublyNode
    Size int
}

func NewDoublyLinkedList() *DoublyLinkedList {
    return &DoublyLinkedList{Head: nil, Tail: nil, Size: 0}
}

func (dll *DoublyLinkedList) InsertAtEnd(data interface{}) {
    newNode := &DoublyNode{Data: data, Prev: nil, Next: nil}
    
    if dll.Head == nil {
        dll.Head = newNode
        dll.Tail = newNode
    } else {
        newNode.Prev = dll.Tail
        dll.Tail.Next = newNode
        dll.Tail = newNode
    }
    dll.Size++
}
```

### 1.2 栈 (Stack)

```go
package stack

import (
    "fmt"
    "sync"
)

// Stack 栈接口
type Stack interface {
    Push(data interface{})
    Pop() (interface{}, error)
    Peek() (interface{}, error)
    IsEmpty() bool
    Size() int
}

// ArrayStack 基于数组的栈
type ArrayStack struct {
    data []interface{}
    top  int
    mu   sync.RWMutex
}

// NewArrayStack 创建数组栈
func NewArrayStack() *ArrayStack {
    return &ArrayStack{
        data: make([]interface{}, 0),
        top:  -1,
    }
}

// Push 入栈
func (s *ArrayStack) Push(data interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.data = append(s.data, data)
    s.top++
}

// Pop 出栈
func (s *ArrayStack) Pop() (interface{}, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.IsEmpty() {
        return nil, fmt.Errorf("stack is empty")
    }
    
    data := s.data[s.top]
    s.data = s.data[:s.top]
    s.top--
    return data, nil
}

// Peek 查看栈顶元素
func (s *ArrayStack) Peek() (interface{}, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if s.IsEmpty() {
        return nil, fmt.Errorf("stack is empty")
    }
    
    return s.data[s.top], nil
}

// IsEmpty 检查栈是否为空
func (s *ArrayStack) IsEmpty() bool {
    return s.top == -1
}

// Size 获取栈大小
func (s *ArrayStack) Size() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.top + 1
}

// LinkedStack 基于链表的栈
type LinkedStack struct {
    head *Node
    size int
    mu   sync.RWMutex
}

// NewLinkedStack 创建链表栈
func NewLinkedStack() *LinkedStack {
    return &LinkedStack{head: nil, size: 0}
}

// Push 入栈
func (s *LinkedStack) Push(data interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    newNode := &Node{Data: data, Next: s.head}
    s.head = newNode
    s.size++
}

// Pop 出栈
func (s *LinkedStack) Pop() (interface{}, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.IsEmpty() {
        return nil, fmt.Errorf("stack is empty")
    }
    
    data := s.head.Data
    s.head = s.head.Next
    s.size--
    return data, nil
}

// Peek 查看栈顶元素
func (s *LinkedStack) Peek() (interface{}, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if s.IsEmpty() {
        return nil, fmt.Errorf("stack is empty")
    }
    
    return s.head.Data, nil
}

// IsEmpty 检查栈是否为空
func (s *LinkedStack) IsEmpty() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.head == nil
}

// Size 获取栈大小
func (s *LinkedStack) Size() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.size
}
```

### 1.3 队列 (Queue)

```go
package queue

import (
    "fmt"
    "sync"
)

// Queue 队列接口
type Queue interface {
    Enqueue(data interface{})
    Dequeue() (interface{}, error)
    Front() (interface{}, error)
    IsEmpty() bool
    Size() int
}

// ArrayQueue 基于数组的队列
type ArrayQueue struct {
    data []interface{}
    head int
    tail int
    size int
    mu   sync.RWMutex
}

// NewArrayQueue 创建数组队列
func NewArrayQueue(capacity int) *ArrayQueue {
    return &ArrayQueue{
        data: make([]interface{}, capacity),
        head: 0,
        tail: 0,
        size: 0,
    }
}

// Enqueue 入队
func (q *ArrayQueue) Enqueue(data interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if q.size == len(q.data) {
        // 扩容
        newData := make([]interface{}, len(q.data)*2)
        for i := 0; i < q.size; i++ {
            newData[i] = q.data[(q.head+i)%len(q.data)]
        }
        q.data = newData
        q.head = 0
        q.tail = q.size
    }
    
    q.data[q.tail] = data
    q.tail = (q.tail + 1) % len(q.data)
    q.size++
}

// Dequeue 出队
func (q *ArrayQueue) Dequeue() (interface{}, error) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if q.IsEmpty() {
        return nil, fmt.Errorf("queue is empty")
    }
    
    data := q.data[q.head]
    q.head = (q.head + 1) % len(q.data)
    q.size--
    return data, nil
}

// Front 查看队首元素
func (q *ArrayQueue) Front() (interface{}, error) {
    q.mu.RLock()
    defer q.mu.RUnlock()
    
    if q.IsEmpty() {
        return nil, fmt.Errorf("queue is empty")
    }
    
    return q.data[q.head], nil
}

// IsEmpty 检查队列是否为空
func (q *ArrayQueue) IsEmpty() bool {
    return q.size == 0
}

// Size 获取队列大小
func (q *ArrayQueue) Size() int {
    q.mu.RLock()
    defer q.mu.RUnlock()
    return q.size
}

// LinkedQueue 基于链表的队列
type LinkedQueue struct {
    head *Node
    tail *Node
    size int
    mu   sync.RWMutex
}

// NewLinkedQueue 创建链表队列
func NewLinkedQueue() *LinkedQueue {
    return &LinkedQueue{head: nil, tail: nil, size: 0}
}

// Enqueue 入队
func (q *LinkedQueue) Enqueue(data interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    newNode := &Node{Data: data, Next: nil}
    
    if q.IsEmpty() {
        q.head = newNode
        q.tail = newNode
    } else {
        q.tail.Next = newNode
        q.tail = newNode
    }
    q.size++
}

// Dequeue 出队
func (q *LinkedQueue) Dequeue() (interface{}, error) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if q.IsEmpty() {
        return nil, fmt.Errorf("queue is empty")
    }
    
    data := q.head.Data
    q.head = q.head.Next
    
    if q.head == nil {
        q.tail = nil
    }
    
    q.size--
    return data, nil
}

// Front 查看队首元素
func (q *LinkedQueue) Front() (interface{}, error) {
    q.mu.RLock()
    defer q.mu.RUnlock()
    
    if q.IsEmpty() {
        return nil, fmt.Errorf("queue is empty")
    }
    
    return q.head.Data, nil
}

// IsEmpty 检查队列是否为空
func (q *LinkedQueue) IsEmpty() bool {
    return q.head == nil
}

// Size 获取队列大小
func (q *LinkedQueue) Size() int {
    q.mu.RLock()
    defer q.mu.RUnlock()
    return q.size
}

// PriorityQueue 优先队列
type PriorityQueue struct {
    data []interface{}
    less func(a, b interface{}) bool
    mu   sync.RWMutex
}

// NewPriorityQueue 创建优先队列
func NewPriorityQueue(less func(a, b interface{}) bool) *PriorityQueue {
    return &PriorityQueue{
        data: make([]interface{}, 0),
        less: less,
    }
}

// Enqueue 入队
func (pq *PriorityQueue) Enqueue(data interface{}) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    pq.data = append(pq.data, data)
    pq.heapifyUp(len(pq.data) - 1)
}

// Dequeue 出队
func (pq *PriorityQueue) Dequeue() (interface{}, error) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    if pq.IsEmpty() {
        return nil, fmt.Errorf("priority queue is empty")
    }
    
    data := pq.data[0]
    pq.data[0] = pq.data[len(pq.data)-1]
    pq.data = pq.data[:len(pq.data)-1]
    
    if len(pq.data) > 0 {
        pq.heapifyDown(0)
    }
    
    return data, nil
}

// heapifyUp 向上调整
func (pq *PriorityQueue) heapifyUp(index int) {
    for index > 0 {
        parent := (index - 1) / 2
        if pq.less(pq.data[index], pq.data[parent]) {
            pq.data[index], pq.data[parent] = pq.data[parent], pq.data[index]
            index = parent
        } else {
            break
        }
    }
}

// heapifyDown 向下调整
func (pq *PriorityQueue) heapifyDown(index int) {
    for {
        left := 2*index + 1
        right := 2*index + 2
        smallest := index
        
        if left < len(pq.data) && pq.less(pq.data[left], pq.data[smallest]) {
            smallest = left
        }
        
        if right < len(pq.data) && pq.less(pq.data[right], pq.data[smallest]) {
            smallest = right
        }
        
        if smallest == index {
            break
        }
        
        pq.data[index], pq.data[smallest] = pq.data[smallest], pq.data[index]
        index = smallest
    }
}

// IsEmpty 检查优先队列是否为空
func (pq *PriorityQueue) IsEmpty() bool {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return len(pq.data) == 0
}

// Size 获取优先队列大小
func (pq *PriorityQueue) Size() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return len(pq.data)
}
```

---

## 2. 树形数据结构 (Tree Data Structures)

### 2.1 二叉树 (Binary Tree)

```go
package tree

import (
    "fmt"
    "sync"
)

// TreeNode 树节点
type TreeNode struct {
    Data  interface{}
    Left  *TreeNode
    Right *TreeNode
}

// BinaryTree 二叉树
type BinaryTree struct {
    Root *TreeNode
    mu   sync.RWMutex
}

// NewBinaryTree 创建二叉树
func NewBinaryTree() *BinaryTree {
    return &BinaryTree{Root: nil}
}

// Insert 插入节点
func (bt *BinaryTree) Insert(data interface{}) {
    bt.mu.Lock()
    defer bt.mu.Unlock()
    
    if bt.Root == nil {
        bt.Root = &TreeNode{Data: data, Left: nil, Right: nil}
        return
    }
    
    bt.insertNode(bt.Root, data)
}

// insertNode 递归插入节点
func (bt *BinaryTree) insertNode(node *TreeNode, data interface{}) {
    // 简单的插入策略：左子树小于等于根节点，右子树大于根节点
    if compare(data, node.Data) <= 0 {
        if node.Left == nil {
            node.Left = &TreeNode{Data: data, Left: nil, Right: nil}
        } else {
            bt.insertNode(node.Left, data)
        }
    } else {
        if node.Right == nil {
            node.Right = &TreeNode{Data: data, Left: nil, Right: nil}
        } else {
            bt.insertNode(node.Right, data)
        }
    }
}

// compare 比较函数
func compare(a, b interface{}) int {
    switch a.(type) {
    case int:
        return a.(int) - b.(int)
    case string:
        if a.(string) < b.(string) {
            return -1
        } else if a.(string) > b.(string) {
            return 1
        }
        return 0
    default:
        return 0
    }
}

// Search 搜索节点
func (bt *BinaryTree) Search(data interface{}) (*TreeNode, bool) {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    return bt.searchNode(bt.Root, data)
}

// searchNode 递归搜索节点
func (bt *BinaryTree) searchNode(node *TreeNode, data interface{}) (*TreeNode, bool) {
    if node == nil {
        return nil, false
    }
    
    if compare(data, node.Data) == 0 {
        return node, true
    }
    
    if compare(data, node.Data) < 0 {
        return bt.searchNode(node.Left, data)
    }
    
    return bt.searchNode(node.Right, data)
}

// InorderTraversal 中序遍历
func (bt *BinaryTree) InorderTraversal() []interface{} {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    result := make([]interface{}, 0)
    bt.inorderTraversal(bt.Root, &result)
    return result
}

// inorderTraversal 递归中序遍历
func (bt *BinaryTree) inorderTraversal(node *TreeNode, result *[]interface{}) {
    if node == nil {
        return
    }
    
    bt.inorderTraversal(node.Left, result)
    *result = append(*result, node.Data)
    bt.inorderTraversal(node.Right, result)
}

// PreorderTraversal 前序遍历
func (bt *BinaryTree) PreorderTraversal() []interface{} {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    result := make([]interface{}, 0)
    bt.preorderTraversal(bt.Root, &result)
    return result
}

// preorderTraversal 递归前序遍历
func (bt *BinaryTree) preorderTraversal(node *TreeNode, result *[]interface{}) {
    if node == nil {
        return
    }
    
    *result = append(*result, node.Data)
    bt.preorderTraversal(node.Left, result)
    bt.preorderTraversal(node.Right, result)
}

// PostorderTraversal 后序遍历
func (bt *BinaryTree) PostorderTraversal() []interface{} {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    result := make([]interface{}, 0)
    bt.postorderTraversal(bt.Root, &result)
    return result
}

// postorderTraversal 递归后序遍历
func (bt *BinaryTree) postorderTraversal(node *TreeNode, result *[]interface{}) {
    if node == nil {
        return
    }
    
    bt.postorderTraversal(node.Left, result)
    bt.postorderTraversal(node.Right, result)
    *result = append(*result, node.Data)
}

// LevelOrderTraversal 层序遍历
func (bt *BinaryTree) LevelOrderTraversal() [][]interface{} {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    if bt.Root == nil {
        return [][]interface{}{}
    }
    
    result := make([][]interface{}, 0)
    queue := []*TreeNode{bt.Root}
    
    for len(queue) > 0 {
        level := make([]interface{}, 0)
        levelSize := len(queue)
        
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            
            level = append(level, node.Data)
            
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        
        result = append(result, level)
    }
    
    return result
}

// Height 计算树高度
func (bt *BinaryTree) Height() int {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    return bt.calculateHeight(bt.Root)
}

// calculateHeight 递归计算高度
func (bt *BinaryTree) calculateHeight(node *TreeNode) int {
    if node == nil {
        return -1
    }
    
    leftHeight := bt.calculateHeight(node.Left)
    rightHeight := bt.calculateHeight(node.Right)
    
    if leftHeight > rightHeight {
        return leftHeight + 1
    }
    return rightHeight + 1
}

// Size 计算树大小
func (bt *BinaryTree) Size() int {
    bt.mu.RLock()
    defer bt.mu.RUnlock()
    
    return bt.calculateSize(bt.Root)
}

// calculateSize 递归计算大小
func (bt *BinaryTree) calculateSize(node *TreeNode) int {
    if node == nil {
        return 0
    }
    
    return 1 + bt.calculateSize(node.Left) + bt.calculateSize(node.Right)
}
```

### 2.2 二叉搜索树 (Binary Search Tree)

```go
package bst

import (
    "fmt"
    "sync"
)

// BST 二叉搜索树
type BST struct {
    Root *TreeNode
    mu   sync.RWMutex
}

// NewBST 创建二叉搜索树
func NewBST() *BST {
    return &BST{Root: nil}
}

// Insert 插入节点
func (bst *BST) Insert(data interface{}) {
    bst.mu.Lock()
    defer bst.mu.Unlock()
    
    bst.Root = bst.insertNode(bst.Root, data)
}

// insertNode 递归插入节点
func (bst *BST) insertNode(node *TreeNode, data interface{}) *TreeNode {
    if node == nil {
        return &TreeNode{Data: data, Left: nil, Right: nil}
    }
    
    if compare(data, node.Data) <= 0 {
        node.Left = bst.insertNode(node.Left, data)
    } else {
        node.Right = bst.insertNode(node.Right, data)
    }
    
    return node
}

// Delete 删除节点
func (bst *BST) Delete(data interface{}) {
    bst.mu.Lock()
    defer bst.mu.Unlock()
    
    bst.Root = bst.deleteNode(bst.Root, data)
}

// deleteNode 递归删除节点
func (bst *BST) deleteNode(node *TreeNode, data interface{}) *TreeNode {
    if node == nil {
        return nil
    }
    
    if compare(data, node.Data) < 0 {
        node.Left = bst.deleteNode(node.Left, data)
    } else if compare(data, node.Data) > 0 {
        node.Right = bst.deleteNode(node.Right, data)
    } else {
        // 找到要删除的节点
        if node.Left == nil {
            return node.Right
        } else if node.Right == nil {
            return node.Left
        } else {
            // 有两个子节点，找到右子树的最小值
            minNode := bst.findMin(node.Right)
            node.Data = minNode.Data
            node.Right = bst.deleteNode(node.Right, minNode.Data)
        }
    }
    
    return node
}

// findMin 找到最小值节点
func (bst *BST) findMin(node *TreeNode) *TreeNode {
    for node.Left != nil {
        node = node.Left
    }
    return node
}

// findMax 找到最大值节点
func (bst *BST) findMax(node *TreeNode) *TreeNode {
    for node.Right != nil {
        node = node.Right
    }
    return node
}

// Min 获取最小值
func (bst *BST) Min() (interface{}, error) {
    bst.mu.RLock()
    defer bst.mu.RUnlock()
    
    if bst.Root == nil {
        return nil, fmt.Errorf("tree is empty")
    }
    
    minNode := bst.findMin(bst.Root)
    return minNode.Data, nil
}

// Max 获取最大值
func (bst *BST) Max() (interface{}, error) {
    bst.mu.RLock()
    defer bst.mu.RUnlock()
    
    if bst.Root == nil {
        return nil, fmt.Errorf("tree is empty")
    }
    
    maxNode := bst.findMax(bst.Root)
    return maxNode.Data, nil
}

// IsValid 验证是否为有效的BST
func (bst *BST) IsValid() bool {
    bst.mu.RLock()
    defer bst.mu.RUnlock()
    
    return bst.isValidBST(bst.Root, nil, nil)
}

// isValidBST 递归验证BST
func (bst *BST) isValidBST(node *TreeNode, min, max interface{}) bool {
    if node == nil {
        return true
    }
    
    if min != nil && compare(node.Data, min) <= 0 {
        return false
    }
    
    if max != nil && compare(node.Data, max) >= 0 {
        return false
    }
    
    return bst.isValidBST(node.Left, min, node.Data) &&
           bst.isValidBST(node.Right, node.Data, max)
}
```

---

## 3. 图数据结构 (Graph Data Structures)

### 3.1 邻接表图

```go
package graph

import (
    "fmt"
    "sync"
)

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

// AdjacencyListGraph 邻接表图
type AdjacencyListGraph struct {
    vertices map[int]*Vertex
    edges    map[int]map[int]*Edge
    directed bool
    mu       sync.RWMutex
}

// NewAdjacencyListGraph 创建邻接表图
func NewAdjacencyListGraph(directed bool) *AdjacencyListGraph {
    return &AdjacencyListGraph{
        vertices: make(map[int]*Vertex),
        edges:    make(map[int]map[int]*Edge),
        directed: directed,
    }
}

// AddVertex 添加顶点
func (g *AdjacencyListGraph) AddVertex(id int, data interface{}) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    g.vertices[id] = &Vertex{ID: id, Data: data}
    if g.edges[id] == nil {
        g.edges[id] = make(map[int]*Edge)
    }
}

// AddEdge 添加边
func (g *AdjacencyListGraph) AddEdge(from, to int, weight float64) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    if g.vertices[from] == nil || g.vertices[to] == nil {
        return
    }
    
    edge := &Edge{From: from, To: to, Weight: weight}
    g.edges[from][to] = edge
    
    if !g.directed {
        g.edges[to][from] = &Edge{From: to, To: from, Weight: weight}
    }
}

// RemoveVertex 删除顶点
func (g *AdjacencyListGraph) RemoveVertex(id int) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    delete(g.vertices, id)
    delete(g.edges, id)
    
    // 删除指向该顶点的边
    for _, adjList := range g.edges {
        delete(adjList, id)
    }
}

// RemoveEdge 删除边
func (g *AdjacencyListGraph) RemoveEdge(from, to int) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    if g.edges[from] != nil {
        delete(g.edges[from], to)
    }
    if !g.directed && g.edges[to] != nil {
        delete(g.edges[to], from)
    }
}

// GetVertex 获取顶点
func (g *AdjacencyListGraph) GetVertex(id int) (*Vertex, bool) {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    vertex, exists := g.vertices[id]
    return vertex, exists
}

// GetEdge 获取边
func (g *AdjacencyListGraph) GetEdge(from, to int) (*Edge, bool) {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    if g.edges[from] != nil {
        edge, exists := g.edges[from][to]
        return edge, exists
    }
    return nil, false
}

// Vertices 获取所有顶点
func (g *AdjacencyListGraph) Vertices() []*Vertex {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    vertices := make([]*Vertex, 0, len(g.vertices))
    for _, vertex := range g.vertices {
        vertices = append(vertices, vertex)
    }
    return vertices
}

// Edges 获取所有边
func (g *AdjacencyListGraph) Edges() []*Edge {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    edges := make([]*Edge, 0)
    for _, adjList := range g.edges {
        for _, edge := range adjList {
            edges = append(edges, edge)
        }
    }
    return edges
}

// Adjacent 获取邻接顶点
func (g *AdjacencyListGraph) Adjacent(vertexID int) []*Vertex {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
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

// IsDirected 是否为有向图
func (g *AdjacencyListGraph) IsDirected() bool {
    return g.directed
}

// Degree 计算度数
func (g *AdjacencyListGraph) Degree(vertexID int) (inDegree, outDegree int) {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
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

## 4. 散列表 (Hash Tables)

### 4.1 基本散列表

```go
package hashtable

import (
    "fmt"
    "hash/fnv"
    "sync"
)

// HashTable 散列表
type HashTable struct {
    buckets []*Bucket
    size    int
    mu      sync.RWMutex
}

// Bucket 桶
type Bucket struct {
    key   interface{}
    value interface{}
    next  *Bucket
}

// NewHashTable 创建散列表
func NewHashTable(size int) *HashTable {
    return &HashTable{
        buckets: make([]*Bucket, size),
        size:    size,
    }
}

// hash 计算哈希值
func (ht *HashTable) hash(key interface{}) int {
    h := fnv.New32a()
    fmt.Fprintf(h, "%v", key)
    return int(h.Sum32()) % ht.size
}

// Put 插入键值对
func (ht *HashTable) Put(key, value interface{}) {
    ht.mu.Lock()
    defer ht.mu.Unlock()
    
    index := ht.hash(key)
    
    // 检查是否已存在
    current := ht.buckets[index]
    for current != nil {
        if current.key == key {
            current.value = value
            return
        }
        current = current.next
    }
    
    // 插入新节点
    newBucket := &Bucket{key: key, value: value, next: ht.buckets[index]}
    ht.buckets[index] = newBucket
}

// Get 获取值
func (ht *HashTable) Get(key interface{}) (interface{}, bool) {
    ht.mu.RLock()
    defer ht.mu.RUnlock()
    
    index := ht.hash(key)
    current := ht.buckets[index]
    
    for current != nil {
        if current.key == key {
            return current.value, true
        }
        current = current.next
    }
    
    return nil, false
}

// Delete 删除键值对
func (ht *HashTable) Delete(key interface{}) bool {
    ht.mu.Lock()
    defer ht.mu.Unlock()
    
    index := ht.hash(key)
    current := ht.buckets[index]
    var prev *Bucket
    
    for current != nil {
        if current.key == key {
            if prev == nil {
                ht.buckets[index] = current.next
            } else {
                prev.next = current.next
            }
            return true
        }
        prev = current
        current = current.next
    }
    
    return false
}

// Contains 检查是否包含键
func (ht *HashTable) Contains(key interface{}) bool {
    _, exists := ht.Get(key)
    return exists
}

// Size 获取大小
func (ht *HashTable) Size() int {
    ht.mu.RLock()
    defer ht.mu.RUnlock()
    
    count := 0
    for _, bucket := range ht.buckets {
        current := bucket
        for current != nil {
            count++
            current = current.next
        }
    }
    return count
}

// Keys 获取所有键
func (ht *HashTable) Keys() []interface{} {
    ht.mu.RLock()
    defer ht.mu.RUnlock()
    
    keys := make([]interface{}, 0)
    for _, bucket := range ht.buckets {
        current := bucket
        for current != nil {
            keys = append(keys, current.key)
            current = current.next
        }
    }
    return keys
}

// Values 获取所有值
func (ht *HashTable) Values() []interface{} {
    ht.mu.RLock()
    defer ht.mu.RUnlock()
    
    values := make([]interface{}, 0)
    for _, bucket := range ht.buckets {
        current := bucket
        for current != nil {
            values = append(values, current.value)
            current = current.next
        }
    }
    return values
}
```

---

## 5. 高级数据结构 (Advanced Data Structures)

### 5.1 红黑树

```go
package rbtree

import (
    "fmt"
    "sync"
)

// Color 颜色
type Color bool

const (
    Red   Color = false
    Black Color = true
)

// RBNode 红黑树节点
type RBNode struct {
    Data  interface{}
    Color Color
    Left  *RBNode
    Right *RBNode
    Parent *RBNode
}

// RBTree 红黑树
type RBTree struct {
    Root *RBNode
    mu   sync.RWMutex
}

// NewRBTree 创建红黑树
func NewRBTree() *RBTree {
    return &RBTree{Root: nil}
}

// Insert 插入节点
func (rbt *RBTree) Insert(data interface{}) {
    rbt.mu.Lock()
    defer rbt.mu.Unlock()
    
    node := &RBNode{Data: data, Color: Red}
    rbt.Root = rbt.insertNode(rbt.Root, node)
    rbt.fixInsert(node)
}

// insertNode 插入节点
func (rbt *RBTree) insertNode(root, node *RBNode) *RBNode {
    if root == nil {
        return node
    }
    
    if compare(node.Data, root.Data) <= 0 {
        root.Left = rbt.insertNode(root.Left, node)
        root.Left.Parent = root
    } else {
        root.Right = rbt.insertNode(root.Right, node)
        root.Right.Parent = root
    }
    
    return root
}

// fixInsert 修复插入后的红黑树性质
func (rbt *RBTree) fixInsert(node *RBNode) {
    for node != rbt.Root && node.Parent.Color == Red {
        if node.Parent == node.Parent.Parent.Left {
            uncle := node.Parent.Parent.Right
            if uncle != nil && uncle.Color == Red {
                node.Parent.Color = Black
                uncle.Color = Black
                node.Parent.Parent.Color = Red
                node = node.Parent.Parent
            } else {
                if node == node.Parent.Right {
                    node = node.Parent
                    rbt.leftRotate(node)
                }
                node.Parent.Color = Black
                node.Parent.Parent.Color = Red
                rbt.rightRotate(node.Parent.Parent)
            }
        } else {
            uncle := node.Parent.Parent.Left
            if uncle != nil && uncle.Color == Red {
                node.Parent.Color = Black
                uncle.Color = Black
                node.Parent.Parent.Color = Red
                node = node.Parent.Parent
            } else {
                if node == node.Parent.Left {
                    node = node.Parent
                    rbt.rightRotate(node)
                }
                node.Parent.Color = Black
                node.Parent.Parent.Color = Red
                rbt.leftRotate(node.Parent.Parent)
            }
        }
    }
    
    rbt.Root.Color = Black
}

// leftRotate 左旋
func (rbt *RBTree) leftRotate(node *RBNode) {
    right := node.Right
    node.Right = right.Left
    if right.Left != nil {
        right.Left.Parent = node
    }
    right.Parent = node.Parent
    if node.Parent == nil {
        rbt.Root = right
    } else if node == node.Parent.Left {
        node.Parent.Left = right
    } else {
        node.Parent.Right = right
    }
    right.Left = node
    node.Parent = right
}

// rightRotate 右旋
func (rbt *RBTree) rightRotate(node *RBNode) {
    left := node.Left
    node.Left = left.Right
    if left.Right != nil {
        left.Right.Parent = node
    }
    left.Parent = node.Parent
    if node.Parent == nil {
        rbt.Root = left
    } else if node == node.Parent.Right {
        node.Parent.Right = left
    } else {
        node.Parent.Left = left
    }
    left.Right = node
    node.Parent = left
}

// Search 搜索节点
func (rbt *RBTree) Search(data interface{}) (*RBNode, bool) {
    rbt.mu.RLock()
    defer rbt.mu.RUnlock()
    
    return rbt.searchNode(rbt.Root, data)
}

// searchNode 递归搜索节点
func (rbt *RBTree) searchNode(node *RBNode, data interface{}) (*RBNode, bool) {
    if node == nil {
        return nil, false
    }
    
    if compare(data, node.Data) == 0 {
        return node, true
    }
    
    if compare(data, node.Data) < 0 {
        return rbt.searchNode(node.Left, data)
    }
    
    return rbt.searchNode(node.Right, data)
}

// InorderTraversal 中序遍历
func (rbt *RBTree) InorderTraversal() []interface{} {
    rbt.mu.RLock()
    defer rbt.mu.RUnlock()
    
    result := make([]interface{}, 0)
    rbt.inorderTraversal(rbt.Root, &result)
    return result
}

// inorderTraversal 递归中序遍历
func (rbt *RBTree) inorderTraversal(node *RBNode, result *[]interface{}) {
    if node == nil {
        return
    }
    
    rbt.inorderTraversal(node.Left, result)
    *result = append(*result, node.Data)
    rbt.inorderTraversal(node.Right, result)
}
```

---

## 总结

本文档介绍了基本数据结构的Go语言实现：

1. **线性数据结构** - 链表、栈、队列
2. **树形数据结构** - 二叉树、二叉搜索树、红黑树
3. **图数据结构** - 邻接表图
4. **散列表** - 基本散列表实现
5. **高级数据结构** - 红黑树等

这些数据结构为算法实现和程序开发提供了基础支持。

---

**相关链接**:

- [01-Hello World (Hello World)](01-Hello-World.md)
- [03-算法实现 (Algorithm Implementation)](03-Algorithm-Implementation.md)
- [04-并发编程 (Concurrent Programming)](04-Concurrent-Programming.md)
- [02-应用示例 (Application Examples)](../02-Application-Examples/README.md)
