package form

import (
	"app/pkg/null"
	"app/pkg/router"
	"app/pkg/validator"
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleEditTitle  = "title"
	articleEditAuthor = "author"
	articleEditDraft  = "draft"
	articleEditSlug   = "slug"
	articleEditLayout = "layout"
)

type ParsedArticleEditMetadata struct {
	Title    string
	Author   string
	Draft    null.Nullable[bool]
	Slug     string
	LayoutID int64
}

type ParsedArticleEditMetadataError struct {
	Title    string
	Author   string
	Slug     string
	LayoutID string
}

func (e ParsedArticleEditMetadataError) HasError() bool {
	return e.Title != "" ||
		e.Author != "" ||
		e.Slug != "" ||
		e.LayoutID != ""
}

func (p *ParsedArticleEditMetadata) Bind(check_id func(int64) error) func(b validator.Binder) {
	return func(b validator.Binder) {
		b.RequiredStringVar(articleEditTitle, &p.Title, validator.Length(1, maxTitleLen))
		b.RequiredStringVar(articleEditAuthor, &p.Author, validator.Length(1, maxAuthorLen))

		b.RequiredStringVar(articleEditSlug, &p.Slug,
			validator.Length(1, maxSlugLen),
			validator.Regexp("^"+router.SlugRegexp+"$"),
		)
		b.RequiredInt64Var(articleEditLayout, &p.LayoutID, validator.Min(int64(1)), check_id)
		b.BoolVar(articleEditDraft, &p.Draft, validator.TrueChoice, validator.FalseChoice)
	}
}

func ParseArticleEditMetadata(
	r *http.Request,
	checkID func(context.Context, int64) (int, error),
) (ParsedArticleEditMetadata, ParsedArticleEditMetadataError, error) {

	parsed := ParsedArticleEditMetadata{}
	errs, err := validator.BindWithForm(r, parsed.Bind(func(val int64) error {
		count, err := checkID(r.Context(), val)
		if err != nil || count == 0 {
			return fmt.Errorf("invalid id")
		}
		return nil
	}))
	if err != nil {
		return ParsedArticleEditMetadata{}, ParsedArticleEditMetadataError{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditMetadataError{
		Title:    strings.Join(errs.Errors[articleEditTitle], " ,"),
		Author:   strings.Join(errs.Errors[articleEditAuthor], " ,"),
		Slug:     strings.Join(errs.Errors[articleEditSlug], " ,"),
		LayoutID: strings.Join(errs.Errors[articleEditLayout], " ,"),
	}

	return parsed, resultErr, nil
}
