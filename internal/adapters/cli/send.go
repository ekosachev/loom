package cli

import (
	basicchatui "github.com/ekosachev/loom/internal/adapters/cli/ui/basicChatUi"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func (a *CLIApp) sendCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   `send "message"`,
		Short: "Request completion for current thread",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var msg string
			var err error
			if len(args) < 1 {
				msg, err = pterm.DefaultInteractiveTextInput.WithMultiLine(true).Show()
				if err != nil {
					return err
				}
			} else {
				msg = args[0]
			}

			agentSession, err := a.chatService.ExecuteChat(ctx, msg)
			if err != nil {
				return err
			}

			return basicchatui.RunUI(ctx, *agentSession)
		},
	}

	sendCmd.Flags().Bool("ui", false, "Use fancy full-screen ui to display the response")

	return sendCmd
}
