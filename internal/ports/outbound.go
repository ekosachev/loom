package ports

import (
	"context"

	"github.com/ekosachev/loom/internal/domain/models"
)

type MessageStoragePort interface {
	SaveMessage(ctx context.Context, msg *models.Message) error
	GetMessage(ctx context.Context, id int64) (*models.Message, error)
	GetThreadContext(ctx context.Context, headMessageID int64) ([]models.Message, error)
}

type BranchStoragePort interface {
	SaveBranch(ctx context.Context, branch *models.Branch) error
	GetBranch(ctx context.Context, branchID int64) (*models.Branch, error)
	GetBranchByName(ctx context.Context, workspaceName string, branchName string) (*models.Branch, error)
	UpdateBranchHead(ctx context.Context, branchID int64, messageID int64) error
	GetAllBranches(ctx context.Context, workspaceName string) ([]models.Branch, error)
}

type WorkspaceStorage interface {
	SaveWorkspace(ctx context.Context, ws *models.Workspace) error
	GetWorkspace(ctx context.Context, name string) (*models.Workspace, error)
	GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error)
}

type StateStorage interface {
	GetState(ctx context.Context, key string) (*string, error)
	SetState(ctx context.Context, key string, value string) error
}

type ModelStoragePort interface {
	LoadModels() (map[string]models.Model, error)
	SaveModels(map[string]models.Model) error
}

type ModelRepositoryPort interface {
	GetMetaForModel(ctx context.Context, cfg models.Config, provider string, slug string) (*models.Model, error)
}

type StoragePort interface {
	MessageStoragePort
	BranchStoragePort
	WorkspaceStorage
	StateStorage
}

type ToolStoragePort interface {
	LoadAllTools() ([]models.Tool, error)
	GetToolByName(toolName string) (*models.Tool, error)
}

type LLMPort interface {
	StreamCompletion(ctx context.Context, cfg *models.Config, modelID string, history []models.Message) (<-chan string, <-chan error, error)
}
