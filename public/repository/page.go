package repository

import (
	"app/model/adminmodel"
	"app/pkg/stacktrace"
	"context"
)

type Page struct {
	ID      int64
	Name    string
	Content string
}

func (r *Repository) SelectPageByID(ctx context.Context, id int64) (Page, bool, error) {
	list, err := adminmodel.New(r.Db).SelectPageByID(ctx, id)
	if err != nil {
		return Page{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Page{}, false, nil
	}
	fromModel := func(l adminmodel.Page) Page {
		return Page{
			Name:    l.Name,
			Content: l.Content,
			ID:      l.ID,
		}
	}
	return fromModel(list[0]), true, nil
}
