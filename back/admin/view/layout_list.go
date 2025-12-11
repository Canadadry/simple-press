package view

import (
	"app/pkg/http/httpresponse"
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

func LayoutsListOk(w http.ResponseWriter, a LayoutsListData) error {
	return httpresponse.Ok(w, a)
}
