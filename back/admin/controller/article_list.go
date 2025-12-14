package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticleList(w http.ResponseWriter, r *http.Request) error {
	const MinPage = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)
	query := ""

	count, err := c.Repository.CountArticleLikeTitle(r.Context(), query)
	if err != nil {
		return fmt.Errorf("cannot count articles : %w", err)
	}

	list, err := c.Repository.GetArticleList(r.Context(), query, limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list articles : %w", err)
	}

	articles := []view.ArticleListData{}
	for _, a := range list {
		articles = append(articles, view.ArticleListData{
			Title:   a.Title,
			Date:    a.Date,
			Author:  a.Author,
			Content: a.Content,
			Slug:    a.Slug,
			Draft:   a.Draft,
		})
	}

	l := view.ArticlesListData{
		Total: count,
		Limit: limit,
		Page:  page,
		Items: articles,
	}

	return view.ArticlesListOk(w, l)
}
