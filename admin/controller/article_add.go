package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticleAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.ArticleAdd(view.ArticleAddData{}, view.ArticleAddError{}))
}

func (c *Controller) PostArticleAdd(w http.ResponseWriter, r *http.Request) error {

	a, errors, err := form.ParseArticleAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return c.render(w, r, view.ArticleAdd(view.ArticleAddData(a), view.ArticleAddError(errors)))
	}

	slug, err := c.Repository.CreateArticle(r.Context(), repository.CreateArticleParams(a))
	if err != nil {
		return fmt.Errorf("cannot create article : %w", err)
	}

	http.Redirect(w, r, "/admin/articles/"+slug+"/edit", http.StatusSeeOther)
	return nil
}
