package services

import "github.com/Foxtrot-14/FitRang/notification-service/repository"

type MessageService struct {
	Repo *repository.MessageRepository
}

func NewMessageService(r *repository.MessageRepository) *MessageService {
	return&MessageService{
		Repo: r,
	}
}
