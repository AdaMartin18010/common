# 01-操作语义 (Operational Semantics)

## 目录

- [01-操作语义 (Operational Semantics)](#01-操作语义-operational-semantics)
  - [目录](#目录)
  - [1. 操作语义基础](#1-操作语义基础)
    - [1.1 操作语义定义](#11-操作语义定义)
    - [1.2 小步语义](#12-小步语义)
    - [1.3 大步语义](#13-大步语义)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 语法定义](#21-语法定义)
    - [2.2 语义规则](#22-语义规则)
    - [2.3 求值关系](#23-求值关系)
  - [3. Go语言实现](#3-go语言实现)
    - [3.1 表达式求值器](#31-表达式求值器)
    - [3.2 语句执行器](#32-语句执行器)
    - [3.3 程序解释器](#33-程序解释器)
  - [4. 应用场景](#4-应用场景)
    - [4.1 语言实现](#41-语言实现)
    - [4.2 程序验证](#42-程序验证)
    - [4.3 编译器设计](#43-编译器设计)
  - [5. 数学证明](#5-数学证明)
    - [5.1 确定性定理](#51-确定性定理)
    - [5.2 终止性定理](#52-终止性定理)
    - [5.3 等价性定理](#53-等价性定理)

---

## 1. 操作语义基础

### 1.1 操作语义定义

操作语义通过描述程序执行的具体步骤来定义编程语言的语义。它关注程序如何从初始状态转换到最终状态。

**定义 1.1**: 操作语义是一个三元组 $(\Sigma, \rightarrow, \Sigma_f)$，其中：
- $\Sigma$ 是配置集合
- $\rightarrow \subseteq \Sigma \times \Sigma$ 是转换关系
- $\Sigma_f \subseteq \Sigma$ 是最终配置集合

### 1.2 小步语义

小步语义描述程序执行的每个小步骤，通过一系列转换达到最终状态。

**定义 1.2**: 小步语义的转换关系 $\rightarrow$ 满足：
- 如果 $\sigma \rightarrow \sigma'$，则 $\sigma'$ 是 $\sigma$ 的一步转换
- 如果 $\sigma \rightarrow^* \sigma'$，则 $\sigma'$ 是 $\sigma$ 的多步转换

### 1.3 大步语义

大步语义直接描述程序从初始状态到最终状态的完整执行过程。

**定义 1.3**: 大步语义的求值关系 $\Downarrow$ 满足：
- $\sigma \Downarrow \sigma'$ 表示从配置 $\sigma$ 求值到配置 $\sigma'$

## 2. 形式化定义

### 2.1 语法定义

**定义 2.1**: 简单编程语言的语法：

$$\begin{align}
e &::= n \mid x \mid e_1 + e_2 \mid e_1 - e_2 \mid e_1 \times e_2 \\
s &::= \text{skip} \mid x := e \mid s_1; s_2 \mid \text{if } e \text{ then } s_1 \text{ else } s_2 \mid \text{while } e \text{ do } s
\end{align}$$

其中：
- $n$ 是数字
- $x$ 是变量
- $e$ 是表达式
- $s$ 是语句

### 2.2 语义规则

**定义 2.2**: 小步语义规则：

**表达式求值**:
$$\frac{n_1 + n_2 = n_3}{n_1 + n_2 \rightarrow n_3}$$

$$\frac{e_1 \rightarrow e_1'}{e_1 + e_2 \rightarrow e_1' + e_2}$$

$$\frac{e_2 \rightarrow e_2'}{n_1 + e_2 \rightarrow n_1 + e_2'}$$

**语句执行**:
$$\frac{e \rightarrow e'}{x := e \rightarrow x := e'}$$

$$\frac{}{x := n \rightarrow \text{skip}}$$

$$\frac{s_1 \rightarrow s_1'}{s_1; s_2 \rightarrow s_1'; s_2}$$

$$\frac{}{\text{skip}; s \rightarrow s}$$

### 2.3 求值关系

**定义 2.3**: 大步语义规则：

**表达式求值**:
$$\frac{e_1 \Downarrow n_1 \quad e_2 \Downarrow n_2 \quad n_1 + n_2 = n_3}{e_1 + e_2 \Downarrow n_3}$$

**语句执行**:
$$\frac{e \Downarrow n}{x := e \Downarrow [x \mapsto n]}$$

$$\frac{s_1 \Downarrow \sigma_1 \quad s_2 \Downarrow \sigma_2}{s_1; s_2 \Downarrow \sigma_2 \circ \sigma_1}$$

## 3. Go语言实现

### 3.1 表达式求值器

```go
package operationalsemantics

import (
	"fmt"
	"strconv"
)

// Expression 表示表达式
type Expression interface {
	String() string
	Evaluate(env Environment) (Value, error)
	Step(env Environment) (Expression, Environment, error)
}

// Value 表示值
type Value interface {
	String() string
	Type() string
}

// Number 数字值
type Number struct {
	Value int
}

func (n *Number) String() string {
	return strconv.Itoa(n.Value)
}

func (n *Number) Type() string {
	return "number"
}

// Variable 变量
type Variable struct {
	Name string
}

func (v *Variable) String() string {
	return v.Name
}

func (v *Variable) Evaluate(env Environment) (Value, error) {
	if value, exists := env.Get(v.Name); exists {
		return value, nil
	}
	return nil, fmt.Errorf("undefined variable: %s", v.Name)
}

func (v *Variable) Step(env Environment) (Expression, Environment, error) {
	if value, exists := env.Get(v.Name); exists {
		if num, ok := value.(*Number); ok {
			return &Number{Value: num.Value}, env, nil
		}
	}
	return nil, env, fmt.Errorf("undefined variable: %s", v.Name)
}

// Addition 加法表达式
type Addition struct {
	Left  Expression
	Right Expression
}

func (a *Addition) String() string {
	return fmt.Sprintf("(%s + %s)", a.Left, a.Right)
}

func (a *Addition) Evaluate(env Environment) (Value, error) {
	leftVal, err := a.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	rightVal, err := a.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	leftNum, ok1 := leftVal.(*Number)
	rightNum, ok2 := rightVal.(*Number)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("cannot add non-numeric values")
	}
	
	return &Number{Value: leftNum.Value + rightNum.Value}, nil
}

func (a *Addition) Step(env Environment) (Expression, Environment, error) {
	// 左操作数求值
	if leftNum, ok := a.Left.(*Number); ok {
		// 右操作数求值
		if rightNum, ok := a.Right.(*Number); ok {
			// 两个操作数都是数字，执行加法
			return &Number{Value: leftNum.Value + rightNum.Value}, env, nil
		} else {
			// 右操作数需要进一步求值
			newRight, newEnv, err := a.Right.Step(env)
			if err != nil {
				return nil, env, err
			}
			return &Addition{Left: a.Left, Right: newRight}, newEnv, nil
		}
	} else {
		// 左操作数需要进一步求值
		newLeft, newEnv, err := a.Left.Step(env)
		if err != nil {
			return nil, env, err
		}
		return &Addition{Left: newLeft, Right: a.Right}, newEnv, nil
	}
}

// Subtraction 减法表达式
type Subtraction struct {
	Left  Expression
	Right Expression
}

func (s *Subtraction) String() string {
	return fmt.Sprintf("(%s - %s)", s.Left, s.Right)
}

func (s *Subtraction) Evaluate(env Environment) (Value, error) {
	leftVal, err := s.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	rightVal, err := s.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	leftNum, ok1 := leftVal.(*Number)
	rightNum, ok2 := rightVal.(*Number)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("cannot subtract non-numeric values")
	}
	
	return &Number{Value: leftNum.Value - rightNum.Value}, nil
}

func (s *Subtraction) Step(env Environment) (Expression, Environment, error) {
	// 左操作数求值
	if leftNum, ok := s.Left.(*Number); ok {
		// 右操作数求值
		if rightNum, ok := s.Right.(*Number); ok {
			// 两个操作数都是数字，执行减法
			return &Number{Value: leftNum.Value - rightNum.Value}, env, nil
		} else {
			// 右操作数需要进一步求值
			newRight, newEnv, err := s.Right.Step(env)
			if err != nil {
				return nil, env, err
			}
			return &Subtraction{Left: s.Left, Right: newRight}, newEnv, nil
		}
	} else {
		// 左操作数需要进一步求值
		newLeft, newEnv, err := s.Left.Step(env)
		if err != nil {
			return nil, env, err
		}
		return &Subtraction{Left: newLeft, Right: s.Right}, newEnv, nil
	}
}

// Multiplication 乘法表达式
type Multiplication struct {
	Left  Expression
	Right Expression
}

func (m *Multiplication) String() string {
	return fmt.Sprintf("(%s * %s)", m.Left, m.Right)
}

func (m *Multiplication) Evaluate(env Environment) (Value, error) {
	leftVal, err := m.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	rightVal, err := m.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	
	leftNum, ok1 := leftVal.(*Number)
	rightNum, ok2 := rightVal.(*Number)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("cannot multiply non-numeric values")
	}
	
	return &Number{Value: leftNum.Value * rightNum.Value}, nil
}

func (m *Multiplication) Step(env Environment) (Expression, Environment, error) {
	// 左操作数求值
	if leftNum, ok := m.Left.(*Number); ok {
		// 右操作数求值
		if rightNum, ok := m.Right.(*Number); ok {
			// 两个操作数都是数字，执行乘法
			return &Number{Value: leftNum.Value * rightNum.Value}, env, nil
		} else {
			// 右操作数需要进一步求值
			newRight, newEnv, err := m.Right.Step(env)
			if err != nil {
				return nil, env, err
			}
			return &Multiplication{Left: m.Left, Right: newRight}, newEnv, nil
		}
	} else {
		// 左操作数需要进一步求值
		newLeft, newEnv, err := m.Left.Step(env)
		if err != nil {
			return nil, env, err
		}
		return &Multiplication{Left: newLeft, Right: m.Right}, newEnv, nil
	}
}
```

### 3.2 语句执行器

```go
// Statement 表示语句
type Statement interface {
	String() string
	Execute(env Environment) (Environment, error)
	Step(env Environment) (Statement, Environment, error)
}

// Skip 空语句
type Skip struct{}

func (s *Skip) String() string {
	return "skip"
}

func (s *Skip) Execute(env Environment) (Environment, error) {
	return env, nil
}

func (s *Skip) Step(env Environment) (Statement, Environment, error) {
	return nil, env, fmt.Errorf("skip statement cannot be stepped")
}

// Assignment 赋值语句
type Assignment struct {
	Variable string
	Value    Expression
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%s := %s", a.Variable, a.Value)
}

func (a *Assignment) Execute(env Environment) (Environment, error) {
	value, err := a.Value.Evaluate(env)
	if err != nil {
		return env, err
	}
	
	newEnv := env.Copy()
	newEnv.Set(a.Variable, value)
	return newEnv, nil
}

func (a *Assignment) Step(env Environment) (Statement, Environment, error) {
	// 如果表达式已经是值，执行赋值
	if _, ok := a.Value.(*Number); ok {
		newEnv := env.Copy()
		newEnv.Set(a.Variable, a.Value)
		return &Skip{}, newEnv, nil
	}
	
	// 否则，进一步求值表达式
	newValue, newEnv, err := a.Value.Step(env)
	if err != nil {
		return nil, env, err
	}
	
	return &Assignment{Variable: a.Variable, Value: newValue}, newEnv, nil
}

// Sequence 顺序语句
type Sequence struct {
	First  Statement
	Second Statement
}

func (s *Sequence) String() string {
	return fmt.Sprintf("%s; %s", s.First, s.Second)
}

func (s *Sequence) Execute(env Environment) (Environment, error) {
	env1, err := s.First.Execute(env)
	if err != nil {
		return env, err
	}
	
	return s.Second.Execute(env1)
}

func (s *Sequence) Step(env Environment) (Statement, Environment, error) {
	// 如果第一个语句是skip，执行第二个语句
	if _, ok := s.First.(*Skip); ok {
		return s.Second, env, nil
	}
	
	// 否则，执行第一个语句的一步
	newFirst, newEnv, err := s.First.Step(env)
	if err != nil {
		return nil, env, err
	}
	
	return &Sequence{First: newFirst, Second: s.Second}, newEnv, nil
}

// If 条件语句
type If struct {
	Condition Expression
	Then      Statement
	Else      Statement
}

func (i *If) String() string {
	return fmt.Sprintf("if %s then %s else %s", i.Condition, i.Then, i.Else)
}

func (i *If) Execute(env Environment) (Environment, error) {
	condition, err := i.Condition.Evaluate(env)
	if err != nil {
		return env, err
	}
	
	if num, ok := condition.(*Number); ok {
		if num.Value != 0 {
			return i.Then.Execute(env)
		} else {
			return i.Else.Execute(env)
		}
	}
	
	return env, fmt.Errorf("condition must evaluate to a number")
}

func (i *If) Step(env Environment) (Statement, Environment, error) {
	// 如果条件已经是值，选择分支
	if num, ok := i.Condition.(*Number); ok {
		if num.Value != 0 {
			return i.Then, env, nil
		} else {
			return i.Else, env, nil
		}
	}
	
	// 否则，进一步求值条件
	newCondition, newEnv, err := i.Condition.Step(env)
	if err != nil {
		return nil, env, err
	}
	
	return &If{Condition: newCondition, Then: i.Then, Else: i.Else}, newEnv, nil
}

// While 循环语句
type While struct {
	Condition Expression
	Body      Statement
}

func (w *While) String() string {
	return fmt.Sprintf("while %s do %s", w.Condition, w.Body)
}

func (w *While) Execute(env Environment) (Environment, error) {
	condition, err := w.Condition.Evaluate(env)
	if err != nil {
		return env, err
	}
	
	if num, ok := condition.(*Number); ok {
		if num.Value != 0 {
			// 条件为真，执行循环体
			env1, err := w.Body.Execute(env)
			if err != nil {
				return env, err
			}
			// 递归执行循环
			return w.Execute(env1)
		} else {
			// 条件为假，结束循环
			return env, nil
		}
	}
	
	return env, fmt.Errorf("condition must evaluate to a number")
}

func (w *While) Step(env Environment) (Statement, Environment, error) {
	// 如果条件已经是值，选择分支
	if num, ok := w.Condition.(*Number); ok {
		if num.Value != 0 {
			// 条件为真，展开为循环体后跟循环
			return &Sequence{First: w.Body, Second: w}, env, nil
		} else {
			// 条件为假，结束循环
			return &Skip{}, env, nil
		}
	}
	
	// 否则，进一步求值条件
	newCondition, newEnv, err := w.Condition.Step(env)
	if err != nil {
		return nil, env, err
	}
	
	return &While{Condition: newCondition, Body: w.Body}, newEnv, nil
}
```

### 3.3 程序解释器

```go
// Environment 环境
type Environment interface {
	Get(name string) (Value, bool)
	Set(name string, value Value)
	Copy() Environment
}

// SimpleEnvironment 简单环境实现
type SimpleEnvironment struct {
	variables map[string]Value
}

// NewEnvironment 创建新环境
func NewEnvironment() Environment {
	return &SimpleEnvironment{
		variables: make(map[string]Value),
	}
}

// Get 获取变量值
func (e *SimpleEnvironment) Get(name string) (Value, bool) {
	value, exists := e.variables[name]
	return value, exists
}

// Set 设置变量值
func (e *SimpleEnvironment) Set(name string, value Value) {
	e.variables[name] = value
}

// Copy 复制环境
func (e *SimpleEnvironment) Copy() Environment {
	newEnv := &SimpleEnvironment{
		variables: make(map[string]Value),
	}
	
	for name, value := range e.variables {
		newEnv.variables[name] = value
	}
	
	return newEnv
}

// Interpreter 解释器
type Interpreter struct {
	env Environment
}

// NewInterpreter 创建新解释器
func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

// EvaluateExpression 求值表达式
func (i *Interpreter) EvaluateExpression(expr Expression) (Value, error) {
	return expr.Evaluate(i.env)
}

// ExecuteStatement 执行语句
func (i *Interpreter) ExecuteStatement(stmt Statement) error {
	newEnv, err := stmt.Execute(i.env)
	if err != nil {
		return err
	}
	i.env = newEnv
	return nil
}

// StepExpression 表达式单步求值
func (i *Interpreter) StepExpression(expr Expression) (Expression, error) {
	newExpr, newEnv, err := expr.Step(i.env)
	if err != nil {
		return nil, err
	}
	i.env = newEnv
	return newExpr, nil
}

// StepStatement 语句单步执行
func (i *Interpreter) StepStatement(stmt Statement) (Statement, error) {
	newStmt, newEnv, err := stmt.Step(i.env)
	if err != nil {
		return nil, err
	}
	i.env = newEnv
	return newStmt, nil
}

// GetEnvironment 获取当前环境
func (i *Interpreter) GetEnvironment() Environment {
	return i.env
}

// SetEnvironment 设置环境
func (i *Interpreter) SetEnvironment(env Environment) {
	i.env = env
}

// TraceExecution 跟踪执行
func (i *Interpreter) TraceExecution(stmt Statement) []string {
	trace := make([]string, 0)
	currentStmt := stmt
	
	for {
		trace = append(trace, fmt.Sprintf("State: %s", currentStmt))
		
		newStmt, err := i.StepStatement(currentStmt)
		if err != nil {
			trace = append(trace, fmt.Sprintf("Error: %s", err))
			break
		}
		
		if newStmt == nil {
			trace = append(trace, "Terminated")
			break
		}
		
		currentStmt = newStmt
	}
	
	return trace
}
```

## 4. 应用场景

### 4.1 语言实现

```go
// LanguageImplementation 语言实现
type LanguageImplementation struct {
	interpreter *Interpreter
}

// NewLanguageImplementation 创建语言实现
func NewLanguageImplementation() *LanguageImplementation {
	return &LanguageImplementation{
		interpreter: NewInterpreter(),
	}
}

// ParseAndExecute 解析并执行程序
func (li *LanguageImplementation) ParseAndExecute(program string) error {
	// 简化的解析器，实际实现需要更复杂的词法和语法分析
	stmt, err := li.parseProgram(program)
	if err != nil {
		return err
	}
	
	return li.interpreter.ExecuteStatement(stmt)
}

// parseProgram 解析程序
func (li *LanguageImplementation) parseProgram(program string) (Statement, error) {
	// 简化的解析实现
	// 实际实现需要词法分析器和语法分析器
	
	// 示例：解析 "x := 5; y := x + 3"
	if program == "x := 5; y := x + 3" {
		stmt1 := &Assignment{
			Variable: "x",
			Value:    &Number{Value: 5},
		}
		
		stmt2 := &Assignment{
			Variable: "y",
			Value: &Addition{
				Left:  &Variable{Name: "x"},
				Right: &Number{Value: 3},
			},
		}
		
		return &Sequence{First: stmt1, Second: stmt2}, nil
	}
	
	return nil, fmt.Errorf("unsupported program: %s", program)
}

// EvaluateExpression 求值表达式
func (li *LanguageImplementation) EvaluateExpression(exprStr string) (Value, error) {
	expr, err := li.parseExpression(exprStr)
	if err != nil {
		return nil, err
	}
	
	return li.interpreter.EvaluateExpression(expr)
}

// parseExpression 解析表达式
func (li *LanguageImplementation) parseExpression(exprStr string) (Expression, error) {
	// 简化的表达式解析
	// 实际实现需要更复杂的解析逻辑
	
	// 示例：解析 "x + 5"
	if exprStr == "x + 5" {
		return &Addition{
			Left:  &Variable{Name: "x"},
			Right: &Number{Value: 5},
		}, nil
	}
	
	return nil, fmt.Errorf("unsupported expression: %s", exprStr)
}
```

### 4.2 程序验证

```go
// ProgramVerifier 程序验证器
type ProgramVerifier struct {
	interpreter *Interpreter
}

// NewProgramVerifier 创建程序验证器
func NewProgramVerifier() *ProgramVerifier {
	return &ProgramVerifier{
		interpreter: NewInterpreter(),
	}
}

// VerifyTermination 验证程序终止性
func (pv *ProgramVerifier) VerifyTermination(stmt Statement) bool {
	// 简化的终止性验证
	// 实际实现需要更复杂的静态分析
	
	// 检查是否包含无限循环
	return pv.checkInfiniteLoop(stmt)
}

// checkInfiniteLoop 检查无限循环
func (pv *ProgramVerifier) checkInfiniteLoop(stmt Statement) bool {
	switch s := stmt.(type) {
	case *While:
		// 检查循环条件是否可能永远为真
		return pv.isAlwaysTrue(s.Condition)
	case *Sequence:
		return pv.checkInfiniteLoop(s.First) || pv.checkInfiniteLoop(s.Second)
	case *If:
		return pv.checkInfiniteLoop(s.Then) || pv.checkInfiniteLoop(s.Else)
	default:
		return false
	}
}

// isAlwaysTrue 检查表达式是否永远为真
func (pv *ProgramVerifier) isAlwaysTrue(expr Expression) bool {
	// 简化的实现，实际需要更复杂的分析
	if num, ok := expr.(*Number); ok {
		return num.Value != 0
	}
	return false
}

// VerifyCorrectness 验证程序正确性
func (pv *ProgramVerifier) VerifyCorrectness(stmt Statement, precondition, postcondition func(Environment) bool) bool {
	// 简化的正确性验证
	// 实际实现需要更复杂的程序逻辑
	
	// 检查前置条件
	if !precondition(pv.interpreter.GetEnvironment()) {
		return false
	}
	
	// 执行程序
	err := pv.interpreter.ExecuteStatement(stmt)
	if err != nil {
		return false
	}
	
	// 检查后置条件
	return postcondition(pv.interpreter.GetEnvironment())
}

// GenerateTestCases 生成测试用例
func (pv *ProgramVerifier) GenerateTestCases(stmt Statement) []Environment {
	testCases := make([]Environment, 0)
	
	// 生成不同的初始环境
	env1 := NewEnvironment()
	env1.Set("x", &Number{Value: 0})
	testCases = append(testCases, env1)
	
	env2 := NewEnvironment()
	env2.Set("x", &Number{Value: 5})
	testCases = append(testCases, env2)
	
	env3 := NewEnvironment()
	env3.Set("x", &Number{Value: -1})
	testCases = append(testCases, env3)
	
	return testCases
}
```

### 4.3 编译器设计

```go
// Compiler 编译器
type Compiler struct {
	interpreter *Interpreter
}

// NewCompiler 创建编译器
func NewCompiler() *Compiler {
	return &Compiler{
		interpreter: NewInterpreter(),
	}
}

// CompileToGo 编译为Go代码
func (c *Compiler) CompileToGo(stmt Statement) string {
	return c.generateGoCode(stmt, 0)
}

// generateGoCode 生成Go代码
func (c *Compiler) generateGoCode(stmt Statement, indent int) string {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "    "
	}
	
	switch s := stmt.(type) {
	case *Skip:
		return indentStr + "// skip\n"
		
	case *Assignment:
		exprCode := c.generateExpressionCode(s.Value)
		return fmt.Sprintf("%s%s = %s\n", indentStr, s.Variable, exprCode)
		
	case *Sequence:
		firstCode := c.generateGoCode(s.First, indent)
		secondCode := c.generateGoCode(s.Second, indent)
		return firstCode + secondCode
		
	case *If:
		conditionCode := c.generateExpressionCode(s.Condition)
		thenCode := c.generateGoCode(s.Then, indent+1)
		elseCode := c.generateGoCode(s.Else, indent+1)
		
		return fmt.Sprintf("%sif %s != 0 {\n%s%s} else {\n%s%s}\n",
			indentStr, conditionCode, thenCode, indentStr, elseCode, indentStr)
		
	case *While:
		conditionCode := c.generateExpressionCode(s.Condition)
		bodyCode := c.generateGoCode(s.Body, indent+1)
		
		return fmt.Sprintf("%sfor %s != 0 {\n%s%s}\n",
			indentStr, conditionCode, bodyCode, indentStr)
		
	default:
		return indentStr + "// unknown statement\n"
	}
}

// generateExpressionCode 生成表达式代码
func (c *Compiler) generateExpressionCode(expr Expression) string {
	switch e := expr.(type) {
	case *Number:
		return fmt.Sprintf("%d", e.Value)
		
	case *Variable:
		return e.Name
		
	case *Addition:
		leftCode := c.generateExpressionCode(e.Left)
		rightCode := c.generateExpressionCode(e.Right)
		return fmt.Sprintf("(%s + %s)", leftCode, rightCode)
		
	case *Subtraction:
		leftCode := c.generateExpressionCode(e.Left)
		rightCode := c.generateExpressionCode(e.Right)
		return fmt.Sprintf("(%s - %s)", leftCode, rightCode)
		
	case *Multiplication:
		leftCode := c.generateExpressionCode(e.Left)
		rightCode := c.generateExpressionCode(e.Right)
		return fmt.Sprintf("(%s * %s)", leftCode, rightCode)
		
	default:
		return "unknown_expression"
	}
}

// CompileToBytecode 编译为字节码
func (c *Compiler) CompileToBytecode(stmt Statement) []byte {
	bytecode := make([]byte, 0)
	c.generateBytecode(stmt, &bytecode)
	return bytecode
}

// generateBytecode 生成字节码
func (c *Compiler) generateBytecode(stmt Statement, bytecode *[]byte) {
	switch s := stmt.(type) {
	case *Skip:
		*bytecode = append(*bytecode, 0x00) // NOP
		
	case *Assignment:
		// 生成表达式字节码
		c.generateExpressionBytecode(s.Value, bytecode)
		// 生成赋值字节码
		*bytecode = append(*bytecode, 0x01) // STORE
		*bytecode = append(*bytecode, byte(len(s.Variable)))
		*bytecode = append(*bytecode, []byte(s.Variable)...)
		
	case *Sequence:
		c.generateBytecode(s.First, bytecode)
		c.generateBytecode(s.Second, bytecode)
		
	case *If:
		// 生成条件字节码
		c.generateExpressionBytecode(s.Condition, bytecode)
		*bytecode = append(*bytecode, 0x02) // JZ
		
		// 记录跳转位置
		jumpPos := len(*bytecode)
		*bytecode = append(*bytecode, 0, 0, 0, 0) // 占位符
		
		// 生成then分支字节码
		thenStart := len(*bytecode)
		c.generateBytecode(s.Then, bytecode)
		thenEnd := len(*bytecode)
		
		// 生成else分支字节码
		elseStart := len(*bytecode)
		c.generateBytecode(s.Else, bytecode)
		elseEnd := len(*bytecode)
		
		// 更新跳转地址
		(*bytecode)[jumpPos] = byte(elseStart & 0xFF)
		(*bytecode)[jumpPos+1] = byte((elseStart >> 8) & 0xFF)
		(*bytecode)[jumpPos+2] = byte((elseStart >> 16) & 0xFF)
		(*bytecode)[jumpPos+3] = byte((elseStart >> 24) & 0xFF)
		
	case *While:
		loopStart := len(*bytecode)
		
		// 生成条件字节码
		c.generateExpressionBytecode(s.Condition, bytecode)
		*bytecode = append(*bytecode, 0x03) // JZ (跳出循环)
		
		// 记录跳出位置
		breakPos := len(*bytecode)
		*bytecode = append(*bytecode, 0, 0, 0, 0) // 占位符
		
		// 生成循环体字节码
		bodyStart := len(*bytecode)
		c.generateBytecode(s.Body, bytecode)
		
		// 跳回循环开始
		*bytecode = append(*bytecode, 0x04) // JMP
		*bytecode = append(*bytecode, byte(loopStart & 0xFF))
		*bytecode = append(*bytecode, byte((loopStart >> 8) & 0xFF))
		*bytecode = append(*bytecode, byte((loopStart >> 16) & 0xFF))
		*bytecode = append(*bytecode, byte((loopStart >> 24) & 0xFF))
		
		// 更新跳出地址
		breakAddr := len(*bytecode)
		(*bytecode)[breakPos] = byte(breakAddr & 0xFF)
		(*bytecode)[breakPos+1] = byte((breakAddr >> 8) & 0xFF)
		(*bytecode)[breakPos+2] = byte((breakAddr >> 16) & 0xFF)
		(*bytecode)[breakPos+3] = byte((breakAddr >> 24) & 0xFF)
	}
}

// generateExpressionBytecode 生成表达式字节码
func (c *Compiler) generateExpressionBytecode(expr Expression, bytecode *[]byte) {
	switch e := expr.(type) {
	case *Number:
		*bytecode = append(*bytecode, 0x10) // PUSH
		*bytecode = append(*bytecode, byte(e.Value & 0xFF))
		*bytecode = append(*bytecode, byte((e.Value >> 8) & 0xFF))
		*bytecode = append(*bytecode, byte((e.Value >> 16) & 0xFF))
		*bytecode = append(*bytecode, byte((e.Value >> 24) & 0xFF))
		
	case *Variable:
		*bytecode = append(*bytecode, 0x11) // LOAD
		*bytecode = append(*bytecode, byte(len(e.Name)))
		*bytecode = append(*bytecode, []byte(e.Name)...)
		
	case *Addition:
		c.generateExpressionBytecode(e.Left, bytecode)
		c.generateExpressionBytecode(e.Right, bytecode)
		*bytecode = append(*bytecode, 0x20) // ADD
		
	case *Subtraction:
		c.generateExpressionBytecode(e.Left, bytecode)
		c.generateExpressionBytecode(e.Right, bytecode)
		*bytecode = append(*bytecode, 0x21) // SUB
		
	case *Multiplication:
		c.generateExpressionBytecode(e.Left, bytecode)
		c.generateExpressionBytecode(e.Right, bytecode)
		*bytecode = append(*bytecode, 0x22) // MUL
	}
}
```

## 5. 数学证明

### 5.1 确定性定理

**定理 5.1** (求值确定性): 对于任意表达式 $e$ 和环境 $\sigma$，如果 $e \rightarrow e_1$ 且 $e \rightarrow e_2$，则 $e_1 = e_2$。

**证明**:
1. 通过结构归纳法证明
2. 对每种表达式类型，证明求值规则是确定的
3. 由于规则是语法导向的，每个表达式最多匹配一个规则

**定理 5.2** (执行确定性): 对于任意语句 $s$ 和环境 $\sigma$，如果 $s \rightarrow s_1$ 且 $s \rightarrow s_2$，则 $s_1 = s_2$。

**证明**:
1. 通过结构归纳法证明
2. 对每种语句类型，证明执行规则是确定的
3. 由于规则是语法导向的，每个语句最多匹配一个规则

### 5.2 终止性定理

**定理 5.3** (表达式终止性): 对于任意表达式 $e$，存在有限步求值序列 $e \rightarrow^* v$，其中 $v$ 是值。

**证明**:
1. 定义表达式的大小度量
2. 证明每步求值都减少表达式大小
3. 由于大小有限，求值必然终止

**定理 5.4** (语句终止性): 对于任意语句 $s$，如果 $s$ 不包含无限循环，则存在有限步执行序列 $s \rightarrow^* \text{skip}$。

**证明**:
1. 定义语句的复杂度度量
2. 证明每步执行都减少复杂度
3. 由于复杂度有限，执行必然终止

### 5.3 等价性定理

**定理 5.5** (小步与大步等价性): 对于任意表达式 $e$ 和值 $v$，$e \rightarrow^* v$ 当且仅当 $e \Downarrow v$。

**证明**:
1. 证明小步语义蕴含大步语义
2. 证明大步语义蕴含小步语义
3. 通过结构归纳法完成证明

**定理 5.6** (上下文等价性): 对于任意表达式 $e_1, e_2$，如果 $e_1 \rightarrow e_2$，则对任意上下文 $C$，$C[e_1] \rightarrow C[e_2]$。

**证明**:
1. 定义上下文的概念
2. 证明求值规则在上下文中保持
3. 通过结构归纳法完成证明

---

**总结**: 操作语义为编程语言提供了精确的执行模型，通过Go语言实现，我们可以构建实用的解释器和编译器，用于语言实现、程序验证和编译器设计。 