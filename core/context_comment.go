package core

import (
	"errors"

	"github.com/google/uuid"
)

func (session *Session) SaveContextComment(contextId string, comment Comment) error {
	ctx, exists := session.State.Contexts[contextId]
	if !exists {
		return errors.New("context not found")
	}
	if ctx.Comments == nil {
		ctx.Comments = make(map[string]Comment)
	}
	if comment.Id == "" {
		comment.Id = uuid.New().String()
	}
	ctx.Comments[comment.Id] = comment
	session.State.Contexts[contextId] = ctx
	return nil
}

func (session *Session) DeleteContextComment(contextId string, commentId string) error {
	ctx, exists := session.State.Contexts[contextId]
	if !exists {
		return errors.New("context not found")
	}
	delete(ctx.Comments, commentId)
	session.State.Contexts[contextId] = ctx
	return nil
}
