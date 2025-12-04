package form

import (
	"app/pkg/validator"
	"fmt"
	"io"
	"net/http"
)

const (
	FileAddContent = "content"
)

type File struct {
	Name    string
	Content []byte
}

type FileError struct {
	Content string
	Raw     validator.Errors
}

func (le FileError) HasError() bool {
	if le.Content != "" {
		return true
	}
	return false
}

func ParseFileAdd(r *http.Request) (File, FileError, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		return File{}, FileError{}, fmt.Errorf("cannot parse multipart form: %w", err)
	}

	errors := FileError{}

	// Récupération du fichier
	file, header, err := r.FormFile(FileAddContent)
	if err != nil {
		errors.Content = errorCannotBeEmpty
	} else {
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			return File{}, FileError{}, fmt.Errorf("cannot read file content: %w", err)
		}

		return File{
			Name:    header.Filename,
			Content: content,
		}, errors, nil
	}

	return File{}, errors, nil
}
