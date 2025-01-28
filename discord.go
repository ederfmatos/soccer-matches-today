package main

import (
	"fmt"
)

type (
	Notificator interface {
		SendMessage(message string) error
	}

	Discord struct {
		httpClient *HttpClient
	}

	DiscordMessage struct {
		Content string `json:"content"`
	}
)

func NewNotificator() Notificator {
	return &Discord{
		httpClient: NewHttpClient("https://discord.com/api/webhooks/1332046251893456967/wNWKwY1MVUM-hlPDJiB0TmsRrEE8AJthbT2qpIC31w0fXZLB2kDtjHTroOxpZC3c-scj"),
	}
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
