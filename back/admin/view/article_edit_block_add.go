package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type ArticleAddBlockData struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Data     map[string]any `json:"data"`
	Position int            `json:"position"`
}

func BlockDataAddCreated(w http.ResponseWriter, a ArticleAddBlockData) error {
	return httpresponse.Created(w, a)
}
