package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	layoutAddName = "name"
)

type Template struct {
	Name string
}

type TemplateError struct {
	Name string
}

func (le TemplateError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	return false
}

func ParseTemplateAdd(r *http.Request) (Template, TemplateError, error) {
	err := r.ParseForm()
	if err != nil {
		return Template{}, TemplateError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	l := Template{
		Name: r.PostForm.Get(layoutAddName),
	}
	errors := TemplateError{}
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
