package cli

import (
	"fmt"

	"github.com/ekosachev/loom/internal/domain/models"
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

			agentSession, err := a.chatService.ExecuteChat(ctx, msg)
			if err != nil {
				return err
			}

			for {
				event, ok := <-agentSession.Events
				if !ok {
					break
				}

				switch event.Type {
				case models.EventText:
					fmt.Print(event.Text)
				case models.EventDone:
					fmt.Println()
				case models.EventError:
					fmt.Printf("Error! %s\n", event.Err.Error())
				case models.EventToolCall:
					fmt.Printf("tool %s called\n", event.ToolCall.Name)
				case models.EventToolCallRequest:
					request := event.ToolCall
					result, _ := pterm.DefaultInteractiveConfirm.Show(fmt.Sprintf("Model wants to run tool %s with arguments %s. Approve?", request.Name, request.Arguments))
					agentSession.Approvals <- models.ApprovalResponse{
						ID:       request.ID,
						Approved: result,
					}
				}
			}

			return nil
		},
	}
}
