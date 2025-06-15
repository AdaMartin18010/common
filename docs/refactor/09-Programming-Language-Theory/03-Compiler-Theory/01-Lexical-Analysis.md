# 01-词法分析 (Lexical Analysis)

## 目录

- [1. 概述](#1-概述)
- [2. 正则表达式](#2-正则表达式)
- [3. 有限自动机](#3-有限自动机)
- [4. 词法分析器](#4-词法分析器)
- [5. Go语言实现](#5-go语言实现)
- [6. 应用实例](#6-应用实例)

## 1. 概述

### 1.1 词法分析定义

词法分析是将源代码字符流转换为词法单元（token）序列的过程。

**形式化定义**：

```latex
\text{词法分析器} = (\Sigma, Q, \delta, q_0, F, T)
```

其中：

- $\Sigma$ 是输入字母表
- $Q$ 是状态集合
- $\delta$ 是转移函数
- $q_0$ 是初始状态
- $F$ 是接受状态集合
- $T$ 是词法单元类型集合

### 1.2 核心概念

#### 1.2.1 词法单元

```go
// Token 词法单元
type Token struct {
    Type    TokenType
    Value   string
    Line    int
    Column  int
}

type TokenType int

const (
    TokenEOF TokenType = iota
    TokenIdentifier
    TokenNumber
    TokenString
    TokenKeyword
    TokenOperator
    TokenDelimiter
    TokenComment
)
```

## 2. 正则表达式

### 2.1 正则表达式定义

```latex
R ::= \epsilon \mid a \mid R_1 \cdot R_2 \mid R_1 + R_2 \mid R^*
```

**Go语言实现**：

```go
// Regex 正则表达式接口
type Regex interface {
    Match(input string) bool
    String() string
}

// EmptyRegex 空正则表达式
type EmptyRegex struct{}

func (er *EmptyRegex) Match(input string) bool {
    return input == ""
}

func (er *EmptyRegex) String() string {
    return "ε"
}

// CharRegex 字符正则表达式
type CharRegex struct {
    Char rune
}

func (cr *CharRegex) Match(input string) bool {
    return len(input) == 1 && rune(input[0]) == cr.Char
}

func (cr *CharRegex) String() string {
    return string(cr.Char)
}

// ConcatRegex 连接正则表达式
type ConcatRegex struct {
    Left, Right Regex
}

func (concr *ConcatRegex) Match(input string) bool {
    for i := 0; i <= len(input); i++ {
        if concr.Left.Match(input[:i]) && concr.Right.Match(input[i:]) {
            return true
        }
    }
    return false
}

func (concr *ConcatRegex) String() string {
    return fmt.Sprintf("(%s·%s)", concr.Left, concr.Right)
}

// UnionRegex 并集正则表达式
type UnionRegex struct {
    Left, Right Regex
}

func (ur *UnionRegex) Match(input string) bool {
    return ur.Left.Match(input) || ur.Right.Match(input)
}

func (ur *UnionRegex) String() string {
    return fmt.Sprintf("(%s+%s)", ur.Left, ur.Right)
}

// StarRegex 星号正则表达式
type StarRegex struct {
    Inner Regex
}

func (sr *StarRegex) Match(input string) bool {
    if input == "" {
        return true
    }
    
    for i := 1; i <= len(input); i++ {
        if sr.Inner.Match(input[:i]) && sr.Match(input[i:]) {
            return true
        }
    }
    return false
}

func (sr *StarRegex) String() string {
    return fmt.Sprintf("(%s)*", sr.Inner)
}
```

## 3. 有限自动机

### 3.1 确定性有限自动机 (DFA)

```latex
\text{DFA} = (Q, \Sigma, \delta, q_0, F)
```

**Go语言实现**：

```go
// DFA 确定性有限自动机
type DFA struct {
    States     map[string]bool
    Alphabet   map[rune]bool
    Transitions map[string]map[rune]string
    StartState string
    AcceptStates map[string]bool
}

func (dfa *DFA) Accept(input string) bool {
    currentState := dfa.StartState
    
    for _, char := range input {
        if !dfa.Alphabet[char] {
            return false
        }
        
        nextState, exists := dfa.Transitions[currentState][char]
        if !exists {
            return false
        }
        
        currentState = nextState
    }
    
    return dfa.AcceptStates[currentState]
}

func (dfa *DFA) AddTransition(from, to string, char rune) {
    if dfa.Transitions[from] == nil {
        dfa.Transitions[from] = make(map[rune]string)
    }
    dfa.Transitions[from][char] = to
    dfa.Alphabet[char] = true
}
```

### 3.2 非确定性有限自动机 (NFA)

```go
// NFA 非确定性有限自动机
type NFA struct {
    States     map[string]bool
    Alphabet   map[rune]bool
    Transitions map[string]map[rune][]string
    StartState string
    AcceptStates map[string]bool
}

func (nfa *NFA) Accept(input string) bool {
    currentStates := map[string]bool{nfa.StartState: true}
    
    for _, char := range input {
        nextStates := make(map[string]bool)
        
        for state := range currentStates {
            if transitions, exists := nfa.Transitions[state][char]; exists {
                for _, nextState := range transitions {
                    nextStates[nextState] = true
                }
            }
        }
        
        if len(nextStates) == 0 {
            return false
        }
        
        currentStates = nextStates
    }
    
    // 检查是否有接受状态
    for state := range currentStates {
        if nfa.AcceptStates[state] {
            return true
        }
    }
    
    return false
}
```

## 4. 词法分析器

### 4.1 词法分析器实现

```go
// Lexer 词法分析器
type Lexer struct {
    Input   string
    Position int
    Line    int
    Column  int
    Tokens  []Token
}

func (l *Lexer) Tokenize() []Token {
    l.Tokens = []Token{}
    
    for l.Position < len(l.Input) {
        // 跳过空白字符
        l.skipWhitespace()
        
        if l.Position >= len(l.Input) {
            break
        }
        
        // 识别词法单元
        token := l.nextToken()
        if token.Type != TokenEOF {
            l.Tokens = append(l.Tokens, token)
        }
    }
    
    // 添加EOF标记
    l.Tokens = append(l.Tokens, Token{
        Type: TokenEOF,
        Value: "",
        Line: l.Line,
        Column: l.Column,
    })
    
    return l.Tokens
}

func (l *Lexer) skipWhitespace() {
    for l.Position < len(l.Input) {
        char := rune(l.Input[l.Position])
        if char == ' ' || char == '\t' {
            l.Position++
            l.Column++
        } else if char == '\n' {
            l.Position++
            l.Line++
            l.Column = 1
        } else {
            break
        }
    }
}

func (l *Lexer) nextToken() Token {
    char := rune(l.Input[l.Position])
    startPos := l.Position
    startLine := l.Line
    startColumn := l.Column
    
    // 识别标识符或关键字
    if isLetter(char) {
        return l.readIdentifier()
    }
    
    // 识别数字
    if isDigit(char) {
        return l.readNumber()
    }
    
    // 识别字符串
    if char == '"' {
        return l.readString()
    }
    
    // 识别操作符
    if isOperator(char) {
        return l.readOperator()
    }
    
    // 识别分隔符
    if isDelimiter(char) {
        return l.readDelimiter()
    }
    
    // 未知字符
    l.Position++
    l.Column++
    return Token{
        Type: TokenEOF,
        Value: string(char),
        Line: startLine,
        Column: startColumn,
    }
}

func (l *Lexer) readIdentifier() Token {
    startPos := l.Position
    startLine := l.Line
    startColumn := l.Column
    
    for l.Position < len(l.Input) {
        char := rune(l.Input[l.Position])
        if !isLetter(char) && !isDigit(char) && char != '_' {
            break
        }
        l.Position++
        l.Column++
    }
    
    value := l.Input[startPos:l.Position]
    
    // 检查是否为关键字
    if isKeyword(value) {
        return Token{
            Type: TokenKeyword,
            Value: value,
            Line: startLine,
            Column: startColumn,
        }
    }
    
    return Token{
        Type: TokenIdentifier,
        Value: value,
        Line: startLine,
        Column: startColumn,
    }
}

func (l *Lexer) readNumber() Token {
    startPos := l.Position
    startLine := l.Line
    startColumn := l.Column
    
    for l.Position < len(l.Input) {
        char := rune(l.Input[l.Position])
        if !isDigit(char) && char != '.' {
            break
        }
        l.Position++
        l.Column++
    }
    
    value := l.Input[startPos:l.Position]
    
    return Token{
        Type: TokenNumber,
        Value: value,
        Line: startLine,
        Column: startColumn,
    }
}

func (l *Lexer) readString() Token {
    startLine := l.Line
    startColumn := l.Column
    
    // 跳过开始的引号
    l.Position++
    l.Column++
    
    startPos := l.Position
    
    for l.Position < len(l.Input) {
        char := rune(l.Input[l.Position])
        if char == '"' {
            l.Position++
            l.Column++
            break
        } else if char == '\n' {
            l.Line++
            l.Column = 1
        }
        l.Position++
        l.Column++
    }
    
    value := l.Input[startPos:l.Position-1]
    
    return Token{
        Type: TokenString,
        Value: value,
        Line: startLine,
        Column: startColumn,
    }
}

func (l *Lexer) readOperator() Token {
    startLine := l.Line
    startColumn := l.Column
    
    char := rune(l.Input[l.Position])
    l.Position++
    l.Column++
    
    // 检查双字符操作符
    if l.Position < len(l.Input) {
        nextChar := rune(l.Input[l.Position])
        if isDoubleOperator(char, nextChar) {
            l.Position++
            l.Column++
            return Token{
                Type: TokenOperator,
                Value: string(char) + string(nextChar),
                Line: startLine,
                Column: startColumn,
            }
        }
    }
    
    return Token{
        Type: TokenOperator,
        Value: string(char),
        Line: startLine,
        Column: startColumn,
    }
}

func (l *Lexer) readDelimiter() Token {
    startLine := l.Line
    startColumn := l.Column
    
    char := rune(l.Input[l.Position])
    l.Position++
    l.Column++
    
    return Token{
        Type: TokenDelimiter,
        Value: string(char),
        Line: startLine,
        Column: startColumn,
    }
}
```

### 4.2 辅助函数

```go
func isLetter(char rune) bool {
    return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char rune) bool {
    return char >= '0' && char <= '9'
}

func isOperator(char rune) bool {
    operators := []rune{'+', '-', '*', '/', '=', '<', '>', '!', '&', '|', '^', '%'}
    for _, op := range operators {
        if char == op {
            return true
        }
    }
    return false
}

func isDelimiter(char rune) bool {
    delimiters := []rune{'(', ')', '[', ']', '{', '}', ';', ',', '.', ':'}
    for _, delim := range delimiters {
        if char == delim {
            return true
        }
    }
    return false
}

func isKeyword(value string) bool {
    keywords := []string{
        "func", "var", "const", "type", "struct", "interface",
        "if", "else", "for", "range", "switch", "case", "default",
        "return", "break", "continue", "defer", "go", "select",
        "chan", "map", "package", "import",
    }
    
    for _, keyword := range keywords {
        if value == keyword {
            return true
        }
    }
    return false
}

func isDoubleOperator(first, second rune) bool {
    doubleOps := map[rune][]rune{
        '=': {'='},
        '!': {'='},
        '<': {'=', '<'},
        '>': {'=', '>'},
        '&': {'&'},
        '|': {'|'},
        '+': {'+', '='},
        '-': {'-', '='},
        '*': {'='},
        '/': {'=', '/'},
        '%': {'='},
        '^': {'='},
    }
    
    if ops, exists := doubleOps[first]; exists {
        for _, op := range ops {
            if second == op {
                return true
            }
        }
    }
    return false
}
```

## 5. Go语言实现

### 5.1 完整的词法分析器

```go
// GoLexer Go语言词法分析器
type GoLexer struct {
    *Lexer
    Keywords map[string]bool
    Operators map[string]bool
}

func NewGoLexer(input string) *GoLexer {
    lexer := &GoLexer{
        Lexer: &Lexer{
            Input: input,
            Position: 0,
            Line: 1,
            Column: 1,
        },
        Keywords: make(map[string]bool),
        Operators: make(map[string]bool),
    }
    
    // 初始化Go语言关键字
    goKeywords := []string{
        "break", "case", "chan", "const", "continue", "default", "defer",
        "else", "fallthrough", "for", "func", "go", "goto", "if", "import",
        "interface", "map", "package", "range", "return", "select", "struct",
        "switch", "type", "var",
    }
    
    for _, keyword := range goKeywords {
        lexer.Keywords[keyword] = true
    }
    
    // 初始化操作符
    goOperators := []string{
        "+", "-", "*", "/", "%", "&", "|", "^", "<<", ">>", "&^",
        "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>=", "&^=",
        "&&", "||", "<-", "++", "--", "==", "<", ">", "!=", "<=", ">=",
        "=", ":=", "...",
    }
    
    for _, op := range goOperators {
        lexer.Operators[op] = true
    }
    
    return lexer
}

func (gl *GoLexer) Tokenize() []Token {
    return gl.Lexer.Tokenize()
}
```

## 6. 应用实例

### 6.1 简单表达式词法分析

```go
func ExampleExpressionLexing() {
    input := "x = 42 + y * 3"
    
    lexer := NewGoLexer(input)
    tokens := lexer.Tokenize()
    
    for _, token := range tokens {
        fmt.Printf("Token: %s, Type: %d, Line: %d, Column: %d\n",
            token.Value, token.Type, token.Line, token.Column)
    }
}
```

### 6.2 Go代码词法分析

```go
func ExampleGoCodeLexing() {
    input := `
package main

import "fmt"

func main() {
    x := 42
    fmt.Println("Hello, World!")
}
`
    
    lexer := NewGoLexer(input)
    tokens := lexer.Tokenize()
    
    for _, token := range tokens {
        if token.Type != TokenEOF {
            fmt.Printf("%s ", token.Value)
        }
    }
    fmt.Println()
}
```

## 总结

词法分析是编译器的第一个阶段，负责将源代码字符流转换为词法单元序列。通过正则表达式和有限自动机的理论基础，我们可以构建高效准确的词法分析器。

**关键要点**：

1. **理论基础**：基于正则表达式和有限自动机的形式化理论
2. **实现技术**：DFA和NFA的转换和优化
3. **实际应用**：在编译器、解释器、文本处理等领域广泛应用
4. **Go语言特性**：充分利用Go语言的字符串处理和map数据结构

**激情澎湃的持续构建** <(￣︶￣)↗[GO!] **词法分析理论完成！** 🚀
