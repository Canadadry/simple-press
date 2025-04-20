package view

import (
	"io"
)

type TemplateAddData struct {
	Name string
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
			"template/pages/layout_add.html",
			TemplateData("LAYOUT_ADD.page_title", viewData{a, errors}),
		)
	}
}
