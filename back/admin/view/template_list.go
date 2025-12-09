package view

import (
	"app/pkg/http/httpresponse"
	"app/pkg/paginator"
	"io"
	"math"
	"net/http"
)

const (
	MaxTemplatePaginationItem = 5
)

type TemplatesListData struct {
	Items []TemplateListData `json:"items"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type TemplateListData struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
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
				Templates: data.Items,
			}),
		)
	}
}

func TemplatesListOk(w http.ResponseWriter, a TemplatesListData) error {
	return httpresponse.Ok(w, a)
}
