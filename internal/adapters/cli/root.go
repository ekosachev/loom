package cli

import (
	"fmt"
	"os"

	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/domain/services"
	"github.com/ekosachev/loom/internal/ports"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CLIApp struct {
	rootCmd          *cobra.Command
	chatService      ports.ChatServicePort
	workspaceService ports.WorkspaceServicePort
	branchService    ports.BranchServicePort
	messageService   ports.MessageServicePort
	modelService     ports.ModelServicePort
	toolService      ports.ToolServicePort
	config           *models.Config
}

func NewCLIApp(
	cs *services.ChatService,
	workspaceService ports.WorkspaceServicePort,
	branchService ports.BranchServicePort,
	messageService ports.MessageServicePort,
	modelService ports.ModelServicePort,
	toolService ports.ToolServicePort,
	cfg *models.Config,
) *CLIApp {
	app := &CLIApp{
		rootCmd: &cobra.Command{
			Use:   "loom",
			Short: "Loom - a CLI tool to manage yor LLM chats",
		},
		chatService:      cs,
		workspaceService: workspaceService,
		branchService:    branchService,
		messageService:   messageService,
		modelService:     modelService,
		toolService:      toolService,
		config:           cfg,
	}

	app.rootCmd.AddCommand(app.initWorkspaceCmd())
	app.rootCmd.AddCommand(app.statusCmd())
	app.rootCmd.AddCommand(app.sendCmd())
	app.rootCmd.AddCommand(app.logCmd())
	app.rootCmd.AddCommand(app.checkoutCmd())
	app.rootCmd.AddCommand(app.initModelCmd())
	app.rootCmd.AddCommand(app.initBranchCmd())
	app.rootCmd.AddCommand(app.initToolCmd())

	return app
}

func (app *CLIApp) Execute() {
	if err := app.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (app *CLIApp) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display current workspace and branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			activeWs, err := app.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active workspace: %w", err)
			}

			if activeWs == nil {
				pterm.Warning.Println("No active workspace. Use `loom workspace init [name]` to create one")
				return nil
			}

			activeBranch, err := app.branchService.GetActiveBranch(ctx, activeWs.Name)
			if err != nil {
				return fmt.Errorf("failed to get active branch for workspace %s: %w", activeWs.Name, err)
			}

			if activeBranch == nil {
				return fmt.Errorf("no active branch for workspace %s", activeWs.Name)
			}

			pterm.Println("Workspace: " + activeWs.Name)
			pterm.Println("Branch: " + activeBranch.Name)

			return nil
		},
	}
}
