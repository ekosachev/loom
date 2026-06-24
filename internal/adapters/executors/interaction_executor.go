package executors

import (
	"fmt"

	"github.com/ekosachev/loom/internal/domain/models"
)

type InteractionExecutor struct{}

func NewInteractionExecutor() *InteractionExecutor {
	return &InteractionExecutor{}
}

func (e *InteractionExecutor) ExecuteTool(toolCall models.ToolExecutionRequest) (*models.ToolResponse, error) {
	config := toolCall.Runtime.InteractionConfig

	result, ok := toolCall.Arguments[config.ResultField]
	if !ok {
		result = "Error: Tool did not return the result as expected"
	}

	return &models.ToolResponse{
		Content: fmt.Sprintf("%v", result),
	}, nil
}
