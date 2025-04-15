package form

import (
	"context"
	"fmt"
	"net/http"
)

const (
	sendEmailField = "email"
)

type Article struct {
	Title  string
	Author string
	Draft  bool
}

type ArticleError struct {
	Title  string
	Author string
	Draft  string
}

func (ae ArticleError) HasError() bool {
	if ae.Title == "" {
		return true
	}
	if ae.Author == "" {
		return true
	}
	if ae.Draft == "" {
		return true
	}
	return false
}

func ParseArticleAdd(r *http.Request, countBySlug func(context.Context, string) (int, error)) (Article, ArticleError, error) {
	a := Article{}
	err := r.ParseForm()
	if err != nil {
		return a, ArticleError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	return Article{}, ArticleError{}, fmt.Errorf("ParseArticleAdd not impl")
}
