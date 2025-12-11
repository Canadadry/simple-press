package view

import (
	"app/pkg/http/httpresponse"
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

func ArticleOk(w http.ResponseWriter, a ArticleEditData) error {
	return httpresponse.Ok(w, a)
}
