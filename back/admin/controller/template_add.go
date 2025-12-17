package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) PostTemplateAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseTemplateAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	id, err := c.Repository.CreateTemplate(r.Context(), repository.CreateTemplateParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Template : %w", err)
	}

	return view.TemplateCreated(w, view.TemplateAddData{
		Name: l.Name,
		ID:   id,
	})
}
