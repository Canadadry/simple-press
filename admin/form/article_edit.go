package form

import (
	"app/pkg/router"
	"fmt"
	"net/http"
	"regexp"
)

const (
	articleEditTitle   = "title"
	articleEditAuthor  = "author"
	articleEditDraft   = "draft"
	articleEditContent = "content"
	articleEditSlug    = "slug"
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
	if ae.Title != "" {
		return true
	}
	if ae.Author != "" {
		return true
	}
	if ae.Content != "" {
		return true
	}
	if ae.Slug != "" {
		return true
	}
	return false
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
		errors.Title = errorCannotBeEmpty
	}
	if a.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if a.Slug == "" {
		errors.Slug = errorCannotBeEmpty
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
	if len(a.Slug) > maxSlugLen {
		errors.Slug = errorTagetToBig
	}
	if len(a.Content) > maxContentLen {
		errors.Author = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.SlugRegexp + "$")
	if !re.Match([]byte(a.Slug)) {
		errors.Slug = errorNotASlug
	}
	return a, errors, nil
}
