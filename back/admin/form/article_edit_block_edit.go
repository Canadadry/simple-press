package form

import (
	"app/pkg/data"
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleEditBlockEditedID       = "block_id"
	articleEditBlockEditedData     = "block_data"
	articleEditBlockEditedPosition = "block_position"
)

type ParsedArticleEditBlockEdit struct {
	EditedBlockID   int64
	EditedBlockData map[string]any
	// EditedBlockPosition int
}

type ParsedArticleEditErrorBlockEdit struct {
	EditedBlockID   string
	EditedBlockData string
	// EditedBlockPosition string
	Raw validator.Errors
}

func (e ParsedArticleEditErrorBlockEdit) HasError() bool {
	return e.EditedBlockID != "" ||
		e.EditedBlockData != "" //||
	// e.EditedBlockPosition != ""
}

func (p *ParsedArticleEditBlockEdit) Bind(exist func(int64) bool) func(b validator.Binder) {
	return func(b validator.Binder) {
		b.RequiredInt64Var(articleEditBlockEditedID, &p.EditedBlockID,
			validator.Min(int64(1)),
			validator.Exist(exist),
		)
		b.RequiredMapVar(articleEditBlockEditedData, &p.EditedBlockData)
	}
}

func ParseArticleEditBlockEdit(
	r *http.Request,
	get_previous_data func() map[string]any,
	is_id_valid func(int64) bool,
) (ParsedArticleEditBlockEdit, ParsedArticleEditErrorBlockEdit, error) {

	parsed := ParsedArticleEditBlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind(is_id_valid))
	if err != nil {
		return parsed,
			ParsedArticleEditErrorBlockEdit{},
			fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditErrorBlockEdit{
		EditedBlockID: strings.Join(errs.Errors[articleEditBlockEditedID], ", "),
		Raw:           errs,
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
