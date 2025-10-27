package cmd

import (
	"fmt"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
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
			for _, migration := range manager.MigrationManager.CreateMigrationChain(core.Version{Major: 1, Minor: 0, Patch: 0}, core.Version{Major: 3, Minor: 2, Patch: 0}) {
				util.Checkm(manager.WithSession(func(session core.Session) error {
					return migration.Migrate(session)
				}), "Migration failed")

				util.Checkm(manager.WithContextArchiver(func(archiver core.Archiver[core.Context]) error {
					return migration.MigrateArchive(archiver)
				}), "Archive migration failed")
			}
			fmt.Println("Migration completed")
		},
	}
}

func init() {
	admCmd.AddCommand(newMigrateCmd(bootstrap.CreateManager()))
}
