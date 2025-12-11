package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type LayoutEditData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func LayoutOk(w http.ResponseWriter, a LayoutEditData) error {
	return httpresponse.Ok(w, a)
}
