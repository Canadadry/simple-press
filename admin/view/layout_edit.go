package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type LayoutEditData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type LayoutEditError struct {
	Name    string
	Content string
}

func LayoutEdit(a LayoutEditData, errors LayoutEditError) ViewFunc {
	type viewData struct {
		Layout LayoutEditData
		Errors LayoutEditError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/layout_edit.html",
			TemplateData("LAYOUT_EDIT.page_title", viewData{a, errors}),
		)
	}
}

func LayoutOk(w http.ResponseWriter, a LayoutEditData) error {
	return httpresponse.Ok(w, a)
}
