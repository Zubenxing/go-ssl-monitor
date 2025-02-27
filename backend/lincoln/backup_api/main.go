package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type BackupLog struct {
	Id            uint64 `db:"Id" json:"id"`
	Ip            string `db:"Ip" json:"ip"`
	ServerName    string `db:"ServerName" json:"server_name"`
	StartTime     string `db:"StartTime" json:"start_time"`
	EndTime       string `db:"EndTime" json:"end_time"`
	BackupStatus  int    `db:"BackupStatus" json:"backup_status"`
	AlertStatus   int    `db:"AlertStatus" json:"alert_status"`
	ScriptVersion string `db:"ScriptVersion" json:"script_version"`
}

var db *sqlx.DB

func initDB() {
	var err error
	dsn := "code_om_backup_status:M2DKM^N6s7b#iumz@tcp(mysql03-prod-modem.mysql.database.azure.com:3306)/om_backup?parseTime=true"

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
}

func main() {
	initDB()
	defer db.Close()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/backupLogs", getBackupLogs)
	router.POST("/backupLogs", createBackupLog)
	router.PUT("/backupLogs/:id", updateBackupLogEndTime)
	router.DELETE("/backupLogs/:id", deleteBackupLog)
	router.GET("/getId/:ip", getLastBackupLogIdByIp)
	router.GET("/getStatus/:ip", getLastBackupLogStatus)

	router.Run(":8080")
}

// 获取所有备份日志
func getBackupLogs(c *gin.Context) {
	var logs []BackupLog
	err := db.Select(&logs, "SELECT * FROM BackupLogs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("error: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, logs)
}

func createBackupLog(c *gin.Context) {
	var log BackupLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("绑定 JSON 错误: %v\n", err)
		return
	}

	// 设置当前时间为 StartTime
	//log.StartTime = time.Now().Format("2006-01-02 15:04:05")

	// 插入数据到 BackupLogs 表
	query := `
        INSERT INTO BackupLogs (Ip, ServerName, StartTime, BackupStatus, AlertStatus, ScriptVersion)
        VALUES (:Ip, :ServerName, :StartTime, :BackupStatus, :AlertStatus, :ScriptVersion)
    `
	// 使用 NamedExec 插入数据
	result, err := db.NamedExec(query, log)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("数据库插入错误: %v\n", err)
		return
	}

	// 获取插入的 ID
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取插入的 ID"})
		fmt.Printf("获取 LastInsertId 错误: %v\n", err)
		return
	}

	// 将 ID 赋值给 log 并返回
	log.Id = uint64(id)
	c.JSON(http.StatusOK, log)
}

// 更新备份日志
func updateBackupLogEndTime(c *gin.Context) {

	id := c.Param("id")
	var log BackupLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logId, err := strconv.ParseUint(id, 10, 64) // 10是基数，64是结果类型
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	log.Ip = ""
	log.ServerName = ""
	log.StartTime = ""

	log.Id = logId // Use id from path parameter
	query := `
        UPDATE BackupLogs
        SET EndTime = :EndTime,
            BackupStatus = :BackupStatus, AlertStatus = :AlertStatus
        WHERE Id = :Id
    `
	_, err = db.NamedExec(query, log)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, log)
}

// 删除备份日志
func deleteBackupLog(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM BackupLogs WHERE Id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Backup log deleted"})

}

// 根据 IP 地址查找最后一条数据的 ID
func getLastBackupLogIdByIp(c *gin.Context) {
	ip := c.Param("ip") // 从请求中获取 IP 参数
	fmt.Printf("ip:%v\n", ip)
	fmt.Printf("ips:%v\n", c.Params)
	var log BackupLog
	query := `SELECT * FROM BackupLogs WHERE Ip = ? ORDER BY StartTime DESC LIMIT 1`

	// 执行查询，查找匹配 IP 的最新记录
	err := db.Get(&log, query, ip)
	if err != nil {
		fmt.Printf("getLastBackupLogIdByIp: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回查询到的 ID
	c.JSON(http.StatusOK, gin.H{"id": log.Id})
}

// 根据 IP 地址查找最后一条数据的 BackupStatus
func getLastBackupLogStatus(c *gin.Context) {
	ip := c.Param("ip") // 从请求中获取 IP 参数
	fmt.Printf("ip:%v\n", ip)
	fmt.Printf("ips:%v\n", c.Params)
	var log BackupLog
	query := `SELECT * FROM BackupLogs WHERE Ip = ? ORDER BY StartTime DESC LIMIT 1`

	// 执行查询，查找匹配 IP 的最新记录
	err := db.Get(&log, query, ip)
	if err != nil {
		fmt.Printf("getLastBackupLogIdByIp: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回查询到的 BackupStatus
	c.JSON(http.StatusOK, gin.H{"backup_status": log.BackupStatus})
}
