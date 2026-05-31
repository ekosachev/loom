package services

import "github.com/ekosachev/loom/internal/domain/models"

type ToolService struct{}

func NewToolService() *ToolService {
	return &ToolService{}
}

func (s *ToolService) ListTools() []models.Tool {
	return []models.Tool{
		{Name: "test-tool-1"},
		{Name: "test-tool-2"},
		{Name: "test-tool-3"},
	}
}
