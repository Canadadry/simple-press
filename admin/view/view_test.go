package view

import (
	"app/pkg/flash"
	"bytes"
	"testing"
)

func fakeTr(key string) string {
	return key
}

func TestView(t *testing.T) {
	tests := map[string]ViewFunc{
		"article_list": ArticlesList(ArticlesListData{
			Articles: []ArticleListData{ArticleListData{}},
		}),
		"article_add":             ArticleAdd(ArticleAddData{}, ArticleAddError{}),
		"article_add with error":  ArticleAdd(ArticleAddData{}, ArticleAddError{"test1", "test2"}),
		"article_edit":            ArticleEdit(ArticleEditData{}, ArticleEditError{}),
		"article_edit with error": ArticleEdit(ArticleEditData{}, ArticleEditError{"test1", "test2", "test3", "test4"}),
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt(buf, fakeTr, flash.Message{})
			if err != nil {
				t.Fatalf("failed : %v", err)
			}
		})
	}
}
