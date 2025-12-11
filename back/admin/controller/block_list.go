package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetBlockList(w http.ResponseWriter, r *http.Request) error {
	const MinPage = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountBlocks(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count block : %w", err)
	}

	list, err := c.Repository.GetBlockList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list block : %w", err)
	}

	blocks := []view.BlockListData{}
	for _, a := range list {
		blocks = append(blocks, view.BlockListData{
			Name: a.Name,
		})
	}

	l := view.BlocksListData{
		Total: count,
		Limit: limit,
		Page:  page,
		Items: blocks,
	}

	return view.BlocksListOk(w, l)
}
