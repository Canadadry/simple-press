package view

import (
	"app/pkg/flash"
	"app/pkg/paginator"
	"io"
	"math"
)

const (
	MaxMayoutPaginationItem = 5
)

type LayoutsListData struct {
	Layouts []LayoutListData
	Total   int
	Page    int
	Limit   int
}

type LayoutListData struct {
	Name string
}

func LayoutsList(data LayoutsListData) ViewFunc {
	type viewData struct {
		Total   int
		Pages   paginator.Pages
		Layouts []LayoutListData
	}
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		lastPage := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastPage, MaxMayoutPaginationItem, "/layouts?page=%page%")
		return render(w, tr,
			"template/pages/layout_list.tmpl",
			TemplateData(msg, viewData{
				Total:   data.Total,
				Pages:   p,
				Layouts: data.Layouts,
			}),
		)
	}
}
