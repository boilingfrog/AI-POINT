package main

import "os"

// Config 是服务的运行配置。demo 只有监听地址一项。
type Config struct {
	Addr string
}

// loadConfig 从环境变量读取配置。PORT 缺省为 8080。
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{Addr: ":" + port}
}
