package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetFileList(w http.ResponseWriter, r *http.Request) error {
	const MinPage = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountFiles(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count file : %w", err)
	}

	list, err := c.Repository.GetFileList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list files : %w", err)
	}

	if len(list) == 0 && count > 0 {
		http.Redirect(w, r, "/admin/templates", http.StatusFound)
		return nil
	}

	files := []view.FileListData{}
	for _, a := range list {
		files = append(files, view.FileListData{
			Name: a.Name,
		})
	}

	l := view.FilesListData{
		Total: count,
		Limit: limit,
		Page:  page,
		Files: files,
	}

	return c.render(w, r, view.FilesList(l))
}
