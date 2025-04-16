package view

import (
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
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/layout_add.html",
			TemplateData("LAYOUT_ADD.page_title", viewData{a, errors}),
		)
	}
}
