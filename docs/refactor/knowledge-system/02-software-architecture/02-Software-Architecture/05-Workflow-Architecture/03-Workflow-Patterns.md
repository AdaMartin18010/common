# 工作流模式 (Workflow Patterns)

## 目录
- [工作流模式 (Workflow Patterns)](#工作流模式-workflow-patterns)
  - [目录](#目录)
  - [1. 工作流模式理论基础](#1-工作流模式理论基础)
    - [1.1 模式分类体系](#11-模式分类体系)
    - [1.2 模式形式化定义](#12-模式形式化定义)
  - [2. 基础工作流模式](#2-基础工作流模式)
    - [2.1 顺序执行模式 (Sequence)](#21-顺序执行模式-sequence)
    - [2.2 并行执行模式 (Parallel Split/Join)](#22-并行执行模式-parallel-splitjoin)
  - [3. 高级工作流模式](#3-高级工作流模式)
    - [3.1 选择模式 (Exclusive Choice)](#31-选择模式-exclusive-choice)
    - [3.2 循环模式 (Arbitrary Cycles)](#32-循环模式-arbitrary-cycles)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 模式接口定义](#41-模式接口定义)
    - [4.2 组合模式实现](#42-组合模式实现)

---

## 1. 工作流模式理论基础

### 1.1 模式分类体系

**定义 1.1 (工作流模式)**:
工作流模式是一个三元组 $\mathcal{P} = (S, T, C)$，其中：

-   $S$ 是结构模式集合
-   $T$ 是任务模式集合
-   $C$ 是控制模式集合

**模式分类**:
工作流模式可以大致分为以下几类：
$$
\text{WorkflowPatterns} = \text{Sequential} \cup \text{Parallel} \cup \text{Choice} \cup \text{Iteration} \cup \text{Exception}
$$

### 1.2 模式形式化定义

**定义 1.2 (顺序模式)**:
顺序模式是一个二元组 $\text{Seq}(t_1, t_2)$，满足：
$$
\text{Seq}(t_1, t_2) = \{ s_1 \to s_2 \to s_3 \mid s_1 \xrightarrow{t_1} s_2 \xrightarrow{t_2} s_3 \}
$$

**定义 1.3 (并行模式)**:
并行模式是一个二元组 $\text{Par}(t_1, t_2)$，满足：
$$
\text{Par}(t_1, t_2) = \{ s_1 \to s_2 \mid s_1 \xrightarrow{t_1 \parallel t_2} s_2 \}
$$

**定义 1.4 (选择模式)**:
选择模式是一个三元组 $\text{Choice}(c, t_1, t_2)$，满足：
$$
\text{Choice}(c, t_1, t_2) = \{ s_1 \to s_2 \mid s_1 \xrightarrow{\text{if } c \text{ then } t_1 \text{ else } t_2} s_2 \}
$$

---

## 2. 基础工作流模式

### 2.1 顺序执行模式 (Sequence)

```go
package patterns

// SequentialPattern 顺序执行模式
type SequentialPattern struct {
    tasks []Task
}

// NewSequentialPattern 创建顺序执行模式
func NewSequentialPattern(tasks ...Task) *SequentialPattern {
    return &SequentialPattern{tasks: tasks}
}

// Execute 执行顺序模式
func (p *SequentialPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    currentInput := input
    for _, task := range p.tasks {
        result, err := task.Execute(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("task %s failed: %w", task.ID(), err)
        }
        currentInput = mergeContexts(currentInput, result)
    }
    return currentInput, nil
}
```

### 2.2 并行执行模式 (Parallel Split/Join)

```go
package patterns

import "sync"

// ParallelPattern 并行执行模式
type ParallelPattern struct {
    tasks []Task
}

// NewParallelPattern 创建并行执行模式
func NewParallelPattern(tasks ...Task) *ParallelPattern {
    return &ParallelPattern{tasks: tasks}
}

// Execute 执行并行模式
func (p *ParallelPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    var wg sync.WaitGroup
    resultsChan := make(chan map[string]interface{}, len(p.tasks))
    errChan := make(chan error, len(p.tasks))

    for _, task := range p.tasks {
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()
            result, err := t.Execute(ctx, input)
            if err != nil {
                errChan <- err
                return
            }
            resultsChan <- result
        }(task)
    }

    wg.Wait()
    close(resultsChan)
    close(errChan)

    if len(errChan) > 0 {
        return nil, <-errChan // 返回第一个错误
    }

    finalResult := make(map[string]interface{})
    for result := range resultsChan {
        finalResult = mergeContexts(finalResult, result)
    }

    return finalResult, nil
}
```

---

## 3. 高级工作流模式

### 3.1 选择模式 (Exclusive Choice)

```go
package patterns

// ChoicePattern 选择模式
type ChoicePattern struct {
    branches map[string]Task // 条件 -> 任务
    defaultBranch Task
}

// NewChoicePattern 创建选择模式
func NewChoicePattern() *ChoicePattern {
    return &ChoicePattern{branches: make(map[string]Task)}
}

// AddBranch 添加分支
func (p *ChoicePattern) AddBranch(condition string, task Task) {
    p.branches[condition] = task
}

// SetDefault 设置默认分支
func (p *ChoicePattern) SetDefault(task Task) {
    p.defaultBranch = task
}

// Execute 执行选择模式
func (p *ChoicePattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    condition, _ := input["_condition"].(string)
    if task, ok := p.branches[condition]; ok {
        return task.Execute(ctx, input)
    }
    if p.defaultBranch != nil {
        return p.defaultBranch.Execute(ctx, input)
    }
    return nil, fmt.Errorf("no branch matched for condition: %s", condition)
}
```

### 3.2 循环模式 (Arbitrary Cycles)

```go
package patterns

// LoopPattern 循环模式
type LoopPattern struct {
    body Task
    exitCondition func(context map[string]interface{}) bool
    maxIterations int
}

// NewLoopPattern 创建循环模式
func NewLoopPattern(body Task, exitCondition func(map[string]interface{}) bool, max int) *LoopPattern {
    return &LoopPattern{body: body, exitCondition: exitCondition, maxIterations: max}
}

// Execute 执行循环模式
func (p *LoopPattern) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
    currentInput := input
    var err error
    
    for i := 0; i < p.maxIterations; i++ {
        if p.exitCondition(currentInput) {
            break
        }
        currentInput, err = p.body.Execute(ctx, currentInput)
        if err != nil {
            return nil, fmt.Errorf("loop body failed at iteration %d: %w", i, err)
        }
    }
    return currentInput, nil
}
```

---

## 4. Go语言实现

### 4.1 模式接口定义

```go
package patterns

import "context"

// WorkflowPattern 定义了所有工作流模式的基础接口
type WorkflowPattern interface {
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}

// Task 定义了可执行的任务单元
type Task interface {
    WorkflowPattern // 任务本身也是一种模式
    ID() string
}
```

### 4.2 组合模式实现

使用组合设计模式，可以将简单模式组合成复杂的工作流。

```go
func main() {
    // 定义基础任务
    taskA := NewSimpleTask("A", ...)
    taskB := NewSimpleTask("B", ...)
    taskC := NewSimpleTask("C", ...)
    taskD := NewSimpleTask("D", ...)

    // 组合模式
    // 并行执行 A 和 B，然后顺序执行 C，最后是 D
    parallelAB := NewParallelPattern(taskA, taskB)
    sequence := NewSequentialPattern(parallelAB, taskC, taskD)

    // 执行整个工作流
    initialContext := make(map[string]interface{})
    finalResult, err := sequence.Execute(context.Background(), initialContext)
    
    // ...
}
``` 