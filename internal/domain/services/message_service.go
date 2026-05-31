package services

import (
	"context"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type MessageService struct {
	messageStorage ports.MessageStoragePort
}

func NewMessageService(messageStorage ports.MessageStoragePort) *MessageService {
	return &MessageService{
		messageStorage: messageStorage,
	}
}

func (s *MessageService) GetThreadContext(ctx context.Context, headID int64) ([]models.Message, error) {
	return s.messageStorage.GetThreadContext(ctx, headID)
}
