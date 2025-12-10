package controller

import (
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (c *Controller) GetFile(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	f, ok, err := c.Repository.DownloadFile(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select file : %w", err)
	}
	if !ok {
		if IsJsonRequest(r) {
			return httpresponse.NotFound(w)
		}
		return c.renderWithStatus(w, r, http.StatusNotFound, view.PageNotFound)
	}

	return httpresponse.File(w, io.NopCloser(bytes.NewReader(f.Content)))
}
