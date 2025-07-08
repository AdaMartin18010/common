package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 配置结构体
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	APIKey     string `json:"api_key"`
	Port       int    `json:"port"`
}

// 加载配置文件
func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// 校验配置有效性
func validateConfig(cfg *Config) error {
	if cfg.DBUser == "" || cfg.DBPassword == "" {
		return fmt.Errorf("数据库用户名或密码不能为空")
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("API密钥不能为空")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("端口号无效")
	}
	return nil
}

func main() {
	// 假设有config.json文件
	cfg, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("加载配置失败:", err)
		os.Exit(1)
	}
	if err := validateConfig(cfg); err != nil {
		fmt.Println("配置校验失败:", err)
		os.Exit(1)
	}
	fmt.Println("配置加载与校验成功:", cfg)
}
