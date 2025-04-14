package view

import (
	"app/pkg/flash"
	"app/pkg/paginator"
	"io"
	"math"
	"time"
)

const (
	MaxPaginationItem = 5
)

type ArticlesListData struct {
	Articles []ArticleData
	Total    int
	Page     int
	Limit    int
}

type ArticleData struct {
	Title  string
	Date   time.Time
	Author string
	Slug   string
	Draft  bool
}

func ArticlesList(data ArticlesListData) ViewFunc {
	type viewData struct {
		Total    int
		Pages    paginator.Pages
		Articles []ArticleData
	}
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		return render(w, tr,
			"template/pages/article_list.tmpl",
			TemplateData(msg, viewData{
				Total:    data.Total,
				Pages:    paginator.New(data.Page, int(math.Ceil(float64(data.Total)/float64(data.Limit))), MaxPaginationItem, "/articles?page=%page%"),
				Articles: data.Articles,
			}),
		)
	}
}
