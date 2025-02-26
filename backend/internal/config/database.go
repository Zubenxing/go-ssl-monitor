package config

import (
	"fmt"
	"log"

	"github.com/go-ssl-monitor/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// NewDatabaseConfig 创建数据库配置
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "192.168.131.130",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "ssl_monitor",
	}
}

// InitDB 初始化数据库连接
func InitDB(config *DatabaseConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&model.Domain{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Printf("Successfully connected to database: %s", config.DBName)
} 