package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) GetBlockAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.BlockAdd(view.BlockAddData{}, view.BlockAddError{}))
}

func (c *Controller) PostBlockAdd(w http.ResponseWriter, r *http.Request) error {

	b, errors, err := form.ParseBlockAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}

		return c.render(w, r, view.BlockAdd(
			view.BlockAddData{Name: b.Name},
			view.BlockAddError{Name: errors.Name}))
	}

	err = c.Repository.CreateBlock(r.Context(), repository.CreateBlockParams(b))
	if err != nil {
		return fmt.Errorf("cannot create Block : %w", err)
	}
	if IsJsonRequest(r) {
		return view.BlockCreated(w, view.BlockAddData{
			Name: b.Name,
		})
	}

	http.Redirect(w, r, "/admin/block/"+b.Name+"/edit", http.StatusSeeOther)
	return nil
}
