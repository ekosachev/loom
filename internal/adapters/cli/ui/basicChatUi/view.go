package basicchatui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	header := m.renderHeader()
	messageChunks := m.renderMessageChunks()
	footer := m.renderFooter()
	var body string

	switch m.state {
	case stateStreaming:
		body = m.accumulatedText
	case stateDone:
		body = m.renderedMarkdown
	case stateInteraction:
		body = m.renderInteraction()
	default:
		body = ""
	}

	return tea.NewView(fmt.Sprintf("\n%s\n", strings.Join([]string{
		header,
		messageChunks,
		body,
		footer,
	}, "\n")))
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
	case stateInteraction:
		spinnerLabel = "Waiting for interaction..."
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
	model := m.session.Model.Name
	branch := fmt.Sprintf("%s/%s", m.session.Workspace, m.session.Branch)

	usage := "? => ? tk @ $?"
	if usageInfo := m.usageInfo; usageInfo != nil {
		usage = fmt.Sprintf("%d => %d tk @ $%.3f", usageInfo.PromptTokens, usageInfo.CompletionTokens, usageInfo.Cost)
	}

	lineLength := m.width - 9 - len([]rune(branch)) - len([]rune(model)) - len([]rune(usage))

	return fmt.Sprintf("  %s @ %s %s %s  ", model, branch, strings.Repeat("─", max(lineLength, 0)), usage)
}

func (m model) renderMessageChunks() string {
	return strings.Join(m.messageBlocks, "\n")
}
