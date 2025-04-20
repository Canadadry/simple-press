package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"time"
)

type Article struct {
	ID       int64
	Title    string
	Date     time.Time
	Author   string
	Slug     string
	Content  string
	Draft    bool
	LayoutID int64
}

func (r *Repository) CountArticle(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountArticle(ctx)
	return int(c), err
}

func (r *Repository) CountArticleBySlug(ctx context.Context, slug string) (int, error) {
	c, err := adminmodel.New(r.Db).CountArticleBySlug(ctx, slug)
	return int(c), err
}

type CreateArticleParams struct {
	Title    string
	Author   string
	Draft    bool
	LayoutID int64
}

func (r *Repository) CreateArticle(ctx context.Context, a CreateArticleParams) (string, error) {
	var draft int64
	if a.Draft {
		draft = 1
	}
	slug := slugify(a.Title)
	_, err := adminmodel.New(r.Db).CreateArticle(ctx, adminmodel.CreateArticleParams{
		Title:    a.Title,
		Date:     r.Clock.Now(),
		Author:   a.Author,
		Slug:     slug,
		Draft:    draft,
		LayoutID: a.LayoutID,
	})
	if err != nil {
		return "", stacktrace.From(err)
	}
	return slug, nil
}

func (r *Repository) DeleteArticle(ctx context.Context, slug string) error {
	err := adminmodel.New(r.Db).DeleteArticle(ctx, slug)
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}

func (r *Repository) GetArticleList(ctx context.Context, limit, offset int) ([]Article, error) {
	list, err := adminmodel.New(r.Db).GetArticleList(ctx, adminmodel.GetArticleListParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		return nil, stacktrace.From(err)
	}
	return sqlutil.Map(list, func(a adminmodel.GetArticleListRow) Article {
		return Article{
			Title:  a.Title,
			Date:   a.Date,
			Author: a.Author,
			Slug:   a.Slug,
			Draft:  a.Draft == 1,
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
			ID:       a.ID,
			Title:    a.Title,
			Date:     a.Date,
			Author:   a.Author,
			Content:  a.Content,
			Slug:     a.Slug,
			Draft:    a.Draft == 1,
			LayoutID: a.LayoutID,
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdateArticle(ctx context.Context, slug string, a Article) error {
	var draft int64
	if a.Draft {
		draft = 1
	}
	err := adminmodel.New(r.Db).UpdateArticle(ctx, adminmodel.UpdateArticleParams{
		Title:    a.Title,
		Date:     r.Clock.Now(),
		Author:   a.Author,
		Content:  a.Content,
		Slug:     a.Slug,
		Draft:    draft,
		Slug_2:   slug,
		LayoutID: a.LayoutID,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
