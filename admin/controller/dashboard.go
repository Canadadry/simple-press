package controller

import (
	"app/admin/view"
	"net/http"
)

func (c *Controller) GetDashboard(w http.ResponseWriter, r *http.Request) error {
	return c.render(w, r, view.Dashboard())
}
