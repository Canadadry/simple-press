package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"net/http"
)

const (
	templateEditName    = "name"
	templateEditContent = "content"
)

type TemplateEdit struct {
	Name    string
	Content string
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

func ParseTemplateEdit(r *http.Request) (TemplateEdit, validator.Errors, error) {
	parsed := TemplateEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return TemplateEdit{}, validator.Errors{}, err
	}
	return parsed, errs, nil
}
