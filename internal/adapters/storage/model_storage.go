package storage

import (
	"errors"
	"os"

	"github.com/ekosachev/loom/internal/domain/models"
	"go.yaml.in/yaml/v3"
)

type ModelStorage struct {
	filePath string
}

func NewModelStorage(filePath string) (*ModelStorage, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.Create(filePath)
		} else {
			return nil, err
		}
	}

	return &ModelStorage{
		filePath: filePath,
	}, nil
}

func (s *ModelStorage) LoadModels() (map[string]models.Model, error) {
	bytes, err := os.ReadFile(s.filePath)
	if err != nil {
		return map[string]models.Model{}, err
	}

	var modelsMap map[string]models.Model

	err = yaml.Unmarshal(bytes, &modelsMap)
	if err != nil {
		return map[string]models.Model{}, err
	}

	if modelsMap == nil {
		modelsMap = map[string]models.Model{}
	}

	return modelsMap, nil
}

func (s *ModelStorage) SaveModels(models map[string]models.Model) error {
	bytes, err := yaml.Marshal(models)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, bytes, 0644)
}
