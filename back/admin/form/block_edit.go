package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
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

type BlockEditError struct {
	Name       string
	Content    string
	Definition string
	Raw        validator.Errors
}

func (e BlockEditError) HasError() bool {
	return e.Name != "" ||
		e.Content != "" ||
		e.Definition != ""
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

func ParseBlockEdit(r *http.Request) (BlockEdit, BlockEditError, error) {
	parsed := BlockEdit{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return BlockEdit{}, BlockEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := BlockEditError{
		Name:       strings.Join(errs.Errors[blockEditName], ", "),
		Content:    strings.Join(errs.Errors[blockEditContent], ", "),
		Definition: strings.Join(errs.Errors[blockEditDefinition], ", "),
		Raw:        errs,
	}

	return parsed, resultErr, nil
}
