package controller

import (
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
	"strings"
)

func (c *Controller) TreeFile(w http.ResponseWriter, r *http.Request) error {
	path := router.GetField(r, 0)
	if strings.Contains(path, "..") {
		return httpresponse.NotFound(w)
	}
	query := path
	if len(query) > 0 && query[len(query)-1] != '/' {
		query = query + "/"
	}
	files, folders, err := c.Repository.SelectFileTree(r.Context(), query)
	if err != nil {
		return fmt.Errorf("cannot fetch file tree : %w", err)
	}
	if len(files) == 0 && len(folders) == 0 {
		return httpresponse.NotFound(w)
	}
	return view.FilesTreeOk(w, view.FileTreeData{
		Path:    path,
		Files:   files,
		Folders: folders,
	})
}
