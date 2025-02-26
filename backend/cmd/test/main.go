package main

import (
	"fmt"
	"log"

	"github.com/go-ssl-monitor/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(192.168.131.130:3306)/ssl_monitor?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 创建一个新的测试用户
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成密码哈希失败:", err)
	}

	// 更新管理员密码
	result := db.Model(&model.User{}).Where("username = ?", "admin").Update("password", string(hashedPassword))
	if result.Error != nil {
		log.Fatal("更新密码失败:", result.Error)
	}
	fmt.Printf("更新了 %d 条记录\n", result.RowsAffected)

	// 验证密码
	var user model.User
	if err := db.Where("username = ?", "admin").First(&user).Error; err != nil {
		log.Fatal("查找用户失败:", err)
	}

	fmt.Printf("找到用户: %+v\n", user)
	
	if user.ComparePassword(password) {
		fmt.Println("密码验证成功!")
	} else {
		fmt.Println("密码验证失败!")
	}
} 