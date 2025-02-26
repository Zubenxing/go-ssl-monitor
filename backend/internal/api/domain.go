package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/internal/model"
	"github.com/go-ssl-monitor/pkg/ssl"
	"gorm.io/gorm"
)

// GetDomains 获取所有域名
func GetDomains(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var domains []model.Domain
	if err := db.Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取域名列表失败"})
		return
	}
	c.JSON(http.StatusOK, domains)
}

// AddDomain 添加新域名
func AddDomain(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var domain model.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查域名是否已存在
	var existingDomain model.Domain
	if err := db.Where("domain_name = ?", domain.DomainName).First(&existingDomain).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名已存在"})
		return
	}

	// 检查证书状态
	certInfo, err := ssl.CheckCertificate(domain.DomainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查证书失败"})
		return
	}

	domain.CertificateStatus = "VALID"
	if !certInfo.IsValid {
		domain.CertificateStatus = "ERROR"
	}
	domain.CertificateIssuer = certInfo.Issuer
	domain.CertificateExpiryDate = certInfo.NotAfter
	domain.LastChecked = time.Now()

	if err := db.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加域名失败"})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// UpdateDomain 更新域名信息
func UpdateDomain(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var domain model.Domain
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "域名不存在"})
		return
	}

	var updateData model.Domain
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	domain.NotificationEmail = updateData.NotificationEmail
	domain.AutoRenewal = updateData.AutoRenewal

	if err := db.Save(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新域名失败"})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// DeleteDomain 删除域名
func DeleteDomain(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if err := db.Delete(&model.Domain{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除域名失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// CheckDomainCertificate 检查指定域名的证书状态
func CheckDomainCertificate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var domain model.Domain
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "域名不存在"})
		return
	}

	certInfo, err := ssl.CheckCertificate(domain.DomainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查证书失败"})
		return
	}

	domain.CertificateStatus = "VALID"
	if !certInfo.IsValid {
		domain.CertificateStatus = "ERROR"
	}
	domain.CertificateIssuer = certInfo.Issuer
	domain.CertificateExpiryDate = certInfo.NotAfter
	domain.LastChecked = time.Now()

	if err := db.Save(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新证书状态失败"})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// ToggleAutoRenewal 切换自动续期状态
func ToggleAutoRenewal(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var domain model.Domain
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "域名不存在"})
		return
	}

	domain.AutoRenewal = !domain.AutoRenewal
	if err := db.Save(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新自动续期状态失败"})
		return
	}

	c.JSON(http.StatusOK, domain)
} 