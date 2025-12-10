package view

import (
	"app/pkg/http/httpresponse"
	"app/pkg/paginator"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	MaxPaginationItem = 5
)

type ArticlesListData struct {
	Items []ArticleListData `json:"items"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type ArticleListData struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Author  string    `json:"author"`
	Slug    string    `json:"slug"`
	Draft   bool      `json:"draft"`
	Content string    `json:"content"`
}

func ArticlesList(data ArticlesListData) ViewFunc {
	type viewData struct {
		Total    int
		Pages    paginator.Pages
		Articles []ArticleListData
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/article_list.html",
			TemplateData("ARTICLE_LIST.page_title", viewData{
				Total:    data.Total,
				Pages:    paginator.New(data.Page, int(math.Ceil(float64(data.Total)/float64(data.Limit))), MaxPaginationItem, "/articles?page=%page%"),
				Articles: data.Items,
			}),
		)
	}
}

func ArticlesListOk(w http.ResponseWriter, a ArticlesListData) error {
	return httpresponse.Ok(w, a)
}
