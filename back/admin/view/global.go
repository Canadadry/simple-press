package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

func GlobalDefinitionOk(w http.ResponseWriter, data map[string]any) error {
	return httpresponse.Ok(w, map[string]any{
		"definition": data,
	})
}

func GlobalDataOk(w http.ResponseWriter, data map[string]any) error {
	return httpresponse.Ok(w, map[string]any{
		"data": data,
	})
}
