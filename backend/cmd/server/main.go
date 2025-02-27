package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/internal/api"
	"github.com/go-ssl-monitor/internal/config"
)

func main() {
	// 加载配置文件
	configPath := filepath.Join("configs", "config.yaml")
	config.LoadConfig(configPath)

	// 初始化数据库连接
	config.InitDB()
	
	// 创建gin实例
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 添加数据库中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", config.DB)
		c.Next()
	})

	// 配置CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API路由
	apiGroup := r.Group("/api")
	{
		// 认证路由
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/login", api.Login)
		}

		// 需要认证的路由
		protected := apiGroup.Group("")
		protected.Use(api.AuthMiddleware())
		{
			// 域名管理
			protected.GET("/domains", api.GetDomains)
			protected.POST("/domains", api.AddDomain)
			protected.PUT("/domains/:id", api.UpdateDomain)
			protected.DELETE("/domains/:id", api.DeleteDomain)
			protected.POST("/domains/:id/check", api.CheckDomainCertificate)
			protected.PUT("/domains/:id/auto-renewal", api.ToggleAutoRenewal)

			// 备份日志相关路由
			protected.GET("/backupLogs", api.GetBackupLogs)
			protected.POST("/backupLogs", api.CreateBackupLog)
			protected.PUT("/backupLogs/:id", api.UpdateBackupLog)
			protected.GET("/getId/:ip", api.GetLastBackupLogByIP)
			protected.GET("/getStatus/:ip", api.GetBackupStatusByIP)
		}
	}

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(r.Run(serverAddr))
}