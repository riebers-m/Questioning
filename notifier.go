package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func getBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func getChatId() int64 {
	if v := os.Getenv("TELEGRAM_CHATID"); v != "" {
		ID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Fatal("invalid chat id:", err)
			return -1
		}
		return ID
	}
	return -1
}

type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessage(text string) error {
	botToken := getBotToken()
	chatID := getChatId()

	if botToken == "" {
		return fmt.Errorf("bot token empty")
	}

	if chatID == -1 {
		return fmt.Errorf("invalid chat id")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	msg := Message{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
