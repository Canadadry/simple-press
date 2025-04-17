package view

import (
	"io"
)

func PageNotFound(w io.Writer, tr func(string) string) error {
	return renderStatic(w, tr, "template/static/page_not_found.html", "404.page_title")
}

func InternalServerError(w io.Writer, tr func(string) string) error {
	return renderStatic(w, tr, "template/static/internal_server_error.html", "500.page_title")
}
