package basicchatui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	switch m.state {
	case stateWaiting:
		header := m.renderHeader()
		return tea.NewView(header)
	case stateStreaming:
		header := m.renderHeader()
		footer := m.renderFooter()
		return tea.NewView(fmt.Sprintf("\n%s\n%s\n%s\n", header, m.accumulatedText, footer))
	case stateDone:
		header := m.renderHeader()
		footer := m.renderFooter()
		return tea.NewView(fmt.Sprintf("\n%s\n%s\n%s\n", header, m.renderedMarkdown, footer))
	default:
		return tea.NewView("")
	}
}

func (m model) renderHeader() string {
	var spinnerLabel string

	switch m.state {
	case stateWaiting:
		spinnerLabel = "Waiting for the response stream to start..."
	case stateStreaming:
		spinnerLabel = "Capturing response stream..."
	case stateDone:
		spinnerLabel = "Streaming complete"
	}

	var spinnerText string
	if m.state == stateDone {
		spinnerText = "✓"
	} else {
		spinnerText = m.spinner.View()
	}

	spinner := fmt.Sprintf("  %s %s ", spinnerText, spinnerLabel)
	ttft := fmt.Sprintf(" ttft: %.3fs  ", m.ttft.Seconds())

	lineLength := m.width - len([]rune(spinner)) - len([]rune(ttft))

	return lipgloss.JoinHorizontal(0, spinner, strings.Repeat("─", max(lineLength, 0)), ttft)
}

func (m model) renderFooter() string {
	return "  " + strings.Repeat("─", max(m.width-4, 0)) + "  "
}
