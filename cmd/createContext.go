package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewCreateContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		name string
	)
	createContextCmd := &cobra.Command{
		Use:   "context",
		Short: "Create a new context",
		Run: func(cmd *cobra.Command, args []string) {
			context := &core.Context{
				Name: name,
			}
			_, err := manager.ContextRepository.Save(context)
			if err != nil {
				cmd.PrintErrln("Error creating context:", err)
				return
			}
			cmd.Println("Context created successfully")
		},
	}
	createContextCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the context")
	createContextCmd.MarkFlagRequired("name")
	return createContextCmd
}

func init() {
	createCmd.AddCommand(NewCreateContextCmd(bootstrap.CreateManager()))
}
