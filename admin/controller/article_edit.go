package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticleEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:    a.Title,
		Author:   a.Author,
		Slug:     a.Slug,
		Content:  a.Content,
		Draft:    a.Draft,
		LayoutID: a.LayoutID,
		Layouts:  layoutSelector,
		Blocks:   blockSelector,
	}, view.ArticleEditError{}))
}

func (c *Controller) PostArticleEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	a, errors, err := form.ParseArticleEdit(
		r,
		c.Repository.CountLayoutByID,
		c.Repository.CountBlockByID,
		c.Repository.CountBlockDataByID,
	)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	switch a.Action {
	case form.ArticleEditActionMetadata:
		article.Title = a.Title
		article.Author = a.Author
		article.Slug = a.Slug
		article.Draft = a.Draft
		article.LayoutID = a.LayoutID
	case form.ArticleEditActionContent:
		article.Content = a.Content
	}
	fmt.Println("will patch", a.Action, article)

	if !errors.HasError() {
		switch a.Action {
		case form.ArticleEditActionMetadata, form.ArticleEditActionContent:
			err := c.Repository.UpdateArticle(r.Context(), slug, article)
			if err != nil {
				return fmt.Errorf("cannot update %s article : %w", slug, err)
			}
		case form.ArticleEditActionBlockAdd:
			err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
				ArticleID: article.ID,
				Block:     repository.Block{ID: a.BlockID},
				Position:  0,
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}

		case form.ArticleEditActionBlockEdit:
			err := c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
				ID:       a.EditedBlockID,
				Data:     a.EditedBlockData,
				Position: int64(a.EditedBlockPosition),
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}
		}
	}
	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:    article.Title,
		Author:   article.Author,
		Slug:     article.Slug,
		Content:  article.Content,
		Draft:    article.Draft,
		LayoutID: article.LayoutID,
		Layouts:  layoutSelector,
		Blocks:   blockSelector,
	}, view.ArticleEditError(errors)))
}
