package cli

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) initModelCmd() *cobra.Command {
	modelCmd := &cobra.Command{
		Use:   "model",
		Short: "Manage available models",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			currentModel, err := a.modelService.GetCurrentModel(ctx)

			if err != nil {
				return fmt.Errorf("failed to get current model: %w", err)
			}

			pterm.Printfln("%s (%s/%s)", currentModel.Name, currentModel.Provider, currentModel.Slug)
			return nil
		},
	}

	modelCmd.AddCommand(a.initModelAddCmd())
	modelCmd.AddCommand(a.initModelListCmd())
	modelCmd.AddCommand(a.initModelSetCmd())

	return modelCmd
}

func (a *CLIApp) initModelAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [search-term]",
		Short: "Add a new model from OpenRouter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			searchTerm := args[0]

			err := a.modelService.AddModel(ctx, *a.config, searchTerm)
			if err != nil {
				return fmt.Errorf("failed to add model: %w", err)
			}
			return nil
		},
	}
}

func boolToMark(val bool) string {
	if val {
		return pterm.FgLightGreen.Sprint("✓")
	}
	return pterm.FgLightRed.Sprint("✗")
}

func cursorString(val bool) string {
	if val {
		return "=>"
	}
	return "  "
}

func (a *CLIApp) initModelListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available models",
		RunE: func(cmd *cobra.Command, args []string) error {
			models, err := a.modelService.ListModels()
			ctx := cmd.Context()
			if err != nil {
				return err
			}
			td := pterm.TableData{
				[]string{"", "Name", "OpenRouter ID", "Supports tools"},
			}

			currentModel, _ := a.modelService.GetCurrentModel(ctx)
			var currentModelName string
			if currentModel == nil {
				currentModelName = ""
			} else {
				currentModelName = currentModel.Name
			}

			for _, model := range models {
				td = append(td, []string{
					cursorString(currentModelName == model.Name),
					model.Name,
					fmt.Sprintf("%s/%s", model.Provider, model.Slug),
					boolToMark(model.SupportsTools),
				})
			}

			printer := pterm.DefaultTable.WithData(td).WithHasHeader(true).WithSeparator("│")
			printer.Render()
			return nil
		},
	}
}

func (a *CLIApp) initModelSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [model-name]",
		Short: "Set current model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			modelName := args[0]
			return a.modelService.SetCurrentModel(ctx, modelName)
		},
	}
}
