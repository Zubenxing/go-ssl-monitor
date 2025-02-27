package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	
	db := c.MustGet("db").(*gorm.DB)
	var user model.User
	
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"username": user.Username,
			"email": user.Email,
			"avatar": user.Avatar,
		},
	})
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		Email string `json:"email"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	
	db := c.MustGet("db").(*gorm.DB)
	result := db.Model(&model.User{}).Where("id = ?", userID).Update("email", input.Email)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdatePassword 修改密码
func UpdatePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword    string `json:"newPassword"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	
	db := c.MustGet("db").(*gorm.DB)
	var user model.User
	
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码错误"})
		return
	}
	
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	
	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	userID, _ := c.Get("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件"})
		return
	}
	
	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("avatar_%v%s", userID, ext)
	filepath := filepath.Join("uploads/avatars", filename)
	
	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	
	// 更新用户头像路径
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Model(&model.User{}).Where("id = ?", userID).Update("avatar", filepath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url": "/" + filepath,
		},
	})
} 