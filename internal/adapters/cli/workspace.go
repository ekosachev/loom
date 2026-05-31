package cli

import (
	"fmt"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (app *CLIApp) initWorkspaceCmd() *cobra.Command {
	workspaceCmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage isolated workspaces",
	}

	initCmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsName := args[0]
			ctx := cmd.Context()

			ws := &models.Workspace{Name: wsName}
			if _, err := app.workspaceService.CreateWorkspace(ctx, ws); err != nil {
				return fmt.Errorf("failed to create a workspace: %w", err)
			}

			if err := app.workspaceService.ActivateWorkspace(ctx, wsName); err != nil {
				return fmt.Errorf("failed to activate workspace: %w", err)
			}

			pterm.Info.Printfln("Successfully createad and switched to workspace: %s", wsName)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			workspaces, err := app.workspaceService.GetAllWorkspaces(ctx)
			if err != nil {
				return fmt.Errorf("failed to list workspaces: %w", err)
			}

			activeWorkspace, err := app.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active workspace: %w", err)
			}

			for _, ws := range workspaces {
				activeIndicator := "  "
				if activeWorkspace.Name == ws.Name {
					activeIndicator = "=>"
				}
				pterm.FgGreen.Print(activeIndicator)
				pterm.Printfln("\t%s", ws.Name)
			}
			return nil
		},
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Activate a workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wsName := args[0]

			if err := app.workspaceService.ActivateWorkspace(ctx, wsName); err != nil {
				return fmt.Errorf("failed to activate workspace: %w", err)
			}

			pterm.Info.Printfln("acttivated workspace %s", wsName)
			return nil
		},
	}

	workspaceCmd.AddCommand(initCmd)
	workspaceCmd.AddCommand(listCmd)
	workspaceCmd.AddCommand(setCmd)
	return workspaceCmd
}
