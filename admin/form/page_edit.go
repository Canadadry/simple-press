package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	pageEditName    = "name"
	pageEditContent = "content"
)

type PageEdit struct {
	Name    string
	Content string
}

type PageEditError struct {
	Name    string
	Content string
}

func (le PageEditError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Content != "" {
		return true
	}
	return false
}

func ParsePageEdit(r *http.Request) (PageEdit, PageEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return PageEdit{}, PageEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := PageEdit{
		Name:    r.PostForm.Get(pageEditName),
		Content: r.PostForm.Get(pageEditContent),
	}
	errors := PageEditError{}
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
