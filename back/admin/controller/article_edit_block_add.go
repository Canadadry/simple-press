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

func (c *Controller) PostArticleEditBlockAdd(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	article, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), article.ID)
	if err != nil {
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []view.BlockData{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, view.BlockData{ID: p.ID, Name: p.BlockName, Data: p.Data})
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

	if errors.HasError() {
		return httpresponse.BadRequest(w, errors.Raw)
	}
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
	blockName := ""
	for _, block := range blockSelector {
		if block.Value == id {
			blockName = block.Name
		}
	}

	blockDataView = append(blockDataView, view.BlockData{ID: id, Name: blockName, Data: def})
	layouts, err := c.Repository.GetAllLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	layoutSelector := []view.LayoutSelector{}
	for _, p := range layouts {
		layoutSelector = append(layoutSelector, view.LayoutSelector{Name: p.Name, Value: p.ID})
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
