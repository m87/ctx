package cmd

import (
	"fmt"
	"strings"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewMergeContextCmd(manager *core.ContextManager) *cobra.Command {
	var (
		sourceID     string
		targetID     string
		deleteSource bool
	)

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Merge source context into target context",
		RunE: func(cmd *cobra.Command, args []string) error {
			source := strings.TrimSpace(sourceID)
			target := strings.TrimSpace(targetID)
			if source == "" || target == "" {
				return fmt.Errorf("source-id and target-id are required")
			}
			if source == target {
				return fmt.Errorf("source-id and target-id must be different")
			}

			moved := 0
			if resolveRemoteAddr() != "" {
				intervals, err := remoteListContextIntervals(source)
				if err != nil {
					return err
				}
				for _, interval := range intervals {
					if interval == nil || strings.TrimSpace(interval.Id) == "" {
						continue
					}
					if err := remoteMoveInterval(interval.Id, target); err != nil {
						return err
					}
					moved++
				}
				if deleteSource {
					if err := remoteDeleteContext(source); err != nil {
						return err
					}
				}
			} else {
				intervals, err := manager.IntervalRepository.ListByContextId(source)
				if err != nil {
					return err
				}
				for _, interval := range intervals {
					if interval == nil {
						continue
					}
					interval.ContextId = target
					if _, err := manager.IntervalRepository.Save(interval); err != nil {
						return err
					}
					moved++
				}
				if deleteSource {
					if err := manager.ContextRepository.Delete(source); err != nil {
						return err
					}
				}
			}

			result := map[string]any{
				"sourceId":       source,
				"targetId":       target,
				"movedIntervals": moved,
				"deletedSource":  deleteSource,
			}
			return printOutput(cmd, result, func() string {
				return fmt.Sprintf("Merged contexts successfully, moved %d intervals", moved)
			}, nil)
		},
	}

	cmd.Flags().StringVar(&sourceID, "source-id", "", "Source context ID")
	cmd.Flags().StringVar(&targetID, "target-id", "", "Target context ID")
	cmd.Flags().BoolVar(&deleteSource, "delete-source", true, "Delete source context after merge")
	_ = cmd.MarkFlagRequired("source-id")
	_ = cmd.MarkFlagRequired("target-id")
	return cmd
}

func init() {
	mergeCmd.AddCommand(NewMergeContextCmd(bootstrap.CreateManager()))
}
