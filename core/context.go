package core

type Context struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
}

type ContextMapper struct {
}

func NewContextMapper() *ContextMapper {
	return &ContextMapper{}
}
