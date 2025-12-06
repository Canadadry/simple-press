package form

import (
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleEditContent = "content"
)

type ParsedArticleEditContent struct {
	Content string
}

type ParsedArticleEditErrorContent struct {
	Content string
	Raw     validator.Errors
}

func (pe ParsedArticleEditErrorContent) HasError() bool {
	return pe.Content != ""
}

func (p *ParsedArticleEditContent) Bind(b validator.Binder) {
	b.RequiredStringVar(articleEditContent, &p.Content, validator.Length(1, maxContentLen))
}

func ParseArticleEditContent(r *http.Request) (ParsedArticleEditContent, ParsedArticleEditErrorContent, error) {
	parsed := ParsedArticleEditContent{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return ParsedArticleEditContent{}, ParsedArticleEditErrorContent{}, fmt.Errorf("cannot parse form : %w", err)
	}

	resultErr := ParsedArticleEditErrorContent{
		Content: strings.Join(errs.Errors[articleEditContent], ", "),
		Raw:     errs,
	}

	return parsed, resultErr, nil
}
