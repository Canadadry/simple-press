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
	article.Draft = a.Draft
	article.LayoutID = a.LayoutID

	err = c.Repository.UpdateArticle(r.Context(), slug, article)
	if err != nil {
		return fmt.Errorf("cannot update %s article : %w", slug, err)
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

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Data: p.Data})
	}

	a, errors, err := form.ParseArticleEdit(form.ParseArticleEditParam{
		Request:       r,
		CheckLayoutID: c.Repository.CountLayoutByID,
		CheckBlockID:  c.Repository.CountBlockByID,
		GetPreviousData: func(id int64) (map[string]any, bool) {
			for _, bd := range blockDatas {
				if bd.ID == id {
					return bd.Data, true
				}
			}
			return nil, false
		},
	})
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

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	if !errors.HasError() {
		switch a.Action {
		case form.ArticleEditActionMetadata, form.ArticleEditActionContent:
			err := c.Repository.UpdateArticle(r.Context(), slug, article)
			if err != nil {
				return fmt.Errorf("cannot update %s article : %w", slug, err)
			}
		case form.ArticleEditActionBlockAdd:
			def := map[string]any{}
			for _, b := range blocks {
				if b.ID == a.BlockID {
					def = b.Definition
				}
			}
			id, err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
				ArticleID: article.ID,
				Block:     repository.Block{ID: a.BlockID, Definition: def},
				Position:  0,
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}

			blockDataView = append(blockDataView, view.BlockData{ID: id, Data: def})

		case form.ArticleEditActionBlockEdit:
			err := c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
				ID:       a.EditedBlockID,
				Data:     a.EditedBlockData,
				Position: int64(a.EditedBlockPosition),
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}
			for i, p := range blockDataView {
				if p.ID == a.EditedBlockID {
					blockDataView[i].Data = a.EditedBlockData
				}
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
	}, view.ArticleEditError(errors)))
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

	a, errors, err := form.ParseArticleEdit(form.ParseArticleEditParam{
		Request:       r,
		CheckLayoutID: c.Repository.CountLayoutByID,
		CheckBlockID:  c.Repository.CountBlockByID,
		GetPreviousData: func(id int64) (map[string]any, bool) {
			for _, bd := range blockDatas {
				if bd.ID == id {
					return bd.Data, true
				}
			}
			return nil, false
		},
	})
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

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	if !errors.HasError() {
		switch a.Action {
		case form.ArticleEditActionMetadata, form.ArticleEditActionContent:
			err := c.Repository.UpdateArticle(r.Context(), slug, article)
			if err != nil {
				return fmt.Errorf("cannot update %s article : %w", slug, err)
			}
		case form.ArticleEditActionBlockAdd:
			def := map[string]any{}
			for _, b := range blocks {
				if b.ID == a.BlockID {
					def = b.Definition
				}
			}
			id, err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
				ArticleID: article.ID,
				Block:     repository.Block{ID: a.BlockID, Definition: def},
				Position:  0,
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}

			blockDataView = append(blockDataView, view.BlockData{ID: id, Data: def})

		case form.ArticleEditActionBlockEdit:
			err := c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
				ID:       a.EditedBlockID,
				Data:     a.EditedBlockData,
				Position: int64(a.EditedBlockPosition),
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}
			for i, p := range blockDataView {
				if p.ID == a.EditedBlockID {
					blockDataView[i].Data = a.EditedBlockData
				}
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
	}, view.ArticleEditError(errors)))
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

	a, errors, err := form.ParseArticleEdit(form.ParseArticleEditParam{
		Request:       r,
		CheckLayoutID: c.Repository.CountLayoutByID,
		CheckBlockID:  c.Repository.CountBlockByID,
		GetPreviousData: func(id int64) (map[string]any, bool) {
			for _, bd := range blockDatas {
				if bd.ID == id {
					return bd.Data, true
				}
			}
			return nil, false
		},
	})
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

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := []view.LayoutSelector{}
	for _, b := range blocks {
		blockSelector = append(blockSelector, view.LayoutSelector{Name: b.Name, Value: b.ID})
	}

	if !errors.HasError() {
		switch a.Action {
		case form.ArticleEditActionMetadata, form.ArticleEditActionContent:
			err := c.Repository.UpdateArticle(r.Context(), slug, article)
			if err != nil {
				return fmt.Errorf("cannot update %s article : %w", slug, err)
			}
		case form.ArticleEditActionBlockAdd:
			def := map[string]any{}
			for _, b := range blocks {
				if b.ID == a.BlockID {
					def = b.Definition
				}
			}
			id, err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
				ArticleID: article.ID,
				Block:     repository.Block{ID: a.BlockID, Definition: def},
				Position:  0,
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}

			blockDataView = append(blockDataView, view.BlockData{ID: id, Data: def})

		case form.ArticleEditActionBlockEdit:
			err := c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
				ID:       a.EditedBlockID,
				Data:     a.EditedBlockData,
				Position: int64(a.EditedBlockPosition),
			})
			if err != nil {
				return fmt.Errorf("cannot add block %v to article : %w", a.BlockID, err)
			}
			for i, p := range blockDataView {
				if p.ID == a.EditedBlockID {
					blockDataView[i].Data = a.EditedBlockData
				}
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
	}, view.ArticleEditError(errors)))
}
