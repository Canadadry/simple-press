package view

import (
	"io"
)

type PageAddData struct {
	Name string
}

type PageAddError struct {
	Name string
}

func PageAdd(a PageAddData, errors PageAddError) ViewFunc {
	type viewData struct {
		Page   PageAddData
		Errors PageAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/page_add.html",
			TemplateData("LAYOUT_ADD.page_title", viewData{a, errors}),
		)
	}
}
