package ports

import (
	"context"

	"github.com/ekosachev/loom/internal/domain/models"
)

type WorkspaceServicePort interface {
	CreateWorkspace(ctx context.Context, ws *models.Workspace) (string, error)
	ActivateWorkspace(ctx context.Context, name string) error
	GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error)
	GetActiveWorkspace(ctx context.Context) (*models.Workspace, error)
	GetWorkspace(ctx context.Context, name string) (*models.Workspace, error)
}

type BranchServicePort interface {
	GetActiveBranch(ctx context.Context, workspaceName string) (*models.Branch, error)
	GetBranchByName(ctx context.Context, workspaceName string, branchName string) (*models.Branch, error)
	SaveBranch(ctx context.Context, branch *models.Branch) error
	ActivateBranch(ctx context.Context, workspaceName string, branchID int64) error
	GetAllBranches(ctx context.Context, workspaceName string) ([]models.Branch, error)
}

type MessageServicePort interface {
	GetThreadContext(ctx context.Context, headID int64) ([]models.Message, error)
}

type ModelServicePort interface {
	AddModel(ctx context.Context, cfg models.Config, searchItem string) error
	ListModels() (map[string]models.Model, error)
	GetCurrentModel(ctx context.Context) (*models.Model, error)
	SetCurrentModel(ctx context.Context, name string) error
}

type ToolServicePort interface {
	ListTools() ([]models.Tool, error)
	GetToolByName(toolName string) (*models.Tool, error)
}
