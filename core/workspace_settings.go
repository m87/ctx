package core

type WorkspaceSettings struct {
	LinkRules []LinkRule `json:"linkRules,omitempty"`
}

type LinkRule struct {
	Regexp string `json:"regexp"`
	Link   string `json:"link"`
}
