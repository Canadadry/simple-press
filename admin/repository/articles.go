package repository

import (
	"app/model/adminmodel"
	"app/pkg/sqlutil"
	"app/pkg/stacktrace"
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"
	"unicode"
)

type Article struct {
	Title   string
	Date    time.Time
	Author  string
	Slug    string
	Content string
	Draft   bool
}

func (r *Repository) CountArticles(ctx context.Context) (int, error) {
	c, err := adminmodel.New(r.Db).CountArticles(ctx)
	return int(c), err
}

func (r *Repository) CountArticlesBySlug(ctx context.Context, slug string) (int, error) {
	c, err := adminmodel.New(r.Db).CountArticlesBySlug(ctx, slug)
	return int(c), err
}

type CreateArticleParams struct {
	Title  string
	Author string
	Draft  bool
}

func (r *Repository) CreateArticle(ctx context.Context, a CreateArticleParams) (string, error) {
	slug := slugify(a.Title)
	_, err := adminmodel.New(r.Db).CreateArticle(ctx, adminmodel.CreateArticleParams{
		Title:  a.Title,
		Date:   r.Clock.Now(),
		Author: a.Author,
		Slug:   slug,
		Draft:  sql.NullInt64{Int64: 1, Valid: a.Draft},
	})
	if err != nil {
		return "", stacktrace.From(err)
	}
	return slug, nil
}

func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = removeAccents(slug)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

func removeAccents(s string) string {
	var result []rune
	for _, r := range s {
		switch r {
		case 'à', 'á', 'â', 'ã', 'ä', 'å':
			result = append(result, 'a')
		case 'è', 'é', 'ê', 'ë':
			result = append(result, 'e')
		case 'ì', 'í', 'î', 'ï':
			result = append(result, 'i')
		case 'ò', 'ó', 'ô', 'õ', 'ö':
			result = append(result, 'o')
		case 'ù', 'ú', 'û', 'ü':
			result = append(result, 'u')
		case 'ç':
			result = append(result, 'c')
		case 'ñ':
			result = append(result, 'n')
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				result = append(result, r)
			} else {
				result = append(result, ' ')
			}
		}
	}
	return string(result)
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
			Title:   a.Title,
			Date:    a.Date,
			Author:  a.Author,
			Content: a.Content,
			Slug:    a.Slug,
			Draft:   a.Draft.Valid,
		}
	}
	return fromModel(list[0]), true, nil
}

func (r *Repository) UpdateArticle(ctx context.Context, slug string, a Article) error {
	err := adminmodel.New(r.Db).UpdateArticle(ctx, adminmodel.UpdateArticleParams{
		Title:   a.Title,
		Date:    r.Clock.Now(),
		Author:  a.Author,
		Content: a.Content,
		Slug:    a.Slug,
		Draft:   sql.NullInt64{Int64: 1, Valid: a.Draft},
		Slug_2:  slug,
	})
	if err != nil {
		return stacktrace.From(err)
	}
	return nil
}
