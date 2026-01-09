package form

import (
	"app/pkg/validator"
	"fmt"
	"net/http"
)

type GobalDataEdit struct {
	Data map[string]any
}

func (b *GobalDataEdit) Bind(bind validator.Binder) {
	bind.RequiredMapVar("data", &b.Data)
}

func ParseGlobalDataEdit(r *http.Request) (GobalDataEdit, validator.Errors, error) {
	parsed := GobalDataEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return GobalDataEdit{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}

	return parsed, errs, nil
}
