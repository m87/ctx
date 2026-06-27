package core

type IntegrityIssue struct {
	EntityType  string                 `json:"entityType"`
	EntityId    string                 `json:"entityId"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Repairable  bool                   `json:"repairable"`
	Details     *IntegrityIssueDetails `json:"details,omitempty"`
}

type IntegrityIssueDetails struct {
	Name        string     `json:"name,omitempty"`
	ContextId   string     `json:"contextId,omitempty"`
	WorkspaceId string     `json:"workspaceId,omitempty"`
	Start       *ZonedTime `json:"start,omitempty"`
	End         *ZonedTime `json:"end,omitempty"`
}

type IntegrityReport struct {
	Healthy        bool              `json:"healthy"`
	WorkspaceCount int               `json:"workspaceCount"`
	ContextCount   int               `json:"contextCount"`
	IntervalCount  int               `json:"intervalCount"`
	Issues         []*IntegrityIssue `json:"issues"`
}

type IntegrityContextOption struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	WorkspaceId   string `json:"workspaceId"`
	WorkspaceName string `json:"workspaceName"`
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

func contextIntegrityIssue(context *Context, code, description string) *IntegrityIssue {
	issue := integrityIssue("context", context.Id, code, description)
	issue.Repairable = true
	issue.Details = &IntegrityIssueDetails{
		Name:        context.Name,
		WorkspaceId: context.WorkspaceId,
	}
	return issue
}

func intervalIntegrityIssue(interval *Interval, code, description string, repairable bool) *IntegrityIssue {
	issue := integrityIssue("interval", interval.Id, code, description)
	issue.Repairable = repairable
	issue.Details = &IntegrityIssueDetails{
		ContextId:   interval.ContextId,
		WorkspaceId: interval.WorkspaceId,
		Start:       &interval.Start,
		End:         &interval.End,
	}
	return issue
}

func zonedTimeIsSet(value ZonedTime) bool {
	return !value.Time.IsZero() && !value.IsZero
}

func intervalIsInactive(interval *Interval) bool {
	return interval != nil && interval.Status != "active"
}

func intervalIsOpenActive(interval *Interval) bool {
	return interval != nil && interval.Status == "active" && !zonedTimeIsSet(interval.End)
}
