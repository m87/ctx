package flags

import (
	"errors"

	"github.com/m87/ctx/util"
)

func GetStringArg(args []string, i int, name string) (string, error) {
	if len(args) <= i {
		return "", errors.New("missing " + name)
	}
	return args[i], nil
}

func ResolveArgumentAsContextId(args []string, i int, name string) (string, error) {
	if len(args) <= i {
		return "", errors.New("missing " + name)
	}

	id, err := util.Id(args[i], false)
	if err != nil {
		return "", err
	}
	return id, nil
}

func ConditionalIndexProvider(ctxIdExists bool) func(int) int {
	return func(baseIndex int) int {
		if ctxIdExists {
			return baseIndex - 1
		}
		return baseIndex
	}
}
