package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/m87/ctx/core"
	localstorage "github.com/m87/ctx/storage/local"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate data to a new format or structure",
		Long:  "This command is used to migrate data from the old format to the new format.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Migration process started...")

			registry, err := localstorage.LoadLocalMigrationRegistry(filepath.Join(viper.GetString("storePath"), "migrations"))
			if err != nil {
				util.Check(err)
			}

			manager := &localstorage.LocalStorageMigrationManager{StatePath: filepath.Join(viper.GetString("storePath"), "state"), ArchivePath: filepath.Join(viper.GetString("storePath"), "archive")}
			util.Check(core.Migrate(manager, registry))

			fmt.Println("Migration completed")
		},
	}
}

func init() {
	admCmd.AddCommand(newMigrateCmd())
}
