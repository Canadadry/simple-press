package form

import (
	"app/pkg/data"
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

func ParseGlobalDataEdit(r *http.Request, definition map[string]any) (GobalDataEdit, validator.Errors, error) {
	parsed := GobalDataEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return GobalDataEdit{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}
	form_data, err := data.ParseFormData(parsed.Data, definition)
	//TODO check err
	parsed.Data = form_data

	return parsed, errs, nil
}
