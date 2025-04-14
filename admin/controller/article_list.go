package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

const MinPage = 0
const MinLimit = 10

func (c *Controller) GetArticleList(w http.ResponseWriter, r *http.Request) error {
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountArticles(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count articles : %w", err)
	}

	list, err := c.Repository.GetArticlesList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list articles : %w", err)
	}

	if len(list) == 0 && count > 0 {
		http.Redirect(w, r, "/admin/articles", http.StatusFound)
		return nil
	}

	articles := []view.ArticleData{}
	for _, a := range list {
		articles = append(articles, view.ArticleData(a))
	}

	l := view.ArticlesListData{
		Total:    count,
		Limit:    limit,
		Page:     page,
		Articles: articles,
	}

	return c.render(w, r, view.ArticlesList(l))
}
