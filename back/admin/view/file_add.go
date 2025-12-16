package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type FileAddListData struct {
	Items []FileAddData `json:"items"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type FileAddData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func FileAddCreated(w http.ResponseWriter, fa FileAddData) error {
	return httpresponse.Created(w, fa)
}

func FileAddListCreated(w http.ResponseWriter, fa FileAddListData) error {
	return httpresponse.Created(w, fa)
}
