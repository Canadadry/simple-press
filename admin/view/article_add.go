package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type ArticleAddData struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Draft  bool   `json:"draft"`
	Slug   string `json:"slug"`
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

func ArticleCreated(w http.ResponseWriter, a ArticleAddData) error {
	return httpresponse.Created(w, a)
}
