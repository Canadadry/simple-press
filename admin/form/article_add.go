package form

import (
	"fmt"
	"net/http"
)

const (
	articleAddTitle        = "title"
	articleAddAuthor       = "author"
	articleAddDraft        = "draft"
	errorTageCannotBeEmpty = "empty"
	errorTagetToBig        = "too_big"
	maxTitleLen            = 255
	maxAuthorLen           = 255
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
		return false
	}
	if ae.Author == "" {
		return false
	}
	if ae.Draft == "" {
		return false
	}
	return true
}

func ParseArticleAdd(r *http.Request) (Article, ArticleError, error) {
	err := r.ParseForm()
	if err != nil {
		return Article{}, ArticleError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := Article{
		Title:  r.PostForm.Get(articleAddTitle),
		Author: r.PostForm.Get(articleAddAuthor),
		Draft:  r.PostForm.Get(articleAddDraft) != "",
	}
	errors := ArticleError{}
	if a.Title == "" {
		errors.Title = errorTageCannotBeEmpty
	}
	if a.Author == "" {
		errors.Author = errorTageCannotBeEmpty
	}
	if len(a.Author) > maxTitleLen {
		errors.Author = errorTagetToBig
	}
	if len(a.Author) > maxTitleLen {
		errors.Author = errorTagetToBig
	}
	return a, errors, nil
}
