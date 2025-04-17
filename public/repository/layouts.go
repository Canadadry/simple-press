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

func (r *Repository) SelectBaseLayout(ctx context.Context) ([]Layout, error) {
	list, err := publicmodel.New(r.Db).SelectBaseLayout(ctx)
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

func (r *Repository) SelectLayout(ctx context.Context, name string) (Layout, bool, error) {
	list, err := publicmodel.New(r.Db).SelectLayout(ctx, name)
	if err != nil {
		return Layout{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Layout{}, false, nil
	}
	fromModel := func(l publicmodel.Layout) Layout {
		return Layout{
			Name:    l.Name,
			Content: l.Content,
		}
	}
	return fromModel(list[0]), true, nil
}
