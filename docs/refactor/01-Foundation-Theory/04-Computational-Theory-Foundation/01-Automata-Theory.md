# 01-自动机理论 (Automata Theory)

## 目录

- [01-自动机理论 (Automata Theory)](#01-自动机理论-automata-theory)
  - [目录](#目录)
  - [1. 基本概念](#1-基本概念)
    - [1.1 有限自动机](#11-有限自动机)
    - [1.2 下推自动机](#12-下推自动机)
    - [1.3 图灵机](#13-图灵机)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 DFA定义](#21-dfa定义)
    - [2.2 NFA定义](#22-nfa定义)
    - [2.3 PDA定义](#23-pda定义)
  - [3. 重要定理](#3-重要定理)
    - [3.1 泵引理](#31-泵引理)
    - [3.2 最小化定理](#32-最小化定理)
    - [3.3 等价性定理](#33-等价性定理)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 自动机数据结构](#41-自动机数据结构)
    - [4.2 NFA实现](#42-nfa实现)
    - [4.3 PDA实现](#43-pda实现)
  - [5. 应用示例](#5-应用示例)
    - [5.1 词法分析器](#51-词法分析器)
    - [5.2 模式匹配](#52-模式匹配)
    - [5.3 状态机](#53-状态机)
  - [总结](#总结)

## 1. 基本概念

### 1.1 有限自动机

**定义 1.1**: 有限自动机 (Finite Automaton)

**有限自动机** $M = (Q, \Sigma, \delta, q_0, F)$ 由以下部分组成：

1. **状态集** $Q$: 有限状态集合
2. **字母表** $\Sigma$: 输入符号的有限集合
3. **转移函数** $\delta$: $Q \times \Sigma \rightarrow Q$ (DFA) 或 $Q \times \Sigma \rightarrow 2^Q$ (NFA)
4. **初始状态** $q_0 \in Q$
5. **接受状态集** $F \subseteq Q$

**定义 1.2**: 自动机的类型

1. **确定性有限自动机 (DFA)**: 转移函数是确定性的
2. **非确定性有限自动机 (NFA)**: 转移函数是非确定性的
3. **ε-NFA**: 允许ε转移的非确定性自动机

### 1.2 下推自动机

**定义 1.3**: 下推自动机 (Pushdown Automaton)

**下推自动机** $P = (Q, \Sigma, \Gamma, \delta, q_0, Z_0, F)$ 由以下部分组成：

1. **状态集** $Q$: 有限状态集合
2. **输入字母表** $\Sigma$: 输入符号的有限集合
3. **栈字母表** $\Gamma$: 栈符号的有限集合
4. **转移函数** $\delta$: $Q \times (\Sigma \cup \{\varepsilon\}) \times \Gamma \rightarrow 2^{Q \times \Gamma^*}$
5. **初始状态** $q_0 \in Q$
6. **初始栈符号** $Z_0 \in \Gamma$
7. **接受状态集** $F \subseteq Q$

### 1.3 图灵机

**定义 1.4**: 图灵机 (Turing Machine)

**图灵机** $T = (Q, \Sigma, \Gamma, \delta, q_0, B, F)$ 由以下部分组成：

1. **状态集** $Q$: 有限状态集合
2. **输入字母表** $\Sigma$: 输入符号的有限集合
3. **带字母表** $\Gamma$: 带符号的有限集合，$\Sigma \subseteq \Gamma$
4. **转移函数** $\delta$: $Q \times \Gamma \rightarrow Q \times \Gamma \times \{L, R, N\}$
5. **初始状态** $q_0 \in Q$
6. **空白符号** $B \in \Gamma \setminus \Sigma$
7. **接受状态集** $F \subseteq Q$

## 2. 形式化定义

### 2.1 DFA定义

**定义 2.1**: 确定性有限自动机

**DFA** $M = (Q, \Sigma, \delta, q_0, F)$ 的转移函数 $\delta: Q \times \Sigma \rightarrow Q$ 满足：

对于任意状态 $q \in Q$ 和输入符号 $a \in \Sigma$，$\delta(q, a)$ 是唯一确定的状态。

**定义 2.2**: DFA的扩展转移函数

**扩展转移函数** $\hat{\delta}: Q \times \Sigma^* \rightarrow Q$ 递归定义：

1. **基础**: $\hat{\delta}(q, \varepsilon) = q$
2. **归纳**: $\hat{\delta}(q, wa) = \delta(\hat{\delta}(q, w), a)$

其中 $w \in \Sigma^*$，$a \in \Sigma$。

**定义 2.3**: DFA的语言

DFA $M$ 接受的语言：
$$L(M) = \{w \in \Sigma^* \mid \hat{\delta}(q_0, w) \in F\}$$

### 2.2 NFA定义

**定义 2.4**: 非确定性有限自动机

**NFA** $M = (Q, \Sigma, \delta, q_0, F)$ 的转移函数 $\delta: Q \times \Sigma \rightarrow 2^Q$ 满足：

对于任意状态 $q \in Q$ 和输入符号 $a \in \Sigma$，$\delta(q, a)$ 是 $Q$ 的子集。

**定义 2.5**: NFA的扩展转移函数

**扩展转移函数** $\hat{\delta}: Q \times \Sigma^* \rightarrow 2^Q$ 递归定义：

1. **基础**: $\hat{\delta}(q, \varepsilon) = \{q\}$
2. **归纳**: $\hat{\delta}(q, wa) = \bigcup_{p \in \hat{\delta}(q, w)} \delta(p, a)$

**定义 2.6**: NFA的语言

NFA $M$ 接受的语言：
$$L(M) = \{w \in \Sigma^* \mid \hat{\delta}(q_0, w) \cap F \neq \emptyset\}$$

### 2.3 PDA定义

**定义 2.7**: 下推自动机的配置

**配置** $(q, w, \gamma)$ 表示：

- $q \in Q$: 当前状态
- $w \in \Sigma^*$: 剩余输入串
- $\gamma \in \Gamma^*$: 栈内容

**定义 2.8**: PDA的转移关系

**转移关系** $\vdash$ 定义：
$$(q, aw, Z\gamma) \vdash (p, w, \beta\gamma)$$

当且仅当 $(p, \beta) \in \delta(q, a, Z)$。

**定义 2.9**: PDA的语言

PDA $P$ 接受的语言：
$$L(P) = \{w \in \Sigma^* \mid (q_0, w, Z_0) \vdash^* (q, \varepsilon, \gamma) \text{ 且 } q \in F\}$$

## 3. 重要定理

### 3.1 泵引理

**定理 3.1**: 正则语言的泵引理

如果 $L$ 是正则语言，则存在常数 $n$，使得对于任意字符串 $w \in L$ 且 $|w| \geq n$，$w$ 可以分解为 $w = xyz$，满足：

1. $|xy| \leq n$
2. $|y| > 0$
3. 对于所有 $i \geq 0$，$xy^iz \in L$

**证明**: 设 $M$ 是接受 $L$ 的DFA，状态数为 $n$。对于长度至少为 $n$ 的字符串 $w$，在 $M$ 上运行 $w$ 时，根据鸽巢原理，至少有一个状态被访问两次。设 $y$ 是导致循环的子串，则 $xy^iz$ 都被接受。

### 3.2 最小化定理

**定理 3.2**: DFA最小化

对于任意DFA $M$，存在唯一的（在同构意义下）最小DFA $M'$，使得 $L(M) = L(M')$。

**证明**: 通过等价类构造最小DFA。两个状态 $p$ 和 $q$ 等价，当且仅当对于所有输入串 $w$，$\hat{\delta}(p, w) \in F$ 当且仅当 $\hat{\delta}(q, w) \in F$。

### 3.3 等价性定理

**定理 3.3**: DFA和NFA等价性

对于任意NFA $M$，存在DFA $M'$ 使得 $L(M) = L(M')$。

**证明**: 使用子集构造法。DFA的状态是NFA状态的子集，转移函数通过NFA的转移函数定义。

## 4. Go语言实现

### 4.1 自动机数据结构

```go
// State 状态
type State struct {
    ID   string
    Name string
}

// NewState 创建状态
func NewState(id, name string) *State {
    return &State{
        ID:   id,
        Name: name,
    }
}

// Transition 转移
type Transition struct {
    From     *State
    To       *State
    Symbol   string
    StackPop string  // 用于PDA
    StackPush string // 用于PDA
}

// NewTransition 创建转移
func NewTransition(from, to *State, symbol string) *Transition {
    return &Transition{
        From:   from,
        To:     to,
        Symbol: symbol,
    }
}

// NewPDATransition 创建PDA转移
func NewPDATransition(from, to *State, symbol, stackPop, stackPush string) *Transition {
    return &Transition{
        From:      from,
        To:        to,
        Symbol:    symbol,
        StackPop:  stackPop,
        StackPush: stackPush,
    }
}

// DFA 确定性有限自动机
type DFA struct {
    States       map[string]*State
    Alphabet     map[string]bool
    Transitions  map[string]map[string]*State // (state,symbol) -> state
    InitialState *State
    AcceptStates map[string]*State
}

// NewDFA 创建DFA
func NewDFA() *DFA {
    return &DFA{
        States:       make(map[string]*State),
        Alphabet:     make(map[string]bool),
        Transitions:  make(map[string]map[string]*State),
        AcceptStates: make(map[string]*State),
    }
}

// AddState 添加状态
func (dfa *DFA) AddState(state *State) {
    dfa.States[state.ID] = state
    dfa.Transitions[state.ID] = make(map[string]*State)
}

// AddSymbol 添加符号
func (dfa *DFA) AddSymbol(symbol string) {
    dfa.Alphabet[symbol] = true
}

// AddTransition 添加转移
func (dfa *DFA) AddTransition(from, symbol, to string) {
    if dfa.Transitions[from] == nil {
        dfa.Transitions[from] = make(map[string]*State)
    }
    dfa.Transitions[from][symbol] = dfa.States[to]
}

// SetInitialState 设置初始状态
func (dfa *DFA) SetInitialState(stateID string) {
    dfa.InitialState = dfa.States[stateID]
}

// AddAcceptState 添加接受状态
func (dfa *DFA) AddAcceptState(stateID string) {
    dfa.AcceptStates[stateID] = dfa.States[stateID]
}

// GetTransition 获取转移
func (dfa *DFA) GetTransition(stateID, symbol string) *State {
    if dfa.Transitions[stateID] == nil {
        return nil
    }
    return dfa.Transitions[stateID][symbol]
}

// Accepts 检查是否接受输入串
func (dfa *DFA) Accepts(input string) bool {
    currentState := dfa.InitialState
    
    for _, symbol := range input {
        symbolStr := string(symbol)
        if !dfa.Alphabet[symbolStr] {
            return false // 输入符号不在字母表中
        }
        
        nextState := dfa.GetTransition(currentState.ID, symbolStr)
        if nextState == nil {
            return false // 没有转移
        }
        
        currentState = nextState
    }
    
    // 检查最终状态是否为接受状态
    _, isAccept := dfa.AcceptStates[currentState.ID]
    return isAccept
}

// GetReachableStates 获取可达状态
func (dfa *DFA) GetReachableStates() map[string]*State {
    reachable := make(map[string]*State)
    visited := make(map[string]bool)
    
    var dfs func(state *State)
    dfs = func(state *State) {
        if visited[state.ID] {
            return
        }
        
        visited[state.ID] = true
        reachable[state.ID] = state
        
        // 遍历所有转移
        for symbol := range dfa.Alphabet {
            if nextState := dfa.GetTransition(state.ID, symbol); nextState != nil {
                dfs(nextState)
            }
        }
    }
    
    if dfa.InitialState != nil {
        dfs(dfa.InitialState)
    }
    
    return reachable
}

// Minimize 最小化DFA
func (dfa *DFA) Minimize() *DFA {
    // 移除不可达状态
    reachable := dfa.GetReachableStates()
    
    // 计算等价类
    equivalenceClasses := dfa.computeEquivalenceClasses(reachable)
    
    // 构建最小DFA
    return dfa.buildMinimalDFA(equivalenceClasses)
}

// computeEquivalenceClasses 计算等价类
func (dfa *DFA) computeEquivalenceClasses(reachable map[string]*State) map[string]string {
    // 初始分类：接受状态和非接受状态
    classes := make(map[string]string)
    classCounter := 0
    
    for stateID := range reachable {
        if _, isAccept := dfa.AcceptStates[stateID]; isAccept {
            classes[stateID] = "accept"
        } else {
            classes[stateID] = "non-accept"
        }
    }
    
    // 迭代细化等价类
    for {
        newClasses := make(map[string]string)
        classMap := make(map[string][]string)
        
        // 按当前类分组
        for stateID, class := range classes {
            classMap[class] = append(classMap[class], stateID)
        }
        
        // 为每个类分配新类名
        for class, states := range classMap {
            if len(states) == 1 {
                newClasses[states[0]] = class
                continue
            }
            
            // 检查状态是否等价
            subClasses := dfa.refineClass(states, classes)
            for stateID, subClass := range subClasses {
                newClasses[stateID] = fmt.Sprintf("%s_%s", class, subClass)
            }
        }
        
        // 检查是否收敛
        if reflect.DeepEqual(classes, newClasses) {
            break
        }
        
        classes = newClasses
    }
    
    return classes
}

// refineClass 细化等价类
func (dfa *DFA) refineClass(states []string, classes map[string]string) map[string]string {
    subClasses := make(map[string]string)
    
    for i, stateID := range states {
        // 检查与其他状态的等价性
        equivalent := false
        for j := 0; j < i; j++ {
            if dfa.statesEquivalent(stateID, states[j], classes) {
                subClasses[stateID] = subClasses[states[j]]
                equivalent = true
                break
            }
        }
        
        if !equivalent {
            subClasses[stateID] = fmt.Sprintf("sub_%d", i)
        }
    }
    
    return subClasses
}

// statesEquivalent 检查两个状态是否等价
func (dfa *DFA) statesEquivalent(state1, state2 string, classes map[string]string) bool {
    // 检查所有输入符号的转移是否等价
    for symbol := range dfa.Alphabet {
        next1 := dfa.GetTransition(state1, symbol)
        next2 := dfa.GetTransition(state2, symbol)
        
        if next1 == nil && next2 == nil {
            continue
        }
        
        if next1 == nil || next2 == nil {
            return false
        }
        
        if classes[next1.ID] != classes[next2.ID] {
            return false
        }
    }
    
    return true
}

// buildMinimalDFA 构建最小DFA
func (dfa *DFA) buildMinimalDFA(classes map[string]string) *DFA {
    minimalDFA := NewDFA()
    
    // 创建新状态
    classStates := make(map[string]*State)
    for stateID, class := range classes {
        if _, exists := classStates[class]; !exists {
            newState := NewState(class, fmt.Sprintf("Class_%s", class))
            minimalDFA.AddState(newState)
            classStates[class] = newState
            
            // 检查是否为接受状态
            if _, isAccept := dfa.AcceptStates[stateID]; isAccept {
                minimalDFA.AddAcceptState(class)
            }
        }
    }
    
    // 复制字母表
    for symbol := range dfa.Alphabet {
        minimalDFA.AddSymbol(symbol)
    }
    
    // 添加转移
    for stateID, class := range classes {
        for symbol := range dfa.Alphabet {
            if nextState := dfa.GetTransition(stateID, symbol); nextState != nil {
                nextClass := classes[nextState.ID]
                minimalDFA.AddTransition(class, symbol, nextClass)
            }
        }
    }
    
    // 设置初始状态
    if dfa.InitialState != nil {
        initialClass := classes[dfa.InitialState.ID]
        minimalDFA.SetInitialState(initialClass)
    }
    
    return minimalDFA
}
```

### 4.2 NFA实现

```go
// NFA 非确定性有限自动机
type NFA struct {
    States       map[string]*State
    Alphabet     map[string]bool
    Transitions  map[string]map[string][]*State // (state,symbol) -> []state
    InitialState *State
    AcceptStates map[string]*State
}

// NewNFA 创建NFA
func NewNFA() *NFA {
    return &NFA{
        States:       make(map[string]*State),
        Alphabet:     make(map[string]bool),
        Transitions:  make(map[string]map[string][]*State),
        AcceptStates: make(map[string]*State),
    }
}

// AddState 添加状态
func (nfa *NFA) AddState(state *State) {
    nfa.States[state.ID] = state
    nfa.Transitions[state.ID] = make(map[string][]*State)
}

// AddSymbol 添加符号
func (nfa *NFA) AddSymbol(symbol string) {
    nfa.Alphabet[symbol] = true
}

// AddTransition 添加转移
func (nfa *NFA) AddTransition(from, symbol, to string) {
    if nfa.Transitions[from] == nil {
        nfa.Transitions[from] = make(map[string][]*State)
    }
    if nfa.Transitions[from][symbol] == nil {
        nfa.Transitions[from][symbol] = []*State{}
    }
    nfa.Transitions[from][symbol] = append(nfa.Transitions[from][symbol], nfa.States[to])
}

// SetInitialState 设置初始状态
func (nfa *NFA) SetInitialState(stateID string) {
    nfa.InitialState = nfa.States[stateID]
}

// AddAcceptState 添加接受状态
func (nfa *NFA) AddAcceptState(stateID string) {
    nfa.AcceptStates[stateID] = nfa.States[stateID]
}

// GetTransitions 获取转移
func (nfa *NFA) GetTransitions(stateID, symbol string) []*State {
    if nfa.Transitions[stateID] == nil {
        return nil
    }
    return nfa.Transitions[stateID][symbol]
}

// Accepts 检查是否接受输入串
func (nfa *NFA) Accepts(input string) bool {
    currentStates := map[string]*State{nfa.InitialState.ID: nfa.InitialState}
    
    for _, symbol := range input {
        symbolStr := string(symbol)
        if !nfa.Alphabet[symbolStr] {
            return false
        }
        
        nextStates := make(map[string]*State)
        for _, state := range currentStates {
            transitions := nfa.GetTransitions(state.ID, symbolStr)
            for _, nextState := range transitions {
                nextStates[nextState.ID] = nextState
            }
        }
        
        if len(nextStates) == 0 {
            return false
        }
        
        currentStates = nextStates
    }
    
    // 检查是否有接受状态
    for _, state := range currentStates {
        if _, isAccept := nfa.AcceptStates[state.ID]; isAccept {
            return true
        }
    }
    
    return false
}

// ToDFA 转换为DFA
func (nfa *NFA) ToDFA() *DFA {
    dfa := NewDFA()
    
    // 复制字母表
    for symbol := range nfa.Alphabet {
        dfa.AddSymbol(symbol)
    }
    
    // 计算ε闭包
    initialClosure := nfa.epsilonClosure([]*State{nfa.InitialState})
    
    // 创建状态映射
    stateMap := make(map[string][]*State)
    stateMap["q0"] = initialClosure
    
    // 使用队列进行广度优先搜索
    queue := []string{"q0"}
    processed := make(map[string]bool)
    
    for len(queue) > 0 {
        currentStateID := queue[0]
        queue = queue[1:]
        
        if processed[currentStateID] {
            continue
        }
        processed[currentStateID] = true
        
        currentStates := stateMap[currentStateID]
        
        // 创建DFA状态
        dfaState := NewState(currentStateID, fmt.Sprintf("DFA_%s", currentStateID))
        dfa.AddState(dfaState)
        
        // 检查是否为接受状态
        for _, state := range currentStates {
            if _, isAccept := nfa.AcceptStates[state.ID]; isAccept {
                dfa.AddAcceptState(currentStateID)
                break
            }
        }
        
        // 计算转移
        for symbol := range nfa.Alphabet {
            nextStates := nfa.move(currentStates, symbol)
            if len(nextStates) > 0 {
                nextClosure := nfa.epsilonClosure(nextStates)
                nextStateID := nfa.getStateID(nextClosure, stateMap)
                
                if nextStateID == "" {
                    nextStateID = fmt.Sprintf("q%d", len(stateMap))
                    stateMap[nextStateID] = nextClosure
                    queue = append(queue, nextStateID)
                }
                
                dfa.AddTransition(currentStateID, symbol, nextStateID)
            }
        }
    }
    
    dfa.SetInitialState("q0")
    return dfa
}

// epsilonClosure 计算ε闭包
func (nfa *NFA) epsilonClosure(states []*State) []*State {
    closure := make(map[string]*State)
    stack := make([]*State, len(states))
    copy(stack, states)
    
    for _, state := range states {
        closure[state.ID] = state
    }
    
    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        // 获取ε转移
        epsilonTransitions := nfa.GetTransitions(current.ID, "ε")
        for _, nextState := range epsilonTransitions {
            if _, exists := closure[nextState.ID]; !exists {
                closure[nextState.ID] = nextState
                stack = append(stack, nextState)
            }
        }
    }
    
    result := make([]*State, 0, len(closure))
    for _, state := range closure {
        result = append(result, state)
    }
    
    return result
}

// move 计算move操作
func (nfa *NFA) move(states []*State, symbol string) []*State {
    result := make(map[string]*State)
    
    for _, state := range states {
        transitions := nfa.GetTransitions(state.ID, symbol)
        for _, nextState := range transitions {
            result[nextState.ID] = nextState
        }
    }
    
    moveResult := make([]*State, 0, len(result))
    for _, state := range result {
        moveResult = append(moveResult, state)
    }
    
    return moveResult
}

// getStateID 获取状态ID
func (nfa *NFA) getStateID(states []*State, stateMap map[string][]*State) string {
    stateSet := make(map[string]bool)
    for _, state := range states {
        stateSet[state.ID] = true
    }
    
    for id, existingStates := range stateMap {
        if len(existingStates) != len(states) {
            continue
        }
        
        match := true
        for _, state := range existingStates {
            if !stateSet[state.ID] {
                match = false
                break
            }
        }
        
        if match {
            return id
        }
    }
    
    return ""
}
```

### 4.3 PDA实现

```go
// PDA 下推自动机
type PDA struct {
    States       map[string]*State
    InputAlphabet map[string]bool
    StackAlphabet map[string]bool
    Transitions   map[string][]*Transition // state -> []transition
    InitialState  *State
    InitialStack  string
    AcceptStates  map[string]*State
}

// NewPDA 创建PDA
func NewPDA() *PDA {
    return &PDA{
        States:        make(map[string]*State),
        InputAlphabet: make(map[string]bool),
        StackAlphabet: make(map[string]bool),
        Transitions:   make(map[string][]*Transition),
        AcceptStates:  make(map[string]*State),
    }
}

// AddState 添加状态
func (pda *PDA) AddState(state *State) {
    pda.States[state.ID] = state
    pda.Transitions[state.ID] = []*Transition{}
}

// AddInputSymbol 添加输入符号
func (pda *PDA) AddInputSymbol(symbol string) {
    pda.InputAlphabet[symbol] = true
}

// AddStackSymbol 添加栈符号
func (pda *PDA) AddStackSymbol(symbol string) {
    pda.StackAlphabet[symbol] = true
}

// AddTransition 添加转移
func (pda *PDA) AddTransition(from, symbol, stackPop, stackPush, to string) {
    transition := NewPDATransition(
        pda.States[from],
        pda.States[to],
        symbol,
        stackPop,
        stackPush,
    )
    
    pda.Transitions[from] = append(pda.Transitions[from], transition)
}

// SetInitialState 设置初始状态
func (pda *PDA) SetInitialState(stateID string) {
    pda.InitialState = pda.States[stateID]
}

// SetInitialStack 设置初始栈符号
func (pda *PDA) SetInitialStack(symbol string) {
    pda.InitialStack = symbol
}

// AddAcceptState 添加接受状态
func (pda *PDA) AddAcceptState(stateID string) {
    pda.AcceptStates[stateID] = pda.States[stateID]
}

// Accepts 检查是否接受输入串
func (pda *PDA) Accepts(input string) bool {
    // 初始配置
    initialConfig := &PDAConfig{
        State: pda.InitialState,
        Input: input,
        Stack: []string{pda.InitialStack},
    }
    
    return pda.acceptsFromConfig(initialConfig)
}

// PDAConfig PDA配置
type PDAConfig struct {
    State *State
    Input string
    Stack []string
}

// acceptsFromConfig 从配置开始检查是否接受
func (pda *PDA) acceptsFromConfig(config *PDAConfig) bool {
    // 如果输入为空且状态为接受状态，则接受
    if config.Input == "" {
        _, isAccept := pda.AcceptStates[config.State.ID]
        return isAccept
    }
    
    // 获取当前符号和栈顶符号
    currentSymbol := string(config.Input[0])
    var stackTop string
    if len(config.Stack) > 0 {
        stackTop = config.Stack[len(config.Stack)-1]
    } else {
        return false
    }
    
    // 尝试所有可能的转移
    transitions := pda.Transitions[config.State.ID]
    for _, transition := range transitions {
        // 检查转移条件
        if (transition.Symbol == currentSymbol || transition.Symbol == "ε") &&
           transition.StackPop == stackTop {
            
            // 创建新配置
            newConfig := &PDAConfig{
                State: transition.To,
                Input: config.Input,
                Stack: make([]string, len(config.Stack)),
            }
            copy(newConfig.Stack, config.Stack)
            
            // 更新栈
            if len(newConfig.Stack) > 0 {
                newConfig.Stack = newConfig.Stack[:len(newConfig.Stack)-1]
            }
            
            // 压入新符号
            for i := len(transition.StackPush) - 1; i >= 0; i-- {
                newConfig.Stack = append(newConfig.Stack, string(transition.StackPush[i]))
            }
            
            // 更新输入
            if transition.Symbol != "ε" {
                newConfig.Input = newConfig.Input[1:]
            }
            
            // 递归检查
            if pda.acceptsFromConfig(newConfig) {
                return true
            }
        }
    }
    
    return false
}
```

## 5. 应用示例

### 5.1 词法分析器

```go
// LexicalAnalyzer 词法分析器
type LexicalAnalyzer struct {
    dfa *DFA
}

// NewLexicalAnalyzer 创建词法分析器
func NewLexicalAnalyzer() *LexicalAnalyzer {
    dfa := NewDFA()
    
    // 创建状态
    start := NewState("start", "Start")
    identifier := NewState("identifier", "Identifier")
    number := NewState("number", "Number")
    string_lit := NewState("string", "String")
    comment := NewState("comment", "Comment")
    
    dfa.AddState(start)
    dfa.AddState(identifier)
    dfa.AddState(number)
    dfa.AddState(string_lit)
    dfa.AddState(comment)
    
    // 添加符号
    for c := 'a'; c <= 'z'; c++ {
        dfa.AddSymbol(string(c))
    }
    for c := 'A'; c <= 'Z'; c++ {
        dfa.AddSymbol(string(c))
    }
    for c := '0'; c <= '9'; c++ {
        dfa.AddSymbol(string(c))
    }
    dfa.AddSymbol("_")
    dfa.AddSymbol("\"")
    dfa.AddSymbol("/")
    dfa.AddSymbol("*")
    
    // 添加转移
    // 标识符：字母或下划线开头，后跟字母、数字或下划线
    for c := 'a'; c <= 'z'; c++ {
        dfa.AddTransition("start", string(c), "identifier")
    }
    for c := 'A'; c <= 'Z'; c++ {
        dfa.AddTransition("start", string(c), "identifier")
    }
    dfa.AddTransition("start", "_", "identifier")
    
    for c := 'a'; c <= 'z'; c++ {
        dfa.AddTransition("identifier", string(c), "identifier")
    }
    for c := 'A'; c <= 'Z'; c++ {
        dfa.AddTransition("identifier", string(c), "identifier")
    }
    for c := '0'; c <= '9'; c++ {
        dfa.AddTransition("identifier", string(c), "identifier")
    }
    dfa.AddTransition("identifier", "_", "identifier")
    
    // 数字：数字序列
    for c := '0'; c <= '9'; c++ {
        dfa.AddTransition("start", string(c), "number")
        dfa.AddTransition("number", string(c), "number")
    }
    
    // 字符串：双引号包围
    dfa.AddTransition("start", "\"", "string")
    for c := 'a'; c <= 'z'; c++ {
        dfa.AddTransition("string", string(c), "string")
    }
    for c := 'A'; c <= 'Z'; c++ {
        dfa.AddTransition("string", string(c), "string")
    }
    for c := '0'; c <= '9'; c++ {
        dfa.AddTransition("string", string(c), "string")
    }
    dfa.AddTransition("string", "\"", "string")
    
    // 注释：/* ... */
    dfa.AddTransition("start", "/", "comment")
    dfa.AddTransition("comment", "*", "comment")
    dfa.AddTransition("comment", "/", "comment")
    for c := 'a'; c <= 'z'; c++ {
        dfa.AddTransition("comment", string(c), "comment")
    }
    for c := 'A'; c <= 'Z'; c++ {
        dfa.AddTransition("comment", string(c), "comment")
    }
    for c := '0'; c <= '9'; c++ {
        dfa.AddTransition("comment", string(c), "comment")
    }
    
    dfa.SetInitialState("start")
    dfa.AddAcceptState("identifier")
    dfa.AddAcceptState("number")
    dfa.AddAcceptState("string")
    dfa.AddAcceptState("comment")
    
    return &LexicalAnalyzer{dfa: dfa}
}

// Tokenize 词法分析
func (la *LexicalAnalyzer) Tokenize(input string) []string {
    var tokens []string
    current := ""
    
    for _, char := range input {
        symbol := string(char)
        current += symbol
        
        if !la.dfa.Accepts(current) {
            if len(current) > 1 {
                // 回退一个字符
                current = current[:len(current)-1]
                if la.dfa.Accepts(current) {
                    tokens = append(tokens, current)
                    current = symbol
                }
            }
        }
    }
    
    if la.dfa.Accepts(current) {
        tokens = append(tokens, current)
    }
    
    return tokens
}

// LexicalAnalyzerExample 词法分析器示例
func LexicalAnalyzerExample() {
    analyzer := NewLexicalAnalyzer()
    
    input := "int x = 42; /* comment */ string name = \"hello\";"
    tokens := analyzer.Tokenize(input)
    
    fmt.Println("词法分析器示例")
    fmt.Printf("输入: %s\n", input)
    fmt.Printf("词法单元: %v\n", tokens)
}
```

### 5.2 模式匹配

```go
// PatternMatcher 模式匹配器
type PatternMatcher struct {
    nfa *NFA
}

// NewPatternMatcher 创建模式匹配器
func NewPatternMatcher(pattern string) *PatternMatcher {
    nfa := NewNFA()
    
    // 创建状态
    for i := 0; i <= len(pattern); i++ {
        state := NewState(fmt.Sprintf("q%d", i), fmt.Sprintf("State_%d", i))
        nfa.AddState(state)
    }
    
    // 添加符号
    for _, char := range pattern {
        nfa.AddSymbol(string(char))
    }
    
    // 构建转移
    for i, char := range pattern {
        symbol := string(char)
        nfa.AddTransition(fmt.Sprintf("q%d", i), symbol, fmt.Sprintf("q%d", i+1))
    }
    
    nfa.SetInitialState("q0")
    nfa.AddAcceptState(fmt.Sprintf("q%d", len(pattern)))
    
    return &PatternMatcher{nfa: nfa}
}

// Match 模式匹配
func (pm *PatternMatcher) Match(text string) []int {
    var matches []int
    
    for i := 0; i <= len(text)-len(pm.pattern); i++ {
        substring := text[i : i+len(pm.pattern)]
        if pm.nfa.Accepts(substring) {
            matches = append(matches, i)
        }
    }
    
    return matches
}

// PatternMatcherExample 模式匹配示例
func PatternMatcherExample() {
    pattern := "abc"
    matcher := NewPatternMatcher(pattern)
    
    text := "abcabcabc"
    matches := matcher.Match(text)
    
    fmt.Println("模式匹配示例")
    fmt.Printf("模式: %s\n", pattern)
    fmt.Printf("文本: %s\n", text)
    fmt.Printf("匹配位置: %v\n", matches)
}
```

### 5.3 状态机

```go
// StateMachine 状态机
type StateMachine struct {
    dfa *DFA
    currentState *State
}

// NewStateMachine 创建状态机
func NewStateMachine() *StateMachine {
    dfa := NewDFA()
    
    // 创建状态
    idle := NewState("idle", "Idle")
    running := NewState("running", "Running")
    paused := NewState("paused", "Paused")
    stopped := NewState("stopped", "Stopped")
    
    dfa.AddState(idle)
    dfa.AddState(running)
    dfa.AddState(paused)
    dfa.AddState(stopped)
    
    // 添加事件
    dfa.AddSymbol("start")
    dfa.AddSymbol("pause")
    dfa.AddSymbol("resume")
    dfa.AddSymbol("stop")
    
    // 添加转移
    dfa.AddTransition("idle", "start", "running")
    dfa.AddTransition("running", "pause", "paused")
    dfa.AddTransition("running", "stop", "stopped")
    dfa.AddTransition("paused", "resume", "running")
    dfa.AddTransition("paused", "stop", "stopped")
    dfa.AddTransition("stopped", "start", "running")
    
    dfa.SetInitialState("idle")
    
    return &StateMachine{
        dfa:          dfa,
        currentState: idle,
    }
}

// ProcessEvent 处理事件
func (sm *StateMachine) ProcessEvent(event string) bool {
    nextState := sm.dfa.GetTransition(sm.currentState.ID, event)
    if nextState != nil {
        sm.currentState = nextState
        return true
    }
    return false
}

// GetCurrentState 获取当前状态
func (sm *StateMachine) GetCurrentState() string {
    return sm.currentState.Name
}

// StateMachineExample 状态机示例
func StateMachineExample() {
    sm := NewStateMachine()
    
    events := []string{"start", "pause", "resume", "stop", "start"}
    
    fmt.Println("状态机示例")
    fmt.Printf("初始状态: %s\n", sm.GetCurrentState())
    
    for _, event := range events {
        if sm.ProcessEvent(event) {
            fmt.Printf("事件 '%s' -> 状态: %s\n", event, sm.GetCurrentState())
        } else {
            fmt.Printf("事件 '%s' -> 无效转移\n", event)
        }
    }
}
```

## 总结

自动机理论是计算机科学的基础理论，提供了：

1. **形式化模型**: 用于描述和验证计算过程
2. **语言理论**: 研究形式语言的性质和分类
3. **计算模型**: 从有限自动机到图灵机的层次结构
4. **广泛应用**: 在编译原理、自然语言处理、人工智能等领域有重要应用

通过Go语言的实现，我们展示了：

- DFA、NFA、PDA的数据结构表示
- 自动机的构造和运行
- 最小化、等价性转换等算法
- 词法分析、模式匹配、状态机等应用

这为后续的形式语言、计算复杂性等更高级的理论奠定了基础。
