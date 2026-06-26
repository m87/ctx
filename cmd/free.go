package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/spf13/cobra"
)

func NewFreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "free",
		Short: "Free active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

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
	rootCmd.AddCommand(NewFreeCmd())
}
