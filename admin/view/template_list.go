package view

import (
	"app/pkg/paginator"
	"io"
	"math"
)

const (
	MaxTemplatePaginationItem = 5
)

type TemplatesListData struct {
	Templates []TemplateListData
	Total     int
	Page      int
	Limit     int
}

type TemplateListData struct {
	Name string
}

func TemplatesList(data TemplatesListData) ViewFunc {
	type viewData struct {
		Total     int
		Pages     paginator.Pages
		Templates []TemplateListData
	}
	return func(w io.Writer, tr func(string) string) error {
		lastPage := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastPage, MaxTemplatePaginationItem, "/template?page=%page%")
		return render(w, tr,
			"template/pages/template_list.html",
			TemplateData("TEMPLATE_LIST.page_title", viewData{
				Total:     data.Total,
				Pages:     p,
				Templates: data.Templates,
			}),
		)
	}
}
