package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

func emitErr(err error) models.StreamEvent {
	return models.StreamEvent{
		Type: models.EventError,
		Err:  err,
	}
}

type ChatService struct {
	workspaceService ports.WorkspaceServicePort
	branchService    ports.BranchServicePort
	messageService   ports.MessageServicePort
	modelService     ports.ModelServicePort
	toolService      ports.ToolServicePort
	llm              ports.LLMPort
}

func NewChatService(
	workspaceService ports.WorkspaceServicePort,
	branchService ports.BranchServicePort,
	messageService ports.MessageServicePort,
	modelService ports.ModelServicePort,
	toolService ports.ToolServicePort,
	llm ports.LLMPort,
) *ChatService {
	return &ChatService{
		workspaceService: workspaceService,
		branchService:    branchService,
		messageService:   messageService,
		modelService:     modelService,
		toolService:      toolService,
		llm:              llm,
	}
}

func (cs *ChatService) ExecuteChat(
	ctx context.Context,
	message string,
) (*models.AgentSession, error) {
	currentWorkspace, err := cs.workspaceService.GetActiveWorkspace(ctx)
	if err != nil {
		return nil, err
	}

	currentBranch, err := cs.branchService.GetActiveBranch(ctx, currentWorkspace.Name)
	if err != nil {
		return nil, err
	}

	userMsg := &models.Message{
		ParentID: currentBranch.CurrentMessageID,
		Role:     "user",
		Content:  message,
	}

	if err = cs.messageService.SaveMessage(ctx, currentBranch.ID, userMsg); err != nil {
		return nil, err
	}

	model, err := cs.modelService.GetCurrentModel(ctx)
	if err != nil {
		return nil, err
	}

	tools, err := cs.toolService.ListTools()
	if err != nil {
		return nil, err
	}

	resutltCh := make(chan models.StreamEvent)
	approveCh := make(chan models.ApprovalResponse)
	toolCalls := []models.ToolCall{}
	headID := userMsg.ID

	go func() {
		defer close(resutltCh)
		defer close(approveCh)
		for {

			history, err := cs.messageService.GetThreadContext(ctx, headID)
			if err != nil {
				resutltCh <- emitErr(err)
				return
			}
			clear(toolCalls)
			toolCalls = toolCalls[:0]

			eventCh, err := cs.llm.StreamCompletion(ctx, models.CompletionRequest{
				ThreadHistory: history,
				Model:         *model,
				Tools:         tools,
			})
			if err != nil {
				resutltCh <- emitErr(err)
				return
			}

			var messageContent strings.Builder

			for {
				event, ok := <-eventCh
				if !ok {
					break
				}

				resutltCh <- event
				switch event.Type {
				case models.EventToolCall:
					toolCalls = append(toolCalls, *event.ToolCall)
				case models.EventText:
					messageContent.Write([]byte(event.Text))
				}
			}

			modelMessage := &models.Message{
				ParentID: &headID,
				Role:     "assistant",
				Content:  messageContent.String(),
			}
			err = cs.messageService.SaveMessage(ctx, currentBranch.ID, modelMessage)
			if err != nil {
				resutltCh <- emitErr(err)
			}
			headID = modelMessage.ID

			for _, toolCall := range toolCalls {
				resutltCh <- models.StreamEvent{
					Type:     models.EventToolCallRequest,
					ToolCall: &toolCall,
				}

				approval, ok := <-approveCh
				if !ok {
					resutltCh <- emitErr(fmt.Errorf("failed to read from approveCh"))
					return
				}

				var toolResult *models.ToolResponse
				if approval.Approved {
					toolResult, err = cs.toolService.ExecuteToolCall(ctx, toolCall)
					if err != nil {
						resutltCh <- emitErr(err)
						return
					}
				} else {
					toolResult = &models.ToolResponse{
						ID:      toolCall.ID,
						Content: "User has DENIED the tool call.",
					}
				}

				toolMessage := &models.Message{
					ParentID:   &headID,
					Role:       "tool",
					ToolCallID: &toolResult.ID,
					Content:    toolResult.Content,
				}

				err = cs.messageService.SaveMessage(ctx, currentBranch.ID, toolMessage)
				if err != nil {
					resutltCh <- emitErr(err)
					return
				}
				headID = toolMessage.ID
			}

			err = cs.branchService.UpdateBranchHead(ctx, currentBranch.ID, headID)
			if err != nil {
				resutltCh <- emitErr(err)
				return
			}

			if !model.SupportsTools || len(toolCalls) == 0 {
				resutltCh <- models.StreamEvent{Type: models.EventLoopComplete}
				return
			}
		}
	}()

	return &models.AgentSession{
		Events:    resutltCh,
		Approvals: approveCh,
		Workspace: currentWorkspace.Name,
		Branch:    currentBranch.Name,
		Model:     *model,
	}, nil
}
