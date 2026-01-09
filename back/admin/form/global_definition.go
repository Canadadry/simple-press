package form

import (
	"app/pkg/validator"
	"fmt"
	"net/http"
)

type GobalDefinitionEdit struct {
	Definition map[string]any
}

func (b *GobalDefinitionEdit) Bind(bind validator.Binder) {
	bind.RequiredMapVar("definition", &b.Definition)
}

func ParseGlobalDefinitionEdit(r *http.Request) (GobalDefinitionEdit, validator.Errors, error) {
	parsed := GobalDefinitionEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return GobalDefinitionEdit{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}

	return parsed, errs, nil
}
