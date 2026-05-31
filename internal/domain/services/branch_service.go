package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/ports"
)

type BranchService struct {
	branchStorage ports.BranchStoragePort
	stateStorage  ports.StateStorage
}

func NewBranchServcie(branchStorage ports.BranchStoragePort, stateStorage ports.StateStorage) *BranchService {
	return &BranchService{
		branchStorage: branchStorage,
		stateStorage:  stateStorage,
	}
}

func (s *BranchService) GetActiveBranch(ctx context.Context, workspaceName string) (*models.Branch, error) {
	activeBranchId, err := s.stateStorage.GetState(ctx, workspaceName+"/head")

	if err != nil {
		return nil, err
	}

	if activeBranchId == nil {
		return nil, fmt.Errorf("no active branch for workspace %s", workspaceName)
	}

	branchIdNum, err := strconv.ParseInt(*activeBranchId, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.branchStorage.GetBranch(ctx, branchIdNum)
}

func (s *BranchService) GetBranchByName(ctx context.Context, workspaceName string, branchName string) (*models.Branch, error) {
	return s.branchStorage.GetBranchByName(ctx, workspaceName, branchName)
}

func (s *BranchService) SaveBranch(ctx context.Context, branch *models.Branch) error {
	return s.branchStorage.SaveBranch(ctx, branch)
}

func (s *BranchService) ActivateBranch(ctx context.Context, workspaceName string, branchID int64) error {
	return s.stateStorage.SetState(ctx, workspaceName+"/head", strconv.FormatInt(branchID, 10))
}

func (s *BranchService) GetAllBranches(ctx context.Context, workspaceName string) ([]models.Branch, error) {
	return s.branchStorage.GetAllBranches(ctx, workspaceName)
}
