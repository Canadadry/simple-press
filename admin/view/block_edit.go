package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type BlockEditData struct {
	Name       string         `json:"name"`
	Content    string         `json:"content"`
	Definition map[string]any `json:"definition"`
}

type BlockEditError struct {
	Name       string
	Content    string
	Definition string
}

func BlockEdit(a BlockEditData, errors BlockEditError) ViewFunc {
	type viewData struct {
		Block  BlockEditData
		Errors BlockEditError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/block_edit.html",
			TemplateData("BLOCK_EDIT.page_title", viewData{a, errors}),
		)
	}
}
func BlockOk(w http.ResponseWriter, a BlockEditData) error {
	return httpresponse.Ok(w, a)
}
