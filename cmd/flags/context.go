package flags

import (
	"errors"
	"strings"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

type ContextId struct {
	Id          string
	Description string
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



func AddCustomContextFlag(cmd *cobra.Command, name string, short string, description string) {
	cmd.Flags().StringP(name, short, "", description+" description")
	cmd.Flags().StringP(name+"-id", strings.ToUpper(short), "", description+" id")
}

func AddContextIdFlags(cmd *cobra.Command, ctxId *string) {
	cmd.Flags().StringVar(ctxId, "ctx-id", "", "context id")
}

func ResolveContextId(positional []string, ctxId string) (ContextId, error) {
	if len(positional) == 0 {
		return ContextId{}, errors.New("either positional argument or --ctx-id must be provided")
	}
	switch {
	case ctxId != "":
		return ContextId{Id: ctxId, Description: ""}, nil

	default:
		return ContextId{Id: util.GenerateId(strings.TrimSpace(positional[0])), Description: strings.TrimSpace(positional[0])}, nil
	}
}

func AddPrefixedContextIdFlags(cmd *cobra.Command, ctxId *string, prefix string, docPrefix string) {
	cmd.Flags().StringVar(ctxId, prefix+"ctx-id", "", docPrefix+"context id")
}

func ResolveArgument(args []string, index int, property string, name string) (string, error) {
	if property != "" {
		return property, nil
	}

	if len(args) > index {
		return args[index], nil
	}

	return "", errors.New("missing required argument: " + name)
}
