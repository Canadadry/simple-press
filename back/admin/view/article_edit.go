package view

import (
	"app/pkg/http/httpresponse"
	"io"
	"net/http"
)

type LayoutSelector struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type ArticleEditData struct {
	Title      string           `json:"title"`
	Author     string           `json:"author"`
	Slug       string           `json:"slug"`
	Content    string           `json:"content"`
	Draft      bool             `json:"draft"`
	LayoutID   int64            `json:"layout_id"`
	Layouts    []LayoutSelector `json:"layouts"`
	Blocks     []LayoutSelector `json:"blocks"`
	BlockDatas []BlockData      `json:"block_datas"`
}

type BlockData struct {
	ID   int64          `json:"id"`
	Data map[string]any `json:"data"`
}

type ArticleEditError struct {
	Title               string
	Author              string
	Slug                string
	Content             string
	LayoutID            string
	EditedBlockID       string
	EditedBlockData     string
	EditedBlockPosition string
	AddedBlockID        string
	Action              string
	Form                string
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

func ArticleOk(w http.ResponseWriter, a ArticleEditData) error {
	return httpresponse.Ok(w, a)
}
