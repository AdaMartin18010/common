# 03-静态分析 (Static Analysis)

## 1. 概述

### 1.1 定义与目标

**静态分析**是在不执行程序的情况下，通过分析源代码来发现潜在错误、安全漏洞和代码质量问题的方法。

**形式化定义**：
设 $P$ 为程序，$S$ 为程序状态空间，静态分析函数 $A: P \rightarrow 2^S$ 满足：
$$A(P) = \{s \in S | \exists \text{ execution path } \pi: s \in \pi\}$$

### 1.2 理论基础

#### 1.2.1 抽象解释理论

抽象解释是静态分析的理论基础，通过抽象域来近似程序行为：

**定义**：抽象域 $(D, \sqsubseteq, \sqcup, \sqcap)$ 是一个完全格，其中：
- $D$ 是抽象值的集合
- $\sqsubseteq$ 是偏序关系
- $\sqcup$ 是上确界操作
- $\sqcap$ 是下确界操作

**Galois连接**：
$$(\alpha, \gamma): \mathcal{P}(S) \leftrightarrow D$$
其中 $\alpha$ 是抽象函数，$\gamma$ 是具体化函数。

## 2. 数据流分析

### 2.1 可达定义分析

#### 2.1.1 理论基础

**定义**：变量 $v$ 在程序点 $p$ 处的定义 $d$ 可达，当且仅当存在从 $d$ 到 $p$ 的路径，且该路径上 $v$ 没有被重新定义。

**形式化定义**：
$$RD(p) = \bigcup_{q \in pred(p)} (RD(q) \setminus kill(q) \cup gen(q))$$

其中：
- $RD(p)$ 是程序点 $p$ 处的可达定义集合
- $pred(p)$ 是 $p$ 的前驱节点集合
- $kill(q)$ 是在 $q$ 处被杀死（覆盖）的定义
- $gen(q)$ 是在 $q$ 处生成的新定义

#### 2.1.2 Go语言实现

```go
package staticanalysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Definition 表示变量定义
type Definition struct {
	Variable string
	Line     int
	Position token.Pos
}

// ProgramPoint 表示程序点
type ProgramPoint struct {
	Node ast.Node
	Line int
}

// ReachingDefinitionAnalysis 可达定义分析器
type ReachingDefinitionAnalysis struct {
	definitions map[string][]Definition
	reaching    map[ProgramPoint]map[string]Definition
}

// NewReachingDefinitionAnalysis 创建新的可达定义分析器
func NewReachingDefinitionAnalysis() *ReachingDefinitionAnalysis {
	return &ReachingDefinitionAnalysis{
		definitions: make(map[string][]Definition),
		reaching:    make(map[ProgramPoint]map[string]Definition),
	}
}

// Analyze 分析Go源代码
func (rda *ReachingDefinitionAnalysis) Analyze(sourceCode string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析失败: %v", err)
	}

	// 收集所有定义
	rda.collectDefinitions(node, fset)
	
	// 计算可达定义
	rda.computeReachingDefinitions(node, fset)
	
	return nil
}

// collectDefinitions 收集所有变量定义
func (rda *ReachingDefinitionAnalysis) collectDefinitions(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			for i, lhs := range x.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					def := Definition{
						Variable: ident.Name,
						Line:     fset.Position(x.Pos()).Line,
						Position: x.Pos(),
					}
					rda.definitions[ident.Name] = append(rda.definitions[ident.Name], def)
				}
			}
		case *ast.ValueSpec:
			for _, name := range x.Names {
				def := Definition{
					Variable: name.Name,
					Line:     fset.Position(x.Pos()).Line,
					Position: x.Pos(),
				}
				rda.definitions[name.Name] = append(rda.definitions[name.Name], def)
			}
		}
		return true
	})
}

// computeReachingDefinitions 计算可达定义
func (rda *ReachingDefinitionAnalysis) computeReachingDefinitions(node ast.Node, fset *token.FileSet) {
	// 初始化所有程序点
	rda.initializeProgramPoints(node, fset)
	
	// 迭代计算直到收敛
	for rda.iterate() {
		// 继续迭代
	}
}

// initializeProgramPoints 初始化程序点
func (rda *ReachingDefinitionAnalysis) initializeProgramPoints(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			point := ProgramPoint{
				Node: n,
				Line: fset.Position(n.Pos()).Line,
			}
			rda.reaching[point] = make(map[string]Definition)
		}
		return true
	})
}

// iterate 执行一次迭代
func (rda *ReachingDefinitionAnalysis) iterate() bool {
	changed := false
	// 实现迭代逻辑
	return changed
}

// GetReachingDefinitions 获取指定程序点的可达定义
func (rda *ReachingDefinitionAnalysis) GetReachingDefinitions(variable string, line int) []Definition {
	var result []Definition
	for point, defs := range rda.reaching {
		if point.Line == line {
			if def, exists := defs[variable]; exists {
				result = append(result, def)
			}
		}
	}
	return result
}

// 泛型实现
type AnalysisResult[T any] struct {
	Value    T
	Confidence float64
	Metadata   map[string]interface{}
}

// GenericStaticAnalyzer 泛型静态分析器
type GenericStaticAnalyzer[T any] struct {
	analysisFunc func(ast.Node) T
	results      map[ProgramPoint]AnalysisResult[T]
}

// NewGenericStaticAnalyzer 创建泛型静态分析器
func NewGenericStaticAnalyzer[T any](analysisFunc func(ast.Node) T) *GenericStaticAnalyzer[T] {
	return &GenericStaticAnalyzer[T]{
		analysisFunc: analysisFunc,
		results:      make(map[ProgramPoint]AnalysisResult[T]),
	}
}

// Analyze 执行泛型分析
func (g *GenericStaticAnalyzer[T]) Analyze(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			point := ProgramPoint{Node: n}
			result := g.analysisFunc(n)
			g.results[point] = AnalysisResult[T]{
				Value:      result,
				Confidence: 1.0,
				Metadata:   make(map[string]interface{}),
			}
		}
		return true
	})
}
```

### 2.2 活跃变量分析

#### 2.2.1 理论基础

**定义**：变量 $v$ 在程序点 $p$ 处活跃，当且仅当存在从 $p$ 开始的路径，在该路径上 $v$ 被使用且在此之前 $v$ 没有被重新定义。

**形式化定义**：
$$LV(p) = \bigcup_{q \in succ(p)} (LV(q) \setminus def(q) \cup use(q))$$

其中：
- $LV(p)$ 是程序点 $p$ 处的活跃变量集合
- $succ(p)$ 是 $p$ 的后继节点集合
- $def(q)$ 是在 $q$ 处定义的变量
- $use(q)$ 是在 $q$ 处使用的变量

#### 2.2.2 Go语言实现

```go
// LiveVariableAnalysis 活跃变量分析器
type LiveVariableAnalysis struct {
	liveVars map[ProgramPoint]map[string]bool
	use      map[ProgramPoint]map[string]bool
	def      map[ProgramPoint]map[string]bool
}

// NewLiveVariableAnalysis 创建活跃变量分析器
func NewLiveVariableAnalysis() *LiveVariableAnalysis {
	return &LiveVariableAnalysis{
		liveVars: make(map[ProgramPoint]map[string]bool),
		use:      make(map[ProgramPoint]map[string]bool),
		def:      make(map[ProgramPoint]map[string]bool),
	}
}

// Analyze 分析活跃变量
func (lva *LiveVariableAnalysis) Analyze(sourceCode string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析失败: %v", err)
	}

	// 收集use和def信息
	lva.collectUseDef(node, fset)
	
	// 计算活跃变量
	lva.computeLiveVariables()
	
	return nil
}

// collectUseDef 收集use和def信息
func (lva *LiveVariableAnalysis) collectUseDef(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			point := ProgramPoint{
				Node: n,
				Line: fset.Position(n.Pos()).Line,
			}
			
			lva.use[point] = make(map[string]bool)
			lva.def[point] = make(map[string]bool)
			
			switch x := n.(type) {
			case *ast.AssignStmt:
				// 收集def（左值）
				for _, lhs := range x.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						lva.def[point][ident.Name] = true
					}
				}
				// 收集use（右值）
				for _, rhs := range x.Rhs {
					lva.collectUses(rhs, point)
				}
			case *ast.ExprStmt:
				lva.collectUses(x.X, point)
			}
		}
		return true
	})
}

// collectUses 收集表达式中的变量使用
func (lva *LiveVariableAnalysis) collectUses(expr ast.Expr, point ProgramPoint) {
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			lva.use[point][ident.Name] = true
		}
		return true
	})
}

// computeLiveVariables 计算活跃变量
func (lva *LiveVariableAnalysis) computeLiveVariables() {
	// 初始化
	for point := range lva.use {
		lva.liveVars[point] = make(map[string]bool)
	}
	
	// 迭代计算直到收敛
	changed := true
	for changed {
		changed = false
		for point := range lva.use {
			oldSize := len(lva.liveVars[point])
			
			// 计算新的活跃变量集合
			newLiveVars := make(map[string]bool)
			
			// 添加后继节点的活跃变量
			for succ := range lva.getSuccessors(point) {
				for varName := range lva.liveVars[succ] {
					newLiveVars[varName] = true
				}
			}
			
			// 移除当前节点定义的变量
			for varName := range lva.def[point] {
				delete(newLiveVars, varName)
			}
			
			// 添加当前节点使用的变量
			for varName := range lva.use[point] {
				newLiveVars[varName] = true
			}
			
			lva.liveVars[point] = newLiveVars
			
			if len(newLiveVars) != oldSize {
				changed = true
			}
		}
	}
}

// getSuccessors 获取后继节点（简化实现）
func (lva *LiveVariableAnalysis) getSuccessors(point ProgramPoint) map[ProgramPoint]bool {
	// 简化实现，实际应该基于控制流图
	return make(map[ProgramPoint]bool)
}

// IsVariableLive 检查变量在指定程序点是否活跃
func (lva *LiveVariableAnalysis) IsVariableLive(variable string, line int) bool {
	for point, liveVars := range lva.liveVars {
		if point.Line == line {
			return liveVars[variable]
		}
	}
	return false
}
```

## 3. 控制流分析

### 3.1 控制流图构建

#### 3.1.1 理论基础

**控制流图 (CFG)** 是一个有向图 $G = (V, E)$，其中：
- $V$ 是基本块的集合
- $E$ 是基本块之间的控制流边

**基本块** 是程序中的线性指令序列，只有一个入口点和一个出口点。

#### 3.1.2 Go语言实现

```go
// BasicBlock 基本块
type BasicBlock struct {
	ID       int
	StartPos token.Pos
	EndPos   token.Pos
	Nodes    []ast.Node
	Pred     []*BasicBlock
	Succ     []*BasicBlock
}

// ControlFlowGraph 控制流图
type ControlFlowGraph struct {
	Blocks    []*BasicBlock
	Entry     *BasicBlock
	Exit      *BasicBlock
	BlockMap  map[token.Pos]*BasicBlock
}

// NewControlFlowGraph 创建控制流图
func NewControlFlowGraph() *ControlFlowGraph {
	return &ControlFlowGraph{
		Blocks:   make([]*BasicBlock, 0),
		BlockMap: make(map[token.Pos]*BasicBlock),
	}
}

// BuildCFG 构建控制流图
func (cfg *ControlFlowGraph) BuildCFG(node ast.Node, fset *token.FileSet) {
	// 识别基本块
	cfg.identifyBasicBlocks(node, fset)
	
	// 建立控制流边
	cfg.buildControlFlowEdges()
}

// identifyBasicBlocks 识别基本块
func (cfg *ControlFlowGraph) identifyBasicBlocks(node ast.Node, fset *token.FileSet) {
	var currentBlock *BasicBlock
	blockID := 0
	
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		
		// 检查是否是块开始
		if cfg.isBlockStart(n) {
			if currentBlock != nil {
				cfg.Blocks = append(cfg.Blocks, currentBlock)
			}
			
			currentBlock = &BasicBlock{
				ID:       blockID,
				StartPos: n.Pos(),
				Nodes:    []ast.Node{n},
				Pred:     make([]*BasicBlock, 0),
				Succ:     make([]*BasicBlock, 0),
			}
			blockID++
			cfg.BlockMap[n.Pos()] = currentBlock
		} else if currentBlock != nil {
			currentBlock.Nodes = append(currentBlock.Nodes, n)
			currentBlock.EndPos = n.End()
		}
		
		return true
	})
	
	// 添加最后一个块
	if currentBlock != nil {
		cfg.Blocks = append(cfg.Blocks, currentBlock)
	}
}

// isBlockStart 检查是否是基本块开始
func (cfg *ControlFlowGraph) isBlockStart(n ast.Node) bool {
	switch n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt:
		return true
	}
	return false
}

// buildControlFlowEdges 建立控制流边
func (cfg *ControlFlowGraph) buildControlFlowEdges() {
	for _, block := range cfg.Blocks {
		cfg.addControlFlowEdges(block)
	}
}

// addControlFlowEdges 为基本块添加控制流边
func (cfg *ControlFlowGraph) addControlFlowEdges(block *BasicBlock) {
	if len(block.Nodes) == 0 {
		return
	}
	
	lastNode := block.Nodes[len(block.Nodes)-1]
	
	switch x := lastNode.(type) {
	case *ast.IfStmt:
		// if语句有两个后继：then分支和else分支
		if thenBlock := cfg.findBlockByPos(x.Body.Pos()); thenBlock != nil {
			cfg.addEdge(block, thenBlock)
		}
		if x.Else != nil {
			if elseBlock := cfg.findBlockByPos(x.Else.Pos()); elseBlock != nil {
				cfg.addEdge(block, elseBlock)
			}
		}
	case *ast.ForStmt:
		// for循环有后继：循环体
		if bodyBlock := cfg.findBlockByPos(x.Body.Pos()); bodyBlock != nil {
			cfg.addEdge(block, bodyBlock)
		}
	}
}

// addEdge 添加控制流边
func (cfg *ControlFlowGraph) addEdge(from, to *BasicBlock) {
	from.Succ = append(from.Succ, to)
	to.Pred = append(to.Pred, from)
}

// findBlockByPos 根据位置查找基本块
func (cfg *ControlFlowGraph) findBlockByPos(pos token.Pos) *BasicBlock {
	for _, block := range cfg.Blocks {
		if block.StartPos <= pos && pos <= block.EndPos {
			return block
		}
	}
	return nil
}

// DominanceAnalysis 支配分析
type DominanceAnalysis struct {
	cfg       *ControlFlowGraph
	dominators map[*BasicBlock]map[*BasicBlock]bool
}

// NewDominanceAnalysis 创建支配分析器
func NewDominanceAnalysis(cfg *ControlFlowGraph) *DominanceAnalysis {
	return &DominanceAnalysis{
		cfg:        cfg,
		dominators: make(map[*BasicBlock]map[*BasicBlock]bool),
	}
}

// ComputeDominators 计算支配关系
func (da *DominanceAnalysis) ComputeDominators() {
	// 初始化
	for _, block := range da.cfg.Blocks {
		da.dominators[block] = make(map[*BasicBlock]bool)
		for _, other := range da.cfg.Blocks {
			da.dominators[block][other] = true
		}
	}
	
	// 入口块只支配自己
	if da.cfg.Entry != nil {
		for _, block := range da.cfg.Blocks {
			if block != da.cfg.Entry {
				delete(da.dominators[da.cfg.Entry], block)
			}
		}
	}
	
	// 迭代计算
	changed := true
	for changed {
		changed = false
		for _, block := range da.cfg.Blocks {
			if block == da.cfg.Entry {
				continue
			}
			
			oldDominators := make(map[*BasicBlock]bool)
			for dom := range da.dominators[block] {
				oldDominators[dom] = true
			}
			
			// 计算新的支配者集合
			newDominators := make(map[*BasicBlock]bool)
			first := true
			
			for _, pred := range block.Pred {
				if first {
					for dom := range da.dominators[pred] {
						newDominators[dom] = true
					}
					first = false
				} else {
					// 交集
					for dom := range newDominators {
						if !da.dominators[pred][dom] {
							delete(newDominators, dom)
						}
					}
				}
			}
			
			// 添加自己
			newDominators[block] = true
			
			da.dominators[block] = newDominators
			
			// 检查是否改变
			if len(newDominators) != len(oldDominators) {
				changed = true
			} else {
				for dom := range newDominators {
					if !oldDominators[dom] {
						changed = true
						break
					}
				}
			}
		}
	}
}

// Dominates 检查block1是否支配block2
func (da *DominanceAnalysis) Dominates(block1, block2 *BasicBlock) bool {
	return da.dominators[block2][block1]
}
```

## 4. 类型检查

### 4.1 类型系统

#### 4.1.1 理论基础

**类型系统** 是一组规则，用于验证程序中的类型使用是否正确。

**类型推导** 是从表达式中推断类型的过程：
$$\frac{\Gamma \vdash e_1 : \tau_1 \quad \Gamma \vdash e_2 : \tau_2}{\Gamma \vdash e_1 + e_2 : \tau_1 \sqcup \tau_2}$$

#### 4.1.2 Go语言实现

```go
// Type 表示类型
type Type interface {
	String() string
	IsAssignableTo(other Type) bool
}

// BasicType 基本类型
type BasicType struct {
	Name string
}

func (bt *BasicType) String() string {
	return bt.Name
}

func (bt *BasicType) IsAssignableTo(other Type) bool {
	if otherBT, ok := other.(*BasicType); ok {
		return bt.Name == otherBT.Name
	}
	return false
}

// FunctionType 函数类型
type FunctionType struct {
	Params []Type
	Return Type
}

func (ft *FunctionType) String() string {
	paramStrs := make([]string, len(ft.Params))
	for i, param := range ft.Params {
		paramStrs[i] = param.String()
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(paramStrs, ", "), ft.Return.String())
}

func (ft *FunctionType) IsAssignableTo(other Type) bool {
	if otherFT, ok := other.(*FunctionType); ok {
		if len(ft.Params) != len(otherFT.Params) {
			return false
		}
		for i, param := range ft.Params {
			if !param.IsAssignableTo(otherFT.Params[i]) {
				return false
			}
		}
		return ft.Return.IsAssignableTo(otherFT.Return)
	}
	return false
}

// TypeChecker 类型检查器
type TypeChecker struct {
	types    map[ast.Expr]Type
	env      map[string]Type
	errors   []string
}

// NewTypeChecker 创建类型检查器
func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		types:  make(map[ast.Expr]Type),
		env:    make(map[string]Type),
		errors: make([]string, 0),
	}
}

// Check 检查类型
func (tc *TypeChecker) Check(node ast.Node) error {
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			tc.checkNode(n)
		}
		return true
	})
	
	if len(tc.errors) > 0 {
		return fmt.Errorf("类型检查错误: %s", strings.Join(tc.errors, "; "))
	}
	return nil
}

// checkNode 检查单个节点
func (tc *TypeChecker) checkNode(n ast.Node) {
	switch x := n.(type) {
	case *ast.Ident:
		tc.checkIdent(x)
	case *ast.BinaryExpr:
		tc.checkBinaryExpr(x)
	case *ast.CallExpr:
		tc.checkCallExpr(x)
	case *ast.AssignStmt:
		tc.checkAssignStmt(x)
	}
}

// checkIdent 检查标识符
func (tc *TypeChecker) checkIdent(ident *ast.Ident) {
	if typ, exists := tc.env[ident.Name]; exists {
		tc.types[ident] = typ
	} else {
		tc.errors = append(tc.errors, fmt.Sprintf("未定义的变量: %s", ident.Name))
	}
}

// checkBinaryExpr 检查二元表达式
func (tc *TypeChecker) checkBinaryExpr(expr *ast.BinaryExpr) {
	leftType := tc.getType(expr.X)
	rightType := tc.getType(expr.Y)
	
	if leftType == nil || rightType == nil {
		return
	}
	
	switch expr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO:
		if tc.isNumericType(leftType) && tc.isNumericType(rightType) {
			tc.types[expr] = tc.unifyNumericTypes(leftType, rightType)
		} else {
			tc.errors = append(tc.errors, fmt.Sprintf("操作符 %s 不能用于类型 %s 和 %s", 
				expr.Op, leftType.String(), rightType.String()))
		}
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		if leftType.IsAssignableTo(rightType) || rightType.IsAssignableTo(leftType) {
			tc.types[expr] = &BasicType{Name: "bool"}
		} else {
			tc.errors = append(tc.errors, fmt.Sprintf("比较操作符不能用于类型 %s 和 %s", 
				leftType.String(), rightType.String()))
		}
	}
}

// checkCallExpr 检查函数调用
func (tc *TypeChecker) checkCallExpr(call *ast.CallExpr) {
	funcType := tc.getType(call.Fun)
	if funcType == nil {
		return
	}
	
	if ft, ok := funcType.(*FunctionType); ok {
		if len(call.Args) != len(ft.Params) {
			tc.errors = append(tc.errors, fmt.Sprintf("函数调用参数数量不匹配: 期望 %d, 实际 %d", 
				len(ft.Params), len(call.Args)))
			return
		}
		
		for i, arg := range call.Args {
			argType := tc.getType(arg)
			if argType != nil && !argType.IsAssignableTo(ft.Params[i]) {
				tc.errors = append(tc.errors, fmt.Sprintf("参数 %d 类型不匹配: 期望 %s, 实际 %s", 
					i, ft.Params[i].String(), argType.String()))
			}
		}
		
		tc.types[call] = ft.Return
	} else {
		tc.errors = append(tc.errors, fmt.Sprintf("尝试调用非函数类型: %s", funcType.String()))
	}
}

// checkAssignStmt 检查赋值语句
func (tc *TypeChecker) checkAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Lhs) != len(assign.Rhs) {
		tc.errors = append(tc.errors, "赋值语句左右值数量不匹配")
		return
	}
	
	for i, lhs := range assign.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			rhsType := tc.getType(assign.Rhs[i])
			if rhsType != nil {
				tc.env[ident.Name] = rhsType
			}
		}
	}
}

// getType 获取表达式类型
func (tc *TypeChecker) getType(expr ast.Expr) Type {
	if typ, exists := tc.types[expr]; exists {
		return typ
	}
	return nil
}

// isNumericType 检查是否是数值类型
func (tc *TypeChecker) isNumericType(typ Type) bool {
	if bt, ok := typ.(*BasicType); ok {
		switch bt.Name {
		case "int", "int8", "int16", "int32", "int64",
			 "uint", "uint8", "uint16", "uint32", "uint64",
			 "float32", "float64":
			return true
		}
	}
	return false
}

// unifyNumericTypes 统一数值类型
func (tc *TypeChecker) unifyNumericTypes(t1, t2 Type) Type {
	// 简化实现，实际应该考虑类型提升规则
	if bt1, ok := t1.(*BasicType); ok {
		return bt1
	}
	return t1
}
```

## 5. 安全分析

### 5.1 缓冲区溢出检测

#### 5.1.1 理论基础

**缓冲区溢出** 是当程序试图在预分配的固定长度缓冲区边界之外写入数据时发生的安全漏洞。

**形式化定义**：
设 $B$ 为缓冲区，$size(B)$ 为缓冲区大小，$access(B, i)$ 为对位置 $i$ 的访问，则：
$$\forall i: access(B, i) \Rightarrow 0 \leq i < size(B)$$

#### 5.1.2 Go语言实现

```go
// BufferOverflowDetector 缓冲区溢出检测器
type BufferOverflowDetector struct {
	buffers map[string]int
	accesses []BufferAccess
	errors   []string
}

// BufferAccess 缓冲区访问
type BufferAccess struct {
	Buffer string
	Index  ast.Expr
	Line   int
}

// NewBufferOverflowDetector 创建缓冲区溢出检测器
func NewBufferOverflowDetector() *BufferOverflowDetector {
	return &BufferOverflowDetector{
		buffers:  make(map[string]int),
		accesses: make([]BufferAccess, 0),
		errors:   make([]string, 0),
	}
}

// Analyze 分析缓冲区溢出
func (bod *BufferOverflowDetector) Analyze(sourceCode string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析失败: %v", err)
	}

	// 识别缓冲区
	bod.identifyBuffers(node, fset)
	
	// 检测访问
	bod.detectAccesses(node, fset)
	
	// 检查溢出
	bod.checkOverflows()
	
	if len(bod.errors) > 0 {
		return fmt.Errorf("缓冲区溢出检测错误: %s", strings.Join(bod.errors, "; "))
	}
	return nil
}

// identifyBuffers 识别缓冲区
func (bod *BufferOverflowDetector) identifyBuffers(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ArrayType:
			if lit, ok := x.Len.(*ast.BasicLit); ok {
				if size, err := strconv.Atoi(lit.Value); err == nil {
					// 简化实现，实际应该跟踪变量名
					bod.buffers["array"] = size
				}
			}
		case *ast.CallExpr:
			if fun, ok := x.Fun.(*ast.Ident); ok && fun.Name == "make" {
				if len(x.Args) >= 2 {
					if lit, ok := x.Args[1].(*ast.BasicLit); ok {
						if size, err := strconv.Atoi(lit.Value); err == nil {
							bod.buffers["slice"] = size
						}
					}
				}
			}
		}
		return true
	})
}

// detectAccesses 检测缓冲区访问
func (bod *BufferOverflowDetector) detectAccesses(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.IndexExpr:
			access := BufferAccess{
				Buffer: bod.getBufferName(x.X),
				Index:  x.Index,
				Line:   fset.Position(x.Pos()).Line,
			}
			bod.accesses = append(bod.accesses, access)
		}
		return true
	})
}

// getBufferName 获取缓冲区名称（简化实现）
func (bod *BufferOverflowDetector) getBufferName(expr ast.Expr) string {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return "unknown"
}

// checkOverflows 检查溢出
func (bod *BufferOverflowDetector) checkOverflows() {
	for _, access := range bod.accesses {
		if size, exists := bod.buffers[access.Buffer]; exists {
			if !bod.isIndexSafe(access.Index, size) {
				bod.errors = append(bod.errors, 
					fmt.Sprintf("第 %d 行: 可能的缓冲区溢出访问", access.Line))
			}
		}
	}
}

// isIndexSafe 检查索引是否安全
func (bod *BufferOverflowDetector) isIndexSafe(index ast.Expr, size int) bool {
	// 简化实现，实际应该进行符号执行
	if lit, ok := index.(*ast.BasicLit); ok {
		if i, err := strconv.Atoi(lit.Value); err == nil {
			return i >= 0 && i < size
		}
	}
	// 无法静态确定，假设不安全
	return false
}
```

### 5.2 空指针检测

#### 5.2.1 理论基础

**空指针解引用** 是当程序试图访问空指针指向的内存时发生的错误。

**形式化定义**：
设 $p$ 为指针，$*p$ 为解引用操作，则：
$$*p \Rightarrow p \neq null$$

#### 5.2.2 Go语言实现

```go
// NullPointerDetector 空指针检测器
type NullPointerDetector struct {
	nullChecks map[string]bool
	derefs     []Dereference
	errors     []string
}

// Dereference 解引用操作
type Dereference struct {
	Pointer string
	Line    int
}

// NewNullPointerDetector 创建空指针检测器
func NewNullPointerDetector() *NullPointerDetector {
	return &NullPointerDetector{
		nullChecks: make(map[string]bool),
		derefs:     make([]Dereference, 0),
		errors:     make([]string, 0),
	}
}

// Analyze 分析空指针
func (npd *NullPointerDetector) Analyze(sourceCode string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析失败: %v", err)
	}

	// 检测空指针检查
	npd.detectNullChecks(node, fset)
	
	// 检测解引用
	npd.detectDereferences(node, fset)
	
	// 检查潜在的空指针解引用
	npd.checkNullDereferences()
	
	if len(npd.errors) > 0 {
		return fmt.Errorf("空指针检测错误: %s", strings.Join(npd.errors, "; "))
	}
	return nil
}

// detectNullChecks 检测空指针检查
func (npd *NullPointerDetector) detectNullChecks(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.IfStmt:
			if binExpr, ok := x.Cond.(*ast.BinaryExpr); ok {
				if binExpr.Op == token.NEQ || binExpr.Op == token.EQL {
					if ident, ok := binExpr.X.(*ast.Ident); ok {
						if lit, ok := binExpr.Y.(*ast.Ident); ok && lit.Name == "nil" {
							npd.nullChecks[ident.Name] = true
						}
					}
				}
			}
		}
		return true
	})
}

// detectDereferences 检测解引用
func (npd *NullPointerDetector) detectDereferences(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StarExpr:
			if ident, ok := x.X.(*ast.Ident); ok {
				deref := Dereference{
					Pointer: ident.Name,
					Line:    fset.Position(x.Pos()).Line,
				}
				npd.derefs = append(npd.derefs, deref)
			}
		case *ast.SelectorExpr:
			if ident, ok := x.X.(*ast.Ident); ok {
				deref := Dereference{
					Pointer: ident.Name,
					Line:    fset.Position(x.Pos()).Line,
				}
				npd.derefs = append(npd.derefs, deref)
			}
		}
		return true
	})
}

// checkNullDereferences 检查潜在的空指针解引用
func (npd *NullPointerDetector) checkNullDereferences() {
	for _, deref := range npd.derefs {
		if !npd.nullChecks[deref.Pointer] {
			npd.errors = append(npd.errors, 
				fmt.Sprintf("第 %d 行: 潜在的空指针解引用: %s", deref.Line, deref.Pointer))
		}
	}
}
```

## 6. 性能分析

### 6.1 复杂度分析

#### 6.1.1 理论基础

**算法复杂度** 描述了算法执行时间或空间需求随输入规模增长的变化规律。

**大O记号**：
$$f(n) = O(g(n)) \Leftrightarrow \exists c, n_0: \forall n \geq n_0, f(n) \leq c \cdot g(n)$$

#### 6.1.2 Go语言实现

```go
// ComplexityAnalyzer 复杂度分析器
type ComplexityAnalyzer struct {
	complexities map[string]string
	errors       []string
}

// NewComplexityAnalyzer 创建复杂度分析器
func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		complexities: make(map[string]string),
		errors:       make([]string, 0),
	}
}

// Analyze 分析复杂度
func (ca *ComplexityAnalyzer) Analyze(sourceCode string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("解析失败: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if fun, ok := n.(*ast.FuncDecl); ok {
			complexity := ca.analyzeFunction(fun)
			ca.complexities[fun.Name.Name] = complexity
		}
		return true
	})
	
	return nil
}

// analyzeFunction 分析函数复杂度
func (ca *ComplexityAnalyzer) analyzeFunction(fun *ast.FuncDecl) string {
	// 计算圈复杂度
	cyclomatic := ca.calculateCyclomaticComplexity(fun)
	
	// 分析嵌套循环
	nestedLoops := ca.analyzeNestedLoops(fun)
	
	// 分析递归
	recursive := ca.isRecursive(fun)
	
	// 综合评估
	if cyclomatic > 10 {
		return "高复杂度"
	} else if nestedLoops > 2 {
		return "高复杂度（嵌套循环）"
	} else if recursive {
		return "需要递归分析"
	} else {
		return "低复杂度"
	}
}

// calculateCyclomaticComplexity 计算圈复杂度
func (ca *ComplexityAnalyzer) calculateCyclomaticComplexity(fun *ast.FuncDecl) int {
	complexity := 1 // 基础复杂度
	
	ast.Inspect(fun, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.CaseClause:
			complexity++
		case *ast.BinaryExpr:
			if binExpr, ok := n.(*ast.BinaryExpr); ok {
				if binExpr.Op == token.LOR || binExpr.Op == token.LAND {
					complexity++
				}
			}
		}
		return true
	})
	
	return complexity
}

// analyzeNestedLoops 分析嵌套循环
func (ca *ComplexityAnalyzer) analyzeNestedLoops(fun *ast.FuncDecl) int {
	maxNesting := 0
	currentNesting := 0
	
	ast.Inspect(fun, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			currentNesting++
			if currentNesting > maxNesting {
				maxNesting = currentNesting
			}
		}
		return true
	})
	
	return maxNesting
}

// isRecursive 检查是否递归
func (ca *ComplexityAnalyzer) isRecursive(fun *ast.FuncDecl) bool {
	funcName := fun.Name.Name
	
	recursive := false
	ast.Inspect(fun, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == funcName {
					recursive = true
					return false
				}
			}
		}
		return true
	})
	
	return recursive
}

// GetComplexity 获取函数复杂度
func (ca *ComplexityAnalyzer) GetComplexity(funcName string) string {
	if complexity, exists := ca.complexities[funcName]; exists {
		return complexity
	}
	return "未知"
}
```

## 7. 总结

### 7.1 静态分析的优势

1. **早期错误检测**：在编译时发现潜在问题
2. **全面性**：分析所有可能的执行路径
3. **自动化**：无需人工干预
4. **可扩展性**：支持多种分析类型

### 7.2 局限性

1. **假阳性**：可能报告不存在的错误
2. **假阴性**：可能遗漏某些错误
3. **精度限制**：抽象可能导致精度损失
4. **性能开销**：复杂分析可能耗时较长

### 7.3 最佳实践

1. **组合使用**：结合多种静态分析技术
2. **配置优化**：根据项目需求调整分析参数
3. **持续集成**：将静态分析集成到CI/CD流程
4. **结果验证**：定期验证分析结果的准确性

---

**参考文献**：
1. Aho, A. V., Lam, M. S., Sethi, R., & Ullman, J. D. (2006). Compilers: Principles, Techniques, and Tools
2. Nielson, F., Nielson, H. R., & Hankin, C. (2015). Principles of Program Analysis
3. Cousot, P., & Cousot, R. (1977). Abstract interpretation: A unified lattice model for static analysis of programs by construction or approximation of fixpoints 