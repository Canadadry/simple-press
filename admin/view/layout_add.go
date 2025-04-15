package view

import (
	"app/pkg/flash"
	"io"
)

type LayoutAddData struct {
	Name string
}

type LayoutAddError struct {
	Name string
}

func LayoutAdd(a LayoutAddData, errors LayoutAddError) ViewFunc {
	type viewData struct {
		Layout LayoutAddData
		Errors LayoutAddError
	}
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		return render(w, tr,
			"template/pages/layout_add.tmpl",
			TemplateData(msg, viewData{a, errors}),
		)
	}
}
