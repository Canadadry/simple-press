package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"fmt"
	"net/http"
)

func (c *Controller) GetFileAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.FileAdd(view.FileAddError{}))
}

func (c *Controller) PostFileAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParseFileAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return c.render(w, r, view.FileAdd(
			view.FileAddError{Content: errors.Content},
		))
	}

	err = c.Repository.UploadFile(r.Context(), repository.File{Name: l.Name, Content: l.Content})
	if err != nil {
		return fmt.Errorf("cannot create File : %w", err)
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
	return nil
}
