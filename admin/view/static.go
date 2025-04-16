package view

import (
	"app/pkg/flash"
	"io"
)

func PageNotFound(firstname string) func(w io.Writer, tr func(string) string, msg flash.Message) error {
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		return render(w, tr, "template/static/page_not_found.html", map[string]any{
			"DisplayTopLink": firstname != "",
			"Name":           firstname,
		})
	}
}

func InternalServerError(w io.Writer, tr func(string) string, msg flash.Message) error {
	return renderStatic(w, tr, "template/static/internal_server_error.html")
}
