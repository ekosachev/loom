package cli

import (
	"fmt"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/spf13/cobra"
)

func (a *CLIApp) checkoutCmd() *cobra.Command {
	var createBranch bool

	cmd := &cobra.Command{
		Use:   "checkout [-b] [branch]",
		Short: "Switch to another branch in current workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			branchName := args[0]

			workspace, err := a.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get ative context: %w", err)
			}

			var branch *models.Branch
			if !createBranch {
				branch, err = a.branchService.GetBranchByName(ctx, workspace.Name, branchName)
				if err != nil {
					return fmt.Errorf("failed to get branch %s: %w", branchName, err)
				}
			} else {
				currentBranch, err := a.branchService.GetActiveBranch(ctx, workspace.Name)
				if err != nil {
					return fmt.Errorf("failed to get current branch for workspace %s: %w", workspace.Name, err)
				}

				branch = &models.Branch{
					Name:             branchName,
					CurrentMessageID: currentBranch.CurrentMessageID,
					WorkspaceID:      currentBranch.WorkspaceID,
				}
				err = a.branchService.SaveBranch(ctx, branch)
				if err != nil {
					return fmt.Errorf("failed to create branch: %w", err)
				}
			}

			a.branchService.ActivateBranch(ctx, workspace.Name, branch.ID)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&createBranch, "create-branch", "b", false, "Create a new branch")

	return cmd
}
