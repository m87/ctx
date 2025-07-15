package flags

import (
	"errors"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func ResolveContextId(cmd *cobra.Command) (string, error) {
	flags := cmd.Flags()

	id, _ := flags.GetString("ctx-id")
	description, _ := flags.GetString("ctx")

	if id != "" && description != "" {
		return "", errors.New("both --ctx and -ctx-id provided")
	}

	if id != "" {
		return id, nil
	}

	if description != "" {
		return util.GenerateId(description), nil
	}

	return "", errors.New("either --ctx or --ctx-id must be provided")
}

func AddContxtFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("ctx", "c", "", "Context description")
	cmd.Flags().StringP("ctx-id", "C", "", "Context id")
}
