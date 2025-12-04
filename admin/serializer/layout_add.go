package serializer

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type LayoutAdded struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func LayoutCreated(w http.ResponseWriter, l LayoutAdded) error {
	return httpresponse.Created(w, l)
}
