package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewFreeCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "free",
		Short: "Free active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			if resolveRemoteAddr() != "" {
				if err := remoteFreeContext(); err != nil {
					return err
				}
			} else {
				if err := manager.FreeActiveContext(); err != nil {
					return err
				}
			}

			return printOutput(cmd, map[string]string{"status": "freed"}, func() string {
				return "Active context freed"
			}, nil)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(NewFreeCmd(bootstrap.CreateManager()))
}
