package view

import (
	"io"
)

type BlockEditData struct {
	Name       string
	Content    string
	Definition map[string]any
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
