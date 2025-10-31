package tui

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/m87/ctx/core"
)

func List(session core.Session) string {
	state := session.State
	ids := session.GetSortedContextIds()
	output := ""
	for _, id := range ids {
		v := state.Contexts[id]
		output += fmt.Sprintf("- %s\n", v.Description)
	}

	return output
}

func ListBash(session core.Session) string {
	state := session.State
	ids := session.GetSortedContextIds()
	output := ""
	for _, id := range ids {
		v := state.Contexts[id]
		output += fmt.Sprintf("%s\n", v.Description)
	}

	return output
}

func ListFull(session core.Session) string {
	state := session.State
	ids := session.GetSortedContextIds()
	output := ""
	for _, id := range ids {
		v := state.Contexts[id]
		output += fmt.Sprintf("- [%s] %s\n", id, v.Description)
		output += fmt.Sprintf("  Intervals:\n")
		for _, interval := range v.Intervals {
			output += fmt.Sprintf("\t[%s] %s - %s\n", interval.Id, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
		}

		output += fmt.Sprintf("  Labels:\n")
		for _, label := range v.Labels {
			output += fmt.Sprintf("\t%s\n", label)
		}

		output += fmt.Sprintf("  Comments:\n")
		for _, comment := range v.Comments {
			output += fmt.Sprintf("\t[%s] %s\n", comment.Id, comment.Content)
		}
	}
	return output
}

func ListJson(session core.Session) string {
	state := session.State
	v := []core.Context{}
	for _, c := range state.Contexts {
		v = append(v, c)
	}
	s, _ := json.Marshal(v)

	return string(s)
}
