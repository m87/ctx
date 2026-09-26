package core

import (
	"fmt"
	"strings"
	"text/scanner"
)

// CQL
// text | string | "search text" | ~ = !=
// duration | string | 2h30 | > < >= <= =
// start, end | datetime | 2023-01-01T00:00:00Z | > <	>= <= | relative time: -30d -1y -1m - 1h now -2w ...
// label | [] | ("label1", "label2") | in, not in
// project | string, [] | pr , ("pr1", "pr2")  | ~ = ~= in not in

type ContextSQLQuery struct {
	WorkspaceId string
	WhereClause string
	Arguments   []any
}

type ContextQueryResult struct {
	WorkspaceId   string                   `json:"workspaceId"`
	Query         string                   `json:"query"`
	Contexts      []*Context               `json:"contexts"`
	ContextStats  []*WorkspaceContextStats `json:"contextStats"`
	TotalDuration int64                    `json:"totalDuration"`
	TotalSessions int                      `json:"totalSessions"`
}

type ContextQueryInterpreter interface {
	Interpret(query string) (*ContextSQLQuery, error)
}

type DefaultContextQueryInterpreter struct{}

func (i *DefaultContextQueryInterpreter) Interpret(query string) (*ContextSQLQuery, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(query))
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanInts | scanner.ScanFloats | scanner.SkipComments

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("token: %s\n", s.TokenText())
	}

	return &ContextSQLQuery{}, nil
}

type PassthroughContextQueryInterpreter struct{}

func (i *PassthroughContextQueryInterpreter) Interpret(_ string) (*ContextSQLQuery, error) {
	return &ContextSQLQuery{}, nil
}

type InvalidContextQueryError struct {
	Err error
}

func (e *InvalidContextQueryError) Error() string {
	if e.Err == nil {
		return "invalid context query"
	}
	return fmt.Sprintf("invalid context query: %s", e.Err)
}

func (e *InvalidContextQueryError) Unwrap() error {
	return e.Err
}

func (m *ContextManager) QueryContexts(query string) ([]*Context, error) {
	sqlQuery, err := m.interpretContextQuery(query)
	if err != nil {
		return nil, err
	}
	return m.ContextRepository.Query(sqlQuery)
}

func (m *ContextManager) QueryWorkspaceContexts(workspaceId string, query string) (*ContextQueryResult, error) {
	workspaceId = strings.TrimSpace(workspaceId)
	if workspaceId == "" {
		return nil, &WorkspaceNotFoundError{}
	}

	workspace, err := m.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, &WorkspaceNotFoundError{WorkspaceId: workspaceId}
	}

	sqlQuery, err := m.interpretContextQuery(query)
	if err != nil {
		return nil, err
	}
	sqlQuery.WorkspaceId = workspaceId
	contexts, err := m.ContextRepository.Query(sqlQuery)
	if err != nil {
		return nil, err
	}

	contextStats, totalDuration, totalSessions, err := m.getContextCollectionStats(contexts)
	if err != nil {
		return nil, err
	}
	return &ContextQueryResult{
		WorkspaceId:   workspaceId,
		Query:         query,
		Contexts:      contexts,
		ContextStats:  contextStats,
		TotalDuration: int64(totalDuration),
		TotalSessions: totalSessions,
	}, nil
}

func (m *ContextManager) interpretContextQuery(query string) (*ContextSQLQuery, error) {
	if m.QueryInterpreter == nil {
		return nil, fmt.Errorf("context query interpreter is required")
	}
	sqlQuery, err := m.QueryInterpreter.Interpret(query)
	if err != nil {
		return nil, err
	}
	if sqlQuery == nil {
		return nil, fmt.Errorf("context query interpreter returned an empty query")
	}
	return &ContextSQLQuery{
		WorkspaceId: sqlQuery.WorkspaceId,
		WhereClause: sqlQuery.WhereClause,
		Arguments:   append([]any(nil), sqlQuery.Arguments...),
	}, nil
}
