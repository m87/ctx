package cmd

import (
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive active contexts",
	Long: `Archive single context: ctx archive test
	Archvie all active contexts: 
		ctx archive --all
		ctx archive -a

	`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolP("all", "a", false, "Archive all active contexts")
}
