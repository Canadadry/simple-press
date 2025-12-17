package form

import (
	"app/pkg/null"
	"app/pkg/validator"
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleEditNewBlock      = "new_block"
	articleEditBlockPosition = "position"
)

type ParsedArticleEditBlockAdd struct {
	AddedBlockID int64
	Position     null.Nullable[int]
}

type ParsedArticleEditErrorBlockAdd struct {
	AddedBlockID string
	Raw          validator.Errors
}

func (e ParsedArticleEditErrorBlockAdd) HasError() bool {
	return e.AddedBlockID != ""
}

func (p *ParsedArticleEditBlockAdd) Bind(check_id func(int64) error) func(b validator.Binder) {
	return func(b validator.Binder) {
		b.RequiredInt64Var(articleEditNewBlock, &p.AddedBlockID,
			validator.Min(int64(1)),
			check_id,
		)
		b.IntVar(articleEditBlockPosition, &p.Position)
	}
}

func ParseArticleEditBlockAdd(
	r *http.Request,
	checkID func(context.Context, int64) (int, error),
) (ParsedArticleEditBlockAdd, ParsedArticleEditErrorBlockAdd, error) {

	parsed := ParsedArticleEditBlockAdd{}

	errs, err := validator.BindWithForm(r, parsed.Bind(func(val int64) error {
		count, err := checkID(r.Context(), val)
		if err != nil || count == 0 {
			return fmt.Errorf("invalid id")
		}
		return nil
	}))
	if err != nil {
		return ParsedArticleEditBlockAdd{}, ParsedArticleEditErrorBlockAdd{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditErrorBlockAdd{
		AddedBlockID: strings.Join(errs.Errors[articleEditNewBlock], ", "),
		Raw:          errs,
	}

	return parsed, resultErr, nil
}
