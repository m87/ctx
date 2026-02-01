/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewEditContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		id   string
		name string
	)

	cmd := &cobra.Command{
		Use:   "editContext",
		Short: "Edit an existing context",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("editContext called")
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the context to edit")
	cmd.Flags().StringVarP(&name, "name", "n", "", "New name for the context")

	return cmd
}

func init() {
	editCmd.AddCommand(NewEditContextCmd(bootstrap.CreateManager()))
}
