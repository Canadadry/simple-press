package view

import (
	"io"
)

type BlockAddData struct {
	Name string
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
