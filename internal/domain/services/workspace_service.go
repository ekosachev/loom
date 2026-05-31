package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type WorkspaceService struct {
	workspaceStorage ports.WorkspaceStorage
	branchStorage    ports.BranchStoragePort
	stateStorage     ports.StateStorage
}

func NewWorkspaceService(
	workspaceStorage ports.WorkspaceStorage,
	branchStorage ports.BranchStoragePort,
	stateStorage ports.StateStorage,
) *WorkspaceService {
	return &WorkspaceService{
		workspaceStorage: workspaceStorage,
		branchStorage:    branchStorage,
		stateStorage:     stateStorage,
	}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, ws *models.Workspace) (string, error) {
	err := s.workspaceStorage.SaveWorkspace(ctx, ws)
	if err != nil {
		return "", err
	}

	return ws.Name, nil
}

func (s *WorkspaceService) ActivateWorkspace(ctx context.Context, name string) error {
	existingWs, err := s.workspaceStorage.GetWorkspace(ctx, name)

	if err != nil {
		return err
	}

	if existingWs == nil {
		return fmt.Errorf("workspace %s does not exist", name)
	}

	head, err := s.stateStorage.GetState(ctx, name+"/head")
	if err != nil {
		return err
	}

	if head == nil {
		masterBranch := &models.Branch{
			Name:             "master",
			CurrentMessageID: nil,
			WorkspaceID:      existingWs.Name,
		}

		s.branchStorage.SaveBranch(ctx, masterBranch)
		s.stateStorage.SetState(ctx, name+"/head", strconv.FormatInt(masterBranch.ID, 10))
	}

	return s.stateStorage.SetState(ctx, "active_workspace", name)
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, name string) (*models.Workspace, error) {
	return s.workspaceStorage.GetWorkspace(ctx, name)
}

func (s *WorkspaceService) GetActiveWorkspace(ctx context.Context) (*models.Workspace, error) {
	activeWorkspaceName, err := s.stateStorage.GetState(ctx, "active_workspace")

	if err != nil {
		return nil, err
	}

	if activeWorkspaceName == nil {
		return nil, fmt.Errorf("no active workspace set")
	}

	return s.GetWorkspace(ctx, *activeWorkspaceName)
}

func (s *WorkspaceService) GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error) {
	return s.workspaceStorage.GetAllWorkspaces(ctx)
}
