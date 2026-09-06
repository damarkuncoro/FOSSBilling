package mailer

import (
	"fmt"
	"strings"
)

// MessageBuilder provides a fluent builder pattern for constructing email messages
type MessageBuilder struct {
	msg Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: Message{
			To:          make([]string, 0),
			Attachments: make([]Attachment, 0),
		},
	}
}

func (b *MessageBuilder) From(senderName, senderEmail string) *MessageBuilder {
	if senderName != "" {
		b.msg.From = fmt.Sprintf("%s <%s>", senderName, senderEmail)
	} else {
		b.msg.From = senderEmail
	}
	return b
}

func (b *MessageBuilder) To(recipients ...string) *MessageBuilder {
	b.msg.To = append(b.msg.To, recipients...)
	return b
}

func (b *MessageBuilder) Subject(subject string) *MessageBuilder {
	b.msg.Subject = strings.TrimSpace(subject)
	return b
}

func (b *MessageBuilder) HTMLBody(html string) *MessageBuilder {
	b.msg.HTMLBody = html
	return b
}

func (b *MessageBuilder) TextBody(text string) *MessageBuilder {
	b.msg.TextBody = text
	return b
}

func (b *MessageBuilder) Attach(filename, contentType string, data []byte) *MessageBuilder {
	b.msg.Attachments = append(b.msg.Attachments, Attachment{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	})
	return b
}

func (b *MessageBuilder) AttachPDF(filename string, data []byte) *MessageBuilder {
	return b.Attach(filename, "application/pdf", data)
}

func (b *MessageBuilder) Build() (Message, error) {
	if len(b.msg.To) == 0 {
		return Message{}, fmt.Errorf("message requires at least one recipient in 'To'")
	}
	if b.msg.Subject == "" {
		return Message{}, fmt.Errorf("message requires a non-empty subject")
	}
	return b.msg, nil
}
