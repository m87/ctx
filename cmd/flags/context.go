package flags

import (
	"errors"
	"strings"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func ResolveContextId(cmd *cobra.Command) (string, error) {
	return ResolveCustomContextId(cmd, "ctx")
}

func ResolveCustomContextId(cmd *cobra.Command, name string) (string, error) {
	flags := cmd.Flags()

	id, _ := flags.GetString(name + "-id")
	description, _ := flags.GetString(name)

	if id != "" && description != "" {
		return "", errors.New("both --" + name + " and --" + name + "-id provided")
	}

	if id != "" {
		return id, nil
	}

	if description != "" {
		return util.GenerateId(description), nil
	}

	return "", errors.New("either --" + name + " or --" + name + "-id must be provided")
}

func AddContxtFlag(cmd *cobra.Command) {
	AddCustomContextFlag(cmd, "ctx", "c", "Context")
}

func AddCustomContextFlag(cmd *cobra.Command, name string, short string, description string) {
	cmd.Flags().StringP(name, short, "", description+" description")
	cmd.Flags().StringP(name+"-id", strings.ToUpper(short), "", description+" id")
}
