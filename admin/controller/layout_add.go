package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
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
		return c.render(w, r, view.LayoutAdd(view.LayoutAddData(l), view.LayoutAddError(errors)))
	}

	err = c.Repository.CreateLayout(r.Context(), repository.CreateLayoutParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Layout : %w", err)
	}

	http.Redirect(w, r, "/admin/layout/"+l.Name+"/edit", http.StatusSeeOther)
	return nil
}
