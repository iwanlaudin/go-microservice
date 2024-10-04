package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/iwanlaudin/go-microservice/pkg/common/config"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailSender struct {
	config SMTPConfig
}

func NewEmailSender(cfg *config.Config) (*EmailSender, error) {
	port, err := strconv.Atoi(cfg.SmtpPort)
	if err != nil {
		return nil, err
	}

	return &EmailSender{config: SMTPConfig{
		Host:     cfg.SmtpHost,
		Port:     port,
		Username: cfg.SmtpUsername,
		Password: cfg.SmtpPassword,
	}}, nil
}

func (s *EmailSender) SendEmail(to []string, subject, body string) error {
	from := s.config.Username
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Compose the email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := fmt.Sprintf("Subject: %s\n%s\n\n%s", subject, mime, body)

	// Connect to the SMTP Server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	tlsConfig := &tls.Config{ServerName: s.config.Host}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set the sender and recipient
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return nil
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
