package executors

import (
	"fmt"
	"os"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/vladimirvivien/startype"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type StarlarkExecutor struct{}

func NewStarlarkExecutor() *StarlarkExecutor {
	return &StarlarkExecutor{}
}

func (e *StarlarkExecutor) ExecuteTool(toolCall models.ToolExecutionRequest) (*models.ToolResponse, error) {
	config := toolCall.Runtime.StarlarkConfig

	scriptPath := toolCall.ToolRoot + "\\" + config.File

	predeclared := starlark.StringDict{}
	err := e.loadRequirements(predeclared, config)
	if err != nil {
		return nil, err
	}

	switch config.Parameters {
	case models.StarlarkParamsGlobals:
		err := e.loadParamsAsGlobals(predeclared, toolCall.Arguments)
		if err != nil {
			return nil, err
		}
	}

	globals, err := starlark.ExecFileOptions(&syntax.FileOptions{}, &starlark.Thread{}, scriptPath, nil, predeclared)
	if err != nil {
		return nil, err
	}

	var result string
	switch config.Return.Type {
	case models.StarlarkReturnGlobals:
		resultStarlark, ok := globals[config.Return.Name]
		if !ok {
			return nil, fmt.Errorf("Tool executor error: tool did not provide the result")
		}
		result = resultStarlark.String()
	}

	return &models.ToolResponse{
		Content: result,
	}, nil
}

func (e *StarlarkExecutor) loadRequirements(predeclared starlark.StringDict, config models.StarlarkRuntimeConfig) error {
	for _, req := range config.Requirements {
		switch req {
		case models.StarlarkReadFile:
			predeclared["READ_FILE"] = e.readFilePredeclare()
		case models.StarlarkPwd:
			predeclared["PWD"] = e.pwdPredeclare()
		}
	}

	return nil
}

func (e *StarlarkExecutor) pwdPredeclare() *starlark.Builtin {
	return starlark.NewBuiltin(
		"PWD",
		func(
			thread *starlark.Thread,
			fn *starlark.Builtin,
			args starlark.Tuple,
			kwargs []starlark.Tuple,
		) (starlark.Value, error) {
			pwd, err := os.Getwd()
			if err != nil {
				return starlark.None, err
			}
			return starlark.String(pwd), nil
		},
	)
}

func (e *StarlarkExecutor) readFilePredeclare() *starlark.Builtin {
	return starlark.NewBuiltin(
		"READ_FILE",
		func(
			thread *starlark.Thread,
			fn *starlark.Builtin,
			args starlark.Tuple,
			kwargs []starlark.Tuple,
		) (starlark.Value, error) {
			var filename string
			if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &filename); err != nil {
				return starlark.None, err
			}

			contents, err := os.ReadFile(filename)
			if err != nil {
				return starlark.None, err
			}

			return starlark.String(contents), nil
		},
	)
}

func (e *StarlarkExecutor) loadParamsAsGlobals(predeclared starlark.StringDict, params map[string]any) error {
	for key, value := range params {
		starlarkVal, err := startype.Go(value).ToStarlarkValue()
		if err != nil {
			return err
		}

		predeclared[key] = starlarkVal
	}
	return nil
}
