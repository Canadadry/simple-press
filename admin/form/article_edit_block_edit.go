package form

import (
	"app/pkg/data"
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleEditEditedBlockID = "edited_block_id"
)

type ParsedArticleEditBlockEdit struct {
	EditedBlockID       int64
	EditedBlockData     map[string]any
	EditedBlockPosition int
}

type ParsedArticleEditErrorBlockEdit struct {
	EditedBlockID       string
	EditedBlockData     string
	EditedBlockPosition string
}

func (e ParsedArticleEditErrorBlockEdit) HasError() bool {
	return e.EditedBlockID != "" ||
		e.EditedBlockData != "" ||
		e.EditedBlockPosition != ""
}

func (p *ParsedArticleEditBlockEdit) Bind(b validator.Binder) {
	b.RequiredInt64Var(articleEditEditedBlockID, &p.EditedBlockID, validator.Min(int64(1)))
}

func ParseArticleEditBlockEdit(
	r *http.Request,
	get_previous_data func(int64) (map[string]any, bool),
) (ParsedArticleEditBlockEdit, ParsedArticleEditErrorBlockEdit, error) {

	parsed := ParsedArticleEditBlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return ParsedArticleEditBlockEdit{}, ParsedArticleEditErrorBlockEdit{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditErrorBlockEdit{
		EditedBlockID: strings.Join(errs.Errors[articleEditEditedBlockID], ", "),
	}

	form_data, ok := get_previous_data(parsed.EditedBlockID)
	if !ok {
		resultErr.EditedBlockID = errorInvalidId
		return parsed, resultErr, nil
	}
	form_data, err = data.ParseFormData(r, form_data)
	if err != nil {
		resultErr.EditedBlockData = errorInvalidJson
		return parsed, resultErr, nil
	}
	parsed.EditedBlockData = form_data
	return parsed, resultErr, nil
}
