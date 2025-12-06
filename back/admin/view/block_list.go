package view

import (
	"app/pkg/paginator"
	"io"
	"math"
)

const (
	MaxBlockPaginationItem = 5
)

type BlocksListData struct {
	Blocks []BlockListData
	Total  int
	Page   int
	Limit  int
}

type BlockListData struct {
	Name string
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
				Blocks: data.Blocks,
			}),
		)
	}
}
