package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type TemplateEditData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type TemplateEditError struct {
	Name    string
	Content string
}

func TemplateEdit(a TemplateEditData, errors TemplateEditError) ViewFunc {
	type viewData struct {
		Template TemplateEditData
		Errors   TemplateEditError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/template_edit.html",
			TemplateData("TEMPLATE_EDIT.page_title", viewData{a, errors}),
		)
	}
}

func TemplateOk(w http.ResponseWriter, a TemplateEditData) error {
	return httpresponse.Ok(w, a)
}
