package controller

import (
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetFile(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	f, ok, err := c.Repository.DownloadFile(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select file : %w", err)
	}
	if !ok {
		http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
	}

	return httpresponse.Bytes(w, f.Content)
}
