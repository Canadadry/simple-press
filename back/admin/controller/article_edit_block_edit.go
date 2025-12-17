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
	id, err := router.GetFieldAsInt(r, 0)
	if err != nil {
		return httpresponse.NotFound(w)
	}
	block, ok, err := c.Repository.SelectBlockDataByID(r.Context(), int64(id))
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	get_previous_data := func() map[string]any {
		return block.Data
	}
	a, errors, err := form.ParseArticleEditBlockEdit(r, get_previous_data)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return httpresponse.BadRequest(w, errors.Raw)
	}
	newPos := block.Position
	if a.EditedBlockPosition.Valid {
		newPos = int64(a.EditedBlockPosition.V)
	}
	err = c.Repository.UpdateBlockData(r.Context(), repository.BlockData{
		ID:       int64(id),
		Data:     a.EditedBlockData,
		Position: newPos,
	})
	if err != nil {
		return fmt.Errorf("cannot add block %v to article : %w", id, err)
	}

	return view.BlockDataEditOk(w, view.BlockData{
		ID:       int64(id),
		Name:     block.BlockName,
		Data:     a.EditedBlockData,
		Position: int(newPos),
	})
}
