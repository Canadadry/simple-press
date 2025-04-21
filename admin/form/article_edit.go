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
	articleEditTitle           = "title"
	articleEditAuthor          = "author"
	articleEditDraft           = "draft"
	articleEditContent         = "content"
	articleEditSlug            = "slug"
	articleEditLayout          = "layout"
	articleEditAction          = "action"
	ArticleEditActionMetadata  = "metadata"
	ArticleEditActionContent   = "content"
	ArticleEditActionBlockEdit = "block_edit"
	ArticleEditActionBlockAdd  = "block_add"
)

type ParsedArticleEdit struct {
	Title           string
	Author          string
	Draft           bool
	Slug            string
	LayoutID        int64
	Content         string
	EditedBlockID   int64
	EditedBlockData map[string]any
	BlockID         int64
	Action          string
}

type ParsedArticleEditError struct {
	Title           string
	Author          string
	Slug            string
	Content         string
	LayoutID        string
	EditedBlockID   string
	EditedBlockData string
	AddedBlockID    string
	Action          string
	ActionError     string
}

func (pe ParsedArticleEditError) HasError() bool {
	if pe.ActionError != "" {
		return true
	}
	switch pe.Action {
	case ArticleEditActionMetadata:
		return pe.HasMetadataError()
	case ArticleEditActionContent:
		if pe.Content != "" {
			return true
		}
	case ArticleEditActionBlockEdit:
		return pe.HasBlockDataError()
	case ArticleEditActionBlockAdd:
		if pe.AddedBlockID != "" {
			return true
		}
	}
	return false
}

func (be ParsedArticleEditError) HasBlockDataError() bool {
	if be.EditedBlockID != "" {
		return true
	}
	if be.EditedBlockData != "" {
		return true
	}
	return false
}

func (ae ParsedArticleEditError) HasMetadataError() bool {
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

func ParseArticleEdit(r *http.Request, check_layout_id, check_block_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	err := r.ParseForm()
	if err != nil {
		return ParsedArticleEdit{}, ParsedArticleEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	switch r.PostForm.Get(articleEditAction) {
	case ArticleEditActionMetadata:
		return parseArticleEditMetadata(r, check_layout_id)
	case ArticleEditActionContent:
		return parseArticleEditContent(r)
	case ArticleEditActionBlockEdit:
		return parseArticleEditBlockEdit(r, check_block_id)
	case ArticleEditActionBlockAdd:
		return parseArticleEditBlockAdd(r, check_block_id)
	}
	return ParsedArticleEdit{}, ParsedArticleEditError{
		ActionError: errorInvalidAction,
	}, nil
}

func parseArticleEditMetadata(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditLayout), 10, 64)
	a := ParsedArticleEdit{
		Title:    r.PostForm.Get(articleEditTitle),
		Author:   r.PostForm.Get(articleEditAuthor),
		Slug:     r.PostForm.Get(articleEditSlug),
		LayoutID: id,
		Draft:    r.PostForm.Get(articleEditDraft) != "",
		Action:   ArticleEditActionMetadata,
	}
	errors := ParsedArticleEditError{}
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
	if c, err := check_id(r.Context(), id); c == 0 {
		fmt.Println(r.PostForm.Get(articleEditLayout), id, c, err)
		errors.LayoutID = errorInvalidId
	}
	return a, errors, nil
}

func parseArticleEditContent(r *http.Request) (ParsedArticleEdit, ParsedArticleEditError, error) {
	pae := ParsedArticleEdit{
		Content: r.PostForm.Get(articleEditContent),
		Action:  ArticleEditActionContent,
	}
	errors := ParsedArticleEditError{}
	if pae.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if len(pae.Content) > maxContentLen {
		errors.Content = errorTagetToBig
	}
	return pae, errors, nil
}

func parseArticleEditBlockEdit(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	return ParsedArticleEdit{}, ParsedArticleEditError{}, nil
}

func parseArticleEditBlockAdd(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	return ParsedArticleEdit{}, ParsedArticleEditError{}, nil
}
