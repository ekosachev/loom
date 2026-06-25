package executors

import (
	"os/exec"
	"strings"
	"text/template"

	"github.com/ekosachev/loom/internal/domain/models"
)

type ShellExecutor struct{}

func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) ExecuteTool(toolCall models.ToolExecutionRequest) (*models.ToolResponse, error) {
	config := toolCall.Runtime.ShellConfig

	var script string
	switch config.Location {
	case models.ShellScriptInProp:
		script = config.Script
	default:
		script = ""
	}

	switch config.Parameters {
	case models.ShellParametersAsTemplate:
		var err error
		script, err = fillScriptTemplate(script, toolCall.Arguments)
		if err != nil {
			return nil, err
		}
	}

	result, err := exec.Command(config.Cmd, config.Flags, script).Output()
	if err != nil {
		return nil, err
	}

	return &models.ToolResponse{
		Content: string(result),
	}, nil
}

func fillScriptTemplate(script string, args map[string]any) (string, error) {
	tmpl, err := template.New("script").Parse(script)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, args)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
