package cli

import (
	"fmt"
	"os"
	"strings"

	"charm.land/glamour/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) sendCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `send "message"`,
		Short: "Request completion for current thread",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			msg := args[0]

			activeWs, err := a.workspaceService.GetActiveWorkspace(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active workspace: %w", err)
			}
			activeBranch, err := a.branchService.GetActiveBranch(ctx, activeWs.Name)
			if err != nil {
				return fmt.Errorf("failed to get active branch: %w", err)
			}
			activeModel, err := a.modelService.GetCurrentModel(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active model: %w", err)
			}

			if activeModel == nil {
				return fmt.Errorf("no active model set. Use loom model set [name]")
			}

			var messageBuffer strings.Builder

			pterm.Info.Printfln(
				"Requested completion from %s on branch %s/%s",
				activeModel.Name,
				activeWs.Name,
				activeBranch.Name,
			)

			modelId := fmt.Sprintf("%s/%s", activeModel.Provider, activeModel.Slug)

			err = a.chatService.ExecuteChat(
				ctx,
				activeBranch.ID,
				modelId,
				msg,
				a.config,
				func(s string) {
					messageBuffer.Write([]byte(s))
					fmt.Print(s)
					os.Stdout.Sync()
				},
			)

			if err != nil {
				return err
			}

			fullResponse := messageBuffer.String()

			termWidth, termHeight, _ := pterm.GetTerminalSize()

			linesToClear := min(termHeight-1, calculatePhysicalLines(fullResponse, termWidth))

			for range linesToClear {
				fmt.Print("\033[A\033[2K")
			}

			rendered, err := glamour.Render(fullResponse, "dark")
			if err != nil {
				return err
			}

			fmt.Print(rendered)

			return nil
		},
	}
}

func calculatePhysicalLines(text string, termWidth int) int {
	if text == "" {
		return 0
	}

	lines := strings.Split(text, "\n")
	totalLines := 0

	for _, line := range lines {
		lineLen := len([]rune(line))
		if lineLen == 0 {
			totalLines++
			continue
		}
		totalLines += (lineLen + termWidth - 1) / termWidth
	}

	return totalLines
}
