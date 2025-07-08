# 质量保证

## 概述

质量保证模块涵盖了现代软件系统的质量保证机制，包括自动化测试、代码质量检查、性能监控、安全审计等技术。本章节深入探讨质量保证在Go生态系统中的应用和实践。

## 目录

- [自动化测试](#自动化测试)
- [代码质量检查](#代码质量检查)
- [性能监控](#性能监控)
- [安全审计](#安全审计)
- [持续集成](#持续集成)

## 自动化测试

### 单元测试框架

```go
// 单元测试管理器
package testing

import (
    "context"
    "fmt"
    "testing"
    "time"
)

type TestManager struct {
    tests map[string]TestSuite
}

type TestSuite struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Tests       []Test                 `json:"tests"`
    Setup       func() error           `json:"-"`
    Teardown    func() error           `json:"-"`
    CreatedAt   time.Time              `json:"created_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type Test struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Function    func(t *testing.T)    `json:"-"`
    Timeout     time.Duration          `json:"timeout"`
    Skip        bool                   `json:"skip"`
    Metadata    map[string]interface{} `json:"metadata"`
}

func NewTestManager() *TestManager {
    return &TestManager{
        tests: make(map[string]TestSuite),
    }
}

func (tm *TestManager) AddTestSuite(suite TestSuite) {
    tm.tests[suite.Name] = suite
}

func (tm *TestManager) RunTestSuite(ctx context.Context, suiteName string) (*TestResult, error) {
    suite, exists := tm.tests[suiteName]
    if !exists {
        return nil, fmt.Errorf("test suite not found: %s", suiteName)
    }
    
    result := &TestResult{
        SuiteName: suiteName,
        StartTime: time.Now(),
        Tests:     make([]TestResult, 0),
    }
    
    // 执行测试套件设置
    if suite.Setup != nil {
        if err := suite.Setup(); err != nil {
            return nil, fmt.Errorf("test suite setup failed: %w", err)
        }
    }
    
    // 执行测试
    for _, test := range suite.Tests {
        if test.Skip {
            continue
        }
        
        testResult := tm.runTest(ctx, test)
        result.Tests = append(result.Tests, testResult)
    }
    
    // 执行测试套件清理
    if suite.Teardown != nil {
        if err := suite.Teardown(); err != nil {
            return nil, fmt.Errorf("test suite teardown failed: %w", err)
        }
    }
    
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)
    
    return result, nil
}

func (tm *TestManager) runTest(ctx context.Context, test Test) TestResult {
    result := TestResult{
        Name:      test.Name,
        StartTime: time.Now(),
    }
    
    // 创建测试上下文
    testCtx, cancel := context.WithTimeout(ctx, test.Timeout)
    defer cancel()
    
    // 执行测试
    done := make(chan bool)
    go func() {
        defer close(done)
        // 这里需要实际的测试执行逻辑
        // 由于Go的testing包限制，这里只是示例
    }()
    
    select {
    case <-done:
        result.Status = "passed"
    case <-testCtx.Done():
        result.Status = "timeout"
        result.Error = "test timeout"
    }
    
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)
    
    return result
}

type TestResult struct {
    SuiteName string        `json:"suite_name"`
    StartTime time.Time     `json:"start_time"`
    EndTime   time.Time     `json:"end_time"`
    Duration  time.Duration `json:"duration"`
    Tests     []TestResult  `json:"tests"`
}

type TestResult struct {
    Name      string        `json:"name"`
    Status    string        `json:"status"`
    Error     string        `json:"error,omitempty"`
    StartTime time.Time     `json:"start_time"`
    EndTime   time.Time     `json:"end_time"`
    Duration  time.Duration `json:"duration"`
}
```

### 集成测试框架

```go
// 集成测试管理器
package integration

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "time"
)

type IntegrationTestManager struct {
    db     *sql.DB
    client *http.Client
}

type IntegrationTest struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Setup       func() error           `json:"-"`
    Test        func() error           `json:"-"`
    Teardown    func() error           `json:"-"`
    Timeout     time.Duration          `json:"timeout"`
    Metadata    map[string]interface{} `json:"metadata"`
}

func NewIntegrationTestManager(db *sql.DB, client *http.Client) *IntegrationTestManager {
    return &IntegrationTestManager{
        db:     db,
        client: client,
    }
}

func (itm *IntegrationTestManager) RunIntegrationTest(ctx context.Context, test IntegrationTest) error {
    // 执行测试设置
    if test.Setup != nil {
        if err := test.Setup(); err != nil {
            return fmt.Errorf("test setup failed: %w", err)
        }
    }
    
    // 执行测试
    testCtx, cancel := context.WithTimeout(ctx, test.Timeout)
    defer cancel()
    
    done := make(chan error)
    go func() {
        done <- test.Test()
    }()
    
    select {
    case err := <-done:
        if err != nil {
            return err
        }
    case <-testCtx.Done():
        return fmt.Errorf("test timeout")
    }
    
    // 执行测试清理
    if test.Teardown != nil {
        if err := test.Teardown(); err != nil {
            return fmt.Errorf("test teardown failed: %w", err)
        }
    }
    
    return nil
}

// 示例集成测试
func (itm *IntegrationTestManager) TestUserAPI() IntegrationTest {
    return IntegrationTest{
        Name:        "Test User API",
        Description: "Test user creation and retrieval",
        Timeout:     30 * time.Second,
        Setup: func() error {
            // 清理测试数据
            _, err := itm.db.Exec("DELETE FROM users WHERE email LIKE '%@test.com'")
            return err
        },
        Test: func() error {
            // 创建用户
            resp, err := itm.client.Post("http://localhost:8080/users", "application/json", strings.NewReader(`{
                "name": "Test User",
                "email": "test@test.com"
            }`))
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusCreated {
                return fmt.Errorf("expected status 201, got %d", resp.StatusCode)
            }
            
            // 获取用户
            resp, err = itm.client.Get("http://localhost:8080/users/test@test.com")
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
            }
            
            return nil
        },
        Teardown: func() error {
            // 清理测试数据
            _, err := itm.db.Exec("DELETE FROM users WHERE email = 'test@test.com'")
            return err
        },
    }
}
```

## 代码质量检查

### 静态代码分析

```go
// 代码质量检查器
package quality

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "path/filepath"
    "strings"
)

type CodeQualityChecker struct {
    rules []QualityRule
}

type QualityRule struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Severity    string `json:"severity"` // "error", "warning", "info"
    Check       func(*ast.File) []QualityIssue `json:"-"`
}

type QualityIssue struct {
    Rule      string `json:"rule"`
    Severity  string `json:"severity"`
    Message   string `json:"message"`
    File      string `json:"file"`
    Line      int    `json:"line"`
    Column    int    `json:"column"`
}

func NewCodeQualityChecker() *CodeQualityChecker {
    return &CodeQualityChecker{
        rules: make([]QualityRule, 0),
    }
}

func (cqc *CodeQualityChecker) AddRule(rule QualityRule) {
    cqc.rules = append(cqc.rules, rule)
}

func (cqc *CodeQualityChecker) CheckFile(filePath string) ([]QualityIssue, error) {
    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }
    
    var issues []QualityIssue
    
    for _, rule := range cqc.rules {
        ruleIssues := rule.Check(file)
        for i := range ruleIssues {
            ruleIssues[i].File = filePath
        }
        issues = append(issues, ruleIssues...)
    }
    
    return issues, nil
}

func (cqc *CodeQualityChecker) CheckDirectory(dirPath string) ([]QualityIssue, error) {
    var allIssues []QualityIssue
    
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() && strings.HasSuffix(path, ".go") {
            issues, err := cqc.CheckFile(path)
            if err != nil {
                return err
            }
            allIssues = append(allIssues, issues...)
        }
        
        return nil
    })
    
    return allIssues, err
}

// 预定义的质量规则
func (cqc *CodeQualityChecker) AddDefaultRules() {
    // 函数长度检查
    cqc.AddRule(QualityRule{
        Name:        "function_length",
        Description: "Functions should not exceed 50 lines",
        Severity:    "warning",
        Check: func(file *ast.File) []QualityIssue {
            var issues []QualityIssue
            
            for _, decl := range file.Decls {
                if funcDecl, ok := decl.(*ast.FuncDecl); ok {
                    if funcDecl.Body != nil {
                        lineCount := len(funcDecl.Body.List)
                        if lineCount > 50 {
                            issues = append(issues, QualityIssue{
                                Rule:     "function_length",
                                Severity: "warning",
                                Message:  fmt.Sprintf("Function %s is %d lines long (max 50)", funcDecl.Name.Name, lineCount),
                                Line:     int(funcDecl.Pos()),
                            })
                        }
                    }
                }
            }
            
            return issues
        },
    })
    
    // 变量命名检查
    cqc.AddRule(QualityRule{
        Name:        "variable_naming",
        Description: "Variable names should follow Go naming conventions",
        Severity:    "warning",
        Check: func(file *ast.File) []QualityIssue {
            var issues []QualityIssue
            
            ast.Inspect(file, func(n ast.Node) bool {
                if ident, ok := n.(*ast.Ident); ok {
                    if ident.Obj != nil && ident.Obj.Kind == ast.Var {
                        if strings.Contains(ident.Name, "_") {
                            issues = append(issues, QualityIssue{
                                Rule:     "variable_naming",
                                Severity: "warning",
                                Message:  fmt.Sprintf("Variable name '%s' contains underscores", ident.Name),
                                Line:     int(ident.Pos()),
                            })
                        }
                    }
                }
                return true
            })
            
            return issues
        },
    })
}
```

### 代码覆盖率分析

```go
// 代码覆盖率分析器
package coverage

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type CoverageAnalyzer struct {
    outputDir string
}

type CoverageReport struct {
    TotalLines    int     `json:"total_lines"`
    CoveredLines  int     `json:"covered_lines"`
    Coverage      float64 `json:"coverage"`
    Files         []FileCoverage `json:"files"`
    GeneratedAt   string  `json:"generated_at"`
}

type FileCoverage struct {
    File         string  `json:"file"`
    TotalLines   int     `json:"total_lines"`
    CoveredLines int     `json:"covered_lines"`
    Coverage     float64 `json:"coverage"`
}

func NewCoverageAnalyzer(outputDir string) *CoverageAnalyzer {
    return &CoverageAnalyzer{
        outputDir: outputDir,
    }
}

func (ca *CoverageAnalyzer) RunCoverage(packagePath string) (*CoverageReport, error) {
    // 运行测试并生成覆盖率报告
    cmd := exec.Command("go", "test", "-coverprofile=coverage.out", packagePath)
    cmd.Dir = packagePath
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    
    // 解析覆盖率输出
    coverageFile := filepath.Join(packagePath, "coverage.out")
    report, err := ca.parseCoverageFile(coverageFile)
    if err != nil {
        return nil, err
    }
    
    // 生成HTML报告
    if err := ca.generateHTMLReport(packagePath); err != nil {
        return nil, err
    }
    
    return report, nil
}

func (ca *CoverageAnalyzer) parseCoverageFile(filePath string) (*CoverageReport, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    
    lines := strings.Split(string(data), "\n")
    report := &CoverageReport{
        Files: make([]FileCoverage, 0),
    }
    
    for _, line := range lines {
        if strings.HasPrefix(line, "mode:") {
            continue
        }
        
        parts := strings.Split(line, ":")
        if len(parts) < 4 {
            continue
        }
        
        file := parts[0]
        statements := strings.Split(parts[3], " ")
        
        fileCoverage := FileCoverage{
            File: file,
        }
        
        for _, stmt := range statements {
            if stmt == "" {
                continue
            }
            
            fileCoverage.TotalLines++
            if strings.HasSuffix(stmt, ".1") {
                fileCoverage.CoveredLines++
            }
        }
        
        if fileCoverage.TotalLines > 0 {
            fileCoverage.Coverage = float64(fileCoverage.CoveredLines) / float64(fileCoverage.TotalLines) * 100
        }
        
        report.Files = append(report.Files, fileCoverage)
        report.TotalLines += fileCoverage.TotalLines
        report.CoveredLines += fileCoverage.CoveredLines
    }
    
    if report.TotalLines > 0 {
        report.Coverage = float64(report.CoveredLines) / float64(report.TotalLines) * 100
    }
    
    return report, nil
}

func (ca *CoverageAnalyzer) generateHTMLReport(packagePath string) error {
    cmd := exec.Command("go", "tool", "cover", "-html=coverage.out", "-o=coverage.html")
    cmd.Dir = packagePath
    return cmd.Run()
}
```

## 性能监控

### 性能分析器

```go
// 性能分析器
package performance

import (
    "context"
    "fmt"
    "runtime"
    "runtime/pprof"
    "time"
)

type PerformanceProfiler struct {
    cpuProfile   *pprof.Profile
    memoryProfile *pprof.Profile
    goroutineProfile *pprof.Profile
}

type PerformanceMetrics struct {
    Timestamp       time.Time `json:"timestamp"`
    CPUUsage        float64   `json:"cpu_usage"`
    MemoryUsage     uint64    `json:"memory_usage"`
    GoroutineCount  int       `json:"goroutine_count"`
    GCStats         GCStats   `json:"gc_stats"`
}

type GCStats struct {
    NumGC         uint32  `json:"num_gc"`
    PauseTotalNs  uint64  `json:"pause_total_ns"`
    PauseNs       uint64  `json:"pause_ns"`
    PauseEnd      uint64  `json:"pause_end"`
    LastGC        uint64  `json:"last_gc"`
}

func NewPerformanceProfiler() *PerformanceProfiler {
    return &PerformanceProfiler{
        cpuProfile:      pprof.Lookup("cpu"),
        memoryProfile:   pprof.Lookup("heap"),
        goroutineProfile: pprof.Lookup("goroutine"),
    }
}

func (pp *PerformanceProfiler) StartCPUProfile(writer io.Writer) error {
    return pprof.StartCPUProfile(writer)
}

func (pp *PerformanceProfiler) StopCPUProfile() {
    pprof.StopCPUProfile()
}

func (pp *PerformanceProfiler) WriteMemoryProfile(writer io.Writer) error {
    return pp.memoryProfile.WriteTo(writer, 0)
}

func (pp *PerformanceProfiler) WriteGoroutineProfile(writer io.Writer) error {
    return pp.goroutineProfile.WriteTo(writer, 0)
}

func (pp *PerformanceProfiler) GetMetrics() *PerformanceMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    var gcStats GCStats
    runtime.ReadGCStats(&gcStats)
    
    return &PerformanceMetrics{
        Timestamp:      time.Now(),
        MemoryUsage:    m.Alloc,
        GoroutineCount: runtime.NumGoroutine(),
        GCStats:        gcStats,
    }
}

func (pp *PerformanceProfiler) MonitorPerformance(ctx context.Context, interval time.Duration) <-chan *PerformanceMetrics {
    metricsChan := make(chan *PerformanceMetrics)
    
    go func() {
        defer close(metricsChan)
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                metrics := pp.GetMetrics()
                metricsChan <- metrics
            }
        }
    }()
    
    return metricsChan
}
```

### 性能基准测试

```go
// 性能基准测试框架
package benchmark

import (
    "context"
    "fmt"
    "testing"
    "time"
)

type BenchmarkSuite struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Benchmarks  []Benchmark           `json:"benchmarks"`
    Setup       func() error           `json:"-"`
    Teardown    func() error           `json:"-"`
}

type Benchmark struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Function    func(b *testing.B)    `json:"-"`
    Timeout     time.Duration          `json:"timeout"`
    Iterations  int                    `json:"iterations"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type BenchmarkResult struct {
    Name         string        `json:"name"`
    Operations   int64         `json:"operations"`
    Duration     time.Duration `json:"duration"`
    OperationsPerSecond float64 `json:"operations_per_second"`
    MemoryAlloc  int64         `json:"memory_alloc"`
    MemoryBytes  int64         `json:"memory_bytes"`
}

func NewBenchmarkSuite(name, description string) *BenchmarkSuite {
    return &BenchmarkSuite{
        Name:        name,
        Description: description,
        Benchmarks:  make([]Benchmark, 0),
    }
}

func (bs *BenchmarkSuite) AddBenchmark(benchmark Benchmark) {
    bs.Benchmarks = append(bs.Benchmarks, benchmark)
}

func (bs *BenchmarkSuite) RunBenchmarks() ([]BenchmarkResult, error) {
    var results []BenchmarkResult
    
    // 执行设置
    if bs.Setup != nil {
        if err := bs.Setup(); err != nil {
            return nil, fmt.Errorf("benchmark setup failed: %w", err)
        }
    }
    
    // 执行基准测试
    for _, benchmark := range bs.Benchmarks {
        result := bs.runBenchmark(benchmark)
        results = append(results, result)
    }
    
    // 执行清理
    if bs.Teardown != nil {
        if err := bs.Teardown(); err != nil {
            return nil, fmt.Errorf("benchmark teardown failed: %w", err)
        }
    }
    
    return results, nil
}

func (bs *BenchmarkSuite) runBenchmark(benchmark Benchmark) BenchmarkResult {
    result := BenchmarkResult{
        Name: benchmark.Name,
    }
    
    // 运行基准测试
    b := &testing.B{}
    benchmark.Function(b)
    
    result.Operations = b.N
    result.Duration = b.NsPerOp()
    result.OperationsPerSecond = float64(b.N) / float64(b.NsPerOp()) * 1e9
    result.MemoryAlloc = b.AllocedBytesPerOp()
    result.MemoryBytes = b.TotalAllocBytesPerOp()
    
    return result
}

// 示例基准测试
func BenchmarkStringConcatenation(b *testing.B) {
    var result string
    for i := 0; i < b.N; i++ {
        result = "Hello" + " " + "World"
    }
    _ = result
}

func BenchmarkStringBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        builder.WriteString("Hello")
        builder.WriteString(" ")
        builder.WriteString("World")
        _ = builder.String()
    }
}
```

## 安全审计

### 安全扫描器

```go
// 安全扫描器
package security

import (
    "context"
    "crypto/tls"
    "fmt"
    "net/http"
    "strings"
    "time"
)

type SecurityScanner struct {
    client *http.Client
}

type SecurityVulnerability struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    CVE         string    `json:"cve,omitempty"`
    URL         string    `json:"url"`
    Timestamp   time.Time `json:"timestamp"`
}

func NewSecurityScanner() *SecurityScanner {
    return &SecurityScanner{
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: true, // 仅用于测试
                },
            },
        },
    }
}

func (ss *SecurityScanner) ScanURL(url string) ([]SecurityVulnerability, error) {
    var vulnerabilities []SecurityVulnerability
    
    // 检查HTTPS
    if !strings.HasPrefix(url, "https://") {
        vulnerabilities = append(vulnerabilities, SecurityVulnerability{
            ID:          "SEC-001",
            Type:        "transport_security",
            Severity:    "high",
            Description: "URL does not use HTTPS",
            URL:         url,
            Timestamp:   time.Now(),
        })
    }
    
    // 检查安全头
    resp, err := ss.client.Get(url)
    if err != nil {
        return vulnerabilities, err
    }
    defer resp.Body.Close()
    
    // 检查安全头
    if resp.Header.Get("X-Frame-Options") == "" {
        vulnerabilities = append(vulnerabilities, SecurityVulnerability{
            ID:          "SEC-002",
            Type:        "security_headers",
            Severity:    "medium",
            Description: "Missing X-Frame-Options header",
            URL:         url,
            Timestamp:   time.Now(),
        })
    }
    
    if resp.Header.Get("X-Content-Type-Options") == "" {
        vulnerabilities = append(vulnerabilities, SecurityVulnerability{
            ID:          "SEC-003",
            Type:        "security_headers",
            Severity:    "medium",
            Description: "Missing X-Content-Type-Options header",
            URL:         url,
            Timestamp:   time.Now(),
        })
    }
    
    if resp.Header.Get("X-XSS-Protection") == "" {
        vulnerabilities = append(vulnerabilities, SecurityVulnerability{
            ID:          "SEC-004",
            Type:        "security_headers",
            Severity:    "medium",
            Description: "Missing X-XSS-Protection header",
            URL:         url,
            Timestamp:   time.Now(),
        })
    }
    
    return vulnerabilities, nil
}

func (ss *SecurityScanner) ScanDependencies(goModPath string) ([]SecurityVulnerability, error) {
    var vulnerabilities []SecurityVulnerability
    
    // 读取go.mod文件
    data, err := os.ReadFile(goModPath)
    if err != nil {
        return nil, err
    }
    
    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "\t") {
            // 解析依赖
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                module := parts[0]
                version := parts[1]
                
                // 检查已知漏洞
                if vuln := ss.checkKnownVulnerability(module, version); vuln != nil {
                    vulnerabilities = append(vulnerabilities, *vuln)
                }
            }
        }
    }
    
    return vulnerabilities, nil
}

func (ss *SecurityScanner) checkKnownVulnerability(module, version string) *SecurityVulnerability {
    // 这里应该连接到漏洞数据库进行检查
    // 示例实现
    knownVulnerabilities := map[string]string{
        "github.com/example/vulnerable": "CVE-2023-1234",
    }
    
    if cve, exists := knownVulnerabilities[module]; exists {
        return &SecurityVulnerability{
            ID:          "DEP-001",
            Type:        "dependency_vulnerability",
            Severity:    "high",
            Description: fmt.Sprintf("Vulnerable dependency: %s %s", module, version),
            CVE:         cve,
            Timestamp:   time.Now(),
        }
    }
    
    return nil
}
```

## 持续集成

### CI/CD流水线

```go
// CI/CD流水线管理器
package cicd

import (
    "context"
    "fmt"
    "os/exec"
    "time"
)

type CICDPipeline struct {
    stages []PipelineStage
}

type PipelineStage struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Command     string                 `json:"command"`
    Args        []string               `json:"args"`
    Timeout     time.Duration          `json:"timeout"`
    Required    bool                   `json:"required"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type PipelineResult struct {
    StageName   string        `json:"stage_name"`
    Status      string        `json:"status"`
    Duration    time.Duration `json:"duration"`
    Output      string        `json:"output"`
    Error       string        `json:"error,omitempty"`
    Timestamp   time.Time     `json:"timestamp"`
}

func NewCICDPipeline() *CICDPipeline {
    return &CICDPipeline{
        stages: make([]PipelineStage, 0),
    }
}

func (cp *CICDPipeline) AddStage(stage PipelineStage) {
    cp.stages = append(cp.stages, stage)
}

func (cp *CICDPipeline) Run(ctx context.Context) ([]PipelineResult, error) {
    var results []PipelineResult
    
    for _, stage := range cp.stages {
        result := cp.runStage(ctx, stage)
        results = append(results, result)
        
        if result.Status == "failed" && stage.Required {
            return results, fmt.Errorf("required stage failed: %s", stage.Name)
        }
    }
    
    return results, nil
}

func (cp *CICDPipeline) runStage(ctx context.Context, stage PipelineStage) PipelineResult {
    result := PipelineResult{
        StageName: stage.Name,
        Timestamp: time.Now(),
    }
    
    // 创建命令
    cmd := exec.CommandContext(ctx, stage.Command, stage.Args...)
    
    // 设置超时
    if stage.Timeout > 0 {
        timeoutCtx, cancel := context.WithTimeout(ctx, stage.Timeout)
        defer cancel()
        cmd = exec.CommandContext(timeoutCtx, stage.Command, stage.Args...)
    }
    
    // 执行命令
    output, err := cmd.CombinedOutput()
    result.Duration = time.Since(result.Timestamp)
    result.Output = string(output)
    
    if err != nil {
        result.Status = "failed"
        result.Error = err.Error()
    } else {
        result.Status = "success"
    }
    
    return result
}

// 预定义的CI/CD阶段
func (cp *CICDPipeline) AddDefaultStages() {
    // 代码格式化
    cp.AddStage(PipelineStage{
        Name:        "format",
        Description: "Format Go code",
        Command:     "go",
        Args:        []string{"fmt", "./..."},
        Timeout:     5 * time.Minute,
        Required:    true,
    })
    
    // 代码检查
    cp.AddStage(PipelineStage{
        Name:        "lint",
        Description: "Run linter",
        Command:     "golangci-lint",
        Args:        []string{"run"},
        Timeout:     10 * time.Minute,
        Required:    true,
    })
    
    // 运行测试
    cp.AddStage(PipelineStage{
        Name:        "test",
        Description: "Run tests",
        Command:     "go",
        Args:        []string{"test", "-v", "./..."},
        Timeout:     15 * time.Minute,
        Required:    true,
    })
    
    // 构建应用
    cp.AddStage(PipelineStage{
        Name:        "build",
        Description: "Build application",
        Command:     "go",
        Args:        []string{"build", "-o", "app", "."},
        Timeout:     10 * time.Minute,
        Required:    true,
    })
}
```

## 总结

质量保证为现代软件系统提供了重要的质量保障机制。通过自动化测试、代码质量检查、性能监控、安全审计等技术，我们可以构建高质量、可靠的软件系统。

### 关键要点

1. **自动化测试**: 建立完整的测试体系，包括单元测试、集成测试、性能测试
2. **代码质量**: 使用静态分析工具确保代码质量和一致性
3. **性能监控**: 实时监控系统性能，及时发现和解决问题
4. **安全审计**: 定期进行安全扫描和漏洞检查
5. **持续集成**: 建立自动化的CI/CD流水线

### 实践建议

- 在项目初期就建立质量保证体系
- 自动化所有可重复的质量检查过程
- 建立性能基准和监控机制
- 定期进行安全审计和漏洞扫描
- 持续优化CI/CD流水线
