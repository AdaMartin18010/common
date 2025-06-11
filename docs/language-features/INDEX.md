# Golang 语言特性文档索引

## 快速导航

### 📚 核心特性文档

| 特性 | 描述 | 文档链接 | 状态 |
|------|------|----------|------|
| 🏗️ **类型系统** | 静态类型、接口、泛型、反射 | [类型系统详解](./01-type-system/README.md) | ✅ 完成 |
| 🔄 **控制流** | 条件、循环、switch、defer、panic/recover | [控制流详解](./02-control-flow/README.md) | ✅ 完成 |
| 📊 **数据流** | 变量作用域、内存管理、垃圾回收 | [数据流详解](./03-data-flow/README.md) | ✅ 完成 |
| ⚡ **并发编程** | Goroutines、Channels、Select、同步原语 | [并发编程详解](./04-concurrency/README.md) | ✅ 完成 |
| 🧠 **内存模型** | Happens-before、内存屏障、原子操作 | [内存模型详解](./05-memory-model/README.md) | ✅ 完成 |

### 🎯 学习路径

#### 初学者路径

1. [控制流详解](./02-control-flow/README.md) - 基础语法和控制结构
2. [类型系统详解](./01-type-system/README.md) - 基础类型和接口
3. [数据流详解](./03-data-flow/README.md) - 变量和内存管理
4. [并发编程详解](./04-concurrency/README.md) - Goroutines和Channels基础

#### 进阶者路径

1. [类型系统详解](./01-type-system/README.md) - 泛型和反射
2. [并发编程详解](./04-concurrency/README.md) - 并发模式和Context
3. [内存模型详解](./05-memory-model/README.md) - 内存模型和原子操作

#### 专家级路径

1. [内存模型详解](./05-memory-model/README.md) - 高级内存模型
2. [并发编程详解](./04-concurrency/README.md) - 无锁数据结构和性能优化
3. [类型系统详解](./01-type-system/README.md) - unsafe包和底层编程

### 📋 文档特性

每个文档都包含：

- ✅ **理论讲解** - 核心概念和原理
- ✅ **代码示例** - 实用的代码示例
- ✅ **最佳实践** - 实际开发中的最佳实践
- ✅ **常见陷阱** - 避免常见错误
- ✅ **性能优化** - 性能优化技巧
- ✅ **2025年改进** - 最新特性和改进

### 🔍 快速查找

#### 按主题查找

**并发编程**

- [Goroutines基础](./04-concurrency/README.md#goroutines)
- [Channels通信](./04-concurrency/README.md#channels)
- [Select语句](./04-concurrency/README.md#select语句)
- [同步原语](./04-concurrency/README.md#同步原语)
- [Context包](./04-concurrency/README.md#context包)
- [并发模式](./04-concurrency/README.md#并发模式)

**类型系统**

- [基础类型](./01-type-system/README.md#基础类型)
- [接口系统](./01-type-system/README.md#接口系统)
- [泛型支持](./01-type-system/README.md#泛型支持)
- [反射机制](./01-type-system/README.md#反射机制)
- [类型断言](./01-type-system/README.md#类型断言)

**控制流**

- [条件语句](./02-control-flow/README.md#条件语句)
- [循环语句](./02-control-flow/README.md#循环语句)
- [Switch语句](./02-control-flow/README.md#switch语句)
- [Defer语句](./02-control-flow/README.md#defer语句)
- [Panic/Recovery](./02-control-flow/README.md#panicrecovery机制)

**内存模型**

- [Happens-before关系](./05-memory-model/README.md#happens-before关系)
- [内存屏障](./05-memory-model/README.md#内存屏障)
- [原子操作](./05-memory-model/README.md#原子操作)
- [缓存一致性](./05-memory-model/README.md#缓存一致性)

#### 按难度查找

**入门级**

- [控制流基础](./02-control-flow/README.md)
- [类型系统基础](./01-type-system/README.md)

**进阶级**

- [并发编程](./04-concurrency/README.md)
- [数据流管理](./03-data-flow/README.md)

**专家级**

- [内存模型](./05-memory-model/README.md)
- [高级并发模式](./04-concurrency/README.md#并发模式)

### 📈 学习建议

1. **循序渐进** - 按照学习路径逐步深入
2. **动手实践** - 运行每个代码示例
3. **理解原理** - 不仅要知道怎么做，还要知道为什么
4. **关注性能** - 在实际项目中注意性能优化
5. **持续更新** - 关注Go语言的最新发展

### 🔗 相关资源

- [Go官方文档](https://golang.org/doc/)
- [Go语言规范](https://golang.org/ref/spec)
- [Go内存模型](https://golang.org/ref/mem)
- [Go并发编程](https://golang.org/doc/effective_go.html#concurrency)

### 📝 贡献指南

欢迎社区贡献：

1. **问题反馈** - 报告文档错误或问题
2. **内容补充** - 添加缺失的内容
3. **示例代码** - 提供更好的代码示例
4. **最佳实践** - 分享实践经验

---

*最后更新时间: 2025年1月*
*文档版本: v1.0*
