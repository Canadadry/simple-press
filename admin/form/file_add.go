package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	FileAddName    = "name"
	FileAddContent = "content"
)

type File struct {
	Name    string
	Content []byte
}

type FileError struct {
	Name    string
	Content string
}

func (le FileError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	return false
}

func ParseFileAdd(r *http.Request) (File, FileError, error) {
	err := r.ParseForm()
	if err != nil {
		return File{}, FileError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	l := File{
		Name: r.PostForm.Get(FileAddName),
	}
	errors := FileError{}
	if l.Name == "" {
		errors.Name = errorCannotBeEmpty
	}
	if len(l.Name) > maxTitleLen {
		errors.Name = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.PathRegexp + "$")
	if !re.Match([]byte(l.Name)) {
		errors.Name = errorNotAPath
	}
	return l, errors, nil
}
