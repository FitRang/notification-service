package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Foxtrot-14/FitRang/notification-service/db"
)

func (r *MessageService) CreateMessage(raw []byte) error {
	var message db.BMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	messageBson := db.Message{
		Sender:    message.Sender.Username,
		Message:   message.Message,
		Receiver:  message.Receiver.Username,
		CreatedAt: now,
	}

	_, err = r.Repo.Col.InsertOne(context.Background(), messageBson)
	if err != nil {
		return err
	}
	return nil
}
