package controller

import (
	"app/admin/view"
	"app/pkg/paginator"
	"fmt"
	"net/http"
)

func (c *Controller) GetTemplateList(w http.ResponseWriter, r *http.Request) error {
	const MinPage = 0
	const MinLimit = 10
	page := paginator.PageFromRequest(r, "page", MinPage)
	limit := paginator.PageFromRequest(r, "limit", MinLimit)

	count, err := c.Repository.CountTemplates(r.Context())
	if err != nil {
		return fmt.Errorf("cannot count template : %w", err)
	}

	list, err := c.Repository.GetTemplateList(r.Context(), limit, page*limit)
	if err != nil {
		return fmt.Errorf("cannot list template : %w", err)
	}

	templates := []view.TemplateListData{}
	for _, t := range list {
		templates = append(templates, view.TemplateListData{
			Name:    t.Name,
			Content: t.Content,
		})
	}

	l := view.TemplatesListData{
		Total: count,
		Limit: limit,
		Page:  page,
		Items: templates,
	}
	return view.TemplatesListOk(w, l)
}
