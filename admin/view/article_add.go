package view

import (
	"app/pkg/flash"
	"io"
)

type ArticleAddData struct {
	Title  string
	Author string
	Draft  bool
}

type ArticleAddError struct {
	Title  string
	Author string
	Draft  string
}

func ArticleAdd(a ArticleAddData, errors ArticleAddError) ViewFunc {
	type viewData struct {
		Article ArticleAddData
		Errors  ArticleAddError
	}
	return func(w io.Writer, tr func(string) string, msg flash.Message) error {
		return render(w, tr,
			"template/pages/article_add.tmpl",
			TemplateData(msg, viewData{a, errors}),
		)
	}
}
