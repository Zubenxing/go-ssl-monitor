package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/internal/model"
	"gorm.io/gorm"
)

// GetBackupLogs 获取所有备份日志
func GetBackupLogs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var logs []model.BackupLog
	
	if err := db.Order("Id DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取备份日志失败"})
		return
	}
	
	c.JSON(http.StatusOK, logs)
}

// CreateBackupLog 创建备份日志
func CreateBackupLog(c *gin.Context) {
	var log model.BackupLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建备份日志失败"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// UpdateBackupLog 更新备份日志
func UpdateBackupLog(c *gin.Context) {
	id := c.Param("id")
	var log model.BackupLog
	
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	result := db.Model(&model.BackupLog{}).Where("Id = ?", id).Updates(map[string]interface{}{
		"EndTime": log.EndTime,
		"BackupStatus": log.BackupStatus,
		"AlertStatus": log.AlertStatus,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新备份日志失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到指定的备份日志"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetLastBackupLogByIP 根据IP获取最后一条备份记录
func GetLastBackupLogByIP(c *gin.Context) {
	ip := c.Param("ip")
	var log model.BackupLog
	
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("Ip = ?", ip).Order("Id DESC").First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到备份记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": log.Id})
}

// GetBackupStatusByIP 根据IP获取最后一条备份状态
func GetBackupStatusByIP(c *gin.Context) {
	ip := c.Param("ip")
	var log model.BackupLog
	
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("Ip = ?", ip).Order("Id DESC").First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到备份记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"backup_status": log.BackupStatus})
} 