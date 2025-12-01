package cmd

import "github.com/spf13/cobra"

func newCommentCmd() *cobra.Command {
	return &cobra.Command{
		Use: "comment",
	}
}

var commentCmd = newCommentCmd()

func init() {
	rootCmd.AddCommand(commentCmd)
}
