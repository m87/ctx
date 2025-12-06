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

type ArgCursor struct {
	next int
}

func NewArgCursor(ctxProvided bool) *ArgCursor {
	c := &ArgCursor{}

	if !ctxProvided {
		c.next = 1
	} else {
		c.next = 0
	}

	return c
}

func (c *ArgCursor) Current() int {
	return c.next
}

func (c *ArgCursor) Next() {
	c.next++
}

type ParamSpec struct {
	Default string
	Name    string
}

func ResolveCidWithParams(args []string, ctxId string, params ...ParamSpec) (ContextId, map[string]string, error) {
	cid, err := ResolveContextId(args, ctxId)
	if err != nil {
		return ContextId{}, nil, err
	}
	cursor := NewArgCursor(ctxId != "")
	resolvedParams := make(map[string]string, len(params))
	for _, param := range params {
		idx := cursor.Current()
		val, usedPos, err := ResolveArgument(args, idx, param.Default, param.Name)
		if err != nil {
			return ContextId{}, nil, err
		}
		resolvedParams[param.Name] = val
		if usedPos {
			cursor.Next()
		}
	}
	return cid, resolvedParams, nil
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
	if ctxId != "" {
		return ContextId{Id: ctxId, Description: ""}, nil
	}

	if len(positional) == 0 {
		return ContextId{}, errors.New("either positional argument or --ctx-id must be provided")
	}

	return ContextId{Id: util.GenerateId(strings.TrimSpace(positional[0])), Description: strings.TrimSpace(positional[0])}, nil
}

func AddPrefixedContextIdFlags(cmd *cobra.Command, ctxId *string, prefix string, docPrefix string) {
	cmd.Flags().StringVar(ctxId, prefix+"ctx-id", "", docPrefix+"context id")
}

func ResolveArgument(args []string, index int, property string, name string) (string, bool, error) {
	if property != "" {
		return property, false, nil
	}

	if len(args) > index {
		return args[index], true, nil
	}

	return "", false, errors.New("missing required argument: " + name)
}
