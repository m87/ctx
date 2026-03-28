package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewDeleteIntervalCmd(manager *core.ContextManager) *cobra.Command {
	var intervalID string

	cmd := &cobra.Command{
		Use:   "interval",
		Short: "Delete an interval by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := strings.TrimSpace(intervalID)
			if id == "" {
				return fmt.Errorf("id is required")
			}

			if resolveRemoteAddr() != "" {
				if err := remoteDeleteInterval(id); err != nil {
					return err
				}
			} else {
				if err := manager.IntervalRepository.Delete(id); err != nil {
					return err
				}
			}

			return printOutput(cmd, map[string]string{"id": id, "status": "deleted"}, func() string {
				return "Interval deleted successfully"
			}, nil)
		},
	}

	cmd.Flags().StringVarP(&intervalID, "id", "i", "", "ID of the interval to delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func init() {
	deleteCmd.AddCommand(NewDeleteIntervalCmd(bootstrap.CreateManager()))
}
