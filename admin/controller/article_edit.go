package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
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

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), a.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
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

	if IsJsonRequest(r) {
		return view.ArticleOk(w, view.ArticleEditData{
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

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:      a.Title,
		Author:     a.Author,
		Slug:       a.Slug,
		Content:    a.Content,
		Draft:      a.Draft,
		LayoutID:   a.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	}, view.ArticleEditError{}))
}

func (c *Controller) PostArticleEditMetadata(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
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

	if !errors.HasError() {
		err = c.Repository.UpdateArticle(r.Context(), slug, article)
		if err != nil {
			return fmt.Errorf("cannot update %s article : %w", slug, err)
		}
	} else if IsJsonRequest(r) {
		return httpresponse.BadRequest(w, errors.Raw)
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
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
	}
	if IsJsonRequest(r) {
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

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	}, view.ArticleEditError{
		Title:    errors.Title,
		LayoutID: errors.LayoutID,
		Author:   errors.Author,
		Slug:     errors.Slug,
	}))
}

func (c *Controller) PostArticleEditContent(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	a, errors, err := form.ParseArticleEditContent(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	article.Content = a.Content

	if !errors.HasError() {
		err := c.Repository.UpdateArticle(r.Context(), slug, article)
		if err != nil {
			return fmt.Errorf("cannot update %s article : %w", slug, err)
		}
	} else if IsJsonRequest(r) {
		return httpresponse.BadRequest(w, errors.Raw)
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
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
	}
	if IsJsonRequest(r) {
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

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	}, view.ArticleEditError{
		Content: errors.Content,
	}))
}

func (c *Controller) PostArticleEditBlockEdit(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
	}

	a, errors, err := form.ParseArticleEditBlockEdit(r, func(id int64) (map[string]any, bool) {
		for _, bd := range blockDatas {
			if bd.ID == id {
				return bd.Data, true
			}
		}
		return nil, false
	})
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	if !errors.HasError() {
		err := c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
			ID:       a.EditedBlockID,
			Data:     a.EditedBlockData,
			Position: int64(a.EditedBlockPosition),
		})
		if err != nil {
			return fmt.Errorf("cannot add block %v to article : %w", a.EditedBlockID, err)
		}
		for i, p := range blockDataView {
			if p.ID == a.EditedBlockID {
				blockDataView[i].Data = a.EditedBlockData
			}
		}
	} else if IsJsonRequest(r) {
		return httpresponse.BadRequest(w, errors.Raw)
	}
	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}
	if IsJsonRequest(r) {
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

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	}, view.ArticleEditError{
		EditedBlockID:       errors.EditedBlockID,
		EditedBlockData:     errors.EditedBlockData,
		EditedBlockPosition: errors.EditedBlockPosition,
	}))
}

func (c *Controller) PostArticleEditBlockAdd(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
	}

	a, errors, err := form.ParseArticleEditBlockAdd(r, c.Repository.CountBlockByID)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	if !errors.HasError() {
		def := map[string]any{}
		for _, b := range blocks {
			if b.ID == a.AddedBlockID {
				def = b.Definition
			}
		}
		id, err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
			ArticleID: article.ID,
			Block:     repository.Block{ID: a.AddedBlockID, Definition: def},
			Position:  0,
		})
		if err != nil {
			return fmt.Errorf("cannot add block %v to article : %w", a.AddedBlockID, err)
		}

		blockDataView = append(blockDataView, view.BlockData{ID: id, Data: def})
	} else if IsJsonRequest(r) {
		return httpresponse.BadRequest(w, errors.Raw)
	}
	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
	}
	if IsJsonRequest(r) {
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

	return c.render(w, r, view.ArticleEdit(view.ArticleEditData{
		Title:      article.Title,
		Author:     article.Author,
		Slug:       article.Slug,
		Content:    article.Content,
		Draft:      article.Draft,
		LayoutID:   article.LayoutID,
		Layouts:    layoutSelector,
		Blocks:     blockSelector,
		BlockDatas: blockDataView,
	}, view.ArticleEditError{
		AddedBlockID: errors.AddedBlockID,
	}))
}
