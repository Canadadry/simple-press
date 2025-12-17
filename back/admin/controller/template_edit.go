package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetTemplateEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	l, ok, err := c.Repository.SelectTemplate(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select Template : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/templates", http.StatusSeeOther)
		return nil
	}

	return view.TemplateOk(w, view.TemplateEditData{
		Name:    l.Name,
		Content: l.Content,
	})
}

func (c *Controller) PostTemplateEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	template, ok, err := c.Repository.SelectTemplate(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select template : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	l, errors, err := form.ParseTemplateEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	template.Name = l.Name
	template.Content = l.Content

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}
	err = c.Repository.UpdateTemplate(r.Context(), name, template)
	if err != nil {
		return fmt.Errorf("cannot update %s template : %w", name, err)
	}
	return view.TemplateOk(w, view.TemplateEditData{
		Name:    l.Name,
		Content: l.Content,
	})
}
