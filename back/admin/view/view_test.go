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
			Items: []ArticleListData{ArticleListData{}},
		}),
		"article_add":            ArticleAdd(ArticleAddData{}, ArticleAddError{}),
		"article_add with error": ArticleAdd(ArticleAddData{}, ArticleAddError{"test1", "test2"}),
		"article_edit": ArticleEdit(ArticleEditData{
			Slug:    "slug",
			Layouts: []LayoutSelector{{Name: "test", Value: 1}},
			Blocks:  []LayoutSelector{{Name: "test", Value: 1}},
			BlockDatas: []BlockData{
				{
					ID: 1,
					Data: map[string]any{
						"profile": map[string]any{
							"name": map[string]any{
								"first": "Jane",
								"last":  "Doe",
							},
							"age": 42,
						},
					},
				},
			},
		}, ArticleEditError{}),
		"article_edit with error": ArticleEdit(ArticleEditData{}, ArticleEditError{"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9", "test10", "test11"}),
		"template_list": TemplatesList(TemplatesListData{
			Templates: []TemplateListData{TemplateListData{}},
		}),
		"template_add":             TemplateAdd(TemplateAddData{}, TemplateAddError{}),
		"template_add with error":  TemplateAdd(TemplateAddData{}, TemplateAddError{"test1"}),
		"template_edit":            TemplateEdit(TemplateEditData{}, TemplateEditError{}),
		"template_edit with error": TemplateEdit(TemplateEditData{}, TemplateEditError{"test1", "test2"}),
		"layout_list": LayoutsList(LayoutsListData{
			Items: []LayoutListData{LayoutListData{}},
		}),
		"layout_add":             LayoutAdd(LayoutAddData{}, LayoutAddError{}),
		"layout_add with error":  LayoutAdd(LayoutAddData{}, LayoutAddError{"test1"}),
		"layout_edit":            LayoutEdit(LayoutEditData{}, LayoutEditError{}),
		"layout_edit with error": LayoutEdit(LayoutEditData{}, LayoutEditError{"test1", "test2"}),
		"file_list": FilesList(FilesListData{
			Files: []FileListData{FileListData{}},
		}),
		"file_add":            FileAdd(FileAddError{}),
		"file_add with error": FileAdd(FileAddError{"test1"}),
		"block_list": BlocksList(BlocksListData{
			Blocks: []BlockListData{BlockListData{}},
		}),
		"block_add":             BlockAdd(BlockAddData{}, BlockAddError{}),
		"block_add with error":  BlockAdd(BlockAddData{}, BlockAddError{"test1"}),
		"block_edit":            BlockEdit(BlockEditData{}, BlockEditError{}),
		"block_edit with error": BlockEdit(BlockEditData{}, BlockEditError{"test1", "test2", "test3"}),
		"404":                   PageNotFound,
		"500":                   InternalServerError,
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
