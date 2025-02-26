package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ssl-monitor/pkg/ssl"
)

type CheckDomainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

func CheckDomain(c *gin.Context) {
	var req CheckDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "域名不能为空",
		})
		return
	}

	info, err := ssl.CheckCertificate(req.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
} 