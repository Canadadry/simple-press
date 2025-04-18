package form

import (
	"app/pkg/router"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

const (
	articleEditTitle   = "title"
	articleEditAuthor  = "author"
	articleEditDraft   = "draft"
	articleEditContent = "content"
	articleEditSlug    = "slug"
	articleEditLayout  = "layout"
)

type ArticleEdit struct {
	Title    string
	Author   string
	Draft    bool
	Content  string
	Slug     string
	LayoutID int64
}

type ArticleEditError struct {
	Title    string
	Author   string
	Slug     string
	Content  string
	LayoutID string
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
	if ae.LayoutID != "" {
		return true
	}
	return false
}

func ParseArticleEdit(r *http.Request, check_id func(context.Context, int64) (int, error)) (ArticleEdit, ArticleEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return ArticleEdit{}, ArticleEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditLayout), 10, 64)
	a := ArticleEdit{
		Title:    r.PostForm.Get(articleEditTitle),
		Author:   r.PostForm.Get(articleEditAuthor),
		Content:  r.PostForm.Get(articleEditContent),
		Slug:     r.PostForm.Get(articleEditSlug),
		LayoutID: id,
		Draft:    r.PostForm.Get(articleEditDraft) != "",
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
	if c, err := check_id(r.Context(), id); c == 0 {
		fmt.Println(r.PostForm.Get(articleEditLayout), id, c, err)
		errors.LayoutID = errorInvalidId
	}
	return a, errors, nil
}
