package repository

import (
	"app/model/publicmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"strings"
)

type Template struct {
	Name    string
	Content string
}

func (r *Repository) SelectAllTemplate(ctx context.Context) ([]Template, error) {
	list, err := publicmodel.New(r.Db).SelectAllTemplate(ctx)
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(l publicmodel.Template) Template {
		name, _ := strings.CutPrefix(l.Name, "_layout/") // TODO remove me
		return Template{
			Name:    name,
			Content: l.Content,
		}
	}), nil
}
