package controller

import (
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) DeleteFile(w http.ResponseWriter, r *http.Request) error {
	name := router.GetField(r, 0)
	err := c.Repository.DeleteFile(r.Context(), name)
	if err != nil {
		return fmt.Errorf("cannot select file : %w", err)
	}

	return httpresponse.NoContent(w)
}
