package basicchatui

import (
	"context"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/ekosachev/loom/internal/domain/models"
)

type sessionState int

const (
	stateWaiting sessionState = iota
	stateStreaming
	stateDone
)

type model struct {
	spinner          spinner.Model
	session          models.AgentSession
	accumulatedText  string
	renderedMarkdown string
	state            sessionState
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
		spinner: s,
		session: session,
		state:   stateWaiting,
	})
	_, err := p.Run()
	return err
}
