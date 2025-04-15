package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GeLayoutEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	l, ok, err := c.Repository.SelectLayout(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select Layout : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/layouts", http.StatusSeeOther)
	}

	return c.render(w, r, view.LayoutEdit(view.LayoutEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.LayoutEditError{}))
}

func (c *Controller) PostLayoutEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	layout, ok, err := c.Repository.SelectLayout(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select layout : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/layouts", http.StatusSeeOther)
	}

	l, errors, err := form.ParseLayoutEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	layout.Name = l.Name
	layout.Content = l.Content

	if !errors.HasError() {
		err := c.Repository.UpdateLayout(r.Context(), name, layout)
		if err != nil {
			return fmt.Errorf("cannot update %s layout : %w", name, err)
		}
	}

	return c.render(w, r, view.LayoutEdit(view.LayoutEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.LayoutEditError(errors)))
}
