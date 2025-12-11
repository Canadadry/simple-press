package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type TemplateAddData struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func TemplateCreated(w http.ResponseWriter, l TemplateAddData) error {
	return httpresponse.Created(w, l)
}
