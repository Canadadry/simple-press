package view

import (
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
}

func ArticleAdd(a ArticleAddData, errors ArticleAddError) ViewFunc {
	type viewData struct {
		Article ArticleAddData
		Errors  ArticleAddError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/article_add.html",
			TemplateData("ARTICLE_ADD.page_title", viewData{a, errors}),
		)
	}
}
