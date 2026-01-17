package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticle(w http.ResponseWriter, r *http.Request) error {
	id, _ := router.GetFieldAsInt(r, 0)
	a, ok, err := c.Repository.SelectArticleByID(r.Context(), int64(id))
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), a.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{
			ID:       p.ID,
			Name:     p.BlockName,
			Data:     p.Data,
			Position: int(p.Position),
		})
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

	return view.ArticleOk(w, view.ArticleEditData{
		ID:         int(a.ID),
		Title:      a.Title,
		Author:     a.Author,
		Slug:       a.Slug,
		Content:    a.Content,
		Draft:      a.Draft,
		LayoutID:   a.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	})
}

func (c *Controller) PostArticleEditMetadata(w http.ResponseWriter, r *http.Request) error {
	id, _ := router.GetFieldAsInt(r, 0)
	article, ok, err := c.Repository.SelectArticleByID(r.Context(), int64(id))
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	a, errors, err := form.ParseArticleEditMetadata(r, c.Repository.CountLayoutByID)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	article.Title = a.Title
	article.Author = a.Author
	article.Slug = a.Slug
	article.Draft = a.Draft.V
	article.LayoutID = a.LayoutID

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}
	err = c.Repository.UpdateArticle(r.Context(), article.Slug, article)
	if err != nil {
		return fmt.Errorf("cannot update %s article : %w", article.Slug, err)
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

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, bd := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{
			ID:       bd.ID,
			Name:     bd.BlockName,
			Data:     bd.Data,
			Position: int(bd.Position),
		})
	}
	return view.ArticleOk(w, view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	})
}

func (c *Controller) PostArticleEditContent(w http.ResponseWriter, r *http.Request) error {
	id, _ := router.GetFieldAsInt(r, 0)
	article, ok, err := c.Repository.SelectArticleByID(r.Context(), int64(id))
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	a, errors, err := form.ParseArticleEditContent(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	article.Content = a.Content

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	err = c.Repository.UpdateArticle(r.Context(), article.Slug, article)
	if err != nil {
		return fmt.Errorf("cannot update %s article : %w", article.Slug, err)
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

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Name: p.BlockName, Data: p.Data})
	}
	return view.ArticleOk(w, view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	})

}
