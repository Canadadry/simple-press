package form

import (
	"app/pkg/null"
	"app/pkg/validator"
	"fmt"
	"net/http"
	"strings"
)

const (
	articleAddTitle  = "title"
	articleAddAuthor = "author"
	articleAddDraft  = "draft"
)

type Article struct {
	Title  string
	Author string
	Draft  null.Nullable[bool]
}

func (a *Article) Bind(b validator.Binder) {
	b.RequiredStringVar(articleAddTitle, &a.Title, validator.Length(3, maxTitleLen))
	b.RequiredStringVar(articleAddAuthor, &a.Author, validator.Length(3, maxAuthorLen))
	b.BoolVar(articleAddDraft, &a.Draft, validator.TrueChoice, validator.FalseChoice)
}

type ArticleError struct {
	Title  string
	Author string
}

func (ae ArticleError) HasError() bool {
	if ae.Title != "" {
		return true
	}
	if ae.Author != "" {
		return true
	}
	return false
}

func ParseArticleAdd(r *http.Request) (Article, ArticleError, error) {
	article := Article{}
	errs, err := validator.BindWithForm(r, article.Bind)
	if err != nil {
		return Article{}, ArticleError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	ae := ArticleError{
		Title:  strings.Join(errs.Errors[articleAddTitle], ", "),
		Author: strings.Join(errs.Errors[articleAddAuthor], ", "),
	}
	return article, ae, nil
}
