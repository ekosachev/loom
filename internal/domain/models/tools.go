package models

type (
	ToolType       string
	ParametersType string
	RuntimeType    string

	StarlarkReturnType  string
	StarlarkRequirement string
	StarlarkParameters  string

	ShellScriptLocation   string
	ShellScriptParameters string

	InteractionKind string

	FieldType  string
	WidgetType string
)

const (
	ToolTypeFunction ToolType = "function"
)

const (
	PropParams ParametersType = "properties"
)

const (
	StarlarkRuntime    RuntimeType = "starlark"
	InteractionRuntime RuntimeType = "interaction"
	ShellRuntime       RuntimeType = "shell"
)

const (
	StarlarkReadFile StarlarkRequirement = "read_file"
	StarlarkPwd      StarlarkRequirement = "pwd"
)

const (
	StarlarkParamsGlobals StarlarkParameters = "globals"
)

const (
	StarlarkReturnGlobals StarlarkReturnType = "globals"
)

const (
	InteractionApprove   InteractionKind = "approve"
	InteractionInputText InteractionKind = "inputText"
)

const (
	ShellScriptInProp ShellScriptLocation = "prop"
)

const (
	ShellParametersAsTemplate ShellScriptParameters = "template"
)

const FieldString FieldType = "string"

const (
	WidgetText      WidgetType = "text"
	WidgetTextInput WidgetType = "textInput"
)

type Tool struct {
	Meta struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"meta"`
	Tool struct {
		Type       ToolType `yaml:"type"`
		Parameters struct {
			Type       ParametersType `yaml:"type"`
			Properties map[string]struct {
				Type        string `yaml:"type"`
				Description string `yaml:"description"`
				Default     any    `yaml:"default" binding:"omitempty"`
			} `yaml:"properties"`
			Required []string `yaml:"required"`
		} `yaml:"parameters"`
	} `yaml:"tool"`
	Runtime     ToolRuntime `yaml:"runtime"`
	Interaction struct {
		Title string          `yaml:"title"`
		Kind  InteractionKind `yaml:"kind"`
		Form  []Field         `yaml:"form"`
		Info  string          `yaml:"info"`
	} `yaml:"interaction"`
}

type ToolRuntime struct {
	Type              RuntimeType              `yaml:"type"`
	StarlarkConfig    StarlarkRuntimeConfig    `yaml:"starlark"`
	InteractionConfig InteractionRuntimeConfig `yaml:"interaction"`
	ShellConfig       ShellRuntimeConfig       `yaml:"shell"`
}

type StarlarkRuntimeConfig struct {
	File       string             `yaml:"file"`
	Parameters StarlarkParameters `yaml:"parameters"`
	Return     struct {
		Type StarlarkReturnType `yaml:"type"`
		Name string             `yaml:"name"`
	} `yaml:"return"`
	Requirements []StarlarkRequirement `yaml:"requirements"`
}

type InteractionRuntimeConfig struct {
	ResultField string `yaml:"result"`
}

type ShellRuntimeConfig struct {
	Location   ShellScriptLocation   `yaml:"location"`
	Parameters ShellScriptParameters `yaml:"parameters"`
	Cmd        string                `yaml:"cmd"`
	Flags      string                `yaml:"flags"`
	Script     string                `yaml:"script"`
}

type Interaction struct {
	ToolCall ToolCall
	Title    string
	Kind     InteractionKind
	Form     JSONSchemaForm
	Info     string
}

type InteractionApproval struct {
	ToolCall ToolCall
	Approved bool
}

type JSONSchemaForm struct {
	Fields []Field
}

type Field struct {
	Name        string     `yaml:"name"`
	Type        FieldType  `yaml:"type"`
	Widget      WidgetType `yaml:"widget"`
	Interaction bool       `yaml:"interaction"`
	Value       any        `yaml:"value"`
}

type ToolCall struct {
	Name      string
	ID        string
	Arguments map[string]any
}

type ToolExecutionRequest struct {
	ToolRoot  string
	Runtime   ToolRuntime
	Arguments map[string]any
}

type ToolResponse struct {
	ID      string
	Content string
}
