package view

import (
	"app/pkg/http/httpresponse"
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

func TemplatesListOk(w http.ResponseWriter, a TemplatesListData) error {
	return httpresponse.Ok(w, a)
}
