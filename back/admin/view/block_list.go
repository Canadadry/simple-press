package view

import (
	"app/pkg/http/httpresponse"
	"app/pkg/paginator"
	"io"
	"math"
	"net/http"
)

const (
	MaxBlockPaginationItem = 5
)

type BlocksListData struct {
	Items []BlockListData `json:"items"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type BlockListData struct {
	Name string `json:"name"`
}

func BlocksList(data BlocksListData) ViewFunc {
	type viewData struct {
		Total  int
		Pages  paginator.Pages
		Blocks []BlockListData
	}
	return func(w io.Writer, tr func(string) string) error {
		lastPage := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastPage, MaxBlockPaginationItem, "/block?page=%page%")
		return render(w, tr,
			"template/pages/block_list.html",
			TemplateData("BLOCK_LIST.page_title", viewData{
				Total:  data.Total,
				Pages:  p,
				Blocks: data.Items,
			}),
		)
	}
}

func BlocksListOk(w http.ResponseWriter, a BlocksListData) error {
	return httpresponse.Ok(w, a)
}
