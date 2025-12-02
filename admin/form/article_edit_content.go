package form

import (
	"fmt"
	"net/http"
)

const (
	articleEditContent = "content"
)

type ParsedArticleEditContent struct {
	Content string
}

type ParsedArticleEditErrorContent struct {
	Content string
}

func (pe ParsedArticleEditErrorContent) HasError() bool {
	if pe.Content != "" {
		return true
	}
	return false
}

func ParseArticleEditContent(r *http.Request) (ParsedArticleEditContent, ParsedArticleEditErrorContent, error) {
	err := r.ParseForm()
	if err != nil {
		return ParsedArticleEditContent{}, ParsedArticleEditErrorContent{}, fmt.Errorf("cannot parse form : %w", err)
	}

	pae := ParsedArticleEditContent{
		Content: r.PostForm.Get(articleEditContent),
	}
	errors := ParsedArticleEditErrorContent{}
	if pae.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if len(pae.Content) > maxContentLen {
		errors.Content = errorTagetToBig
	}
	return pae, errors, nil
}
