package core

import (
	"github.com/m87/ctx/util"
)

func (session *Session) LabelContext(ctxId string, label string) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	if !util.Contains(session.State.Contexts[ctxId].Labels, label) {
		ctx.Labels = append(session.State.Contexts[ctxId].Labels, label)
		// manager.PublishContextEvent(ses.Contexts[id], session.TimeProvider.Now(), LABEL_CTX, map[string]string{
		// 	"label": label,
		// })
		session.State.Contexts[ctxId] = ctx
	}
	return nil
}

func (session *Session) DeleteLabelContext(ctxId string, label string) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	if util.Contains(session.State.Contexts[ctxId].Labels, label) {
		ctx.Labels = util.Remove(session.State.Contexts[ctxId].Labels, label)
		// manager.PublishContextEvent(s.Contexts[id], session.TimeProvider.Now(), DELETE_CTX_LABEL, map[string]string{
		// 	"label": label,
		// })
		session.State.Contexts[ctxId] = ctx
	}
	return nil
}
