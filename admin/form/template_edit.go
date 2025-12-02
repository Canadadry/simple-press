package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
	"strings"
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

func (te TemplateEditError) HasError() bool {
	return te.Name != "" || te.Content != ""
}

func (t *TemplateEdit) Bind(b validator.Binder) {
	b.RequiredStringVar(
		templateEditName,
		&t.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
	b.RequiredStringVar(
		templateEditContent,
		&t.Content,
		validator.Length(1, maxContentLen),
	)
}

func ParseTemplateEdit(r *http.Request) (TemplateEdit, TemplateEditError, error) {
	parsed := TemplateEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return TemplateEdit{}, TemplateEditError{}, err
	}

	resultErr := TemplateEditError{
		Name:    strings.Join(errs.Errors[templateEditName], ", "),
		Content: strings.Join(errs.Errors[templateEditContent], ", "),
	}

	return parsed, resultErr, nil
}
