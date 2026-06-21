package basicchatui

import (
	"time"

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
		case "y", "Y":
			if m.state == stateToolCallConfirmation {
				m.emitToolApproval(true)
				m.state = stateWaiting
			}
		case "n", "N":
			if m.state == stateToolCallConfirmation {
				m.emitToolApproval(false)
				m.state = stateWaiting
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.state == stateWaiting {
			m.ttft = time.Since(m.startTime)
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case models.StreamEvent:
		m.currentEvent = &msg
		switch msg.Type {
		case models.EventText:
			if m.state == stateWaiting {
				m.state = stateStreaming
			}
			m.accumulatedText += msg.Text
		case models.EventUsageInfo:
			m.usageInfo = msg.Usage
		case models.EventError:
			panic(msg.Err)
		case models.EventToolCallRequest:
			m.toolCallRequest = msg.ToolCall
			m.state = stateToolCallConfirmation
			m = m.renderMarkdown()
			m = m.pushChunk()
		case models.EventDone:
			m = m.renderMarkdown()
			m = m.pushChunk()
		case models.EventLoopComplete:
			m.state = stateDone
		}
		return m, waitForEvent(m.session.Events)
	}
	return m, nil
}

func (m model) renderMarkdown() model {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(m.width),
		glamour.WithStandardStyle("dark"),
	)
	if err == nil {
		rendered, err := renderer.Render(m.accumulatedText)
		if err == nil {
			m.renderedMarkdown = rendered
			return m
		}
	}

	m.renderedMarkdown = m.accumulatedText
	return m
}

func (m model) pushChunk() model {
	if len(m.renderedMarkdown) > 0 {
		m.messageBlocks = append(m.messageBlocks, m.renderedMarkdown)
		m.renderedMarkdown = ""
		m.accumulatedText = ""
	}
	return m
}

func (m model) emitToolApproval(isApproved bool) model {
	if m.toolCallRequest == nil {
		return m
	}
	m.session.Approvals <- models.ApprovalResponse{
		Approved: isApproved,
		ID:       m.toolCallRequest.ID,
	}

	return m
}
