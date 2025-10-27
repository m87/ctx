package core

import (
	"github.com/m87/ctx/util"
)

func (session *Session) EditContextComment(ctxId string, commentId string, comment string) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}

	ctx := session.MustGetCtx(ctxId)
	ctx.Labels = labels
	session.State.Contexts[ctxId] = ctx

	return nil
}

func (session *Session) CreateContextComment(ctxId string, comment string) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	if !util.Contains(session.State.Contexts[ctxId].Labels, label) {
		ctx.Labels = append(session.State.Contexts[ctxId].Labels, label)
		session.State.Contexts[ctxId] = ctx
	}
	return nil
}

func (session *Session) DeleteContextComment(ctxId string, commentId string) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}

	ctx := session.MustGetCtx(ctxId)
	if _, ok := ctx.Comments[commentId]; ok {

	}
	if util.Contains(session.State.Contexts[ctxId].Labels, label) {
		ctx.Labels = util.Remove(session.State.Contexts[ctxId].Labels, label)
		session.State.Contexts[ctxId] = ctx
	}
	return nil
}
