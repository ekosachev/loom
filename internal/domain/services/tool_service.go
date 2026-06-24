package services

import (
	"context"

	"github.com/ekosachev/loom/internal/adapters/executors"
	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type ToolService struct {
	toolStorage ports.ToolStoragePort
	executors   map[models.RuntimeType]ports.ToolExecutor
}

func NewToolService(toolStorage ports.ToolStoragePort) *ToolService {
	return &ToolService{
		toolStorage: toolStorage,
		executors: map[models.RuntimeType]ports.ToolExecutor{
			models.StarlarkRuntime:    executors.NewStarlarkExecutor(),
			models.InteractionRuntime: executors.NewInteractionExecutor(),
		},
	}
}

func (s *ToolService) ListTools() ([]models.Tool, error) {
	return s.toolStorage.LoadAllTools()
}

func (s *ToolService) GetToolByName(toolName string) (*models.Tool, error) {
	return s.toolStorage.GetToolByName(toolName)
}

func (s *ToolService) ExecuteToolCall(ctx context.Context, toolCall models.ToolCall) (*models.ToolResponse, error) {
	tool, err := s.GetToolByName(toolCall.Name)
	if err != nil {
		return nil, err
	}

	arguments := toolCall.Arguments
	for prop, propConfig := range tool.Tool.Parameters.Properties {
		if _, ok := arguments[prop]; !ok {
			arguments[prop] = propConfig.Default
		}
	}

	executionRequest := models.ToolExecutionRequest{
		ToolRoot:  s.toolStorage.GetToolsRoot() + "\\" + toolCall.Name,
		Runtime:   tool.Runtime,
		Arguments: arguments,
	}
	toolResult, err := s.executors[executionRequest.Runtime.Type].ExecuteTool(executionRequest)
	if toolResult == nil {
		toolResult = &models.ToolResponse{}
	}
	if err != nil {
		toolResult.Content = err.Error()
	}
	toolResult.ID = toolCall.ID
	return toolResult, nil
}
