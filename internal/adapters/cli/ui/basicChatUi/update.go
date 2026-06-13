package basicchatui

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"github.com/ekosachev/loom/internal/domain/models"
)

func waitForEvent(ch <-chan models.StreamEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-ch
		if !ok {
			return models.StreamEvent{Type: models.EventLoopComplete}
		}
		return event
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Ctrl+c":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case models.StreamEvent:
		switch msg.Type {
		case models.EventText:
			if m.state == stateWaiting {
				m.state = stateStreaming
			}
			m.accumulatedText += msg.Text

			return m, waitForEvent(m.session.Events)
		case models.EventLoopComplete:
			m.state = stateDone
			rendered, err := glamour.Render(m.accumulatedText, "dark")
			if err != nil {
				m.renderedMarkdown = m.accumulatedText
			} else {
				m.renderedMarkdown = rendered
			}

			return m, tea.Quit
		}
	}
	return m, nil
}
