package core

import (
	"strconv"

	"github.com/m87/nod"
)

type WorkspaceSettings struct {
	LinkRules []LinkRule `json:"linkRules,omitempty"`
}

type LinkRule struct {
	Regexp string `json:"regexp"`
	Link   string `json:"link"`
}

func FromKV(kv map[string]*nod.NodeKV) *WorkspaceSettings {
	settings := &WorkspaceSettings{
		LinkRules: []LinkRule{},
	}

	if kv == nil {
		return settings
	}

	for i := 0; ; i++ {
		prefix := "linkRule." + strconv.Itoa(i)
		regexpKey := prefix + ".regexp"
		linkKey := prefix + ".link"

		regexpKV, regexpExists := kv[regexpKey]
		linkKV, linkExists := kv[linkKey]

		if !regexpExists || !linkExists {
			break
		}

		rule := LinkRule{
			Regexp: "",
			Link:   "",
		}

		if regexpKV != nil && regexpKV.ValueText != nil {
			rule.Regexp = *regexpKV.ValueText
		}

		if linkKV != nil && linkKV.ValueText != nil {
			rule.Link = *linkKV.ValueText
		}

		settings.LinkRules = append(settings.LinkRules, rule)
	}

	return settings

}

func ToKV(settings *WorkspaceSettings) map[string]*nod.NodeKV {
	kv := make(map[string]*nod.NodeKV)

	if settings == nil {
		return kv
	}

	for i, rule := range settings.LinkRules {
		prefix := "linkRule." + strconv.Itoa(i)
		regexpKey := prefix + ".regexp"
		linkKey := prefix + ".link"

		kv[regexpKey] = &nod.NodeKV{Key: regexpKey, ValueText: nod.Ptr(rule.Regexp)}
		kv[linkKey] = &nod.NodeKV{Key: linkKey, ValueText: nod.Ptr(rule.Link)}
	}

	return kv
}
