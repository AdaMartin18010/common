# 01-Hello World (Hello World)

## 概述

Hello World是编程学习的第一个程序，展示了Go语言的基本语法和程序结构。本文档介绍Go语言的基础概念、程序结构以及各种Hello World变体。

## 目录

1. [基本Hello World程序](#1-基本hello-world程序)
2. [程序结构分析](#2-程序结构分析)
3. [变量和数据类型](#3-变量和数据类型)
4. [函数定义](#4-函数定义)
5. [包管理](#5-包管理)
6. [高级示例](#6-高级示例)

---

## 1. 基本Hello World程序

### 1.1 最简单的Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**程序分析**:
- `package main`: 声明这是一个可执行程序的主包
- `import "fmt"`: 导入格式化输出包
- `func main()`: 程序入口函数
- `fmt.Println()`: 打印并换行

### 1.2 运行程序

```bash
# 编译并运行
go run hello.go

# 编译为可执行文件
go build hello.go
./hello

# 交叉编译
GOOS=linux GOARCH=amd64 go build hello.go
```

### 1.3 Go语言特点

**优势**:
1. **简洁语法** - 自动分号插入，清晰的代码结构
2. **强类型** - 编译时类型检查，减少运行时错误
3. **垃圾回收** - 自动内存管理
4. **并发支持** - goroutine和channel
5. **标准库丰富** - 内置网络、文件、加密等库

---

## 2. 程序结构分析

### 2.1 包声明

```go
package main
```

**包的作用**:
- 组织代码结构
- 控制可见性
- 避免命名冲突
- 支持模块化开发

### 2.2 导入声明

```go
// 单个包导入
import "fmt"

// 多个包导入
import (
    "fmt"
    "os"
    "strings"
)

// 别名导入
import f "fmt"

// 点导入（不推荐）
import . "fmt"

// 仅执行init函数
import _ "database/sql/driver"
```

### 2.3 函数声明

```go
// 基本函数
func main() {
    // 函数体
}

// 带参数的函数
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 带返回值的函数
func getGreeting() string {
    return "Hello, World!"
}

// 多返回值函数
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

---

## 3. 变量和数据类型

### 3.1 变量声明

```go
package main

import "fmt"

func main() {
    // 方式1：var关键字
    var message string = "Hello, World!"
    
    // 方式2：短变量声明
    greeting := "Hello, Go!"
    
    // 方式3：多变量声明
    var (
        name    string = "Alice"
        age     int    = 25
        isAdult bool   = true
    )
    
    // 方式4：类型推断
    var count = 42
    var price = 19.99
    
    fmt.Println(message)
    fmt.Println(greeting)
    fmt.Printf("Name: %s, Age: %d, Adult: %t\n", name, age, isAdult)
    fmt.Printf("Count: %d, Price: %.2f\n", count, price)
}
```

### 3.2 基本数据类型

```go
package main

import "fmt"

func main() {
    // 整数类型
    var int8Val int8 = 127
    var int16Val int16 = 32767
    var int32Val int32 = 2147483647
    var int64Val int64 = 9223372036854775807
    var intVal int = 42 // 平台相关
    
    // 无符号整数
    var uint8Val uint8 = 255
    var uintVal uint = 42
    
    // 浮点数
    var float32Val float32 = 3.14
    var float64Val float64 = 3.14159265359
    
    // 复数
    var complex64Val complex64 = 1 + 2i
    var complex128Val complex128 = 1.5 + 2.5i
    
    // 布尔类型
    var boolVal bool = true
    
    // 字符串
    var stringVal string = "Hello, World!"
    
    // 字节（uint8的别名）
    var byteVal byte = 'A'
    
    // rune（int32的别名，表示Unicode码点）
    var runeVal rune = '中'
    
    fmt.Printf("整数: %d, %d, %d, %d, %d\n", int8Val, int16Val, int32Val, int64Val, intVal)
    fmt.Printf("无符号整数: %d, %d\n", uint8Val, uintVal)
    fmt.Printf("浮点数: %.2f, %.10f\n", float32Val, float64Val)
    fmt.Printf("复数: %v, %v\n", complex64Val, complex128Val)
    fmt.Printf("布尔: %t\n", boolVal)
    fmt.Printf("字符串: %s\n", stringVal)
    fmt.Printf("字节: %c\n", byteVal)
    fmt.Printf("rune: %c\n", runeVal)
}
```

### 3.3 类型转换

```go
package main

import "fmt"

func main() {
    // 显式类型转换
    var intVal int = 42
    var floatVal float64 = float64(intVal)
    var stringVal string = string(intVal) // 转换为Unicode字符
    
    fmt.Printf("int: %d, float64: %.1f, string: %s\n", intVal, floatVal, stringVal)
    
    // 字符串转换
    var str string = "123"
    var num int = 456
    
    // 注意：Go没有自动类型转换，需要显式转换
    // var result = str + num // 编译错误
    
    // 正确的做法
    var result string = str + fmt.Sprintf("%d", num)
    fmt.Println("Result:", result)
}
```

---

## 4. 函数定义

### 4.1 函数基础

```go
package main

import "fmt"

// 基本函数
func sayHello() {
    fmt.Println("Hello, World!")
}

// 带参数的函数
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 带返回值的函数
func getGreeting() string {
    return "Hello, World!"
}

// 带参数的返回值函数
func createGreeting(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

// 多返回值函数
func divideAndRemainder(a, b int) (int, int) {
    return a / b, a % b
}

// 命名返回值
func divideAndRemainderNamed(a, b int) (quotient, remainder int) {
    quotient = a / b
    remainder = a % b
    return // 裸返回
}

// 可变参数函数
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

func main() {
    sayHello()
    greet("Alice")
    
    greeting := getGreeting()
    fmt.Println(greeting)
    
    customGreeting := createGreeting("Bob")
    fmt.Println(customGreeting)
    
    q, r := divideAndRemainder(17, 5)
    fmt.Printf("17 ÷ 5 = %d remainder %d\n", q, r)
    
    q2, r2 := divideAndRemainderNamed(23, 7)
    fmt.Printf("23 ÷ 7 = %d remainder %d\n", q2, r2)
    
    total := sum(1, 2, 3, 4, 5)
    fmt.Printf("Sum: %d\n", total)
}
```

### 4.2 函数类型

```go
package main

import "fmt"

// 函数类型定义
type Greeter func(string) string

// 高阶函数
func createGreeter(prefix string) Greeter {
    return func(name string) string {
        return fmt.Sprintf("%s, %s!", prefix, name)
    }
}

// 函数作为参数
func applyGreeter(greeter Greeter, names []string) {
    for _, name := range names {
        fmt.Println(greeter(name))
    }
}

func main() {
    // 使用函数类型
    var greeter Greeter = func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    }
    
    fmt.Println(greeter("Alice"))
    
    // 高阶函数
    formalGreeter := createGreeter("Good morning")
    casualGreeter := createGreeter("Hey")
    
    names := []string{"Alice", "Bob", "Charlie"}
    
    fmt.Println("Formal greetings:")
    applyGreeter(formalGreeter, names)
    
    fmt.Println("Casual greetings:")
    applyGreeter(casualGreeter, names)
}
```

---

## 5. 包管理

### 5.1 创建自定义包

**greeting/greeting.go**:
```go
package greeting

import "fmt"

// Greet 公开函数（首字母大写）
func Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

// greet 私有函数（首字母小写）
func greet(name string) string {
    return fmt.Sprintf("Hi, %s!", name)
}

// Greeter 结构体
type Greeter struct {
    Name string
}

// NewGreeter 构造函数
func NewGreeter(name string) *Greeter {
    return &Greeter{Name: name}
}

// Greet 方法
func (g *Greeter) Greet() string {
    return fmt.Sprintf("Hello, %s!", g.Name)
}
```

**main.go**:
```go
package main

import (
    "fmt"
    "./greeting"
)

func main() {
    // 使用包中的函数
    message := greeting.Greet("World")
    fmt.Println(message)
    
    // 使用包中的结构体
    greeter := greeting.NewGreeter("Alice")
    fmt.Println(greeter.Greet())
    
    // 注意：greeting.greet("World") 会编译错误，因为greet是私有的
}
```

### 5.2 Go Modules

**go.mod**:
```go
module hello-world

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/sirupsen/logrus v1.9.3
)
```

**使用外部包**:
```go
package main

import (
    "fmt"
    "github.com/sirupsen/logrus"
)

func main() {
    log := logrus.New()
    log.SetLevel(logrus.InfoLevel)
    
    log.Info("Starting Hello World application")
    
    fmt.Println("Hello, World!")
    
    log.Info("Hello World application completed")
}
```

---

## 6. 高级示例

### 6.1 并发Hello World

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func sayHello(name string, wg *sync.WaitGroup) {
    defer wg.Done()
    time.Sleep(100 * time.Millisecond) // 模拟工作
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
    
    var wg sync.WaitGroup
    
    fmt.Println("Starting concurrent greetings...")
    
    for _, name := range names {
        wg.Add(1)
        go sayHello(name, &wg)
    }
    
    wg.Wait()
    fmt.Println("All greetings completed!")
}
```

### 6.2 HTTP Hello World

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/greet", greetHandler)
    
    fmt.Println("Server starting on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 6.3 文件操作Hello World

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
)

func main() {
    // 写入文件
    message := "Hello, World!"
    err := ioutil.WriteFile("hello.txt", []byte(message), 0644)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Message written to hello.txt")
    
    // 读取文件
    content, err := ioutil.ReadFile("hello.txt")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Read from file: %s\n", string(content))
    
    // 使用os包
    file, err := os.Create("hello2.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    _, err = file.WriteString("Hello, World from os package!\n")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Message written to hello2.txt")
}
```

### 6.4 JSON Hello World

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

// Person 结构体
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Message string `json:"message"`
}

func main() {
    // 创建Person实例
    person := Person{
        Name:    "Alice",
        Age:     25,
        Message: "Hello, World!",
    }
    
    // 序列化为JSON
    jsonData, err := json.Marshal(person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("JSON: %s\n", string(jsonData))
    
    // 格式化JSON
    prettyJSON, err := json.MarshalIndent(person, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Pretty JSON:\n%s\n", string(prettyJSON))
    
    // 反序列化JSON
    var newPerson Person
    err = json.Unmarshal(jsonData, &newPerson)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Deserialized: %+v\n", newPerson)
}
```

### 6.5 测试Hello World

**hello_test.go**:
```go
package main

import (
    "testing"
)

func TestGreet(t *testing.T) {
    expected := "Hello, World!"
    result := getGreeting()
    
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}

func TestCreateGreeting(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"Alice", "Alice", "Hello, Alice!"},
        {"Bob", "Bob", "Hello, Bob!"},
        {"", "", "Hello, !"},
    }
    
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := createGreeting(test.input)
            if result != test.expected {
                t.Errorf("Expected %s, got %s", test.expected, result)
            }
        })
    }
}

// 基准测试
func BenchmarkGreet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        getGreeting()
    }
}
```

**运行测试**:
```bash
# 运行所有测试
go test

# 运行测试并显示详细信息
go test -v

# 运行基准测试
go test -bench=.

# 运行测试并生成覆盖率报告
go test -cover
```

---

## 总结

本文档介绍了Go语言Hello World程序的各种变体和基础概念：

1. **基本程序结构** - 包声明、导入、函数定义
2. **变量和数据类型** - 各种数据类型的声明和使用
3. **函数定义** - 基本函数、多返回值、高阶函数
4. **包管理** - 自定义包、Go Modules
5. **高级示例** - 并发、HTTP、文件操作、JSON、测试

这些示例展示了Go语言的基本语法和常用功能，为后续学习奠定基础。

---

**相关链接**:
- [02-数据结构 (Data Structures)](02-Data-Structures.md)
- [03-算法实现 (Algorithm Implementation)](03-Algorithm-Implementation.md)
- [04-并发编程 (Concurrent Programming)](04-Concurrent-Programming.md)
- [02-应用示例 (Application Examples)](../02-Application-Examples/README.md) 