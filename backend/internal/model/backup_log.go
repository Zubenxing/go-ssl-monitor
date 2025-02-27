package model

import (
	"time"
)

type BackupLog struct {
	Id            uint   `json:"id" gorm:"column:Id;primaryKey"`
	Ip            string `json:"ip" gorm:"column:Ip;not null"`
	ServerName    string `json:"server_name" gorm:"column:ServerName;not null"`
	StartTime     string `json:"start_time" gorm:"column:StartTime;not null;default:''"`
	EndTime       string `json:"end_time" gorm:"column:EndTime;default:''"`
	BackupStatus  int    `json:"backup_status" gorm:"column:BackupStatus;not null;default:0"`
	AlertStatus   int    `json:"alert_status" gorm:"column:AlertStatus;not null;default:0"`
	ScriptVersion string `json:"script_version" gorm:"column:ScriptVersion;not null;default:'0'"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 指定表名
func (BackupLog) TableName() string {
	return "backuplogs"
} 