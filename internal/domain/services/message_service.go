package services

import (
	"context"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type MessageService struct {
	messageStorage ports.MessageStoragePort
	branchStorage  ports.BranchStoragePort
}

func NewMessageService(messageStorage ports.MessageStoragePort, branchStorage ports.BranchStoragePort) *MessageService {
	return &MessageService{
		messageStorage: messageStorage,
		branchStorage:  branchStorage,
	}
}

func (s *MessageService) GetThreadContext(ctx context.Context, headID int64) ([]models.Message, error) {
	return s.messageStorage.GetThreadContext(ctx, headID)
}

func (s *MessageService) SaveMessage(ctx context.Context, branchID int64, message *models.Message) error {
	err := s.messageStorage.SaveMessage(ctx, message)
	if err != nil {
		return err
	}

	return s.branchStorage.UpdateBranchHead(ctx, branchID, message.ID)
}
