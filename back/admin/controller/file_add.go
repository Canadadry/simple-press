package controller

import (
	"app/admin/form"
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (c *Controller) PostFileAdd(w http.ResponseWriter, r *http.Request) error {

	f, errors, err := form.ParseFileAdd(r)
	if err != nil {
		return fmt.Errorf("cannot parse form request : %w", err)
	}

	if errors.HasError {
		return httpresponse.BadRequest(w, errors)
	}
	if !f.Archive {
		id, err := c.Repository.UploadFile(r.Context(), repository.File{Name: f.Name, Content: f.Content})
		if err != nil {
			return fmt.Errorf("cannot create File : %w", err)
		}
		return view.FileAddCreated(w, view.FileAddData{ID: id, Name: f.Name})
	}
	zr := bytes.NewReader(f.Content)
	dc, err := zip.NewReader(zr, int64(len(f.Content)))
	if err != nil {
		return httpresponse.BadRequest(w, fmt.Errorf("cannot read archive"))
	}
	added := []view.FileAddData{}
	for _, file := range dc.File {
		if file.Name[len(file.Name)-1] != '/' {
			fr, err := file.Open()
			Content, err := io.ReadAll(fr)
			if err != nil || len(Content) == 0 {
				return httpresponse.BadRequest(w, fmt.Errorf("cannot read archive"))
			}
			id, err := c.Repository.UploadFile(r.Context(), repository.File{Name: file.Name, Content: Content})
			if err != nil {
				return fmt.Errorf("cannot create File : %w", err)
			}
			added = append(added, view.FileAddData{ID: id, Name: file.Name})
		}
	}
	return view.FileAddListCreated(w, view.FileAddListData{
		Total: len(added),
		Items: added,
	})
}
