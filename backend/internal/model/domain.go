package model

import "time"

type Domain struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	DomainName         string    `json:"domainName" gorm:"unique;not null"`
	NotificationEmail   string    `json:"notificationEmail"`
	CertificateStatus  string    `json:"certificateStatus"`
	CertificateIssuer  string    `json:"certificateIssuer"`
	CertificateExpiryDate time.Time `json:"certificateExpiryDate"`
	LastChecked        time.Time `json:"lastChecked"`
	AutoRenewal        bool      `json:"autoRenewal" gorm:"default:true"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
} 