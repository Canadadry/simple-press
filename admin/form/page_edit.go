package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	layoutEditName    = "name"
	layoutEditContent = "content"
)

type LayoutEdit struct {
	Name    string
	Content string
}

type LayoutEditError struct {
	Name    string
	Content string
}

func (le LayoutEditError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Content != "" {
		return true
	}
	return false
}

func ParseLayoutEdit(r *http.Request) (LayoutEdit, LayoutEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return LayoutEdit{}, LayoutEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := LayoutEdit{
		Name:    r.PostForm.Get(layoutEditName),
		Content: r.PostForm.Get(layoutEditContent),
	}
	errors := LayoutEditError{}
	if a.Name == "" {
		errors.Name = errorCannotBeEmpty
	}
	if a.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if len(a.Name) > maxTitleLen {
		errors.Name = errorTagetToBig
	}
	if len(a.Content) > maxContentLen {
		errors.Content = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.PathRegexp + "$")
	if !re.Match([]byte(a.Name)) {
		errors.Name = errorNotAPath
	}
	fmt.Println("content submited -", a.Content, "-")
	return a, errors, nil
}
