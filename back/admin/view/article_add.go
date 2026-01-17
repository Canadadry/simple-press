package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type ArticleAddData struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Draft  bool   `json:"draft"`
	Slug   string `json:"slug"`
}

type ArticleAddError struct {
	Title  string
	Author string
}

func ArticleCreated(w http.ResponseWriter, a ArticleAddData) error {
	return httpresponse.Created(w, a)
}
