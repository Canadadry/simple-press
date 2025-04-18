package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"fmt"
	"net/http"
)

func (c *Controller) GetPageAdd(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.PageAdd(view.PageAddData{}, view.PageAddError{}))
}

func (c *Controller) PostPageAdd(w http.ResponseWriter, r *http.Request) error {

	l, errors, err := form.ParsePageAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError() {
		return c.render(w, r, view.PageAdd(view.PageAddData(l), view.PageAddError(errors)))
	}

	err = c.Repository.CreatePage(r.Context(), repository.CreatePageParams(l))
	if err != nil {
		return fmt.Errorf("cannot create Page : %w", err)
	}

	http.Redirect(w, r, "/admin/pages/"+l.Name+"/edit", http.StatusSeeOther)
	return nil
}
