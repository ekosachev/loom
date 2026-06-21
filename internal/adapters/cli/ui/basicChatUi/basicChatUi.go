package basicchatui

import (
	"context"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/ekosachev/loom/internal/domain/models"
)

type sessionState int

const (
	stateWaiting sessionState = iota
	stateStreaming
	stateToolCallConfirmation
	stateDone
)

type model struct {
	state           sessionState
	session         models.AgentSession
	ttft            time.Duration
	startTime       time.Time
	usageInfo       *models.UsageInfo
	toolCallRequest *models.ToolCall

	accumulatedText  string
	renderedMarkdown string
	messageBlocks    []string
	currentEvent     *models.StreamEvent

	spinner spinner.Model
	width   int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForEvent(m.session.Events),
	)
}

func RunUI(ctx context.Context, session models.AgentSession) error {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := tea.NewProgram(model{
		spinner:   s,
		session:   session,
		state:     stateWaiting,
		startTime: time.Now(),
	})
	_, err := p.Run()
	return err
}
