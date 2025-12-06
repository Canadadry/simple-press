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

	count, err := c.Repository.CountArticle(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count articles : %w", err)
	}

	list, err := c.Repository.GetArticleList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list articles : %w", err)
	}

	if len(list) == 0 && count > 0 {
		http.Redirect(w, r, "/admin/articles", http.StatusFound)
		return nil
	}

	articles := []view.ArticleListData{}
	for _, a := range list {
		articles = append(articles, view.ArticleListData{
			Title:  a.Title,
			Date:   a.Date,
			Author: a.Author,
			Slug:   a.Slug,
			Draft:  a.Draft,
		})
	}

	l := view.ArticlesListData{
		Total:    count,
		Limit:    limit,
		Page:     page,
		Articles: articles,
	}

	if IsJsonRequest(r) {
		return view.ArticlesListOk(w, l)
	}

	return c.render(w, r, view.ArticlesList(l))
}
