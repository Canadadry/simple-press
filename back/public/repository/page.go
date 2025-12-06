package repository

import (
	"app/model/adminmodel"
	"app/pkg/stacktrace"
	"context"
)

type Layout struct {
	ID      int64
	Name    string
	Content string
}

func (r *Repository) SelectLayoutByID(ctx context.Context, id int64) (Layout, bool, error) {
	list, err := adminmodel.New(r.Db).SelectLayoutByID(ctx, id)
	if err != nil {
		return Layout{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Layout{}, false, nil
	}
	fromModel := func(l adminmodel.Layout) Layout {
		return Layout{
			Name:    l.Name,
			Content: l.Content,
			ID:      l.ID,
		}
	}
	return fromModel(list[0]), true, nil
}
