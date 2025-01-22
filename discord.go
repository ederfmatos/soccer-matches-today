package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DiscordMessage struct {
	Content string `json:"content"`
}

func SendMessageToDiscord(message string) error {
	url := os.Getenv("DISCORD_WEBHOOK_URL")

	payload := DiscordMessage{
		Content: message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal discord message: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("fail to send message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to send message, response status: %d", resp.StatusCode)
	}

	return nil
}
