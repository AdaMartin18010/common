package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CustomError 自定义错误类型
type CustomError struct {
	Code    int
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("错误代码: %d, 消息: %s, 原因: %v", e.Code, e.Message, e.Err)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

// NewCustomError 创建自定义错误
func NewCustomError(code int, message string, err error) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Config 配置结构
type Config struct {
	Port     int
	Host     string
	Database string
}

// LoadConfig 加载配置
func LoadConfig(filename string) (*Config, error) {
	// 模拟文件读取
	if filename == "" {
		return nil, NewCustomError(1001, "配置文件路径不能为空", errors.New("invalid filename"))
	}

	// 模拟配置文件不存在
	if filename == "nonexistent.conf" {
		return nil, NewCustomError(1002, "配置文件不存在", os.ErrNotExist)
	}

	// 模拟配置解析错误
	if filename == "invalid.conf" {
		return nil, NewCustomError(1003, "配置文件格式错误", errors.New("invalid format"))
	}

	// 成功情况
	return &Config{
		Port:     8080,
		Host:     "localhost",
		Database: "mydb",
	}, nil
}

// ParsePort 解析端口号
func ParsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("端口号解析失败: %w", err)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("端口号超出范围 (1-65535): %d", port)
	}

	return port, nil
}

// StartServer 启动服务器
func StartServer(config *Config) error {
	if config == nil {
		return errors.New("配置不能为空")
	}

	if config.Port <= 0 {
		return errors.New("端口号无效")
	}

	// 模拟服务器启动
	fmt.Printf("服务器启动成功: %s:%d\n", config.Host, config.Port)
	return nil
}

// ProcessData 处理数据
func ProcessData(data []string) error {
	if len(data) == 0 {
		return errors.New("数据不能为空")
	}

	for i, item := range data {
		if item == "" {
			return fmt.Errorf("第 %d 项数据为空", i+1)
		}
		// 模拟数据处理
		fmt.Printf("处理数据: %s\n", item)
	}

	return nil
}

// RecoverFromPanic 从panic中恢复
func RecoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Printf("从panic中恢复: %v\n", r)
	}
}

// SafeOperation 安全操作包装器
func SafeOperation(operation func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("操作发生panic: %v", r)
		}
	}()

	return operation()
}

func main() {
	// 示例1: 基本错误处理
	fmt.Println("=== 基本错误处理 ===")

	config, err := LoadConfig("")
	if err != nil {
		var customErr *CustomError
		if errors.As(err, &customErr) {
			fmt.Printf("自定义错误: %s\n", customErr.Error())
		} else {
			fmt.Printf("一般错误: %v\n", err)
		}
	}

	// 示例2: 错误包装和展开
	fmt.Println("\n=== 错误包装和展开 ===")

	port, err := ParsePort("invalid")
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		// 展开错误
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			fmt.Printf("原始错误: %v\n", unwrapped)
		}
	}

	// 示例3: 错误类型检查
	fmt.Println("\n=== 错误类型检查 ===")

	config, err = LoadConfig("nonexistent.conf")
	if err != nil {
		var customErr *CustomError
		if errors.As(err, &customErr) {
			fmt.Printf("错误代码: %d\n", customErr.Code)
			fmt.Printf("错误消息: %s\n", customErr.Message)
		}

		// 检查是否为文件不存在错误
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("文件确实不存在")
		}
	}

	// 示例4: 安全操作
	fmt.Println("\n=== 安全操作 ===")

	err = SafeOperation(func() error {
		// 模拟可能panic的操作
		data := []string{"item1", "item2", ""}
		return ProcessData(data)
	})

	if err != nil {
		fmt.Printf("安全操作错误: %v\n", err)
	}

	// 示例5: 成功情况
	fmt.Println("\n=== 成功情况 ===")

	config, err = LoadConfig("app.conf")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	err = StartServer(config)
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	fmt.Println("程序执行完成")
}
