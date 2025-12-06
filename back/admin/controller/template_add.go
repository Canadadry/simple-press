package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) GetTemplateAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.TemplateAdd(view.TemplateAddData{}, view.TemplateAddError{}))
}

func (c *Controller) PostTemplateAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseTemplateAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}

		return c.render(w, r, view.TemplateAdd(view.TemplateAddData{Name: l.Name}, view.TemplateAddError{
			Name: errors.Name,
		}))
	}

	id, err := c.Repository.CreateTemplate(r.Context(), repository.CreateTemplateParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Template : %w", err)
	}

	if IsJsonRequest(r) {
		return view.TemplateCreated(w, view.TemplateAddData{
			Name: l.Name,
			ID:   id,
		})
	}

	http.Redirect(w, r, "/admin/templates/"+l.Name+"/edit", http.StatusSeeOther)
	return nil
}
