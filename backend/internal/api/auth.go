package api

import (
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-ssl-monitor/internal/model"
    "github.com/golang-jwt/jwt/v5"
    "gorm.io/gorm"
)

var jwtSecret = []byte("your-secret-key")

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Login 处理用户登录
func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    log.Printf("Login attempt for user: %s with password: %s", req.Username, req.Password)

    db := c.MustGet("db").(*gorm.DB)
    var user model.User
    
    // 查询用户并打印结果
    result := db.Where("username = ?", req.Username).First(&user)
    if result.Error != nil {
        log.Printf("Database error or user not found: %v", result.Error)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    log.Printf("Found user: %+v", user)

    // 验证密码
    if !user.ComparePassword(req.Password) {
        log.Printf("Password verification failed for user: %s", req.Username)
        log.Printf("Stored password hash: %s", user.Password)
        log.Printf("Provided password: %s", req.Password)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
        return
    }

    log.Printf("Password correct for user: %s, generating token", req.Username)

    // 生成 JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        log.Printf("Failed to generate token: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
        return
    }

    // 更新最后登录时间
    user.LastLogin = time.Now()
    if err := db.Save(&user).Error; err != nil {
        log.Printf("Failed to update last login time: %v", err)
    }

    log.Printf("Login successful for user: %s", req.Username)

    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证token"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        tokenString = strings.TrimSpace(tokenString)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token声明"})
            c.Abort()
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Set("username", claims["username"])

        c.Next()
    }
} 