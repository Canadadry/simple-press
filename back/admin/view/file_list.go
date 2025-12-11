package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

const (
	MaxFilePaginationItem = 5
)

type FilesListData struct {
	Items []FileListData `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type FileListData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func FilesListOk(w http.ResponseWriter, a FilesListData) error {
	return httpresponse.Ok(w, a)
}
