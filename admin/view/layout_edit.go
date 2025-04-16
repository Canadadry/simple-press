package view

import (
	"app/pkg/flash"
	"io"
)

type LayoutEditData struct {
	Name    string
	Content string
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
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		return render(w, tr,
			"template/pages/layout_edit.html",
			TemplateData(msg, viewData{a, errors}),
		)
	}
}
