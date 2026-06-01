package storage

import (
	"os"
	"path/filepath"

	"github.com/ekosachev/loom/internal/domain/models"
	"go.yaml.in/yaml/v3"
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

func (s *ToolStorage) LoadAllTools() ([]models.Tool, error) {
	var manifestPaths []string

	entries, err := os.ReadDir(s.toolsDirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(s.toolsDirPath, entry.Name(), "manifest.yaml")
			if _, err := os.Stat(manifestPath); err == nil {
				manifestPaths = append(manifestPaths, manifestPath)
			}
		}
	}

	tools := []models.Tool{}
	for _, path := range manifestPaths {
		tool, err := loadToolManifest(path)
		if err != nil {
			continue
		}
		tools = append(tools, *tool)
	}

	return tools, nil
}

func loadToolManifest(path string) (*models.Tool, error) {
	var tool models.Tool

	fileContents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileContents, &tool)
	if err != nil {
		return nil, err
	}

	return &tool, nil
}

func (s *ToolStorage) GetToolByName(toolName string) (*models.Tool, error) {
	manifestPath := filepath.Join(s.toolsDirPath, toolName, "manifest.yaml")
	return loadToolManifest(manifestPath)
}
