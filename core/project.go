package core

const (
	ProjectType = "project"
)

type Project struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ParentId    string `json:"parentId,omitempty"`
	WorkspaceId string `json:"workspaceId"`
}
