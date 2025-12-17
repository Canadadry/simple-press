package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) PostBlockAdd(w http.ResponseWriter, r *http.Request) error {

	b, errors, err := form.ParseBlockAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	err = c.Repository.CreateBlock(r.Context(), repository.CreateBlockParams(b))
	if err != nil {
		return fmt.Errorf("cannot create Block : %w", err)
	}
	return view.BlockCreated(w, view.BlockAddData{
		Name: b.Name,
	})
}
