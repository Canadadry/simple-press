package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type TemplateAddData struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type TemplateAddError struct {
	Name string
}

func TemplateAdd(a TemplateAddData, errors TemplateAddError) ViewFunc {
	type viewData struct {
		Template TemplateAddData
		Errors   TemplateAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/template_add.html",
			TemplateData("TEMPLATE_ADD.page_title", viewData{a, errors}),
		)
	}
}

func TemplateCreated(w http.ResponseWriter, l TemplateAddData) error {
	return httpresponse.Created(w, l)
}
