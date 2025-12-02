package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
	"strings"
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
	return le.Name != "" || le.Content != ""
}

func (l *LayoutEdit) Bind(b validator.Binder) {
	b.RequiredStringVar(
		layoutEditName,
		&l.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
	b.RequiredStringVar(
		layoutEditContent,
		&l.Content,
		validator.Length(1, maxContentLen),
	)
}

func ParseLayoutEdit(r *http.Request) (LayoutEdit, LayoutEditError, error) {
	parsed := LayoutEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return LayoutEdit{}, LayoutEditError{}, err
	}

	resultErr := LayoutEditError{
		Name:    strings.Join(errs.Errors[layoutEditName], ", "),
		Content: strings.Join(errs.Errors[layoutEditContent], ", "),
	}

	return parsed, resultErr, nil
}
