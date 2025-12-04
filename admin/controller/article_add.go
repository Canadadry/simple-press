package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/serializer"
	"app/admin/view"
	"app/pkg/http/httpresponse"
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
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}
		return c.render(w, r, view.ArticleAdd(view.ArticleAddData{
			Title:  a.Title,
			Author: a.Author,
			Draft:  a.Draft.V,
		}, view.ArticleAddError{
			Title:  errors.Title,
			Author: errors.Author,
		}))
	}

	layouts, err := c.Repository.GetLayoutList(r.Context(), 1, 0)
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	if len(layouts) == 0 {
		return fmt.Errorf("need at leats one layout to create an article : %w", err)
	}

	slug, err := c.Repository.CreateArticle(r.Context(), repository.CreateArticleParams{
		Title:    a.Title,
		Author:   a.Author,
		Draft:    a.Draft.V,
		LayoutID: layouts[0].ID,
	})
	if err != nil {
		return fmt.Errorf("cannot create article : %w", err)
	}
	if IsJsonRequest(r) {
		return serializer.ArticleCreated(w, serializer.ArticleAdded{
			Title:  a.Title,
			Author: a.Author,
			Draft:  a.Draft.V && a.Draft.Valid,
			Slug:   slug,
		})
	}

	http.Redirect(w, r, "/admin/articles/"+slug+"/edit", http.StatusSeeOther)
	return nil
}
