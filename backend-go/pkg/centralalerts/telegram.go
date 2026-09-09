package centralalerts

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if !s.IsConfigured() {
		return nil
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)
	payload := map[string]string{
		"chat_id":    s.chatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	body, _ := json.Marshal(payload)
	resp, err := s.client.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: status %d", resp.StatusCode)
	}

	return nil
}
