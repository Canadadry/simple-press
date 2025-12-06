package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"fmt"
	"net/http"
)

func (c *Controller) GetLayoutAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.LayoutAdd(view.LayoutAddData{}, view.LayoutAddError{}))
}

func (c *Controller) PostLayoutAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseLayoutAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}
		return c.render(w, r, view.LayoutAdd(
			view.LayoutAddData{Name: l.Name},
			view.LayoutAddError{Name: errors.Name},
		))
	}

	id, err := c.Repository.CreateLayout(r.Context(), repository.CreateLayoutParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Layout : %w", err)
	}

	if IsJsonRequest(r) {
		return view.LayoutCreated(w, view.LayoutAddData{
			Name: l.Name,
			ID:   id,
		})
	}

	http.Redirect(w, r, "/admin/layouts/"+l.Name+"/edit", http.StatusSeeOther)
	return nil
}
