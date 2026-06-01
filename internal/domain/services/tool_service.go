package services

import (
	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type ToolService struct {
	toolStorage ports.ToolStoragePort
}

func NewToolService(toolStorage ports.ToolStoragePort) *ToolService {
	return &ToolService{toolStorage: toolStorage}
}

func (s *ToolService) ListTools() ([]models.Tool, error) {
	return s.toolStorage.LoadAllTools()
}

func (s *ToolService) GetToolByName(toolName string) (*models.Tool, error) {
	return s.toolStorage.GetToolByName(toolName)
}
