package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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
	discord := &Discord{
		httpClient: NewHttpClient(os.Getenv("DISCORD_WEBHOOK_URL")),
	}
	telegram := &Telegram{
		httpClient: NewHttpClient(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_TOKEN"))),
	}
	return &ComposeNotifier{notifiers: []Notifier{discord, telegram}}
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
	var errs = make([]error, 0)
	for _, notifier := range c.notifiers {
		if err := notifier.SendMessage(message); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == len(c.notifiers) {
		return fmt.Errorf("fail to send messages: %v", errors.Join(errs...))
	}

	if len(errs) > 0 {
		log.Printf("Error on send a message: %v. But at least one was sent\n", errors.Join(errs...))
	}

	return nil
}
