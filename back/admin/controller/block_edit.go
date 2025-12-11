package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetBlockEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	l, ok, err := c.Repository.SelectBlock(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select Block : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}
	return view.BlockOk(w, view.BlockEditData{
		Name:       l.Name,
		Content:    l.Content,
		Definition: l.Definition,
	})
}

func (c *Controller) PostBlockEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	block, ok, err := c.Repository.SelectBlock(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select block : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	b, errors, err := form.ParseBlockEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	block.Name = b.Name
	block.Content = b.Content
	block.Definition = b.Definition

	if errors.HasError() {
		return httpresponse.BadRequest(w, errors.Raw)
	}

	err = c.Repository.UpdateBlock(r.Context(), name, block)
	if err != nil {
		return fmt.Errorf("cannot update %s block : %w", name, err)
	}

	return view.BlockOk(w, view.BlockEditData{
		Name:       b.Name,
		Content:    b.Content,
		Definition: b.Definition,
	})

}
