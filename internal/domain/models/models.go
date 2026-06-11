package models

import (
	"time"
)

type Message struct {
	ID         int64  `json:"id"`
	ParentID   *int64 `json:"parent_id"`
	ToolCallID *string
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
}

type Branch struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	CurrentMessageID *int64 `json:"current_message_id"`
	WorkspaceID      string `json:"workspace_id"`
}

type Workspace struct {
	Name string `json:"name"`
}

type Model struct {
	Name          string
	Provider      string `yaml:"provider"`
	Slug          string `yaml:"slug"`
	SupportsTools bool   `yaml:"supports_tools"`
}

type CompletionRequest struct {
	ThreadHistory []Message
	Model         Model
	Tools         []Tool
}

type StreamEvent struct {
	Type StreamEventType

	Text     string
	ToolCall *ToolCall
	Err      error
}

type StreamEventType int

const (
	EventText StreamEventType = iota
	EventToolCall
	EventToolCallRequest
	EventError
	EventDone
	EventLoopComplete
)

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type ParametersType string

const (
	PropParams ParametersType = "properties"
)

type RuntimeType string

const (
	StarlarkRuntime RuntimeType = "starlark"
)

type StarlarkRequirement string

const (
	StarlarkReadFile StarlarkRequirement = "read_file"
	StarlarkPwd      StarlarkRequirement = "pwd"
)

type StarlarkParameters string

const (
	StarlarkParamsGlobals StarlarkParameters = "globals"
)

type StarlarkReturnType string

const (
	StarlarkReturnGlobals StarlarkReturnType = "globals"
)

type Tool struct {
	Meta struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"meta"`
	Tool struct {
		Type       ToolType `yaml:"type"`
		Parameters struct {
			Type       ParametersType `yaml:"type"`
			Properties map[string]struct {
				Type        string `yaml:"type"`
				Description string `yaml:"description"`
				Default     any    `yaml:"default" binding:"omitempty"`
			} `yaml:"properties"`
			Required []string `yaml:"required"`
		} `yaml:"parameters"`
	} `yaml:"tool"`
	Runtime ToolRuntime `yaml:"runtime"`
}

type ToolRuntime struct {
	Type           RuntimeType           `yaml:"type"`
	StarlarkConfig StarlarkRuntimeConfig `yaml:"starlark"`
}

type StarlarkRuntimeConfig struct {
	File       string             `yaml:"file"`
	Parameters StarlarkParameters `yaml:"parameters"`
	Return     struct {
		Type StarlarkReturnType `yaml:"type"`
		Name string             `yaml:"name"`
	} `yaml:"return"`
	Requirements []StarlarkRequirement `yaml:"requirements"`
}

type ToolCall struct {
	Name      string
	ID        string
	Arguments string
}

type ToolExecutionRequest struct {
	ToolRoot  string
	Runtime   ToolRuntime
	Arguments map[string]any
}

type ToolResponse struct {
	ID      string
	Content string
}

type Config struct {
	Openrouter struct {
		Key string `yaml:"key"`
	} `yaml:"openrouter"`

	Models struct {
		Default string `yaml:"default"`
	} `yaml:"models"`
}

type ApprovalResponse struct {
	ID       string
	Approved bool
}

type AgentSession struct {
	Events    <-chan StreamEvent
	Approvals chan<- ApprovalResponse
}
