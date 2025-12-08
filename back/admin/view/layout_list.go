package view

import (
	"app/pkg/http/httpresponse"
	"app/pkg/paginator"
	"io"
	"math"
	"net/http"
)

const (
	MaxLayoutPaginationItem = 5
)

type LayoutsListData struct {
	Items []LayoutListData `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type LayoutListData struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func LayoutsList(data LayoutsListData) ViewFunc {
	type viewData struct {
		Total   int
		Pages   paginator.Pages
		Layouts []LayoutListData
	}
	return func(w io.Writer, tr func(string) string) error {
		lastLayout := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastLayout, MaxLayoutPaginationItem, "/pages?page=%page%")
		return render(w, tr,
			"template/pages/layout_list.html",
			TemplateData("LAYOUT_LIST.page_title", viewData{
				Total:   data.Total,
				Pages:   p,
				Layouts: data.Items,
			}),
		)
	}
}

func LayoutsListOk(w http.ResponseWriter, a LayoutsListData) error {
	return httpresponse.Ok(w, a)
}
