package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构体
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`

	Email struct {
		SMTPHost    string `yaml:"smtp_host"`
		SMTPPort    int    `yaml:"smtp_port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		FromAddress string `yaml:"from_address"`
	} `yaml:"email"`
}

// AppConfig 全局配置变量
var AppConfig Config

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	log.Printf("Configuration loaded successfully from: %s", configPath)
}