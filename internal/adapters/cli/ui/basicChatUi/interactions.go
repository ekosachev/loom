package basicchatui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ekosachev/loom/internal/domain/models"
)

var (
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true)
	fieldNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Italic(true)

	approveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("47"))
	denyTextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	newline = []byte("\n")
)

func (m *model) renderInteraction() string {
	if m.interaction == nil {
		return ""
	}

	i := m.interaction

	var output strings.Builder
	output.Write(newline)

	output.Write([]byte(titleStyle.Inline(true).Render(i.Title)))
	output.Write(newline)

	for _, field := range i.Form.Fields {
		switch field.Widget {
		case models.WidgetText:
			output.Write([]byte(renderTextWidget(field, i.ToolCall.Arguments[field.Name])))
		}
		output.Write(newline)
	}

	output.Write(newline)

	switch i.Kind {
	case models.InteractionApprove:
		output.Write([]byte(renderApprovalInput()))
	}

	output.Write(newline)
	output.Write(newline)

	return output.String()
}

func renderTextWidget(field models.Field, value any) string {
	return fieldNameStyle.Inline(true).Render(field.Name) + "\t" + fmt.Sprintf("%v", value)
}

func renderApprovalInput() string {
	return "Approve? " + approveTextStyle.Inline(true).Render("Yes [y]") + " / " + denyTextStyle.Inline(true).Render("No [n]")
}
