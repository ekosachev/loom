package cli

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) initToolCmd() *cobra.Command {
	toolCmd := &cobra.Command{
		Use:   "tool",
		Short: "Manage available tools",
	}

	toolCmd.AddCommand(a.initToolListCmd())

	return toolCmd
}

func (a *CLIApp) initToolListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Displays a list of available tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools, err := a.toolService.ListTools()
			if err != nil {
				return err
			}

			for _, tool := range tools {
				pterm.Printfln("%s", tool.Meta.Name)
			}

			return nil
		},
	}
}
