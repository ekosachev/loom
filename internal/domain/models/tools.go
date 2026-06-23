package models

type (
	ToolType       string
	ParametersType string
	RuntimeType    string

	StarlarkReturnType  string
	StarlarkRequirement string
	StarlarkParameters  string

	InteractionBodyElementType string
	InteractionActionType      string
)

const (
	ToolTypeFunction ToolType = "function"
)

const (
	PropParams ParametersType = "properties"
)

const (
	StarlarkRuntime RuntimeType = "starlark"
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
	InteractionBodyText InteractionBodyElementType = "text"
)

const (
	InteractionActionApprove InteractionActionType = "approve"
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
	Runtime ToolRuntime `yaml:"runtime"`
}

type ToolRuntime struct {
	Type           RuntimeType           `yaml:"type"`
	StarlarkConfig StarlarkRuntimeConfig `yaml:"starlark"`
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

type ToolInteraction struct {
	Idempotent bool                     `yaml:"idempotent"`
	Title      string                   `yaml:"title"`
	Body       []InteractionBodyElement `yaml:"body"`
	Action     InteractionActionType    `yaml:"action"`
}

type InteractionBodyElement struct {
	Type InteractionBodyElementType `yaml:"type"`
	Text string                     `yaml:"text"`
}

type ToolCall struct {
	Name      string
	ID        string
	Arguments string
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

type ApprovalResponse struct {
	ID       string
	Approved bool
}
