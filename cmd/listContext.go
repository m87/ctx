package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/spf13/cobra"
)

func NewListContextCmd(manager *core.ContextManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "List all contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			var contexts []*core.Context
			var err error
			if resolveRemoteAddr() != "" {
				contexts, err = remoteListContexts()
			} else {
				contexts, err = manager.ContextRepository.List()
			}
			if err != nil {
				return err
			}

			// Non-verbose: keep previous compact listing
			if !Verbose {
				return printOutput(cmd, contexts, func() string {
					if len(contexts) == 0 {
						return "No contexts found"
					}
					lines := make([]string, 0, len(contexts))
					for _, context := range contexts {
						lines = append(lines, "- ID: "+context.Id+", Name: "+context.Name)
					}
					return strings.Join(lines, "\n")
				}, nil)
			}

			// Verbose: include full context info and list intervals for each context
			type ctxWithIntervals struct {
				Context   *core.Context    `json:"context"`
				Intervals []*core.Interval `json:"intervals"`
			}

			verboseList := make([]*ctxWithIntervals, 0, len(contexts))
			for _, context := range contexts {
				var intervals []*core.Interval
				if resolveRemoteAddr() != "" {
					ivs, ierr := remoteListContextIntervals(context.Id)
					if ierr != nil {
						return ierr
					}
					intervals = ivs
				} else {
					ivs, ierr := manager.IntervalRepository.ListByContextId(context.Id)
					if ierr != nil {
						return ierr
					}
					intervals = ivs
				}
				verboseList = append(verboseList, &ctxWithIntervals{Context: context, Intervals: intervals})
			}

			// Text renderer: detailed human-readable with indented intervals
			textRenderer := func() string {
				if len(verboseList) == 0 {
					return "No contexts found"
				}
				var b strings.Builder
				for idx, item := range verboseList {
					c := item.Context
					if idx > 0 {
						b.WriteString("\n")
					}
					b.WriteString(fmt.Sprintf("ID: %s\n", c.Id))
					b.WriteString(fmt.Sprintf("Name: %s\n", c.Name))
					b.WriteString(fmt.Sprintf("ParentId: %s\n", c.ParentId))
					b.WriteString(fmt.Sprintf("Status: %s\n", c.Status))
					if c.Description != "" {
						b.WriteString(fmt.Sprintf("Description: %s\n", c.Description))
					}
					if len(c.Tags) > 0 {
						b.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(c.Tags, ", ")))
					}
					b.WriteString("Intervals:\n")
					if len(item.Intervals) == 0 {
						b.WriteString("  (none)\n")
					} else {
						for _, iv := range item.Intervals {
							var endStr string
							if iv.End.IsZero {
								endStr = "(ongoing)"
							} else {
								endStr = iv.End.Time.In(iv.End.Time.Location()).Format(time.RFC3339)
							}
							b.WriteString(fmt.Sprintf("  - ID: %s, Start: %s, End: %s, Status: %s\n", iv.Id, iv.Start.Time.In(iv.Start.Time.Location()).Format(time.RFC3339), endStr, iv.Status))
						}
					}
				}
				return b.String()
			}

			// For structured outputs (json/yaml/shell), print the verboseList so intervals are included
			return printOutput(cmd, verboseList, textRenderer, nil)
		},
	}
	return cmd
}

func init() {
	listCmd.AddCommand(NewListContextCmd(bootstrap.CreateManager()))
}
