package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

const (
	blockAddName = "name"
)

type Block struct {
	Name string
}

type BlockError struct {
	Name string
}

func (e BlockError) HasError() bool {
	return e.Name != ""
}

func (b *Block) Bind(binder validator.Binder) {
	binder.RequiredStringVar(
		blockAddName,
		&b.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseBlockAdd(r *http.Request) (Block, BlockError, error) {
	parsed := Block{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Block{}, BlockError{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := BlockError{
		Name: strings.Join(errs.Errors[blockAddName], ", "),
	}

	return parsed, resultErr, nil
}
