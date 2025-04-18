package view

import (
	"io"
)

type PageEditData struct {
	Name    string
	Content string
}

type PageEditError struct {
	Name    string
	Content string
}

func PageEdit(a PageEditData, errors PageEditError) ViewFunc {
	type viewData struct {
		Page   PageEditData
		Errors PageEditError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/page_edit.html",
			TemplateData("PAGE_EDIT.page_title", viewData{a, errors}),
		)
	}
}
