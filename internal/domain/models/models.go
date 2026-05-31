package models

import "time"

type Message struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type Branch struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	CurrentMessageID *int64 `json:"current_message_id"`
	WorkspaceID      string `json:"workspace_id"`
}

type Workspace struct {
	Name string `json:"name"`
}

type Model struct {
	Name          string
	Provider      string `yaml:"provider"`
	Slug          string `yaml:"slug"`
	SupportsTools bool   `yaml:"supports_tools"`
}

type Tool struct {
	Meta struct {
		Name string `yaml:"name"`
	} `yaml:"meta"`
}

type Config struct {
	Openrouter struct {
		Key string `yaml:"key"`
	} `yaml:"openrouter"`

	Models struct {
		Default string `yaml:"default"`
	} `yaml:"models"`
}
