package view

import (
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
		"layout_list": LayoutsList(LayoutsListData{
			Layouts: []LayoutListData{LayoutListData{}},
		}),
		"layout_add":             LayoutAdd(LayoutAddData{}, LayoutAddError{}),
		"layout_add with error":  LayoutAdd(LayoutAddData{}, LayoutAddError{"test1"}),
		"layout_edit":            LayoutEdit(LayoutEditData{}, LayoutEditError{}),
		"layout_edit with error": LayoutEdit(LayoutEditData{}, LayoutEditError{"test1", "test2"}),
		"file_list": FilesList(FilesListData{
			Files: []FileListData{FileListData{}},
		}),
		"file_add":            FileAdd(FileAddData{}, FileAddError{}),
		"file_add with error": FileAdd(FileAddData{}, FileAddError{"test1"}),
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := tt(buf, fakeTr)
			if err != nil {
				t.Fatalf("failed : %v", err)
			}
		})
	}
}
