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

type Layout struct {
	Name string
}

type LayoutError struct {
	Name string
}

func (le LayoutError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Name != "" {
		return true
	}
	return false
}

func ParseLayoutAdd(r *http.Request) (Layout, LayoutError, error) {
	err := r.ParseForm()
	if err != nil {
		return Layout{}, LayoutError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	l := Layout{
		Name: r.PostForm.Get(layoutAddName),
	}
	errors := LayoutError{}
	if l.Name == "" {
		errors.Name = errorCannotBeEmpty
	}
	if len(l.Name) > maxTitleLen {
		errors.Name = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.SlugRegexp + "$")
	if !re.Match([]byte(l.Name)) {
		errors.Name = errorNotASlug
	}
	return l, errors, nil
}
