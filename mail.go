package vertmailer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/smtp"
	"os"
)

// Mail represents the email details.
type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// MailServer represents the SMTP server details.
type MailServer struct {
	Host       string
	Port       string
	Username   string
	Password   string
	CACertPath string
}

// Mailer interface defines the methods for sending emails.
type Mailer interface {
	SendMail(mail Mail) error
}

// MailerService implements the Mailer interface.
type MailerService struct {
	Server MailServer
}

// NewMailerService creates a new MailerService.
func NewMailerService(server MailServer) *MailerService {
	return &MailerService{Server: server}
}

// LoadCACert loads the CA certificate from the well-known location.
func LoadCACert(path string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}
	return caCertPool, nil
}

// SendMail sends an email using the specified MailServer.
func (ms *MailerService) SendMail(mail Mail) error {
	if err := ValidateMail(mail); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	mail = SanitizeMail(mail)

	auth := smtp.PlainAuth("", ms.Server.Username, ms.Server.Password, ms.Server.Host)

	header := map[string]string{
		"From":    mail.From,
		"To":      mail.To[0],
		"Subject": mail.Subject,
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + mail.Body

	// Load the CA certificate.
	caCertPool, err := LoadCACert(ms.Server.CACertPath)
	if err != nil {
		return fmt.Errorf("failed to load CA certificate: %w", err)
	}

	// Establish a TLS connection with CA verification.
	tlsconfig := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: ms.Server.Host,
	}

	conn, err := tls.Dial("tcp", ms.Server.Host+":"+ms.Server.Port, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, ms.Server.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication error: %w", err)
	}

	if err = client.Mail(mail.From); err != nil {
		return fmt.Errorf("failed to set mail sender: %w", err)
	}

	for _, to := range mail.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get Data writer: %w", err)
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close Data writer: %w", err)
	}

	if err = client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP client: %w", err)
	}

	return nil
}
