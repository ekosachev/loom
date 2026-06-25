package llm

import (
	"github.com/ekosachev/loom/internal/domain/models"
)

type openRouterMessage struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id" binding:"omitempty"`
}

type openRouterRequestDTO struct {
	Model             string              `json:"model"`
	Messages          []openRouterMessage `json:"messages"`
	Stream            bool                `json:"stream"`
	ToolChoice        string              `json:"tool_choice"`
	Tools             []requestToolDTO    `json:"tools"`
	ParallelToolCalls bool                `json:"parallel_tool_calls"`
}

type requestToolDTO struct {
	Type     string          `json:"type"`
	Function toolFunctionDTO `json:"function"`
}

type toolFunctionDTO struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  toolParametersDTO `json:"parameters"`
}

type toolParametersDTO struct {
	Type       string                 `json:"type"`
	Properties map[string]toolPropDTO `json:"properties"`
	Required   []string               `json:"required"`
}

type toolPropDTO struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     any    `json:"default"`
}

type toolCallDTO struct {
	Type     string `json:"type"`
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func (c *toolCallDTO) combineChunks(other toolCallDTO) {
	if other.Type != "" {
		c.Type = other.Type
	}

	if other.ID != "" {
		c.ID = other.ID
	}

	if other.Function.Name != "" {
		c.Function.Name = other.Function.Name
	}

	if other.Function.Arguments != "" {
		c.Function.Arguments = c.Function.Arguments + other.Function.Arguments
	}
}

func preprocessTool(tool models.Tool) requestToolDTO {
	preprocessedProps := make(map[string]toolPropDTO, len(tool.Tool.Parameters.Properties))
	for i, p := range tool.Tool.Parameters.Properties {
		preprocessedProps[i] = toolPropDTO{
			Type:        p.Type,
			Default:     p.Default,
			Description: p.Description,
		}
	}

	return requestToolDTO{
		Type: string(tool.Tool.Type),
		Function: toolFunctionDTO{
			Name:        tool.Meta.Name,
			Description: tool.Meta.Description,
			Parameters: toolParametersDTO{
				Type:       "object",
				Properties: preprocessedProps,
				Required:   tool.Tool.Parameters.Required,
			},
		},
	}
}
