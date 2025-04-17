package view

import (
	"app/pkg/paginator"
	"io"
	"math"
)

const (
	MaxFilePaginationItem = 5
)

type FilesListData struct {
	Files []FileListData
	Total int
	Page  int
	Limit int
}

type FileListData struct {
	Name string
}

func FilesList(data FilesListData) ViewFunc {
	type viewData struct {
		Total int
		Pages paginator.Pages
		Files []FileListData
	}
	return func(w io.Writer, tr func(string) string) error {
		lastPage := int(math.Ceil(float64(data.Total) / float64(data.Limit)))
		p := paginator.New(data.Page, lastPage, MaxFilePaginationItem, "/files?page=%page%")
		return render(w, tr,
			"template/pages/file_list.html",
			TemplateData("FILE_LIST.page_title", viewData{
				Total: data.Total,
				Pages: p,
				Files: data.Files,
			}),
		)
	}
}
