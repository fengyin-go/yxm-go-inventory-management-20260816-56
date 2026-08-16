// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 服务运行配置。
type Config struct {
	Addr              string // HTTP 监听地址，如 ":8080"
	LowStockThreshold int    // 全局默认低库存预警阈值
	MaxPageSize       int    // 列表接口最大分页大小
}

// Load 从环境变量加载配置，未设置时使用缺省值。
func Load() *Config {
	cfg := &Config{
		Addr:              ":" + getenv("PORT", "8080"),
		LowStockThreshold: getenvInt("LOW_STOCK_THRESHOLD", 10),
		MaxPageSize:       getenvInt("MAX_PAGE_SIZE", 100),
	}
	// ADDR 优先级高于 PORT。
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

// getenv 读取环境变量，为空时返回默认值。
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getenvInt 读取整型环境变量，解析失败或非正数时返回默认值。
func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

// String 返回配置的可读描述（不含敏感信息）。
func (c *Config) String() string {
	return fmt.Sprintf("addr=%s low_stock_threshold=%d max_page_size=%d",
		c.Addr, c.LowStockThreshold, c.MaxPageSize)
}
