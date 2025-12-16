package form

import (
	"app/pkg/data"
	"app/pkg/validator"
	"fmt"
	"net/http"
)

const (
	articleEditBlockEditedData     = "block_data"
	articleEditBlockEditedPosition = "block_position"
)

type ParsedArticleEditBlockEdit struct {
	EditedBlockData map[string]any
	// EditedBlockPosition int
}

type ParsedArticleEditErrorBlockEdit struct {
	EditedBlockData string
	// EditedBlockPosition string
	Raw validator.Errors
}

func (e ParsedArticleEditErrorBlockEdit) HasError() bool {
	return e.EditedBlockData != "" //||
	// e.EditedBlockPosition != ""
}

func (p *ParsedArticleEditBlockEdit) Bind(b validator.Binder) {
	b.RequiredMapVar(articleEditBlockEditedData, &p.EditedBlockData)
}

func ParseArticleEditBlockEdit(
	r *http.Request,
	get_previous_data func() map[string]any,
) (ParsedArticleEditBlockEdit, ParsedArticleEditErrorBlockEdit, error) {

	parsed := ParsedArticleEditBlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return parsed,
			ParsedArticleEditErrorBlockEdit{},
			fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditErrorBlockEdit{
		Raw: errs,
	}

	form_data := get_previous_data()
	form_data, err = data.ParseFormData(parsed.EditedBlockData, form_data)
	if err != nil {
		resultErr.EditedBlockData = errorInvalidJson
		return parsed, resultErr, nil
	}
	parsed.EditedBlockData = form_data
	return parsed, resultErr, nil
}
