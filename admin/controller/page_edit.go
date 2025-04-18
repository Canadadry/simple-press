package controller

import (
	"app/admin/form"
	"app/admin/view"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetPageEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	l, ok, err := c.Repository.SelectPage(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select Page : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/pages", http.StatusSeeOther)
	}

	return c.render(w, r, view.PageEdit(view.PageEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.PageEditError{}))
}

func (c *Controller) PostPageEdit(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	page, ok, err := c.Repository.SelectPage(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select page : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/pages", http.StatusSeeOther)
	}

	l, errors, err := form.ParsePageEdit(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	page.Name = l.Name
	page.Content = l.Content

	if !errors.HasError() {
		err := c.Repository.UpdatePage(r.Context(), name, page)
		if err != nil {
			return fmt.Errorf("cannot update %s page : %w", name, err)
		}
	}

	return c.render(w, r, view.PageEdit(view.PageEditData{
		Name:    l.Name,
		Content: l.Content,
	}, view.PageEditError(errors)))
}
