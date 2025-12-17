package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"fmt"
	"net/http"
)

const (
	blockEditName       = "name"
	blockEditContent    = "content"
	blockEditDefinition = "definition"
)

type BlockEdit struct {
	Name       string
	Content    string
	Definition map[string]any
}

func (b *BlockEdit) Bind(bind validator.Binder) {

	bind.RequiredStringVar(
		blockEditName,
		&b.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)

	bind.RequiredStringVar(
		blockEditContent,
		&b.Content,
		validator.Length(0, maxContentLen),
	)

	bind.RequiredMapVar(blockEditDefinition, &b.Definition)
}

func ParseBlockEdit(r *http.Request) (BlockEdit, validator.Errors, error) {
	parsed := BlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return BlockEdit{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}

	return parsed, errs, nil
}
