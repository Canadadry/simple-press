package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) PostLayoutAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseLayoutAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	id, err := c.Repository.CreateLayout(r.Context(), repository.CreateLayoutParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Layout : %w", err)
	}

	return view.LayoutCreated(w, view.LayoutAddData{
		Name: l.Name,
		ID:   id,
	})

}
