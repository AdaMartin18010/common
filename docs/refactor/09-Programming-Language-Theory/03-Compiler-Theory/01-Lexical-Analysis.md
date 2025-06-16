# 01-词法分析 (Lexical Analysis)

## 目录

- [1. 概述](#1-概述)
- [2. 正则表达式](#2-正则表达式)
- [3. 有限自动机](#3-有限自动机)
- [4. 词法分析器](#4-词法分析器)
- [5. Go语言实现](#5-go语言实现)
- [6. 应用实例](#6-应用实例)

## 1. 概述

词法分析是编译器的第一个阶段，将源代码字符串转换为标记(token)序列。本文档基于对编程语言理论的深度分析，建立了完整的词法分析理论体系。

### 1.1 词法分析定义

**定义 1.1** (词法分析)
词法分析是一个函数 $L: \Sigma^* \rightarrow T^*$，其中：

- $\Sigma$ 是输入字母表
- $T$ 是标记集合
- $L$ 将输入字符串映射为标记序列

**定理 1.1** (词法分析确定性)
对于任意输入字符串 $s \in \Sigma^*$，词法分析器产生唯一的标记序列。

**证明**:

1. 通过最长匹配原则
2. 标记优先级规则
3. 确定性有限自动机

### 1.2 标记定义

**定义 1.2** (标记)
标记是一个三元组 $(type, value, position)$，其中：

- $type$ 是标记类型
- $value$ 是标记值
- $position$ 是位置信息

```go
// 标记定义
type Token struct {
    Type     TokenType
    Value    string
    Position Position
    Line     int
    Column   int
}

type TokenType int

const (
    TOKEN_EOF TokenType = iota
    TOKEN_IDENTIFIER
    TOKEN_NUMBER
    TOKEN_STRING
    TOKEN_KEYWORD
    TOKEN_OPERATOR
    TOKEN_DELIMITER
    TOKEN_COMMENT
    TOKEN_WHITESPACE
)

type Position struct {
    Offset int
    Line   int
    Column int
}

// 标记类型字符串表示
func (tt TokenType) String() string {
    switch tt {
    case TOKEN_EOF:
        return "EOF"
    case TOKEN_IDENTIFIER:
        return "IDENTIFIER"
    case TOKEN_NUMBER:
        return "NUMBER"
    case TOKEN_STRING:
        return "STRING"
    case TOKEN_KEYWORD:
        return "KEYWORD"
    case TOKEN_OPERATOR:
        return "OPERATOR"
    case TOKEN_DELIMITER:
        return "DELIMITER"
    case TOKEN_COMMENT:
        return "COMMENT"
    case TOKEN_WHITESPACE:
        return "WHITESPACE"
    default:
        return "UNKNOWN"
    }
}
```

## 2. 正则表达式

### 2.1 正则表达式定义

**定义 3.1** (正则表达式)
正则表达式的语法定义为：
$$R ::= \varepsilon \mid a \mid R_1 \cdot R_2 \mid R_1 + R_2 \mid R^*$$

其中：

- $\varepsilon$ 是空字符串
- $a$ 是字母表中的符号
- $\cdot$ 是连接操作
- $+$ 是选择操作
- $*$ 是Kleene星号

**定理 3.1** (正则表达式等价性)
正则表达式与有限自动机等价。

## 3. 有限自动机

### 3.1 确定性有限自动机

**定义 2.1** (DFA)
确定性有限自动机是一个五元组 $M = (Q, \Sigma, \delta, q_0, F)$，其中：

- $Q$ 是状态集合
- $\Sigma$ 是输入字母表
- $\delta: Q \times \Sigma \rightarrow Q$ 是转移函数
- $q_0 \in Q$ 是初始状态
- $F \subseteq Q$ 是接受状态集合

**定理 2.1** (DFA接受性)
DFA $M$ 接受字符串 $w$，当且仅当存在状态序列 $q_0, q_1, \ldots, q_n$ 使得：

1. $q_0$ 是初始状态
2. $\delta(q_i, w_i) = q_{i+1}$ 对于 $0 \leq i < n$
3. $q_n \in F$

### 3.2 非确定性有限自动机

**定义 2.2** (NFA)
非确定性有限自动机是一个五元组 $M = (Q, \Sigma, \delta, q_0, F)$，其中：

- $\delta: Q \times \Sigma \rightarrow 2^Q$ 是转移函数
- 其他定义与DFA相同

**定理 2.2** (NFA到DFA转换)
对于任意NFA，存在等价的DFA。

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
