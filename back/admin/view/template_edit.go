package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type TemplateEditData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type TemplateEditError struct {
	Name    string
	Content string
}

func TemplateOk(w http.ResponseWriter, a TemplateEditData) error {
	return httpresponse.Ok(w, a)
}
