package storage

import (
	"os"

	"github.com/ekosachev/loom/internal/domain/models"
)

type ToolStorage struct {
	toolsDirPath string
}

func NewToolStorage(toolsDirPath string) (*ToolStorage, error) {
	err := os.MkdirAll(toolsDirPath, 0755)
	if err != nil {
		return nil, err
	}

	return &ToolStorage{toolsDirPath: toolsDirPath}, nil
}

func (s *ToolStorage) LoadAllTools() []models.Tool {
	return []models.Tool{
		{Name: "test tool 1"},
		{Name: "test tool 2"},
		{Name: "test tool 3"},
	}
}
