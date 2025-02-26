package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 测试密码
	password := "admin123"
	
	// 生成新的哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成哈希失败:", err)
	}
	fmt.Printf("新生成的哈希: %s\n", string(hash))

	// 验证新生成的哈希
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("新哈希验证失败:", err)
	} else {
		fmt.Println("新哈希验证成功!")
	}

	// 验证数据库中的哈希
	storedHash := "$2a$10$QOXc1.1XcW3IGV.mj3Jh/.6.WXe.V3fGWTxkGRgvGtY8xyxAGJBAO"
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Println("存储的哈希验证失败:", err)
	} else {
		fmt.Println("存储的哈希验证成功!")
	}

	// 生成一个新的哈希用于数据库
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成新哈希失败:", err)
	}
	fmt.Printf("建议使用的新哈希: %s\n", string(newHash))
} 