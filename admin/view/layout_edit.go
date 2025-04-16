package view

import (
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
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/layout_edit.html",
			TemplateData("LAYOUT_EDIT.page_title", viewData{a, errors}),
		)
	}
}
