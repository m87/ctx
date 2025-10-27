package localstorage

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/m87/ctx/core"
)

type LocalStoreContextArchiver struct {
	path string
}

func (archiver *LocalStoreContextArchiver) Update(contextsToUpdate []core.Context, session core.Session) error {
	return nil
}

func (archiver *LocalStoreContextArchiver) Archive(contextsToArchvie []core.Context, session core.Session) error {
	for _, context := range contextsToArchvie {
		err := archiver.archiveContext(context)
		if err != nil {
			return err
		}

		err = session.Delete(context.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func merge(src *core.Context, dst *core.Context) *core.Context {
	if src == nil {
		return dst
	}

	src.Id = dst.Id
	src.Description = dst.Description
	src.Duration = src.Duration + dst.Duration

	for _, comment := range dst.Comments {
		src.Comments[comment.Id] = comment
	}

	src.Labels = append(src.Labels, dst.Labels...)
	src.State = dst.State

	if src.Intervals == nil {
		src.Intervals = map[string]core.Interval{}
	}

	for k, v := range dst.Intervals {
		src.Intervals[k] = v
	}

	return src
}

func (archiver *LocalStoreContextArchiver) archiveContext(context core.Context) error {
	contextPath := filepath.Join(archiver.path, context.Id+".ctx")
	archivedContext, err := Load[core.Context](contextPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return Save(merge(&archivedContext, &context), contextPath)
}
