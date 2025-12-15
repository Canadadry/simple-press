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
}

type FileError struct {
	Request string
	Name    string
	Content string
	Raw     validator.Errors
}

func (le FileError) HasError() bool {
	if le.Request != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	if le.Content != "" {
		return true
	}
	return false
}

func invalidRequest(msg string) FileError {
	return FileError{
		Request: msg,
		Raw: validator.Errors{
			Errors: map[string][]string{
				"name":    []string{},
				"content": []string{},
				"request": []string{
					msg,
				},
			},
		},
	}
}

func invalidName(msg string) FileError {
	return FileError{
		Name: msg,
		Raw: validator.Errors{
			Errors: map[string][]string{
				"name":    []string{msg},
				"content": []string{},
				"request": []string{},
			},
		},
	}
}

func invalidContent(msg string) FileError {
	return FileError{
		Content: msg,
		Raw: validator.Errors{
			Errors: map[string][]string{
				"name":    []string{},
				"content": []string{msg},
				"request": []string{},
			},
		},
	}
}

func ParseFileAdd(r *http.Request) (File, FileError, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
		return File{}, invalidRequest(RequestWrongContentType), nil
	}
	err := r.ParseMultipartForm(MaxRequestSize)
	if err != nil {
		return File{}, invalidRequest(RequestInvalidContent), nil
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
	}, FileError{}, nil

}
