package helper

import (
	"net/smtp"

	"github.com/andikabahari/eoplatform/config"
)

func SendEmail(to []string, message string) error {
	smtpConfig := config.LoadSMTPConfig()
	emailConfig := config.LoadEmailConfig()

	auth := smtp.PlainAuth("", emailConfig.Address, emailConfig.Password, smtpConfig.Host)
	addr := smtpConfig.Host + ":" + smtpConfig.Port

	if err := smtp.SendMail(addr, auth, emailConfig.Address, to, []byte(message)); err != nil {
		return err
	}

	return nil
}
