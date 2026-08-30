package storage

import "github.com/m87/ctx/core"

type WorkspaceEntity struct {
	Id          string `gorm:"primaryKey"`
	Name        string
	Description string
	LinkRules   []WorkspaceLinkRuleEntity `gorm:"foreignKey:WorkspaceId;references:Id;constraint:OnDelete:CASCADE"`
}

type WorkspaceLinkRuleEntity struct {
	WorkspaceId string `gorm:"primaryKey"`
	Position    int    `gorm:"primaryKey"`
	Regexp      string `gorm:"not null"`
	Link        string `gorm:"not null"`
}

func (WorkspaceEntity) TableName() string {
	return "workspaces"
}

func (WorkspaceLinkRuleEntity) TableName() string {
	return "workspace_link_rules"
}

func (e *WorkspaceEntity) ToModel() *core.Workspace {
	linkRules := make([]core.LinkRule, len(e.LinkRules))
	for i, rule := range e.LinkRules {
		linkRules[i] = core.LinkRule{Regexp: rule.Regexp, Link: rule.Link}
	}

	return &core.Workspace{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		Properties:  &core.WorkspaceSettings{LinkRules: linkRules},
	}
}

func NewWorkspaceEntityFromModel(workspace *core.Workspace) *WorkspaceEntity {
	entity := &WorkspaceEntity{
		Id:          workspace.Id,
		Name:        workspace.Name,
		Description: workspace.Description,
	}
	if workspace.Properties != nil {
		entity.LinkRules = make([]WorkspaceLinkRuleEntity, len(workspace.Properties.LinkRules))
		for i, rule := range workspace.Properties.LinkRules {
			entity.LinkRules[i] = WorkspaceLinkRuleEntity{
				WorkspaceId: workspace.Id,
				Position:    i,
				Regexp:      rule.Regexp,
				Link:        rule.Link,
			}
		}
	}

	return entity
}
