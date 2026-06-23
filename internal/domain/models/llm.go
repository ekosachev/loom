package models

type StreamEventType int

const (
	EventText StreamEventType = iota
	EventToolCall
	EventInteractionRequired
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

	Text        string
	ToolCall    *ToolCall
	Interaction *Interaction
	Usage       *UsageInfo
	Err         error
}

type AgentSession struct {
	Events    <-chan StreamEvent
	Approvals chan<- InteractionApproval
	Workspace string
	Branch    string
	Model     Model
}
