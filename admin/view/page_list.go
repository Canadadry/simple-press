package view

import (
	"app/pkg/paginator"
	"io"
	"math"
)

const (
	MaxPagePaginationItem = 5
)

type PagesListData struct {
	Pages []PageListData
	Total int
	Page  int
	Limit int
}

type PageListData struct {
	Name string
}

func PagesList(data PagesListData) ViewFunc {
	type viewData struct {
		Total int
		Pages paginator.Pages
		Data  []PageListData
	}
	return func(w io.Writer, tr func(string) string) error {
		lastPage := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastPage, MaxPagePaginationItem, "/pages?page=%page%")
		return render(w, tr,
			"template/pages/page_list.html",
			TemplateData("LAYOUT_LIST.page_title", viewData{
				Total: data.Total,
				Pages: p,
				Data:  data.Pages,
			}),
		)
	}
}
