package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	templateEditName    = "name"
	templateEditContent = "content"
)

type TemplateEdit struct {
	Name    string
	Content string
}

type TemplateEditError struct {
	Name    string
	Content string
}

func (le TemplateEditError) HasError() bool {
	if le.Name != "" {
		return true
	}
	if le.Content != "" {
		return true
	}
	return false
}

func ParseTemplateEdit(r *http.Request) (TemplateEdit, TemplateEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return TemplateEdit{}, TemplateEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := TemplateEdit{
		Name:    r.PostForm.Get(templateEditName),
		Content: r.PostForm.Get(templateEditContent),
	}
	errors := TemplateEditError{}
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
	return a, errors, nil
}
