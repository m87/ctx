package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/server"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print ctx version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version := server.CurrentVersion()
			return printOutput(
				cmd,
				version,
				func() string {
					parts := []string{version.Release}
					if version.Commit != "" {
						parts = append(parts, version.Commit)
					}
					if version.Date != "" {
						parts = append(parts, version.Date)
					}
					return strings.Join(parts, " ")
				},
				func() string {
					return fmt.Sprintf(
						"VERSION=%q\nRELEASE=%q\nCOMMIT=%q\nDATE=%q\n",
						version.Version,
						version.Release,
						version.Commit,
						version.Date,
					)
				},
			)
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}
