package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type BlockAddData struct {
	Name string `json:"name"`
}

type BlockAddError struct {
	Name string
}

func BlockAdd(a BlockAddData, errors BlockAddError) ViewFunc {
	type viewData struct {
		Block  BlockAddData
		Errors BlockAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/block_add.html",
			TemplateData("BLOCK_ADD.page_title", viewData{a, errors}),
		)
	}
}

func BlockCreated(w http.ResponseWriter, a BlockAddData) error {
	return httpresponse.Created(w, a)
}
