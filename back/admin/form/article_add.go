package form

import (
	"app/pkg/null"
	"app/pkg/validator"
	"fmt"
	"net/http"
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

func ParseArticleAdd(r *http.Request) (Article, validator.Errors, error) {
	article := Article{}
	errs, err := validator.BindWithForm(r, article.Bind)
	if err != nil {
		return Article{}, validator.Errors{}, fmt.Errorf("cannot parse form : %w", err)
	}
	return article, errs, nil
}
