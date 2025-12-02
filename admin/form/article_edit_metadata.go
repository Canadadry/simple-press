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
	articleEditTitle  = "title"
	articleEditAuthor = "author"
	articleEditDraft  = "draft"
	articleEditSlug   = "slug"
	articleEditLayout = "layout"
)

type ParsedArticleEditMetadata struct {
	Title    string
	Author   string
	Draft    bool
	Slug     string
	LayoutID int64
}

type ParsedArticleEditMetadataError struct {
	Title    string
	Author   string
	Slug     string
	LayoutID string
}

func (ae ParsedArticleEditMetadataError) HasError() bool {
	if ae.Title != "" {
		return true
	}
	if ae.Author != "" {
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

func ParseArticleEditMetadata(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEditMetadata, ParsedArticleEditMetadataError, error) {
	err := r.ParseForm()
	if err != nil {
		return ParsedArticleEditMetadata{}, ParsedArticleEditMetadataError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditLayout), 10, 64)
	a := ParsedArticleEditMetadata{
		Title:    r.PostForm.Get(articleEditTitle),
		Author:   r.PostForm.Get(articleEditAuthor),
		Slug:     r.PostForm.Get(articleEditSlug),
		LayoutID: id,
		Draft:    r.PostForm.Get(articleEditDraft) != "",
	}
	errors := ParsedArticleEditMetadataError{}
	if a.Title == "" {
		errors.Title = errorCannotBeEmpty
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
	re := regexp.MustCompile("^" + router.SlugRegexp + "$")
	if !re.Match([]byte(a.Slug)) {
		errors.Slug = errorNotASlug
	}
	c, err := check_id(r.Context(), id)
	if err != nil {
		return a, errors, err
	}
	if c == 0 {
		errors.LayoutID = errorInvalidId
	}
	return a, errors, nil
}
