# 01-自动机理论

 (Automata Theory)

## 目录

- [01-自动机理论](#01-自动机理论)
	- [目录](#目录)
	- [1. 自动机基础](#1-自动机基础)
		- [1.1 有限状态自动机](#11-有限状态自动机)
		- [1.2 下推自动机](#12-下推自动机)
		- [1.3 图灵机](#13-图灵机)
	- [2. 形式化定义](#2-形式化定义)
		- [2.1 DFA定义](#21-dfa定义)
		- [2.2 NFA定义](#22-nfa定义)
		- [2.3 PDA定义](#23-pda定义)
	- [3. Go语言实现](#3-go语言实现)
		- [3.1 DFA实现](#31-dfa实现)
		- [3.2 NFA实现](#32-nfa实现)
		- [3.3 PDA实现](#33-pda实现)
	- [4. 应用场景](#4-应用场景)
		- [4.1 词法分析](#41-词法分析)
		- [4.2 模式匹配](#42-模式匹配)
		- [4.3 协议验证](#43-协议验证)
	- [5. 数学证明](#5-数学证明)
		- [5.1 等价性定理](#51-等价性定理)
		- [5.2 最小化算法](#52-最小化算法)
		- [5.3 复杂度分析](#53-复杂度分析)

---

## 1. 自动机基础

### 1.1 有限状态自动机

有限状态自动机是计算理论的基础模型，用于描述具有有限内存的计算过程。

**定义 1.1**: 确定性有限自动机 (DFA) 是一个五元组 ```latex
$M = (Q, \Sigma, \delta, q_0, F)$
```，其中：

- ```latex
$Q$
``` 是有限状态集合
- ```latex
$\Sigma$
``` 是有限输入字母表
- ```latex
$\delta: Q \times \Sigma \rightarrow Q$
``` 是转移函数
- ```latex
$q_0 \in Q$
``` 是初始状态
- ```latex
$F \subseteq Q$
``` 是接受状态集合

### 1.2 下推自动机

下推自动机扩展了有限自动机，增加了栈作为辅助存储。

**定义 1.2**: 下推自动机 (PDA) 是一个七元组 ```latex
$M = (Q, \Sigma, \Gamma, \delta, q_0, Z_0, F)$
```，其中：

- ```latex
$Q$
``` 是有限状态集合
- ```latex
$\Sigma$
``` 是输入字母表
- ```latex
$\Gamma$
``` 是栈字母表
- ```latex
$\delta: Q \times (\Sigma \cup \{\epsilon\}) \times \Gamma \rightarrow 2^{Q \times \Gamma^*}$
``` 是转移函数
- ```latex
$q_0 \in Q$
``` 是初始状态
- ```latex
$Z_0 \in \Gamma$
``` 是初始栈符号
- ```latex
$F \subseteq Q$
``` 是接受状态集合

### 1.3 图灵机

图灵机是最通用的计算模型，能够模拟任何可计算函数。

**定义 1.3**: 图灵机是一个七元组 ```latex
$M = (Q, \Sigma, \Gamma, \delta, q_0, B, F)$
```，其中：

- ```latex
$Q$
``` 是有限状态集合
- ```latex
$\Sigma$
``` 是输入字母表
- ```latex
$\Gamma$
``` 是磁带字母表，```latex
$\Sigma \subseteq \Gamma$
```
- ```latex
$\delta: Q \times \Gamma \rightarrow Q \times \Gamma \times \{L, R\}$
``` 是转移函数
- ```latex
$q_0 \in Q$
``` 是初始状态
- ```latex
$B \in \Gamma \setminus \Sigma$
``` 是空白符号
- ```latex
$F \subseteq Q$
``` 是接受状态集合

## 2. 形式化定义

### 2.1 DFA定义

**定义 2.1**: DFA的扩展转移函数 ```latex
$\hat{\delta}: Q \times \Sigma^* \rightarrow Q$
``` 定义为：

```latex
$$
\begin{align}
\hat{\delta}(q, \epsilon) &= q \\
\hat{\delta}(q, wa) &= \delta(\hat{\delta}(q, w), a)
\end{align}
$$
```

其中 ```latex
$w \in \Sigma^*$
```, ```latex
$a \in \Sigma$
```。

**定义 2.2**: DFA接受的语言 ```latex
$L(M) = \{w \in \Sigma^* \mid \hat{\delta}(q_0, w) \in F\}$
```。

### 2.2 NFA定义

**定义 2.3**: 非确定性有限自动机 (NFA) 是一个五元组 ```latex
$M = (Q, \Sigma, \delta, q_0, F)$
```，其中：

- ```latex
$\delta: Q \times (\Sigma \cup \{\epsilon\}) \rightarrow 2^Q$
``` 是转移函数

**定义 2.4**: NFA的扩展转移函数 ```latex
$\hat{\delta}: 2^Q \times \Sigma^* \rightarrow 2^Q$
``` 定义为：

```latex
$$
\begin{align}
\hat{\delta}(S, \epsilon) &= S \\
\hat{\delta}(S, wa) &= \bigcup_{q \in \hat{\delta}(S, w)} \delta(q, a)
\end{align}
$$
```

### 2.3 PDA定义

**定义 2.5**: PDA的配置是一个三元组 ```latex
$(q, w, \gamma)$
```，其中：

- ```latex
$q \in Q$
``` 是当前状态
- ```latex
$w \in \Sigma^*$
``` 是剩余输入
- ```latex
$\gamma \in \Gamma^*$
``` 是栈内容

**定义 2.6**: PDA的转移关系 ```latex
$\vdash$
``` 定义为：
```latex
$(q, aw, Z\gamma) \vdash (p, w, \beta\gamma)$
``` 当且仅当 ```latex
$(p, \beta) \in \delta(q, a, Z)$
```。

## 3. Go语言实现

### 3.1 DFA实现

```go
package automatatheory

import (
 "fmt"
 "strings"
)

// State 表示自动机状态
type State string

// Symbol 表示输入符号
type Symbol rune

// DFA 确定性有限自动机
type DFA struct {
 States       map[State]bool
 Alphabet     map[Symbol]bool
 Transitions  map[State]map[Symbol]State
 InitialState State
 AcceptStates map[State]bool
}

// NewDFA 创建新的DFA
func NewDFA(initialState State) *DFA {
 return &DFA{
  States:       make(map[State]bool),
  Alphabet:     make(map[Symbol]bool),
  Transitions:  make(map[State]map[Symbol]State),
  InitialState: initialState,
  AcceptStates: make(map[State]bool),
 }
}

// AddState 添加状态
func (dfa *DFA) AddState(state State) {
 dfa.States[state] = true
}

// AddSymbol 添加输入符号
func (dfa *DFA) AddSymbol(symbol Symbol) {
 dfa.Alphabet[symbol] = true
}

// AddTransition 添加转移
func (dfa *DFA) AddTransition(from State, symbol Symbol, to State) {
 if dfa.Transitions[from] == nil {
  dfa.Transitions[from] = make(map[Symbol]State)
 }
 dfa.Transitions[from][symbol] = to
}

// AddAcceptState 添加接受状态
func (dfa *DFA) AddAcceptState(state State) {
 dfa.AcceptStates[state] = true
}

// Transition 执行转移
func (dfa *DFA) Transition(state State, symbol Symbol) (State, bool) {
 if transitions, exists := dfa.Transitions[state]; exists {
  if nextState, exists := transitions[symbol]; exists {
   return nextState, true
  }
 }
 return "", false
}

// ExtendedTransition 扩展转移函数
func (dfa *DFA) ExtendedTransition(state State, input string) (State, bool) {
 currentState := state

 for _, char := range input {
  symbol := Symbol(char)
  if nextState, exists := dfa.Transition(currentState, symbol); exists {
   currentState = nextState
  } else {
   return "", false
  }
 }

 return currentState, true
}

// Accepts 检查是否接受输入字符串
func (dfa *DFA) Accepts(input string) bool {
 finalState, exists := dfa.ExtendedTransition(dfa.InitialState, input)
 if !exists {
  return false
 }
 return dfa.AcceptStates[finalState]
}

// Minimize 最小化DFA
func (dfa *DFA) Minimize() *DFA {
 // 实现Hopcroft算法进行最小化
 equivalenceClasses := dfa.findEquivalenceClasses()
 return dfa.buildMinimizedDFA(equivalenceClasses)
}

// findEquivalenceClasses 找到等价类
func (dfa *DFA) findEquivalenceClasses() map[State]int {
 // 初始划分：接受状态和非接受状态
 partition := make(map[State]int)
 classID := 0

 // 为接受状态分配类ID
 for state := range dfa.AcceptStates {
  partition[state] = classID
 }
 classID++

 // 为非接受状态分配类ID
 for state := range dfa.States {
  if !dfa.AcceptStates[state] {
   partition[state] = classID
  }
 }

 // 迭代细化划分
 for {
  newPartition := dfa.refinePartition(partition)
  if dfa.partitionsEqual(partition, newPartition) {
   break
  }
  partition = newPartition
 }

 return partition
}

// refinePartition 细化划分
func (dfa *DFA) refinePartition(partition map[State]int) map[State]int {
 newPartition := make(map[State]int)
 classID := 0

 // 为每个等价类创建新的划分
 classes := make(map[int][]State)
 for state, class := range partition {
  classes[class] = append(classes[class], state)
 }

 for _, states := range classes {
  subClasses := dfa.splitClass(states, partition)
  for _, subClass := range subClasses {
   for _, state := range subClass {
    newPartition[state] = classID
   }
   classID++
  }
 }

 return newPartition
}

// splitClass 分割等价类
func (dfa *DFA) splitClass(states []State, partition map[State]int) [][]State {
 if len(states) <= 1 {
  return [][]State{states}
 }

 // 根据转移行为分割类
 behaviorMap := make(map[string][]State)

 for _, state := range states {
  behavior := dfa.getStateBehavior(state, partition)
  behaviorMap[behavior] = append(behaviorMap[behavior], state)
 }

 result := make([][]State, 0, len(behaviorMap))
 for _, subClass := range behaviorMap {
  result = append(result, subClass)
 }

 return result
}

// getStateBehavior 获取状态行为
func (dfa *DFA) getStateBehavior(state State, partition map[State]int) string {
 behavior := ""
 for symbol := range dfa.Alphabet {
  if nextState, exists := dfa.Transition(state, symbol); exists {
   behavior += fmt.Sprintf("%c->%d,", symbol, partition[nextState])
  } else {
   behavior += fmt.Sprintf("%c->-1,", symbol)
  }
 }
 return behavior
}

// partitionsEqual 检查划分是否相等
func (dfa *DFA) partitionsEqual(p1, p2 map[State]int) bool {
 if len(p1) != len(p2) {
  return false
 }

 for state, class1 := range p1 {
  if class2, exists := p2[state]; !exists || class1 != class2 {
   return false
  }
 }

 return true
}

// buildMinimizedDFA 构建最小化DFA
func (dfa *DFA) buildMinimizedDFA(equivalenceClasses map[State]int) *DFA {
 // 创建新的最小化DFA
 minimized := NewDFA(State(fmt.Sprintf("q%d", equivalenceClasses[dfa.InitialState])))

 // 添加状态
 classStates := make(map[int]State)
 for state, class := range equivalenceClasses {
  if _, exists := classStates[class]; !exists {
   classState := State(fmt.Sprintf("q%d", class))
   classStates[class] = classState
   minimized.AddState(classState)

   // 如果是接受状态，标记为接受
   if dfa.AcceptStates[state] {
    minimized.AddAcceptState(classState)
   }
  }
 }

 // 添加转移
 for state, class := range equivalenceClasses {
  classState := classStates[class]
  for symbol := range dfa.Alphabet {
   if nextState, exists := dfa.Transition(state, symbol); exists {
    nextClass := equivalenceClasses[nextState]
    nextClassState := classStates[nextClass]
    minimized.AddTransition(classState, symbol, nextClassState)
   }
  }
 }

 // 复制字母表
 for symbol := range dfa.Alphabet {
  minimized.AddSymbol(symbol)
 }

 return minimized
}
```

### 3.2 NFA实现

```go
// NFA 非确定性有限自动机
type NFA struct {
 States       map[State]bool
 Alphabet     map[Symbol]bool
 Transitions  map[State]map[Symbol][]State
 InitialState State
 AcceptStates map[State]bool
}

// NewNFA 创建新的NFA
func NewNFA(initialState State) *NFA {
 return &NFA{
  States:       make(map[State]bool),
  Alphabet:     make(map[Symbol]bool),
  Transitions:  make(map[State]map[Symbol][]State),
  InitialState: initialState,
  AcceptStates: make(map[State]bool),
 }
}

// AddState 添加状态
func (nfa *NFA) AddState(state State) {
 nfa.States[state] = true
}

// AddSymbol 添加输入符号
func (nfa *NFA) AddSymbol(symbol Symbol) {
 nfa.Alphabet[symbol] = true
}

// AddTransition 添加转移
func (nfa *NFA) AddTransition(from State, symbol Symbol, to State) {
 if nfa.Transitions[from] == nil {
  nfa.Transitions[from] = make(map[Symbol][]State)
 }
 if nfa.Transitions[from][symbol] == nil {
  nfa.Transitions[from][symbol] = make([]State, 0)
 }
 nfa.Transitions[from][symbol] = append(nfa.Transitions[from][symbol], to)
}

// AddAcceptState 添加接受状态
func (nfa *NFA) AddAcceptState(state State) {
 nfa.AcceptStates[state] = true
}

// Transition 执行转移
func (nfa *NFA) Transition(state State, symbol Symbol) []State {
 if transitions, exists := nfa.Transitions[state]; exists {
  if nextStates, exists := transitions[symbol]; exists {
   return nextStates
  }
 }
 return []State{}
}

// EpsilonClosure 计算ε闭包
func (nfa *NFA) EpsilonClosure(states []State) []State {
 closure := make(map[State]bool)
 stack := make([]State, 0)

 // 初始化栈
 for _, state := range states {
  closure[state] = true
  stack = append(stack, state)
 }

 // 计算闭包
 for len(stack) > 0 {
  current := stack[len(stack)-1]
  stack = stack[:len(stack)-1]
  
  // 添加ε转移
  for _, nextState := range nfa.Transition(current, Symbol(0)) { // ε用0表示
   if !closure[nextState] {
    closure[nextState] = true
    stack = append(stack, nextState)
   }
  }
 }

 // 转换为切片
 result := make([]State, 0, len(closure))
 for state := range closure {
  result = append(result, state)
 }

 return result
}

// ExtendedTransition 扩展转移函数
func (nfa *NFA) ExtendedTransition(states []State, input string) []State {
 currentStates := states

 for _, char := range input {
  symbol := Symbol(char)
  nextStates := make([]State, 0)
  
  // 计算ε闭包
  currentStates = nfa.EpsilonClosure(currentStates)
  
  // 执行转移
  for _, state := range currentStates {
   nextStates = append(nextStates, nfa.Transition(state, symbol)...)
  }
  
  currentStates = nextStates
 }

 // 最终ε闭包
 return nfa.EpsilonClosure(currentStates)
}

// Accepts 检查是否接受输入字符串
func (nfa *NFA) Accepts(input string) bool {
 finalStates := nfa.ExtendedTransition([]State{nfa.InitialState}, input)

 for _, state := range finalStates {
  if nfa.AcceptStates[state] {
   return true
  }
 }

 return false
}

// ToDFA 转换为DFA
func (nfa *NFA) ToDFA() *DFA {
 dfa := NewDFA("q0")

 // 计算初始状态的ε闭包
 initialStates := nfa.EpsilonClosure([]State{nfa.InitialState})
 stateMap := make(map[string][]State)
 stateMap["q0"] = initialStates

 // 添加初始状态
 dfa.AddState("q0")

 // 检查初始状态是否为接受状态
 for _, state := range initialStates {
  if nfa.AcceptStates[state] {
   dfa.AddAcceptState("q0")
   break
  }
 }

 // 处理所有状态
 processed := make(map[string]bool)
 queue := []string{"q0"}

 for len(queue) > 0 {
  currentDFAState := queue[0]
  queue = queue[1:]
  
  if processed[currentDFAState] {
   continue
  }
  processed[currentDFAState] = true
  
  nfaStates := stateMap[currentDFAState]
  
  // 对每个输入符号计算转移
  for symbol := range nfa.Alphabet {
   if symbol == Symbol(0) { // 跳过ε
    continue
   }

   nextNFAStates := make([]State, 0)
   for _, nfaState := range nfaStates {
    nextNFAStates = append(nextNFAStates, nfa.Transition(nfaState, symbol)...)
   }

   if len(nextNFAStates) > 0 {
    nextNFAStates = nfa.EpsilonClosure(nextNFAStates)

    // 创建新的DFA状态
    nextDFAState := dfa.createStateName(nextNFAStates)
    stateMap[nextDFAState] = nextNFAStates

    // 添加状态和转移
    dfa.AddState(nextDFAState)
    dfa.AddTransition(currentDFAState, symbol, nextDFAState)

    // 检查是否为接受状态
    for _, state := range nextNFAStates {
     if nfa.AcceptStates[state] {
      dfa.AddAcceptState(nextDFAState)
      break
     }
    }

    // 添加到队列
    if !processed[nextDFAState] {
     queue = append(queue, nextDFAState)
    }
   }
  }
 }

 // 复制字母表
 for symbol := range nfa.Alphabet {
  if symbol != Symbol(0) { // 排除ε
   dfa.AddSymbol(symbol)
  }
 }

 return dfa
}

// createStateName 创建状态名称
func (nfa *NFA) createStateName(states []State) string {
 // 简化的状态命名，实际实现可能需要更复杂的命名策略
 return fmt.Sprintf("q%d", len(states))
}
```

### 3.3 PDA实现

```go
// StackSymbol 栈符号
type StackSymbol string

// PDATransition PDA转移
type PDATransition struct {
 FromState State
 Input     Symbol
 StackTop  StackSymbol
 ToState   State
 StackPush []StackSymbol
}

// PDA 下推自动机
type PDA struct {
 States       map[State]bool
 InputAlphabet map[Symbol]bool
 StackAlphabet map[StackSymbol]bool
 Transitions  []PDATransition
 InitialState State
 InitialStack StackSymbol
 AcceptStates map[State]bool
}

// NewPDA 创建新的PDA
func NewPDA(initialState State, initialStack StackSymbol) *PDA {
 return &PDA{
  States:        make(map[State]bool),
  InputAlphabet: make(map[Symbol]bool),
  StackAlphabet: make(map[StackSymbol]bool),
  Transitions:   make([]PDATransition, 0),
  InitialState:  initialState,
  InitialStack:  initialStack,
  AcceptStates:  make(map[State]bool),
 }
}

// AddState 添加状态
func (pda *PDA) AddState(state State) {
 pda.States[state] = true
}

// AddInputSymbol 添加输入符号
func (pda *PDA) AddInputSymbol(symbol Symbol) {
 pda.InputAlphabet[symbol] = true
}

// AddStackSymbol 添加栈符号
func (pda *PDA) AddStackSymbol(symbol StackSymbol) {
 pda.StackAlphabet[symbol] = true
}

// AddTransition 添加转移
func (pda *PDA) AddTransition(from State, input Symbol, stackTop StackSymbol, to State, stackPush []StackSymbol) {
 transition := PDATransition{
  FromState: from,
  Input:     input,
  StackTop:  stackTop,
  ToState:   to,
  StackPush: stackPush,
 }
 pda.Transitions = append(pda.Transitions, transition)
}

// AddAcceptState 添加接受状态
func (pda *PDA) AddAcceptState(state State) {
 pda.AcceptStates[state] = true
}

// PDAConfiguration PDA配置
type PDAConfiguration struct {
 State     State
 Input     string
 Stack     []StackSymbol
 InputPos  int
}

// NewPDAConfiguration 创建PDA配置
func NewPDAConfiguration(state State, input string) *PDAConfiguration {
 return &PDAConfiguration{
  State:    state,
  Input:    input,
  Stack:    []StackSymbol{pda.InitialStack},
  InputPos: 0,
 }
}

// Step 执行一步转移
func (pda *PDA) Step(config *PDAConfiguration) []*PDAConfiguration {
 configs := make([]*PDAConfiguration, 0)

 for _, transition := range pda.Transitions {
  if transition.FromState == config.State {
   // 检查输入符号
   if config.InputPos < len(config.Input) {
    if transition.Input != Symbol(config.Input[config.InputPos]) && transition.Input != Symbol(0) {
     continue
    }
   } else if transition.Input != Symbol(0) {
    continue
   }

   // 检查栈顶
   if len(config.Stack) == 0 || config.Stack[len(config.Stack)-1] != transition.StackTop {
    continue
   }

   // 创建新配置
   newConfig := &PDAConfiguration{
    State: transition.ToState,
    Input: config.Input,
    Stack: make([]StackSymbol, len(config.Stack)),
    InputPos: config.InputPos,
   }
   copy(newConfig.Stack, config.Stack)

   // 更新栈
   newConfig.Stack = newConfig.Stack[:len(newConfig.Stack)-1] // 弹出栈顶
   for i := len(transition.StackPush) - 1; i >= 0; i-- {
    newConfig.Stack = append(newConfig.Stack, transition.StackPush[i])
   }

   // 更新输入位置
   if transition.Input != Symbol(0) {
    newConfig.InputPos++
   }

   configs = append(configs, newConfig)
  }
 }

 return configs
}

// Accepts 检查是否接受输入字符串
func (pda *PDA) Accepts(input string) bool {
 initialConfig := NewPDAConfiguration(pda.InitialState, input)
 configs := []*PDAConfiguration{initialConfig}

 for len(configs) > 0 {
  newConfigs := make([]*PDAConfiguration, 0)
  
  for _, config := range configs {
   // 检查是否接受
   if config.InputPos == len(config.Input) && pda.AcceptStates[config.State] {
    return true
   }

   // 执行转移
   nextConfigs := pda.Step(config)
   newConfigs = append(newConfigs, nextConfigs...)
  }
  
  configs = newConfigs
 }

 return false
}
```

## 4. 应用场景

### 4.1 词法分析

```go
// LexicalAnalyzer 词法分析器
type LexicalAnalyzer struct {
 dfa *DFA
 tokens map[string]string
}

// NewLexicalAnalyzer 创建词法分析器
func NewLexicalAnalyzer() *LexicalAnalyzer {
 return &LexicalAnalyzer{
  dfa:    NewDFA("start"),
  tokens: make(map[string]string),
 }
}

// BuildIdentifierDFA 构建标识符DFA
func (la *LexicalAnalyzer) BuildIdentifierDFA() {
 // 添加状态
 la.dfa.AddState("start")
 la.dfa.AddState("identifier")
 la.dfa.AddState("error")

 // 添加符号
 for c := 'a'; c <= 'z'; c++ {
  la.dfa.AddSymbol(Symbol(c))
 }
 for c := 'A'; c <= 'Z'; c++ {
  la.dfa.AddSymbol(Symbol(c))
 }
 for c := '0'; c <= '9'; c++ {
  la.dfa.AddSymbol(Symbol(c))
 }
 la.dfa.AddSymbol('_')

 // 添加转移
 // 字母或下划线 -> identifier
 for c := 'a'; c <= 'z'; c++ {
  la.dfa.AddTransition("start", Symbol(c), "identifier")
 }
 for c := 'A'; c <= 'Z'; c++ {
  la.dfa.AddTransition("start", Symbol(c), "identifier")
 }
 la.dfa.AddTransition("start", '_', "identifier")

 // identifier -> identifier (字母、数字、下划线)
 for c := 'a'; c <= 'z'; c++ {
  la.dfa.AddTransition("identifier", Symbol(c), "identifier")
 }
 for c := 'A'; c <= 'Z'; c++ {
  la.dfa.AddTransition("identifier", Symbol(c), "identifier")
 }
 for c := '0'; c <= '9'; c++ {
  la.dfa.AddTransition("identifier", Symbol(c), "identifier")
 }
 la.dfa.AddTransition("identifier", '_', "identifier")

 // 其他符号 -> error
 la.dfa.AddTransition("start", Symbol(0), "error")
 la.dfa.AddTransition("identifier", Symbol(0), "error")

 // 设置接受状态
 la.dfa.AddAcceptState("identifier")
}

// Tokenize 词法分析
func (la *LexicalAnalyzer) Tokenize(input string) []string {
 tokens := make([]string, 0)
 current := ""

 for i, char := range input {
  test := current + string(char)
  if la.dfa.Accepts(test) {
   current = test
  } else {
   if current != "" {
    tokens = append(tokens, current)
    current = ""
   }
   if char != ' ' && char != '\t' && char != '\n' {
    current = string(char)
   }
  }
 }

 if current != "" {
  tokens = append(tokens, current)
 }

 return tokens
}
```

### 4.2 模式匹配

```go
// PatternMatcher 模式匹配器
type PatternMatcher struct {
 dfa *DFA
}

// NewPatternMatcher 创建模式匹配器
func NewPatternMatcher(pattern string) *PatternMatcher {
 pm := &PatternMatcher{
  dfa: NewDFA("start"),
 }
 pm.buildPatternDFA(pattern)
 return pm
}

// buildPatternDFA 构建模式DFA
func (pm *PatternMatcher) buildPatternDFA(pattern string) {
 // 为每个位置创建状态
 for i := 0; i <= len(pattern); i++ {
  stateName := fmt.Sprintf("q%d", i)
  pm.dfa.AddState(State(stateName))
  
  if i == len(pattern) {
   pm.dfa.AddAcceptState(State(stateName))
  }
 }

 // 添加转移
 for i := 0; i < len(pattern); i++ {
  currentState := State(fmt.Sprintf("q%d", i))
  nextState := State(fmt.Sprintf("q%d", i+1))
  symbol := Symbol(pattern[i])
  
  pm.dfa.AddSymbol(symbol)
  pm.dfa.AddTransition(currentState, symbol, nextState)
 }
}

// FindAll 查找所有匹配
func (pm *PatternMatcher) FindAll(text string) []int {
 matches := make([]int, 0)

 for i := 0; i < len(text); i++ {
  if pm.dfa.Accepts(text[i:]) {
   matches = append(matches, i)
  }
 }

 return matches
}

// FindFirst 查找第一个匹配
func (pm *PatternMatcher) FindFirst(text string) int {
 for i := 0; i < len(text); i++ {
  if pm.dfa.Accepts(text[i:]) {
   return i
  }
 }
 return -1
}
```

### 4.3 协议验证

```go
// ProtocolVerifier 协议验证器
type ProtocolVerifier struct {
 pda *PDA
}

// NewProtocolVerifier 创建协议验证器
func NewProtocolVerifier() *ProtocolVerifier {
 pv := &ProtocolVerifier{
  pda: NewPDA("idle", "Z0"),
 }
 pv.buildProtocolPDA()
 return pv
}

// buildProtocolPDA 构建协议PDA
func (pv *ProtocolVerifier) buildProtocolPDA() {
 // 添加状态
 pv.pda.AddState("idle")
 pv.pda.AddState("waiting")
 pv.pda.AddState("processing")
 pv.pda.AddState("error")

 // 添加符号
 pv.pda.AddInputSymbol('R') // Request
 pv.pda.AddInputSymbol('A') // Acknowledge
 pv.pda.AddInputSymbol('D') // Data
 pv.pda.AddInputSymbol('E') // End

 pv.pda.AddStackSymbol("Z0")
 pv.pda.AddStackSymbol("R")
 pv.pda.AddStackSymbol("D")

 // 添加转移
 // idle -> waiting (收到请求)
 pv.pda.AddTransition("idle", 'R', "Z0", "waiting", []StackSymbol{"R", "Z0"})

 // waiting -> processing (收到数据)
 pv.pda.AddTransition("waiting", 'D', "R", "processing", []StackSymbol{"D", "R"})

 // processing -> waiting (发送确认)
 pv.pda.AddTransition("processing", 'A', "D", "waiting", []StackSymbol{"R"})

 // waiting -> idle (结束)
 pv.pda.AddTransition("waiting", 'E', "R", "idle", []StackSymbol{"Z0"})

 // 错误处理
 pv.pda.AddTransition("idle", 'A', "Z0", "error", []StackSymbol{"Z0"})
 pv.pda.AddTransition("waiting", 'E', "Z0", "error", []StackSymbol{"Z0"})

 // 设置接受状态
 pv.pda.AddAcceptState("idle")
}

// VerifyProtocol 验证协议
func (pv *ProtocolVerifier) VerifyProtocol(trace string) bool {
 return pv.pda.Accepts(trace)
}

// GenerateTestCases 生成测试用例
func (pv *ProtocolVerifier) GenerateTestCases() []string {
 testCases := []string{
  "RDAE",     // 正常流程
  "RDAE",     // 重复正常流程
  "RA",       // 错误：没有数据就确认
  "RDE",      // 错误：没有确认就结束
  "RDA",      // 错误：没有结束
 }
 return testCases
}
```

## 5. 数学证明

### 5.1 等价性定理

**定理 5.1** (NFA与DFA等价性): 对任意NFA ```latex
$M$
```，存在等价的DFA ```latex
$M'$
```，使得 ```latex
$L(M) = L(M')$
```。

**证明**:
1. 构造DFA ```latex
$M'$
```，其状态是NFA ```latex
$M$
``` 的状态集合的幂集
2. 初始状态是NFA初始状态的ε闭包
3. 转移函数定义为：```latex
$\delta'(S, a) = \bigcup_{q \in S} \delta(q, a)$
``` 的ε闭包
4. 接受状态是包含NFA接受状态的状态集合
5. 证明 ```latex
$L(M) = L(M')$
```

**定理 5.2** (PDA与上下文无关文法等价性): 对任意PDA ```latex
$M$
```，存在等价的上下文无关文法 ```latex
$G$
```，使得 ```latex
$L(M) = L(G)$
```。

**证明**:
1. 构造文法 ```latex
$G$
```，其变元表示PDA的配置
2. 产生式对应PDA的转移
3. 证明 ```latex
$L(M) = L(G)$
```

### 5.2 最小化算法

**定理 5.3** (DFA最小化): 对任意DFA ```latex
$M$
```，存在唯一的最小DFA ```latex
$M'$
```，使得 ```latex
$L(M) = L(M')$
``` 且 ```latex
$M'$
``` 的状态数最少。

**证明**:
1. 使用等价类划分算法
2. 初始划分：接受状态和非接受状态
3. 迭代细化：根据转移行为分割等价类
4. 证明算法终止且结果唯一

### 5.3 复杂度分析

**定理 5.4**:
- DFA最小化的时间复杂度为 ```latex
$O(n \log n)$
```
- NFA到DFA转换的最坏情况时间复杂度为 ```latex
$O(2^n)$
```
- PDA接受性问题是可判定的

**证明**:
1. DFA最小化：使用Hopcroft算法
2. NFA到DFA：状态数可能指数增长
3. PDA接受性：使用配置图搜索

---

**总结**: 自动机理论为计算理论提供了基础模型，通过Go语言实现，我们可以构建实用的自动机框架，用于词法分析、模式匹配和协议验证等应用。
