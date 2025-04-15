package form

import (
	"fmt"
	"net/http"
)

const (
	articleEditTitle   = "title"
	articleEditAuthor  = "author"
	articleEditDraft   = "draft"
	articleEditContent = "content"
	articleEditSlug    = "Slug"
)

type ArticleEdit struct {
	Title   string
	Author  string
	Draft   bool
	Content string
	Slug    string
}

type ArticleEditError struct {
	Title   string
	Author  string
	Slug    string
	Content string
}

func (ae ArticleEditError) HasError() bool {
	if ae.Title == "" {
		return false
	}
	if ae.Author == "" {
		return false
	}
	if ae.Content == "" {
		return false
	}
	if ae.Slug == "" {
		return false
	}
	return true
}

func ParseArticleEdit(r *http.Request) (ArticleEdit, ArticleEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return ArticleEdit{}, ArticleEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	a := ArticleEdit{
		Title:   r.PostForm.Get(articleEditTitle),
		Author:  r.PostForm.Get(articleEditAuthor),
		Content: r.PostForm.Get(articleEditContent),
		Slug:    r.PostForm.Get(articleEditSlug),
		Draft:   r.PostForm.Get(articleEditDraft) != "",
	}
	errors := ArticleEditError{}
	if a.Title == "" {
		errors.Title = errorTageCannotBeEmpty
	}
	if a.Content == "" {
		errors.Content = errorTageCannotBeEmpty
	}
	if a.Slug == "" {
		errors.Slug = errorTageCannotBeEmpty
	}
	if a.Author == "" {
		errors.Author = errorTageCannotBeEmpty
	}
	if len(a.Title) > maxTitleLen {
		errors.Title = errorTagetToBig
	}
	if len(a.Author) > maxAuthorLen {
		errors.Author = errorTagetToBig
	}
	if len(a.Slug) > maxSlugLen {
		errors.Slug = errorTagetToBig
	}
	if len(a.Content) > maxContentLen {
		errors.Author = errorTagetToBig
	}
	return a, errors, nil
}
