package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type BlockEditData struct {
	Name       string         `json:"name"`
	Content    string         `json:"content"`
	Definition map[string]any `json:"definition"`
}

func BlockOk(w http.ResponseWriter, a BlockEditData) error {
	return httpresponse.Ok(w, a)
}
