package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type BlockAddData struct {
	Name string `json:"name"`
}

func BlockCreated(w http.ResponseWriter, a BlockAddData) error {
	return httpresponse.Created(w, a)
}
