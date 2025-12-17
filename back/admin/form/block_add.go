package form

import (
	"app/pkg/router"
	"app/pkg/validator"
	"fmt"
	"net/http"
)

const (
	blockAddName = "name"
)

type Block struct {
	Name string
}

func (b *Block) Bind(binder validator.Binder) {
	binder.RequiredStringVar(
		blockAddName,
		&b.Name,
		validator.Length(1, maxTitleLen),
		validator.Regexp("^"+router.PathRegexp+"$"),
	)
}

func ParseBlockAdd(r *http.Request) (Block, validator.Errors, error) {
	parsed := Block{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return Block{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}
	return parsed, errs, nil
}
