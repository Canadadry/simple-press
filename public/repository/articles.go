package repository

import (
	"app/model/publicmodel"
	"app/pkg/stacktrace"
	"context"
	"time"
)

type Article struct {
	Title   string
	Date    time.Time
	Author  string
	Slug    string
	Content string
}

func (r *Repository) SelectArticleBySlug(ctx context.Context, slug string) (Article, bool, error) {
	list, err := publicmodel.New(r.Db).SelectArticleBySlug(ctx, slug)
	if err != nil {
		return Article{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Article{}, false, nil
	}
	fromModel := func(a publicmodel.Article) Article {
		return Article{
			Title:   a.Title,
			Date:    a.Date,
			Author:  a.Author,
			Content: a.Content,
			Slug:    a.Slug,
		}
	}
	return fromModel(list[0]), true, nil
}
