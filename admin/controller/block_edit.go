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
		http.Redirect(w, r, "/admin/block", http.StatusSeeOther)
	}
	if IsJsonRequest(r) {
		return view.BlockOk(w, view.BlockEditData{
			Name:       l.Name,
			Content:    l.Content,
			Definition: l.Definition,
		})
	}

	return c.render(w, r, view.BlockEdit(view.BlockEditData{
		Name:       l.Name,
		Content:    l.Content,
		Definition: l.Definition,
	}, view.BlockEditError{}))
}

func (c *Controller) PostBlockEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	block, ok, err := c.Repository.SelectBlock(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select block : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/block", http.StatusSeeOther)
	}

	b, errors, err := form.ParseBlockEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	block.Name = b.Name
	block.Content = b.Content
	block.Definition = b.Definition

	if !errors.HasError() {
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}

		err := c.Repository.UpdateBlock(r.Context(), name, block)
		if err != nil {
			return fmt.Errorf("cannot update %s block : %w", name, err)
		}
	}

	return c.render(w, r, view.BlockEdit(view.BlockEditData{
		Name:       b.Name,
		Content:    b.Content,
		Definition: b.Definition,
	}, view.BlockEditError{
		Name:       errors.Name,
		Content:    errors.Content,
		Definition: errors.Definition,
	}))
}
