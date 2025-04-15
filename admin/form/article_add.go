package form

import (
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
	Draft  bool
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
		errors.Title = errorCannotBeEmpty
	}
	if a.Author == "" {
		errors.Author = errorCannotBeEmpty
	}
	if len(a.Title) > maxTitleLen {
		errors.Title = errorTagetToBig
	}
	if len(a.Author) > maxAuthorLen {
		errors.Author = errorTagetToBig
	}
	return a, errors, nil
}
