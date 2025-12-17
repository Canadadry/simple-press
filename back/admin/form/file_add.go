package form

import (
	"app/pkg/validator"
	"io"
	"net/http"
	"strings"
)

const (
	FileAddContent          = "content"
	FileAddName             = "name"
	FileAddArchive          = "archive"
	RequestWrongContentType = "content type should be a multipart/form-data"
	RequestInvalidContent   = "cannot parse multipart/form-data"
	FileNotFound            = "file not found"
	NoFileName              = "file has no name"
	CannotReadContent       = "cannot read content"
	MaxRequestSize          = 10 << 20
)

type File struct {
	Name    string
	Content []byte
	Archive bool
}

func invalidRequest(msg string) validator.Errors {
	return validator.Errors{
		HasError: true,
		Errors: map[string][]string{
			"name":    []string{},
			"content": []string{},
			"request": []string{
				msg,
			},
		},
	}
}

func invalidName(msg string) validator.Errors {
	return validator.Errors{
		HasError: true,
		Errors: map[string][]string{
			"name":    []string{msg},
			"content": []string{},
			"request": []string{},
		},
	}
}

func invalidContent(msg string) validator.Errors {
	return validator.Errors{
		HasError: true,
		Errors: map[string][]string{
			"name":    []string{},
			"content": []string{msg},
			"request": []string{},
		},
	}
}

func ParseFileAdd(r *http.Request) (File, validator.Errors, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
		return File{}, invalidRequest(RequestWrongContentType), nil
	}
	err := r.ParseMultipartForm(MaxRequestSize)
	if err != nil {
		return File{}, invalidRequest(RequestInvalidContent), nil
	}
	archive := r.FormValue(FileAddArchive)
	filename := r.FormValue(FileAddName)
	if filename == "" {
		return File{}, invalidName(NoFileName), nil
	}
	file, _, err := r.FormFile(FileAddContent)
	if err != nil {
		return File{}, invalidContent(FileNotFound), nil
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return File{}, invalidContent(CannotReadContent), nil
	}

	return File{
		Name:    r.FormValue(FileAddName),
		Content: content,
		Archive: archive == "true",
	}, validator.Errors{}, nil

}
