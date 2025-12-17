package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetLayoutEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	l, ok, err := c.Repository.SelectLayout(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select Layout : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	return view.LayoutOk(w, view.LayoutEditData{
		Name:    l.Name,
		Content: l.Content,
	})

}

func (c *Controller) PostLayoutEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	layout, ok, err := c.Repository.SelectLayout(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select layout : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	l, errors, err := form.ParseLayoutEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	layout.Name = l.Name
	layout.Content = l.Content

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}

	err = c.Repository.UpdateLayout(r.Context(), name, layout)
	if err != nil {
		return fmt.Errorf("cannot update %s layout : %w", name, err)
	}

	return view.LayoutOk(w, view.LayoutEditData{
		Name:    l.Name,
		Content: l.Content,
	})
}
