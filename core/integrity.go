package core

type IntegrityIssue struct {
	EntityType  string `json:"entityType"`
	EntityId    string `json:"entityId"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type IntegrityReport struct {
	Healthy        bool              `json:"healthy"`
	WorkspaceCount int               `json:"workspaceCount"`
	ContextCount   int               `json:"contextCount"`
	IntervalCount  int               `json:"intervalCount"`
	Issues         []*IntegrityIssue `json:"issues"`
}

type IntegrityRepairResult struct {
	RepairedCount int              `json:"repairedCount"`
	Report        *IntegrityReport `json:"report"`
}

func integrityIssue(entityType, entityId, code, description string) *IntegrityIssue {
	return &IntegrityIssue{
		EntityType:  entityType,
		EntityId:    entityId,
		Code:        code,
		Description: description,
	}
}
