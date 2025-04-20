package view

import (
	"io"
)

type TemplateEditData struct {
	Name    string
	Content string
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
