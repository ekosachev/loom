package basicchatui

import (
	"strings"
	"time"
	"unicode/utf8"

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
		switch msg.Key().Keystroke() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m.handleKey(msg)
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
		case models.EventInteractionRequired:
			m.interaction = msg.Interaction
			m.state = stateInteraction
			m.renderMarkdown()
			m.pushChunk()
		case models.EventDone:
			m.renderMarkdown()
			m.pushChunk()
		case models.EventLoopComplete:
			m.state = stateDone
		}
		return m, waitForEvent(m.session.Events)
	}
	return m, nil
}

func (m *model) renderMarkdown() {
	text := strings.Trim(m.accumulatedText, "\n\t\r ")
	if text == "" {
		m.renderedMarkdown = ""
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(m.width),
		glamour.WithStandardStyle("dark"),
	)
	if err == nil {
		rendered, err := renderer.Render(text)
		if err == nil {
			m.renderedMarkdown = rendered
		}
	}

	m.renderedMarkdown = text
}

func (m *model) pushChunk() {
	if len(m.renderedMarkdown) > 0 {
		m.messageBlocks = append(m.messageBlocks, m.renderedMarkdown)
		m.renderedMarkdown = ""
		m.accumulatedText = ""
	}
}

func (m *model) emitToolApproval(isApproved bool) {
	if m.interaction == nil {
		return
	}

	for _, field := range m.interaction.Form.Fields {
		if !field.Interaction {
			continue
		}

		m.interaction.ToolCall.Arguments[field.Name] = field.Value
	}

	m.session.Approvals <- models.InteractionApproval{
		Approved: isApproved,
		ToolCall: m.interaction.ToolCall,
	}

	m.accumulatedText = m.renderPostInteractionInfo()
	m.renderMarkdown()
	m.pushChunk()
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateInteraction && m.interaction != nil {
		switch m.interaction.Kind {
		case models.InteractionInputText:
			for i, field := range m.interaction.Form.Fields {
				if field.Widget == models.WidgetTextInput && field.Interaction {
					if field.Value == nil {
						field.Value = ""
					}

					switch msg.Key().Code {
					case tea.KeyBackspace:
						_, size := utf8.DecodeLastRuneInString(field.Value.(string))
						if len(field.Value.(string))-size > 0 {
							field.Value = field.Value.(string)[:len(field.Value.(string))-size]
						}
					case tea.KeyEnter:
						m.emitToolApproval(true)
						m.state = stateWaiting
					default:
						field.Value = field.Value.(string) + msg.Key().Text
					}

					m.interaction.Form.Fields[i] = field
					break
				}
			}

		case models.InteractionApprove:
			switch msg.String() {
			case "y", "Y":
				if m.state == stateInteraction {
					m.emitToolApproval(true)
					m.state = stateWaiting
				}
			case "n", "N":
				if m.state == stateInteraction {
					m.emitToolApproval(false)
					m.state = stateWaiting
				}
			}
		}
	}

	return m, nil
}
