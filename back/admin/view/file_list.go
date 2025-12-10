package view

import (
	"app/pkg/http/httpresponse"
	"app/pkg/paginator"
	"io"
	"math"
	"net/http"
)

const (
	MaxFilePaginationItem = 5
)

type FilesListData struct {
	Items []FileListData `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type FileListData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
				Files: data.Items,
			}),
		)
	}
}

func FilesListOk(w http.ResponseWriter, a FilesListData) error {
	return httpresponse.Ok(w, a)
}
