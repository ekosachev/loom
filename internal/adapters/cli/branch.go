package cli

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) initBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branch",
		Short: "Display a list of branches in the current workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			currentWs, err := a.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get current workspace: %w", err)
			}

			branches, err := a.branchService.GetAllBranches(ctx, currentWs.Name)
			if err != nil {
				return fmt.Errorf("failed to get branches for workspace %s: %w", currentWs.Name, err)
			}

			activeBranch, err := a.branchService.GetActiveBranch(ctx, currentWs.Name)
			if err != nil {
				return fmt.Errorf("failed to get active branch for workspace %s: %w", currentWs.Name, err)
			}

			for _, br := range branches {
				activeIndicator := "  "
				if activeBranch.ID == br.ID {
					activeIndicator = "=>"
				}
				pterm.FgGreen.Print(activeIndicator)
				pterm.Printfln("\t%s", br.Name)
			}

			return nil
		},
	}
}
