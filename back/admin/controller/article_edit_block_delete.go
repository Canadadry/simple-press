package controller

import (
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) DeleteBlockData(w http.ResponseWriter, r *http.Request) error {
	id, err := router.GetFieldAsInt(r, 0)
	if err != nil {
		return httpresponse.NoContent(w)
	}
	err = c.Repository.DeleteBlockData(r.Context(), int64(id))
	if err != nil {
		return fmt.Errorf("cannot select Block : %w", err)
	}
	return httpresponse.NoContent(w)
}
