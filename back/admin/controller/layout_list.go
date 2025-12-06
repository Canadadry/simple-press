package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetLayoutList(w http.ResponseWriter, r *http.Request) error {
	const MinLayout = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinLayout)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountLayout(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count layout : %w", err)
	}

	list, err := c.Repository.GetLayoutList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list layout : %w", err)
	}

	if len(list) == 0 && count > 0 {
		http.Redirect(w, r, "/admin/layouts", http.StatusFound)
		return nil
	}

	layouts := []view.LayoutListData{}
	for _, a := range list {
		layouts = append(layouts, view.LayoutListData{
			Name: a.Name,
		})
	}

	l := view.LayoutsListData{
		Total:   count,
		Limit:   limit,
		Page:    page,
		Layouts: layouts,
	}

	return c.render(w, r, view.LayoutsList(l))
}
