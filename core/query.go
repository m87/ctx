package core

import "fmt"

type ContextSQLQuery struct {
	WhereClause string
	Arguments   []any
}

type ContextQueryInterpreter interface {
	Interpret(query string) (*ContextSQLQuery, error)
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
	return m.ContextRepository.Query(sqlQuery)
}
