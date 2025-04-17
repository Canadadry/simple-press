package controller

import (
	"app/admin/view"
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
		return c.render(w, r, view.PageNotFound)
	}

	return httpresponse.Bytes(w, f.Content)
}
