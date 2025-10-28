package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func newMigrateCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate data to a new format or structure",
		Long:  "This command is used to migrate data from the old format to the new format.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Migration process started...")

			registry, err := localstorage.LoadLocalMigrationRegistry("migrations.json")
			if err != nil {
				util.Check(err)
			}
			util.Check(manager.MigrationManager.CallMigrationChain(core.Version{Major: 1, Minor: 0, Patch: 0}, core.Version{Major: 3, Minor: 2, Patch: 0}, registry))

			fmt.Println("Migration completed")
		},
	}
}

func init() {
	admCmd.AddCommand(newMigrateCmd(bootstrap.CreateManager()))
}
