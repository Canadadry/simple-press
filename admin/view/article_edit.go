package view

import (
	"io"
)

type ArticleEditData struct {
	Title   string
	Author  string
	Slug    string
	Content string
	Draft   bool
}

type ArticleEditError struct {
	Title   string
	Author  string
	Slug    string
	Content string
}

func ArticleEdit(a ArticleEditData, errors ArticleEditError) ViewFunc {
	type viewData struct {
		Article ArticleEditData
		Errors  ArticleEditError
	}
	return func(w io.Writer, tr func(string) string) error {
		return render(w, tr,
			"template/pages/article_edit.html",
			TemplateData("ARTICLE_EDIT.page_title", viewData{a, errors}),
		)
	}
}
