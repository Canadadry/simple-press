package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"database/sql"
	"time"
)

type Article struct {
	Title  string
	Date   time.Time
	Author string
	Slug   string
	Draft  bool
}

func (r *Repository) CountArticles(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountArticles(ctx)
	return int(c), err
}

func (r *Repository) CreateArticle(ctx context.Context, a Article) error {
	_, err := adminmodel.New(r.Db).CreateArticle(ctx, adminmodel.CreateArticleParams{
		Title:  a.Title,
		Date:   r.Clock.Now(),
		Author: a.Author,
		Slug:   a.Slug,
		Draft:  sql.NullInt64{Int64: 1, Valid: a.Draft},
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) DeleteArticle(ctx context.Context, slug string) error {
	err := adminmodel.New(r.Db).DeleteArticle(ctx, slug)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetArticlesList(ctx context.Context, limit, offset int) ([]Article, error) {
	list, err := adminmodel.New(r.Db).GetArticlesList(ctx, adminmodel.GetArticlesListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(a adminmodel.GetArticlesListRow) Article {
		return Article{
			Title:  a.Title,
			Date:   a.Date,
			Author: a.Author,
			Slug:   a.Slug,
			Draft:  a.Draft.Valid,
		}
	}), nil
}

func (r *Repository) SelectArticleBySlug(ctx context.Context, slug string) (Article, bool, error) {
	list, err := adminmodel.New(r.Db).SelectArticleBySlug(ctx, slug)
	if err != nil {
		return Article{}, false, stacktrace.From(err)
	}
	if len(list) == 0 {
		return Article{}, false, nil
	}
	fromModel := func(a adminmodel.Article) Article {
		return Article{
			Title:  a.Title,
			Date:   a.Date,
			Author: a.Author,
			Slug:   a.Slug,
			Draft:  a.Draft.Valid,
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdateArticle(ctx context.Context, slug string, a Article) error {
	err := adminmodel.New(r.Db).UpdateArticle(ctx, adminmodel.UpdateArticleParams{
		Title:  a.Title,
		Date:   r.Clock.Now(),
		Author: a.Author,
		Slug:   a.Slug,
		Draft:  sql.NullInt64{Int64: 1, Valid: a.Draft},
		Slug_2: slug,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
