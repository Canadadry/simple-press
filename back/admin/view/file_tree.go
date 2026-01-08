package view

import (
	"app/pkg/http/httpresponse"
	"net/http"
)

type FileTreeData struct {
	Path    string   `json:"path"`
	Files   []string `json:"files"`
	Folders []string `json:"folders"`
}

func FilesTreeOk(w http.ResponseWriter, ftd FileTreeData) error {
	return httpresponse.Ok(w, ftd)
}
