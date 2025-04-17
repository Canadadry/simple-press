package view

import (
	"io"
)

type FileAddError struct {
	Content string
}

func FileAdd(errors FileAddError) ViewFunc {
	type viewData struct {
		Errors FileAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/file_add.html",
			TemplateData("FILE_ADD.page_title", viewData{errors}),
		)
	}
}
