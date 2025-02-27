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

// InitDB 初始化数据库连接
func InitDB() {
	// 使用 AppConfig 中的配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MySQL.User,
		AppConfig.MySQL.Password,
		AppConfig.MySQL.Host,
		AppConfig.MySQL.Port,
		AppConfig.MySQL.Database)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用默认事务
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 只对 domains 和 users 表进行自动迁移
	err = DB.AutoMigrate(&model.Domain{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Printf("Successfully connected to database: %s", AppConfig.MySQL.Database)
} 