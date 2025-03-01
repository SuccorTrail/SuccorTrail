package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func SendVerificationEmail(email string, token string) error {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	logrus.WithFields(logrus.Fields{
		"smtpHost": smtpHost,
		"smtpPort": smtpPort,
		"from":     from,
		"to":       email,
	}).Debug("Attempting to send verification email")

	// Set up authentication
	auth := smtp.PlainAuth("", from, password, smtpHost+":"+smtpPort)

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	// Connect to server with timeout
	conn, err := net.DialTimeout("tcp", smtpHost+":"+smtpPort, 15*time.Second)
	if err != nil {
		return fmt.Errorf("smtp connection failed: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("smtp client creation failed: %w", err)
	}
	defer client.Close()

	// Start TLS
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("starttls failed: %w", err)
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(email); err != nil {
		return err
	}

	// Send email body
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	// Build message
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
<html>
<body>
    <h2>SuccorTrail Email Verification</h2>
    <p>Please click the button below to verify your email address:</p>
    <a href="%s/verify-email?token=%s" style="
        background-color: #4CAF50;
        border: none;
        color: white;
        padding: 15px 32px;
        text-align: center;
        text-decoration: none;
        display: inline-block;
        font-size: 16px;
        margin: 4px 2px;
        cursor: pointer;
    ">Verify Email</a>
    <p>This link will expire in 24 hours.</p>
</body>
</html>
`, os.Getenv("APP_URL"), token)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, email, subject, body)

	_, err = fmt.Fprintf(wc, msg)
	return err
}
