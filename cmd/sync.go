package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

const syncProgressBarWidth = 20

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync command with remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := bootstrap.CreateManager()
			if err != nil {
				return err
			}

			progressLineOpen := false
			manager.OnSyncProgress = func(progress core.SyncProgress) {
				cmd.Printf("\r%s", formatSyncProgress(progress))
				progressLineOpen = progress.Total > 0 && progress.Current < progress.Total
				if !progressLineOpen {
					cmd.Println()
				}
			}

			cmd.Println("Synchronization started")
			err = manager.Sync(resolveRemoteAddr())
			if err != nil {
				if progressLineOpen {
					cmd.Println()
				}
				return err
			}
			cmd.Println("Synchronization completed")

			return nil
		},
	}

	return cmd
}

func formatSyncProgress(progress core.SyncProgress) string {
	current := progress.Current
	if current < 0 {
		current = 0
	}
	total := progress.Total
	if total < 0 {
		total = 0
	}
	if current > total {
		current = total
	}

	filled := 0
	if total > 0 {
		filled = current * syncProgressBarWidth / total
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", syncProgressBarWidth-filled)
	return fmt.Sprintf("Syncing %-10s [%s] %d/%d", progress.Resource, bar, current, total)
}

func init() {
	rootCmd.AddCommand(NewSyncCmd())
}
