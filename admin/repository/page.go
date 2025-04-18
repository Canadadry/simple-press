package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
)

type Page struct {
	Name    string
	Content string
}

func (r *Repository) CountPages(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountPage(ctx)
	return int(c), err
}

func (r *Repository) CountPageByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountPageByName(ctx, name)
	return int(c), err
}

type CreatePageParams struct {
	Name string
}

func (r *Repository) CreatePage(ctx context.Context, l CreatePageParams) error {
	_, err := adminmodel.New(r.Db).CreatePage(ctx, adminmodel.CreatePageParams{
		Name: l.Name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeletePage(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeletePage(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetPageList(ctx context.Context, limit, offset int) ([]Page, error) {
	list, err := adminmodel.New(r.Db).GePageList(ctx, adminmodel.GePageListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(name string) Page {
		return Page{
			Name: name,
		}
	}), nil
}

func (r *Repository) SelectPage(ctx context.Context, name string) (Page, bool, error) {
	list, err := adminmodel.New(r.Db).SelectPage(ctx, name)
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
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdatePage(ctx context.Context, name string, l Page) error {
	err := adminmodel.New(r.Db).UpdatePage(ctx, adminmodel.UpdatePageParams{
		Name:    l.Name,
		Content: l.Content,
		Name_2:  name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
