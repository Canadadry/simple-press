package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"strings"
)

type Layout struct {
	Name    string
	Content string
}

func (r *Repository) CountLayouts(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountLayout(ctx)
	return int(c), err
}

func (r *Repository) CountLayoutByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountLayoutByName(ctx, name)
	return int(c), err
}

type CreateLayoutParams struct {
	Name string
}

func (r *Repository) CreateLayout(ctx context.Context, l CreateLayoutParams) error {
	_, err := adminmodel.New(r.Db).CreateLayout(ctx, adminmodel.CreateLayoutParams{
		Name: l.Name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteLayout(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeleteLayout(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GeLayoutList(ctx context.Context, limit, offset int) ([]Layout, error) {
	list, err := adminmodel.New(r.Db).GeLayoutList(ctx, adminmodel.GeLayoutListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(name string) Layout {
		return Layout{
			Name: name,
		}
	}), nil
}

func (r *Repository) SelectBaseLayout(ctx context.Context) ([]Layout, error) {
	list, err := adminmodel.New(r.Db).SelectBaseLayout(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(l adminmodel.Layout) Layout {
		name, _ := strings.CutPrefix(l.Name, "_layout/")
		return Layout{
			Name:    name,
			Content: l.Content,
		}
	}), nil
}

func (r *Repository) SelectLayout(ctx context.Context, name string) (Layout, bool, error) {
	list, err := adminmodel.New(r.Db).SelectLayout(ctx, name)
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
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdateLayout(ctx context.Context, name string, l Layout) error {
	err := adminmodel.New(r.Db).UpdateLayout(ctx, adminmodel.UpdateLayoutParams{
		Name:    l.Name,
		Content: l.Content,
		Name_2:  name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
