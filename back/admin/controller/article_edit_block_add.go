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

	a, errors, err := form.ParseArticleEditBlockAdd(r, c.Repository.CountBlockByID)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}
	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	block, ok, err := c.Repository.SelectBlockByID(r.Context(), a.AddedBlockID)
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}

	id, err := c.Repository.CreateBlockData(r.Context(), repository.CreateBlockDataParams{
		ArticleID: article.ID,
		Block:     repository.Block{ID: a.AddedBlockID, Definition: block.Definition},
		Position:  int64(a.Position.V),
	})
	if err != nil {
		return fmt.Errorf("cannot add block %v to article : %w", a.AddedBlockID, err)
	}
	return view.BlockDataAddCreated(w, view.ArticleAddBlockData{
		ID:       id,
		Name:     block.Name,
		Data:     block.Definition,
		Position: a.Position.V,
	})
}
