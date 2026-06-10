package cli

import (
	"fmt"
	"slices"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) initToolCmd() *cobra.Command {
	toolCmd := &cobra.Command{
		Use:   "tool",
		Short: "Manage available tools",
	}

	toolCmd.AddCommand(a.initToolListCmd())
	toolCmd.AddCommand(a.initToolInfoCmd())

	return toolCmd
}

func (a *CLIApp) initToolListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Displays a list of available tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools, err := a.toolService.ListTools()
			if err != nil {
				return fmt.Errorf("failed to load tools: %w", err)
			}

			for _, tool := range tools {
				pterm.Printfln("%s", tool.Meta.Name)
			}

			return nil
		},
	}
}

func (a *CLIApp) initToolInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [tool-name]",
		Short: "Display information about a tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]
			tool, err := a.toolService.GetToolByName(toolName)
			if err != nil {
				return fmt.Errorf("failed to get tool %s: %w", toolName, err)
			}

			pterm.FgCyan.Printfln("Tool: %s", tool.Meta.Name)
			pterm.Println(tool.Meta.Description)

			if len(tool.Tool.Parameters.Properties) == 0 {
				return nil
			}

			for paramName, param := range tool.Tool.Parameters.Properties {
				var requiredStr string
				if slices.Contains(tool.Tool.Parameters.Required[:], paramName) {
					requiredStr = pterm.FgLightRed.Sprint(" * ")
				} else {
					requiredStr = "   "
				}
				pterm.Printfln("%s %s (%s)", requiredStr, paramName, param.Type)
			}

			return nil
		},
	}
}
