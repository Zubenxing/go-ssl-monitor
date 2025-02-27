package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/go-ssl-monitor/internal/config"
)

type EmailSender struct {
	config *config.EmailConfig
}

func NewEmailSender(config *config.EmailConfig) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

func (e *EmailSender) SendAlertEmail(ip, serverName string, backupError string) error {
	// 如果邮件配置未启用，直接返回
	if e.config == nil || e.config.SMTPHost == "" {
		return fmt.Errorf("email configuration not set")
	}

	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.SMTPHost)

	subject := "备份异常告警通知"
	body := fmt.Sprintf(`
服务器备份异常告警：

IP地址: %s
服务器名称: %s
错误信息: %s

请及时检查并处理。

此邮件为系统自动发送，请勿回复。
`, ip, serverName, backupError)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", strings.Join(e.config.ToAddresses, ","), subject, body))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort),
		auth,
		e.config.FromAddress,
		e.config.ToAddresses,
		msg,
	)

	return err
} 