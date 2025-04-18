package repository

import (
	"app/model/publicmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"strings"
)

type Layout struct {
	Name    string
	Content string
}

func (r *Repository) SelectAllLayout(ctx context.Context) ([]Layout, error) {
	list, err := publicmodel.New(r.Db).SelectAllLayout(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(l publicmodel.Layout) Layout {
		name, _ := strings.CutPrefix(l.Name, "_layout/")
		return Layout{
			Name:    name,
			Content: l.Content,
		}
	}), nil
}
