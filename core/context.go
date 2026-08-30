package core

type ProjectMetadata struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Context struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	ProjectId   *string          `json:"projectId"`
	WorkspaceId string           `json:"workspaceId"`
	Status      string           `json:"status"`
	Archived    bool             `json:"archived"`
	Description string           `json:"description,omitempty"`
	Tags        []*Tag           `json:"tags,omitempty"`
	Project     *ProjectMetadata `json:"project"`
}

const ContextType = "context"
