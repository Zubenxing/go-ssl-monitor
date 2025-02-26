package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/internal/api"
	"github.com/go-ssl-monitor/internal/config"
)

func main() {
	// 初始化数据库连接
	dbConfig := config.NewDatabaseConfig()
	config.InitDB(dbConfig)

	r := gin.Default()

	// CORS配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 基础健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由
	apiGroup := r.Group("/api")
	{
		// 域名管理
		apiGroup.GET("/domains", api.GetDomains)
		apiGroup.POST("/domains", api.AddDomain)
		apiGroup.PUT("/domains/:id", api.UpdateDomain)
		apiGroup.DELETE("/domains/:id", api.DeleteDomain)
		apiGroup.POST("/domains/:id/check", api.CheckDomainCertificate)
		apiGroup.PUT("/domains/:id/auto-renewal", api.ToggleAutoRenewal)

		// 单个域名证书检查
		apiGroup.POST("/check-domain", api.CheckDomain)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 