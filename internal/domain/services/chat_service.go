package services

import (
	"context"
	"strings"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type ChatService struct {
	storage ports.StoragePort
	llm     ports.LLMPort
}

func NewChatService(s ports.StoragePort, l ports.LLMPort) *ChatService {
	return &ChatService{storage: s, llm: l}
}

func (cs *ChatService) ExecuteChat(ctx context.Context, branchID int64, modelID string, prompt string, cfg *models.Config, onChunk func(string)) error {
	branch, err := cs.storage.GetBranch(ctx, branchID)
	var parentID *int64
	if err == nil && branch != nil {
		parentID = branch.CurrentMessageID
	}

	userMsg := &models.Message{
		ParentID: parentID,
		Role:     "user",
		Content:  prompt,
	}

	if err := cs.storage.SaveMessage(ctx, userMsg); err != nil {
		return err
	}

	history, err := cs.storage.GetThreadContext(ctx, userMsg.ID)
	if err != nil {
		return err
	}

	outCh, errCh, err := cs.llm.StreamCompletion(ctx, cfg, modelID, history)
	if err != nil {
		return err
	}

	var assistantResponse strings.Builder

	firstChunk := true

	for outCh != nil || errCh != nil {
		select {
		case chunk, ok := <-outCh:
			if !ok {
				outCh = nil
				continue
			}

			if firstChunk {
				chunk = strings.TrimLeft(chunk, "\r\n ")
				if chunk == "" {
					continue
				}

				firstChunk = false
			}

			assistantResponse.WriteString(chunk)
			onChunk(chunk)
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	assistantResponseStr := strings.TrimRight(assistantResponse.String(), "\r\n ")

	assistantMsg := &models.Message{
		ParentID: &userMsg.ID,
		Role:     "assistant",
		Content:  assistantResponseStr,
	}

	if err := cs.storage.SaveMessage(ctx, assistantMsg); err != nil {
		return err
	}

	return cs.storage.UpdateBranchHead(ctx, branchID, assistantMsg.ID)
}
