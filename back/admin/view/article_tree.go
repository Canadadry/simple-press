package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type ArticleTreeData struct {
	Path     string   `json:"path"`
	Articles []string `json:"articles"`
	Folders  []string `json:"folders"`
}

func ArticlesTreeOk(w http.ResponseWriter, atd ArticleTreeData) error {
	return httpresponse.Ok(w, atd)
}
