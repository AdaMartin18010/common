# 02-形式语言 (Formal Languages)

## 目录

- [02-形式语言 (Formal Languages)](#02-形式语言-formal-languages)
  - [目录](#目录)
  - [1. 概念定义](#1-概念定义)
    - [1.1 基本概念](#11-基本概念)
    - [1.2 核心思想](#12-核心思想)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 数学定义](#21-数学定义)
    - [2.2 类型定义](#22-类型定义)
  - [3. 定理证明](#3-定理证明)
    - [3.1 定理陈述](#31-定理陈述)
    - [3.2 证明过程](#32-证明过程)
  - [4. Go语言实现](#4-go语言实现)
    - [4.1 基础实现](#41-基础实现)
    - [4.2 泛型实现](#42-泛型实现)
    - [4.3 并发实现](#43-并发实现)
  - [5. 应用示例](#5-应用示例)
    - [5.1 基础示例](#51-基础示例)
    - [5.2 高级示例](#52-高级示例)
  - [6. 性能分析](#6-性能分析)
    - [6.1 时间复杂度](#61-时间复杂度)
    - [6.2 空间复杂度](#62-空间复杂度)
    - [6.3 基准测试](#63-基准测试)
  - [7. 参考文献](#7-参考文献)

## 1. 概念定义

### 1.1 基本概念

**定义 1.1**: 字母表 (Alphabet)
字母表 $\Sigma$ 是一个有限的符号集合。

**定义 1.2**: 字符串 (String)
字母表 $\Sigma$ 上的字符串是 $\Sigma$ 中符号的有限序列。

**定义 1.3**: 语言 (Language)
字母表 $\Sigma$ 上的语言是 $\Sigma^*$ 的子集，其中 $\Sigma^*$ 表示 $\Sigma$ 上所有字符串的集合。

### 1.2 核心思想

形式语言理论是研究字符串集合的数学理论，它为编译器设计、自然语言处理、密码学等领域提供了理论基础。

## 2. 形式化定义

### 2.1 数学定义

**定义 2.1**: 字符串操作
设 $w = a_1a_2\ldots a_n$ 和 $v = b_1b_2\ldots b_m$ 是字符串：

1. **连接**: $w \cdot v = a_1a_2\ldots a_nb_1b_2\ldots b_m$
2. **长度**: $|w| = n$
3. **空字符串**: $\varepsilon$ 是长度为0的字符串
4. **幂**: $w^0 = \varepsilon$, $w^{n+1} = w^n \cdot w$

**定义 2.2**: 语言操作
设 $L_1, L_2$ 是语言：

1. **并**: $L_1 \cup L_2 = \{w \mid w \in L_1 \text{ 或 } w \in L_2\}$
2. **交**: $L_1 \cap L_2 = \{w \mid w \in L_1 \text{ 且 } w \in L_2\}$
3. **连接**: $L_1 \cdot L_2 = \{w \cdot v \mid w \in L_1, v \in L_2\}$
4. **克林闭包**: $L^* = \bigcup_{i=0}^{\infty} L^i$

**定义 2.3**: 乔姆斯基层次结构

1. **正则语言**: 由正则表达式定义
2. **上下文无关语言**: 由上下文无关文法定义
3. **上下文有关语言**: 由上下文有关文法定义
4. **递归可枚举语言**: 由图灵机定义

### 2.2 类型定义

```go
// Alphabet 表示字母表
type Alphabet struct {
    Symbols map[string]bool
    mu      sync.RWMutex
}

// String 表示字符串
type String struct {
    Symbols []string
    Length  int
}

// Language 表示语言
type Language interface {
    Contains(str String) bool
    IsEmpty() bool
    IsFinite() bool
    Cardinality() int
    String() string
}

// RegularLanguage 正则语言
type RegularLanguage struct {
    Regex string
    DFA   *DeterministicFiniteAutomaton
}

// ContextFreeLanguage 上下文无关语言
type ContextFreeLanguage struct {
    Grammar *ContextFreeGrammar
    PDA     *PushdownAutomaton
}

// DeterministicFiniteAutomaton 确定有限自动机
type DeterministicFiniteAutomaton struct {
    States       map[string]bool
    Alphabet     *Alphabet
    Transitions  map[string]map[string]string // state -> symbol -> state
    StartState   string
    AcceptStates map[string]bool
}

// ContextFreeGrammar 上下文无关文法
type ContextFreeGrammar struct {
    Variables    map[string]bool
    Terminals    *Alphabet
    Productions  []Production
    StartSymbol  string
}

// Production 产生式
type Production struct {
    Left  string   // 左部（单个变量）
    Right []string // 右部（符号序列）
}
```

## 3. 定理证明

### 3.1 定理陈述

**定理 3.1**: 泵引理 (Pumping Lemma)
设 $L$ 是正则语言，则存在常数 $n$，使得对于任意 $w \in L$ 且 $|w| \geq n$，存在分解 $w = xyz$，满足：

1. $|xy| \leq n$
2. $|y| > 0$
3. 对于所有 $i \geq 0$，$xy^iz \in L$

**定理 3.2**: 上下文无关语言的泵引理
设 $L$ 是上下文无关语言，则存在常数 $n$，使得对于任意 $w \in L$ 且 $|w| \geq n$，存在分解 $w = uvxyz$，满足：

1. $|vxy| \leq n$
2. $|vy| > 0$
3. 对于所有 $i \geq 0$，$uv^ixy^iz \in L$

**定理 3.3**: 语言类的包含关系
$$
\text{Regular} \subset \text{Context-Free} \subset \text{Context-Sensitive} \subset \text{Recursively-Enumerable}
$$

### 3.2 证明过程

**证明定理 3.1**: 泵引理

设 $L$ 是正则语言，$M$ 是接受 $L$ 的DFA，有 $n$ 个状态。

对于任意 $w \in L$ 且 $|w| \geq n$，考虑 $M$ 在输入 $w$ 时的状态序列：
$p_0, p_1, \ldots, p_{|w|}$
其中 $p_0$ 是初始状态，$p_{i+1} = \delta(p_i, w_{i+1})$。

由于 $|w| \geq n$，状态序列长度为 $|w|+1 > n$。根据鸽巢原理，序列中必有重复状态。
设 $p_j = p_k$ 是第一个重复状态，其中 $j < k \leq n$。

令 $w = xyz$：
- $x = w_1\ldots w_j$
- $y = w_{j+1}\ldots w_k$
- $z = w_{k+1}\ldots w_{|w|}$

那么：
1. $|xy| = k \leq n$
2. $|y| = k - j > 0$
3. 对于所有 $i \geq 0$，$M$ 在输入 $xy^iz$ 时，状态序列为：
   - 从 $p_0$ 经过 $x$ 到达 $p_j$
   - 经过 $y$ 从 $p_j$ 回到 $p_j$ (循环 $i$ 次)
   - 从 $p_j=p_k$ 经过 $z$ 到达最终接受状态
因此，$xy^iz \in L$。

## 4. Go语言实现

### 4.1 基础实现

```go
package formal_languages

import (
	"fmt"
	"sync"
)

// Alphabet 表示字母表
type Alphabet struct {
	Symbols map[rune]bool
	mu      sync.RWMutex
}

// NewAlphabet 创建一个新的字母表
func NewAlphabet(symbols []rune) *Alphabet {
	s := make(map[rune]bool)
	for _, r := range symbols {
		s[r] = true
	}
	return &Alphabet{Symbols: s}
}

// Contains 检查符号是否在字母表中
func (a *Alphabet) Contains(symbol rune) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.Symbols[symbol]
	return ok
}

// String 表示字符串
type String struct {
	Symbols []rune
}

// NewString 创建一个新的字符串
func NewString(s string) *String {
	return &String{Symbols: []rune(s)}
}

// Length 返回字符串长度
func (s *String) Length() int {
	return len(s.Symbols)
}

// Concat 连接两个字符串
func (s *String) Concat(other *String) *String {
	return &String{Symbols: append(s.Symbols, other.Symbols...)}
}

// Language 表示语言接口
type Language interface {
	Contains(s *String) bool
}

// FiniteLanguage 有限语言
type FiniteLanguage struct {
	Strings map[string]bool
}

// NewFiniteLanguage 创建一个新的有限语言
func NewFiniteLanguage(strings []string) *FiniteLanguage {
	s := make(map[string]bool)
	for _, str := range strings {
		s[str] = true
	}
	return &FiniteLanguage{Strings: s}
}

// Contains 检查字符串是否在语言中
func (l *FiniteLanguage) Contains(s *String) bool {
	_, ok := l.Strings[string(s.Symbols)]
	return ok
}
```

### 4.2 泛型实现

```go
package formal_languages

import (
	"sync"
)

// GenericAlphabet 泛型字母表
type GenericAlphabet[T comparable] struct {
	Symbols map[T]bool
	mu      sync.RWMutex
}

// NewGenericAlphabet 创建一个新的泛型字母表
func NewGenericAlphabet[T comparable](symbols []T) *GenericAlphabet[T] {
	s := make(map[T]bool)
	for _, sym := range symbols {
		s[sym] = true
	}
	return &GenericAlphabet[T]{Symbols: s}
}

// Contains 检查符号是否在字母表中
func (a *GenericAlphabet[T]) Contains(symbol T) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.Symbols[symbol]
	return ok
}

// GenericString 泛型字符串
type GenericString[T comparable] struct {
	Symbols []T
}

// NewGenericString 创建一个新的泛型字符串
func NewGenericString[T comparable](symbols []T) *GenericString[T] {
	return &GenericString[T]{Symbols: symbols}
}

// Length 返回字符串长度
func (s *GenericString[T]) Length() int {
	return len(s.Symbols)
}

// Concat 连接两个泛型字符串
func (s *GenericString[T]) Concat(other *GenericString[T]) *GenericString[T] {
	return &GenericString[T]{Symbols: append(s.Symbols, other.Symbols...)}
}

// GenericLanguage 泛型语言接口
type GenericLanguage[T comparable] interface {
	Contains(s *GenericString[T]) bool
}
```

### 4.3 并发实现

```go
package formal_languages

import "sync"

// ConcurrentFiniteLanguage 并发安全的有限语言
type ConcurrentFiniteLanguage struct {
	Strings map[string]bool
	mu      sync.RWMutex
}

// NewConcurrentFiniteLanguage 创建一个新的并发安全有限语言
func NewConcurrentFiniteLanguage(strings []string) *ConcurrentFiniteLanguage {
	s := make(map[string]bool)
	for _, str := range strings {
		s[str] = true
	}
	return &ConcurrentFiniteLanguage{Strings: s}
}

// Contains 并发安全地检查字符串是否在语言中
func (l *ConcurrentFiniteLanguage) Contains(s *String) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.Strings[string(s.Symbols)]
	return ok
}

// AddString 并发安全地添加字符串
func (l *ConcurrentFiniteLanguage) AddString(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Strings[s] = true
}
```

## 5. 应用示例

### 5.1 基础示例

**示例**: 定义一个语言 $L = \{ a^n b^n \mid n \geq 1 \}$。这不是正则语言，而是上下文无关语言。

```go
package main

import "fmt"

// IsInAnBn 检查字符串是否属于 {a^n b^n | n >= 1}
func IsInAnBn(s string) bool {
    if len(s) == 0 || len(s)%2 != 0 {
        return false
    }
    n := len(s) / 2
    for i := 0; i < n; i++ {
        if s[i] != 'a' {
            return false
        }
    }
    for i := n; i < 2*n; i++ {
        if s[i] != 'b' {
            return false
        }
    }
    return true
}

func main() {
    fmt.Println(IsInAnBn("ab"))      // true
    fmt.Println(IsInAnBn("aabb"))   // true
    fmt.Println(IsInAnBn("aaabbb"))  // true
    fmt.Println(IsInAnBn("abc"))     // false
    fmt.Println(IsInAnBn("a"))       // false
    fmt.Println(IsInAnBn(""))        // false
}
```

### 5.2 高级示例

**示例**: 使用上下文无关文法定义 $L = \{ a^n b^n \mid n \geq 1 \}$

- **变量**: $S$
- **终端**: $a, b$
- **产生式**:
  1. $S \rightarrow ab$
  2. $S \rightarrow aSb$
- **开始符号**: $S$

## 6. 性能分析

### 6.1 时间复杂度

- **DFA/NFA识别**: $O(|w|)$，其中 $|w|$ 是输入字符串的长度。
- **上下文无关文法解析 (CYK算法)**: $O(|w|^3 \cdot |G|)$，其中 $|G|$ 是文法的大小。

### 6.2 空间复杂度

- **DFA/NFA**: $O(1)$ (不包括输入)
- **CYK算法**: $O(|w|^2 \cdot |V|)$，其中 $|V|$ 是变量数量。

### 6.3 基准测试

```go
package formal_languages_test

import "testing"

func IsInAnBn(s string) bool {
    // ... (same as above)
    if len(s) == 0 || len(s)%2 != 0 {
		return false
	}
	n := len(s) / 2
	for i := 0; i < n; i++ {
		if s[i] != 'a' {
			return false
		}
	}
	for i := n; i < 2*n; i++ {
		if s[i] != 'b' {
			return false
		}
	}
	return true
}


func BenchmarkIsInAnBn(b *testing.B) {
    s := "aaaaabbbbb"
    for i := 0; i < b.N; i++ {
        IsInAnBn(s)
    }
}
```

## 7. 参考文献

1. Hopcroft, John E., Rajeev Motwani, and Jeffrey D. Ullman. *Introduction to Automata Theory, Languages, and Computation*. 3rd ed., Pearson, 2006.
2. Sipser, Michael. *Introduction to the Theory of Computation*. 3rd ed., Cengage Learning, 2012.
3. Grune, Dick, and Ceriel J.H. Jacobs. *Parsing Techniques: A Practical Guide*. 2nd ed., Springer, 2008. 