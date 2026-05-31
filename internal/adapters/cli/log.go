package cli

import (
	"fmt"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type logItem struct {
	msg      models.Message
	expanded bool
	lines    []string
	offset   int
	cursor   int
}

type logModel struct {
	branchName string
	wsName     string

	messages []logItem
	cursor   int

	width   int
	height  int
	offsetY int
}

func (a *CLIApp) logCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "Display log of the current branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			activeWs, err := a.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active workspace: %w", err)
			}
			activeBranch, err := a.branchService.GetActiveBranch(ctx, activeWs.Name)
			if err != nil {
				return fmt.Errorf(
					"failed to get active branch for workspace %s: %w",
					activeWs.Name,
					err,
				)
			}

			if activeBranch.CurrentMessageID == nil {
				pterm.Warning.Println("This branch has no messages")
				return nil
			}
			messages, err := a.messageService.GetThreadContext(ctx, *activeBranch.CurrentMessageID)
			if err != nil {
				return fmt.Errorf("failed to get history for branch %s: %w", activeBranch.Name, err)
			}

			if len(messages) == 0 {
				pterm.Warning.Println("This branch has no messages")
				return nil
			}

			slices.Reverse(messages)
			p := tea.NewProgram(initialModel(activeBranch.Name, activeWs.Name, messages))
			_, err = p.Run()
			return err
		},
	}
}

func initialModel(branchName string, wsName string, messages []models.Message) logModel {
	messageItems := make([]logItem, len(messages))
	for i, m := range messages {
		messageItems[i] = logItem{
			msg:      m,
			expanded: false,
			lines:    []string{},
			offset:   0,
		}
	}

	return logModel{
		branchName: branchName,
		wsName:     wsName,
		messages:   messageItems,
		cursor:     0,

		width:   0,
		height:  0,
		offsetY: 0,
	}
}

func (m logModel) Init() tea.Cmd {
	return nil
}

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.messages[m.cursor].expanded {
				m.messages[m.cursor].cursor = max(0, m.messages[m.cursor].cursor-1)
			} else {
				m.cursor = max(0, m.cursor-1)
			}
		case "down", "j":
			if m.messages[m.cursor].expanded {
				m.messages[m.cursor].cursor = min(len(m.messages[m.cursor].lines)-1, m.messages[m.cursor].cursor+1)
			} else {
				m.cursor = min(len(m.messages)-1, m.cursor+1)
			}
		case "enter", "space":
			m.messages[m.cursor].expanded = !m.messages[m.cursor].expanded
			if m.messages[m.cursor].expanded {
				m.messages[m.cursor].lines = renderMarkdown(m.messages[m.cursor].msg.Content, m.width-4-11)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		for i := range m.messages {
			if !m.messages[i].expanded {
				continue
			}

			m.messages[i].lines = renderMarkdown(m.messages[i].msg.Content, m.width-4-11)
			m.messages[i].cursor = min(len(m.messages[i].lines)-1, max(0, m.messages[i].cursor))
		}
	}

	return m, nil
}

func (m logModel) View() tea.View {
	headerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(1).
		Background(lipgloss.Color("#104e64")).
		Foreground(lipgloss.Color("#cefafe")).
		AlignHorizontal(lipgloss.Center)

	footerStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(1).
		Background(lipgloss.Color("#020618")).
		Foreground(lipgloss.Color("#45556c")).
		AlignHorizontal(lipgloss.Left)

	header := headerStyle.Render(fmt.Sprintf("loom log for %s/%s", m.wsName, m.branchName))
	footer := footerStyle.Render("k / j to go up / down | <Enter> or <Space> to expand message | <Ctrl-C> or q to quit")

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, m.renderMessages(), footer))

	v.AltScreen = true
	return v
}

func (m logModel) renderMessages() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	availableHeight := m.height - 2

	type displayLine struct {
		text     string
		isCursor bool
		role     string
		expanded bool
	}

	var allLines []displayLine

	for i, it := range m.messages {
		if !it.expanded {
			text := truncate(getFirstLine(it.msg.Content), m.width-4-11)
			allLines = append(allLines, displayLine{
				text:     text,
				isCursor: i == m.cursor,
				role:     it.msg.Role,
				expanded: false,
			})
		} else {
			if len(it.lines) == 0 {
				allLines = append(allLines, displayLine{
					text:     "*(empty message. how did you get here?)*",
					isCursor: i == m.cursor,
					role:     it.msg.Role,
					expanded: true,
				})
			}

			for j, l := range it.lines {
				allLines = append(allLines, displayLine{
					text:     l,
					isCursor: i == m.cursor && j == m.messages[m.cursor].cursor,
					role:     it.msg.Role,
					expanded: true,
				})
			}
		}
	}

	cursorGlobalIndex := 0
	for i, l := range allLines {
		if l.isCursor {
			cursorGlobalIndex = i
			break
		}
	}

	if cursorGlobalIndex < m.offsetY {
		m.offsetY = cursorGlobalIndex
	} else if cursorGlobalIndex >= m.offsetY+availableHeight {
		m.offsetY = cursorGlobalIndex - availableHeight + 1
	}

	var visibleLines []string
	end := min(m.offsetY+availableHeight, len(allLines))

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	roleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f1f5f9")).Bold(true).AlignHorizontal(lipgloss.Center).Width(11)

	for _, l := range allLines[m.offsetY:end] {
		switch l.role {
		case "user":
			roleStyle = roleStyle.Background(lipgloss.Color("#024a70"))
		case "assistant":
			roleStyle = roleStyle.Background(lipgloss.Color("#861043"))
		}
		var roleText string

		if l.isCursor || !l.expanded {
			roleText = roleStyle.Render(l.role)
		} else {
			roleText = roleStyle.Render("")
		}

		if l.isCursor {
			visibleLines = append(visibleLines, cursorStyle.Render("> ")+roleText+l.text)
		} else {
			visibleLines = append(visibleLines, "  "+roleText+l.text)
		}
	}

	for len(visibleLines) < availableHeight {
		visibleLines = append(visibleLines, "")
	}

	content := strings.Join(visibleLines, "\n")

	return content
}

func getFirstLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ↵ ")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		if max > 3 {
			return string(runes[:max-3]) + "..."
		}
		return string(runes[:max])
	}

	return string(runes)
}

func renderMarkdown(content string, width int) []string {
	wrapWidth := max(width-4, 10)

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(wrapWidth),
	)

	out, _ := r.Render(content)

	out = strings.TrimSuffix(out, "\n")
	return strings.Split(out, "\n")
}
