package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type LayoutAddData struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type LayoutAddError struct {
	Name string
}

func LayoutAdd(a LayoutAddData, errors LayoutAddError) ViewFunc {
	type viewData struct {
		Layout LayoutAddData
		Errors LayoutAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/layout_add.html",
			TemplateData("LAYOUT_ADD.page_title", viewData{a, errors}),
		)
	}
}
func LayoutCreated(w http.ResponseWriter, l LayoutAddData) error {
	return httpresponse.Created(w, l)
}
