package form

import (
	"app/pkg/data"
	"app/pkg/null"
	"app/pkg/validator"
	"fmt"
	"net/http"
)

const (
	articleEditBlockEditedData     = "block_data"
	articleEditBlockEditedPosition = "block_position"
)

type ParsedArticleEditBlockEdit struct {
	EditedBlockData     map[string]any
	EditedBlockPosition null.Nullable[int]
}

func (p *ParsedArticleEditBlockEdit) Bind(b validator.Binder) {
	b.MapVar(articleEditBlockEditedData, &p.EditedBlockData)
	b.IntVar(articleEditBlockEditedPosition, &p.EditedBlockPosition)
}

func ParseArticleEditBlockEdit(
	r *http.Request,
	get_previous_data func() map[string]any,
) (ParsedArticleEditBlockEdit, validator.Errors, error) {

	parsed := ParsedArticleEditBlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return parsed,
			validator.Errors{},
			fmt.Errorf("cannot parse form : %w", err)
	}

	form_data := get_previous_data()
	form_data, err = data.ParseFormData(parsed.EditedBlockData, form_data)
	//TODO check err
	parsed.EditedBlockData = form_data
	return parsed, errs, nil
}
