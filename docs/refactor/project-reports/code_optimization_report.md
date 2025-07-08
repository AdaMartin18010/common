# Go代码性能优化分析报告

## 总体统计

- 总文件数: 98

## 文件分析

### 07-Implementation-Examples\01-Basic-Examples\abac_auth.go

- 代码行数: 65
- 函数数量: 4
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\aes_encrypt.go

- 代码行数: 59
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\api_rate_limit.go

- 代码行数: 77
- 函数数量: 5
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 1
- 错误处理评分: 0/3
- 并发安全评分: 3/3

### 07-Implementation-Examples\01-Basic-Examples\audit_log.go

- 代码行数: 74
- 函数数量: 2
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\channel_worker_pool.go

- 代码行数: 128
- 函数数量: 3
- Goroutine数量: 2
- Channel数量: 2
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 3/3
- 性能问题:
  - 在goroutine中使用fmt.Print可能影响性能

### 07-Implementation-Examples\01-Basic-Examples\concurrent_counter.go

- 代码行数: 26
- 函数数量: 1
- Goroutine数量: 1
- Channel数量: 0
- 互斥锁数量: 1
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\csrf_protection.go

- 代码行数: 62
- 函数数量: 5
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\error_handling.go

- 代码行数: 202
- 函数数量: 8
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\input_validation.go

- 代码行数: 56
- 函数数量: 4
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\jwt_auth.go

- 代码行数: 56
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\oauth2_auth.go

- 代码行数: 48
- 函数数量: 1
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 1/3
- 并发安全评分: 3/3

### 07-Implementation-Examples\01-Basic-Examples\rbac_auth.go

- 代码行数: 59
- 函数数量: 4
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\rsa_encrypt.go

- 代码行数: 64
- 函数数量: 4
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\secure_config.go

- 代码行数: 58
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\security_audit.go

- 代码行数: 107
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\security_monitor.go

- 代码行数: 79
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\sha256_hash.go

- 代码行数: 46
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 2/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\threat_detection.go

- 代码行数: 55
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### 07-Implementation-Examples\01-Basic-Examples\tls_client.go

- 代码行数: 45
- 函数数量: 1
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 1/3
- 并发安全评分: 2/3

### benchmarks\abac_auth_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\abac_auth_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\aes_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\api_rate_limit_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\audit_log_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\audit_log_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\audit_log_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\audit_log_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\channel_worker_pool_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\concurrent_counter_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\csrf_protection_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\csrf_protection_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\csrf_protection_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\csrf_protection_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\error_handling_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\input_validation_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\input_validation_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\input_validation_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\input_validation_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\input_validation_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\jwt_auth_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\oauth2_auth_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\oauth2_auth_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rbac_auth_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rbac_auth_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rbac_auth_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\rsa_encrypt_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\secure_config_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\secure_config_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\secure_config_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\security_audit_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\sha256_hash_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\sha256_hash_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\sha256_hash_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\sha256_hash_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\sha256_hash_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\tls_client_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\tls_client_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\tls_client_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\tls_client_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

### benchmarks\tls_client_benchmark_test_benchmark_test_benchmark_test_benchmark_test_benchmark_test.go

- 代码行数: 37
- 函数数量: 3
- Goroutine数量: 0
- Channel数量: 0
- 互斥锁数量: 0
- 错误处理评分: 0/3
- 并发安全评分: 2/3

## 优化建议

- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议检查并修复潜在的内存泄漏问题
- 性能优化建议: 在goroutine中使用fmt.Print可能影响性能
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议检查并修复潜在的内存泄漏问题
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
- 建议完善错误处理机制，包括错误返回、检查和恢复
