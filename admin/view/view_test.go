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
		"article_edit":            ArticleEdit(ArticleEditData{Pages: []PageSelector{{Name: "test", Value: 1}}}, ArticleEditError{}),
		"article_edit with error": ArticleEdit(ArticleEditData{}, ArticleEditError{"test1", "test2", "test3", "test4", "test5"}),
		"layout_list": TemplatesList(TemplatesListData{
			Templates: []TemplateListData{TemplateListData{}},
		}),
		"layout_add":             TemplateAdd(TemplateAddData{}, TemplateAddError{}),
		"layout_add with error":  TemplateAdd(TemplateAddData{}, TemplateAddError{"test1"}),
		"layout_edit":            TemplateEdit(TemplateEditData{}, TemplateEditError{}),
		"layout_edit with error": TemplateEdit(TemplateEditData{}, TemplateEditError{"test1", "test2"}),
		"page_list": PagesList(PagesListData{
			Pages: []PageListData{PageListData{}},
		}),
		"page_add":             PageAdd(PageAddData{}, PageAddError{}),
		"page_add with error":  PageAdd(PageAddData{}, PageAddError{"test1"}),
		"page_edit":            PageEdit(PageEditData{}, PageEditError{}),
		"page_edit with error": PageEdit(PageEditData{}, PageEditError{"test1", "test2"}),
		"file_list": FilesList(FilesListData{
			Files: []FileListData{FileListData{}},
		}),
		"file_add":            FileAdd(FileAddError{}),
		"file_add with error": FileAdd(FileAddError{"test1"}),
		"404":                 PageNotFound,
		"500":                 InternalServerError,
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
