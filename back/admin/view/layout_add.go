package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type LayoutAddData struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func LayoutCreated(w http.ResponseWriter, l LayoutAddData) error {
	return httpresponse.Created(w, l)
}
