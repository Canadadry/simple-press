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

	if IsJsonRequest(r) {
		return view.TemplateOk(w, view.TemplateEditData{
			Name:    l.Name,
			Content: l.Content,
		})
	}
	return c.render(w, r, view.TemplateEdit(view.TemplateEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.TemplateEditError{}))
}

func (c *Controller) PostTemplateEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	template, ok, err := c.Repository.SelectTemplate(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select template : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/templates", http.StatusSeeOther)
	}

	l, errors, err := form.ParseTemplateEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	template.Name = l.Name
	template.Content = l.Content

	if !errors.HasError() {

		err := c.Repository.UpdateTemplate(r.Context(), name, template)
		if err != nil {
			return fmt.Errorf("cannot update %s template : %w", name, err)
		}
	} else if IsJsonRequest(r) {
		return httpresponse.BadRequest(w, errors.Raw)
	}
	if IsJsonRequest(r) {
		return view.TemplateOk(w, view.TemplateEditData{
			Name:    l.Name,
			Content: l.Content,
		})
	}
	return c.render(w, r, view.TemplateEdit(view.TemplateEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.TemplateEditError{
		Name:    errors.Name,
		Content: errors.Content,
	}))
}
