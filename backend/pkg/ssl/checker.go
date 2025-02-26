package ssl

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"
)

type CertInfo struct {
	Domain            string    `json:"domain"`
	Issuer           string    `json:"issuer"`
	NotBefore        time.Time `json:"not_before"`
	NotAfter         time.Time `json:"not_after"`
	RemainingDays    int       `json:"remaining_days"`
	IsValid          bool      `json:"is_valid"`
	ValidationErrors []string  `json:"validation_errors,omitempty"`
}

func CheckCertificate(domain string) (*CertInfo, error) {
	// 确保域名格式正确
	if !strings.Contains(domain, ":") {
		domain = domain + ":443"
	}

	conn, err := tls.Dial("tcp", domain, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return &CertInfo{
			Domain:            domain,
			IsValid:          false,
			ValidationErrors: []string{fmt.Sprintf("连接失败: %v", err)},
		}, nil
	}
	defer conn.Close()

	// 获取证书信息
	cert := conn.ConnectionState().PeerCertificates[0]
	now := time.Now()

	info := &CertInfo{
		Domain:         domain,
		Issuer:        cert.Issuer.CommonName,
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		RemainingDays: int(cert.NotAfter.Sub(now).Hours() / 24),
		IsValid:       true,
	}

	// 验证证书
	if now.Before(cert.NotBefore) {
		info.IsValid = false
		info.ValidationErrors = append(info.ValidationErrors, "证书还未生效")
	}
	if now.After(cert.NotAfter) {
		info.IsValid = false
		info.ValidationErrors = append(info.ValidationErrors, "证书已过期")
	}

	return info, nil
} 