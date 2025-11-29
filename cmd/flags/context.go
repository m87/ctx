package flags

import (
	"errors"
	"strings"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

type ContextCmdIdentifier struct {
	Id          string
	Description string
}

func ResolveContextIdLegacy(cmd *cobra.Command) (string, error) {
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

func AddContextIdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("ctx-id", "c", "", "context id")
	cmd.Flags().String("ctx", "", "context description")
	cmd.MarkFlagsMutuallyExclusive("ctx-id", "ctx")
}

func resolveContextId(positional string, ctxId string, ctxDescription string) (string, bool) {
	switch {
	case ctxDescription != "":
		return ctxDescription, false

	case ctxId != "":
		return ctxId, true

	default:
		return positional, false
	}
}

func ResolveContextId(positional string, ctxId string, ctxDescription string) (string, string, bool, error) {
	rawId, isId := resolveContextId(positional, ctxId, ctxDescription)
	trimmedId := strings.TrimSpace(rawId)
	if trimmedId == "" {
		return "", "", false, errors.New("context id not provided")
	}

	id, err := util.Id(trimmedId, isId)
	if err != nil {
		return "", "", false, err
	}

	return id, rawId, isId, nil
}

func ResolveContextIdentifier(cmd *cobra.Command, args []string) (ContextCmdIdentifier, error) {
	ctxId, _ := cmd.Flags().GetString("ctx-id")
	ctxDescription, _ := cmd.Flags().GetString("ctx")

	var positional string
	if len(args) > 0 {
		positional = args[0]
	}

	if positional != "" && (ctxId != "" || ctxDescription != "") {
		return ContextCmdIdentifier{}, errors.New("cannot use positional argument together with --ctx or --ctx-id flags")
	}

	if positional != "" {
		return ContextCmdIdentifier{
			Id:          util.GenerateId(positional),
			Description: positional,
		}, nil
	}

	if ctxId != "" && ctxDescription != "" {
		return ContextCmdIdentifier{}, errors.New("both --ctx and --ctx-id provided")
	}

	if ctxDescription != "" {
		generatedId := util.GenerateId(ctxDescription)
		return ContextCmdIdentifier{
			Id:          generatedId,
			Description: ctxDescription,
		}, nil
	}

	if ctxId != "" {
		return ContextCmdIdentifier{
			Id:          ctxId,
			Description: "",
		}, nil
	}

	return ContextCmdIdentifier{}, errors.New("no context identifier provided. Use either positional argument, --ctx or --ctx-id")
}

func AddPrefixedContextIdFlags(cmd *cobra.Command, ctxId *string, ctxDescription *string, prefix string, docPrefix string) {
	cmd.Flags().StringVar(ctxId, prefix+"ctx-id", "", docPrefix+"context id")
	cmd.Flags().StringVar(ctxDescription, prefix+"ctx", "", docPrefix+"context description")
	cmd.MarkFlagsMutuallyExclusive("ctx-id", "ctx")
}
