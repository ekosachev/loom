package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type ModelService struct {
	modelStorage ports.ModelStoragePort
	modelRepo    ports.ModelRepositoryPort
	stateStorage ports.StateStorage
}

func NewModelService(storage ports.ModelStoragePort, repo ports.ModelRepositoryPort, stateStorage ports.StateStorage) *ModelService {
	return &ModelService{
		modelStorage: storage,
		modelRepo:    repo,
		stateStorage: stateStorage,
	}
}

func (s *ModelService) AddModel(ctx context.Context, cfg models.Config, searchItem string) error {
	itemParts := strings.Split(searchItem, "/")
	if len(itemParts) != 2 {
		return fmt.Errorf(`%s appears to not be a valid OpenRouter model ID,
							use {series}/{slug} format, e. g.: anthropic/claude-opus-4.8-fast`, searchItem)
	}

	provider := itemParts[0]
	slug := itemParts[1]

	model, err := s.modelRepo.GetMetaForModel(ctx, cfg, provider, slug)
	if err != nil {
		return err
	}

	currentModels, err := s.modelStorage.LoadModels()
	if err != nil {
		return err
	}

	currentModels[model.Name] = *model
	s.modelStorage.SaveModels(currentModels)

	return nil
}

func (s *ModelService) ListModels() (map[string]models.Model, error) {
	return s.modelStorage.LoadModels()
}

func (s *ModelService) GetCurrentModel(ctx context.Context) (*models.Model, error) {
	models, err := s.modelStorage.LoadModels()
	if err != nil {
		return nil, fmt.Errorf("failed to load models: %w", err)
	}

	currentModelName, err := s.stateStorage.GetState(ctx, "model")
	if err != nil {
		return nil, fmt.Errorf("failed to get current model: %w", err)
	}

	if currentModelName == nil {
		return nil, fmt.Errorf("no current model set")
	}
	model := models[*currentModelName]
	return &model, nil
}

func (s *ModelService) SetCurrentModel(ctx context.Context, name string) error {
	return s.stateStorage.SetState(ctx, "model", name)
}
