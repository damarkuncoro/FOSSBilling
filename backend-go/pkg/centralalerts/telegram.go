package centralalerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TelegramService struct {
	botToken string
	chatID   string
	client   *http.Client
}

func NewTelegramService(token, chatID string) *TelegramService {
	return &TelegramService{
		botToken: token,
		chatID:   chatID,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *TelegramService) IsConfigured() bool {
	return s.botToken != "" && s.chatID != ""
}

func (s *TelegramService) SendMessage(message string) error {
	return s.send(message, "")
}

func (s *TelegramService) SendAlert(title, message string, priority bool) error {
	emoji := "ℹ️"
	if priority {
		emoji = "🔥"
	}

	formatted := fmt.Sprintf("%s <b>%s</b>\n\n%s\n\n<i>Generated at: %s</i>",
		emoji, title, message, time.Now().Format(time.RFC1123))

	return s.send(formatted, "HTML")
}

func (s *TelegramService) send(text, parseMode string) error {
	if !s.IsConfigured() {
		return nil
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)
	payload := map[string]string{
		"chat_id": s.chatID,
		"text":    text,
	}
	if parseMode != "" {
		payload["parse_mode"] = parseMode
	}

	body, _ := json.Marshal(payload)
	resp, err := s.client.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
