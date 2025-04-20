package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"strings"
)

type Template struct {
	Name    string
	Content string
}

func (r *Repository) CountTemplates(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountTemplate(ctx)
	return int(c), err
}

func (r *Repository) CountTemplateByName(ctx context.Context, name string) (int, error) {
	c, err := adminmodel.New(r.Db).CountTemplateByName(ctx, name)
	return int(c), err
}

type CreateTemplateParams struct {
	Name string
}

func (r *Repository) CreateTemplate(ctx context.Context, l CreateTemplateParams) error {
	_, err := adminmodel.New(r.Db).CreateTemplate(ctx, adminmodel.CreateTemplateParams{
		Name: l.Name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteTemplate(ctx context.Context, name string) error {
	err := adminmodel.New(r.Db).DeleteTemplate(ctx, name)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetTemplateList(ctx context.Context, limit, offset int) ([]Template, error) {
	list, err := adminmodel.New(r.Db).GetTemplateList(ctx, adminmodel.GetTemplateListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(name string) Template {
		return Template{
			Name: name,
		}
	}), nil
}

func (r *Repository) SelectAllTemplate(ctx context.Context) ([]Template, error) {
	list, err := adminmodel.New(r.Db).SelectAllTemplate(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(l adminmodel.Template) Template {
		name, _ := strings.CutPrefix(l.Name, "_layout/")
		return Template{
			Name:    name,
			Content: l.Content,
		}
	}), nil
}

func (r *Repository) SelectTemplate(ctx context.Context, name string) (Template, bool, error) {
	list, err := adminmodel.New(r.Db).SelectTemplate(ctx, name)
	if err != nil {
		return Template{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Template{}, false, nil
	}
	fromModel := func(l adminmodel.Template) Template {
		return Template{
			Name:    l.Name,
			Content: l.Content,
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdateTemplate(ctx context.Context, name string, l Template) error {
	err := adminmodel.New(r.Db).UpdateTemplate(ctx, adminmodel.UpdateTemplateParams{
		Name:    l.Name,
		Content: l.Content,
		Name_2:  name,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
