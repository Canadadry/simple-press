package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type FileAddData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func FileAddCreated(w http.ResponseWriter, fa FileAddData) error {
	return httpresponse.Created(w, fa)
}
