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

func (s *ToolService) ListTools() []models.Tool {
	return s.toolStorage.LoadAllTools()
}
