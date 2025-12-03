package serializer

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type ArticleAdded struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Draft  bool   `json:"draft"`
}

func ArticleCreated(w http.ResponseWriter, a ArticleAdded) error {
	return httpresponse.Created(w, a)
}
