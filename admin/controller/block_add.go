package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"fmt"
	"net/http"
)

func (c *Controller) GetBlockAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.BlockAdd(view.BlockAddData{}, view.BlockAddError{}))
}

func (c *Controller) PostBlockAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseBlockAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return c.render(w, r, view.BlockAdd(view.BlockAddData(l), view.BlockAddError(errors)))
	}

	err = c.Repository.CreateBlock(r.Context(), repository.CreateBlockParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Block : %w", err)
	}

	http.Redirect(w, r, "/admin/block/"+l.Name+"/edit", http.StatusSeeOther)
	return nil
}
