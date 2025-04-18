package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	pageAddName = "name"
)

type Page struct {
	Name string
}

type PageError struct {
	Name string
}

func (le PageError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	return false
}

func ParsePageAdd(r *http.Request) (Page, PageError, error) {
	err := r.ParseForm()
	if err != nil {
		return Page{}, PageError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	l := Page{
		Name: r.PostForm.Get(pageAddName),
	}
	errors := PageError{}
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
