package view

import (
	"app/pkg/http/httpresponse"
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

func BlocksListOk(w http.ResponseWriter, a BlocksListData) error {
	return httpresponse.Ok(w, a)
}
