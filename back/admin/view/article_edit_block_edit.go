package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

func BlockDataEditOk(w http.ResponseWriter, a BlockData) error {
	return httpresponse.Ok(w, a)
}
