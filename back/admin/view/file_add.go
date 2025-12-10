package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type FileAddData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

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

func FileAddCreated(w http.ResponseWriter, fa FileAddData) error {
	return httpresponse.Created(w, fa)
}
