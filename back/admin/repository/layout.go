package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
)

type Layout struct {
	ID      int64
	Name    string
	Content string
}

func (r *Repository) CountLayout(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountLayout(ctx)
	return int(c), err
}

func (r *Repository) CountLayoutByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountLayoutByName(ctx, name)
	return int(c), err
}

func (r *Repository) CountLayoutByID(ctx context.Context, id int64) (int, error) {
	c, err := adminmodel.New(r.Db).CountLayoutByID(ctx, id)
	return int(c), err
}

type CreateLayoutParams struct {
	Name string
}

func (r *Repository) CreateLayout(ctx context.Context, l CreateLayoutParams) (int64, error) {
	id, err := adminmodel.New(r.Db).CreateLayout(ctx, adminmodel.CreateLayoutParams{
		Name: l.Name,
	})
	if err != nil {
		return 0, stacktrace.From(err)
	}
	return id, nil
}

func (r *Repository) DeleteLayout(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeleteLayout(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetLayoutList(ctx context.Context, limit, offset int) ([]Layout, error) {
	list, err := adminmodel.New(r.Db).GetLayoutList(ctx, adminmodel.GetLayoutListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(p adminmodel.Layout) Layout {
		return Layout{
			Name: p.Name,
			ID:   p.ID,
		}
	}), nil
}

func (r *Repository) GetAllLayout(ctx context.Context) ([]Layout, error) {
	list, err := adminmodel.New(r.Db).GetAllLayout(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(p adminmodel.Layout) Layout {
		return Layout{
			Name: p.Name,
			ID:   p.ID,
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
			ID:      l.ID,
		}
	}
	return fromModel(list[0]), true, nil
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
