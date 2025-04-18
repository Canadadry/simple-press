package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetPageList(w http.ResponseWriter, r *http.Request) error {
	const MinPage = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountPages(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count pages : %w", err)
	}

	list, err := c.Repository.GetPageList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list pages : %w", err)
	}

	if len(list) == 0 && count > 0 {
		http.Redirect(w, r, "/admin/pages", http.StatusFound)
		return nil
	}

	pages := []view.PageListData{}
	for _, a := range list {
		pages = append(pages, view.PageListData{
			Name: a.Name,
		})
	}

	l := view.PagesListData{
		Total: count,
		Limit: limit,
		Page:  page,
		Pages: pages,
	}

	return c.render(w, r, view.PagesList(l))
}
