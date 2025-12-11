package view

import (
	"app/pkg/http/httpresponse"
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

func ArticlesListOk(w http.ResponseWriter, a ArticlesListData) error {
	return httpresponse.Ok(w, a)
}
