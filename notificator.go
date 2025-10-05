package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ederfmatos/go-concurrency/pkg/concurrency"
)

type (
	Notifier interface {
		SendMessage(message string) error
	}

	Discord struct {
		httpClient *HttpClient
	}

	DiscordMessage struct {
		Content string `json:"content"`
	}

	Telegram struct {
		httpClient *HttpClient
	}

	TelegramMessage struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}

	ComposeNotifier struct {
		notifiers []Notifier
	}
)

func NewNotificator() Notifier {
	var notifiers []Notifier
	if os.Getenv("DISCORD_WEBHOOK_URL") != "" {
		notifiers = append(notifiers, &Discord{
			httpClient: NewHttpClient(os.Getenv("DISCORD_WEBHOOK_URL")),
		})
	}
	if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
		notifiers = append(notifiers, &Telegram{
			httpClient: NewHttpClient(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_TOKEN"))),
		})
	}
	return &ComposeNotifier{notifiers: notifiers}
}

func (d Discord) SendMessage(message string) error {
	payload := DiscordMessage{
		Content: message,
	}

	_, err := d.httpClient.Post("", payload)
	if err != nil {
		return fmt.Errorf("send message to discord: %v", err)
	}

	return nil
}

func (t Telegram) SendMessage(message string) error {
	parsedMessage := message
	parsedMessage = strings.ReplaceAll(parsedMessage, "#", "")
	parsedMessage = strings.ReplaceAll(parsedMessage, "*", "")
	parsedMessage = strings.ReplaceAll(parsedMessage, "\n ", "\n")

	payload := TelegramMessage{
		Text:   parsedMessage,
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
	}

	_, err := t.httpClient.Post("", payload)
	if err != nil {
		return fmt.Errorf("send message to telegram: %v", err)
	}

	return nil
}

func (c ComposeNotifier) SendMessage(message string) error {
	_, err := concurrency.ForEach(c.notifiers, 2, func(notifier Notifier) (any, error) {
		return nil, notifier.SendMessage(message)
	})

	if err != nil {
		return fmt.Errorf("fail to send messages: %v", err)
	}

	return nil
}
