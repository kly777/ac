package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// 配置结构体（需导出字段）
type Config struct {
	DeepSeekAPIKey string `json:"deepseek_api_key"`
}

// Get 读取当前目录下的 config.json
func Get() (*Config, error) {
	// 1. 读取文件内容
	data, err := os.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在: %v", err)
		}
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 2. 解析JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	return &config, nil
}
