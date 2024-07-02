package vertmailer

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	smtpserver "github.com/emersion/go-smtp"
)

type mockBackend struct{}

func (bkd *mockBackend) Login(_ *smtpserver.Conn, username, password string) (smtpserver.Session, error) {
	if username == "user" && password == "pass" {
		return &mockSession{}, nil
	}
	return nil, fmt.Errorf("invalid username or password")
}

func (bkd *mockBackend) AnonymousLogin(_ *smtpserver.Conn) (smtpserver.Session, error) {
	return &mockSession{}, nil
}

func (bkd *mockBackend) NewSession(_ *smtpserver.Conn) (smtpserver.Session, error) {
	return &mockSession{}, nil
}

type mockSession struct{}

func (s *mockSession) Mail(from string, opts *smtpserver.MailOptions) error {
	return nil
}

func (s *mockSession) Rcpt(to string, opts *smtpserver.RcptOptions) error {
	return nil
}

func (s *mockSession) Data(r io.Reader) error {
	// Read and print the email message for verification
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}
	fmt.Println(buf.String())
	return nil
}

func (s *mockSession) Reset() {}

func (s *mockSession) Logout() error {
	return nil
}

func startMockSMTPServer(addr string) *smtpserver.Server {
	be := &mockBackend{}
	s := smtpserver.NewServer(be)

	s.Addr = addr
	s.Domain = "localhost"
	s.AllowInsecureAuth = true

	go func() {
		if err := s.ListenAndServe(); err != nil {
			fmt.Println("Mock SMTP Server error:", err)
		}
	}()

	// Give the server some time to start
	time.Sleep(500 * time.Millisecond)
	return s
}

func TestSendMail(t *testing.T) {
	mockServerAddr := "127.0.0.1:1025"
	mockServer := startMockSMTPServer(mockServerAddr)
	defer mockServer.Close()

	server := MailServer{
		Host:       "127.0.0.1",
		Port:       "1025",
		Username:   "user",
		Password:   "pass",
		CACertPath: "", // Not needed for mock server
		UseTLS:     false,
	}

	mailer := NewMailerService(server)

	mail := Mail{
		From:    "your-email@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    `<h1>This is a test email body.</h1><p style="color:red;">This is a paragraph.</p>`,
	}

	if err := mailer.SendMail(mail); err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}
