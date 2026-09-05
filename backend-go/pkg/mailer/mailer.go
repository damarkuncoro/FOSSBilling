package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
)

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type Message struct {
	From        string
	To          []string
	Subject     string
	HTMLBody    string
	TextBody    string
	Attachments []Attachment
}

type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// MockMailer captures sent emails for testing
type MockMailer struct {
	mu           sync.RWMutex
	SentMessages []Message
}

func NewMockMailer() *MockMailer {
	return &MockMailer{SentMessages: make([]Message, 0)}
}

func (m *MockMailer) Send(ctx context.Context, msg Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentMessages = append(m.SentMessages, msg)
	return nil
}

func (m *MockMailer) GetSentCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.SentMessages)
}

func (m *MockMailer) LastMessage() *Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.SentMessages) == 0 {
		return nil
	}
	last := m.SentMessages[len(m.SentMessages)-1]
	return &last
}

// SMTPMailer handles sending real emails over standard SMTP protocol
type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
	return &SMTPMailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (s *SMTPMailer) Send(ctx context.Context, msg Message) error {
	if msg.From == "" {
		msg.From = s.From
	}

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	headers := make(map[string]string)
	headers["From"] = msg.From
	headers["To"] = strings.Join(msg.To, ", ")
	headers["Subject"] = msg.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	headerStr := ""
	for k, v := range headers {
		headerStr += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	body := []byte(headerStr + "\r\n" + msg.HTMLBody)

	// Note: in local dev/testing without active SMTP server, gracefully handle or return error
	return smtp.SendMail(addr, auth, msg.From, msg.To, body)
}
