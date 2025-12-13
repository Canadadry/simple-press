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

func (c *Controller) PostArticleEditBlockEdit(w http.ResponseWriter, r *http.Request) error {
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

	blockData := repository.BlockData{}
	is_valid := func(id int64) bool {
		for _, bd := range blockDatas {
			if bd.ID == id {
				blockData = bd
				return true
			}
		}
		return false
	}
	get_previous_data := func() map[string]any {
		return blockData.Data
	}
	a, errors, err := form.ParseArticleEditBlockEdit(r, get_previous_data, is_valid)
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
	err = c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
		ID:       a.EditedBlockID,
		Data:     a.EditedBlockData,
		Position: blockData.Position, //int64(a.EditedBlockPosition),
	})
	if err != nil {
		return fmt.Errorf("cannot add block %v to article : %w", a.EditedBlockID, err)
	}
	for i, p := range blockDataView {
		if p.ID == a.EditedBlockID {
			blockDataView[i].Data = a.EditedBlockData
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
