package basicchatui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	switch m.state {
	case stateWaiting:
		return tea.NewView(fmt.Sprintf("%s Waiting for response stream start...", m.spinner.View()))
	case stateStreaming:
		return tea.NewView(m.accumulatedText)
	case stateDone:
		return tea.NewView(m.renderedMarkdown)
	default:
		return tea.NewView("")
	}
}
