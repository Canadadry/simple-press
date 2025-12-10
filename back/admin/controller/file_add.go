package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
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
	fmt.Println("PostFileAdd", errors, err)

	if errors.HasError() {
		if IsJsonRequest(r) {
			return httpresponse.BadRequest(w, errors.Raw)
		}

		return c.render(w, r, view.FileAdd(
			view.FileAddError{Content: errors.Content},
		))
	}

	id, err := c.Repository.UploadFile(r.Context(), repository.File{Name: l.Name, Content: l.Content})
	if err != nil {
		return fmt.Errorf("cannot create File : %w", err)
	}

	if IsJsonRequest(r) {
		return view.FileAddCreated(w, view.FileAddData{ID: id, Name: l.Name})
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
	return nil
}
