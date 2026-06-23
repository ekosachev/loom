package models

type StreamEventType int

const (
	EventText StreamEventType = iota
	EventToolCall
	EventToolCallRequest
	EventUsageInfo
	EventError
	EventDone
	EventLoopComplete
)

type UsageInfo struct {
	PromptTokens     int
	CompletionTokens int
	Cost             float64
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
	Usage    *UsageInfo
	Err      error
}

type AgentSession struct {
	Events    <-chan StreamEvent
	Approvals chan<- ApprovalResponse
	Workspace string
	Branch    string
	Model     Model
}
