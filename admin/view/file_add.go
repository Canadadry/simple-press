package view

import (
	"io"
)

type FileAddData struct {
	Name string
}

type FileAddError struct {
	Name string
}

func FileAdd(a FileAddData, errors FileAddError) ViewFunc {
	type viewData struct {
		File   FileAddData
		Errors FileAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/file_add.html",
			TemplateData("FILE_ADD.page_title", viewData{a, errors}),
		)
	}
}
