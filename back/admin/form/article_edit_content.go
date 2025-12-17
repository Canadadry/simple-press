package form

import (
	"app/pkg/validator"
	"fmt"
	"net/http"
)

const (
	articleEditContent = "content"
)

type ParsedArticleEditContent struct {
	Content string
}

func (p *ParsedArticleEditContent) Bind(b validator.Binder) {
	b.RequiredStringVar(articleEditContent, &p.Content, validator.Length(1, maxContentLen))
}

func ParseArticleEditContent(r *http.Request) (ParsedArticleEditContent, validator.Errors, error) {
	parsed := ParsedArticleEditContent{}

	errs, err := validator.BindWithForm(r, parsed.Bind)
	if err != nil {
		return ParsedArticleEditContent{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}

	return parsed, errs, nil
}
