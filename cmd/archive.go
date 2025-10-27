package cmd

import (
	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func NewArchiveCmd(manager *core.ContextManager) *cobra.Command {
	return &cobra.Command{
		Use:     "archive",
		Aliases: []string{"archive", "a"},
		Short:   "Archive contexts",

		Run: func(cmd *cobra.Command, args []string) {
			util.Check(manager.WithSession(func(session core.Session) error {
				return manager.WithContextArchiver(func(archver core.Archiver[core.Context]) error {
					contextsToArchvie := []core.Context{}

					for _, v := range session.State.Contexts {
						if session.State.CurrentId == v.Id {
							continue
						}

						contextsToArchvie = append(contextsToArchvie, v)
					}

					return archver.Archive(contextsToArchvie, session)
				})
			}))
		},
	}

}

func init() {
	cmd := NewArchiveCmd(bootstrap.CreateManager())

	admCmd.AddCommand(cmd)
}
