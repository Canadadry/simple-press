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
	articleEditActionMetadata  = "metadata"
	articleEditActionContent   = "content"
	articleEditActionBlockEdit = "block_edit"
	articleEditActionBlockAdd  = "block_add"
)

type ParsedArticleEdit struct {
	Metadata  ArticleMetadataEdit
	Content   string
	BlockData BlockData
	BlockID   int64
	Action    string
}

type BlockData struct {
	ID        int64
	BlockData map[string]any
}

type ArticleMetadataEdit struct {
	Title    string
	Author   string
	Draft    bool
	Slug     string
	LayoutID int64
}

type ParsedArticleEditError struct {
	Metadata    ArticleMetadataEditError
	Content     string
	BlockData   BlockDataError
	BlockID     string
	Action      string
	ActionError string
}

func (pe ParsedArticleEditError) HasError() bool {
	switch pe.Action {
	case articleEditActionMetadata:
		return pe.Metadata.HasError()
	case articleEditActionContent:
		if pe.Content != "" {
			return true
		}
	case articleEditActionBlockEdit:
		return pe.BlockData.HasError()
	case articleEditActionBlockAdd:
		if pe.BlockID != "" {
			return true
		}
	}
	return false
}

type BlockDataError struct {
	ID        string
	BlockData string
}

func (be BlockDataError) HasError() bool {
	if be.ID != "" {
		return true
	}
	if be.BlockData != "" {
		return true
	}
	return false
}

type ArticleMetadataEditError struct {
	Title    string
	Author   string
	Slug     string
	LayoutID string
}

func (ae ArticleMetadataEditError) HasError() bool {
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
	case articleEditActionMetadata:
		return parseArticleEditMetadata(r, check_layout_id)
	case articleEditActionContent:
		return parseArticleEditContent(r)
	case articleEditActionBlockEdit:
		return parseArticleEditBlockEdit(r, check_block_id)
	case articleEditActionBlockAdd:
		return parseArticleEditBlockAdd(r, check_block_id)
	}
	return ParsedArticleEdit{}, ParsedArticleEditError{
		ActionError: errorInvalidAction,
	}, nil
}

func parseArticleEditMetadata(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditLayout), 10, 64)
	a := ParsedArticleEdit{
		Metadata: ArticleMetadataEdit{
			Title:    r.PostForm.Get(articleEditTitle),
			Author:   r.PostForm.Get(articleEditAuthor),
			Slug:     r.PostForm.Get(articleEditSlug),
			LayoutID: id,
			Draft:    r.PostForm.Get(articleEditDraft) != "",
		},
		Action: articleEditActionMetadata,
	}
	errors := ParsedArticleEditError{}
	if a.Metadata.Title == "" {
		errors.Metadata.Title = errorCannotBeEmpty
	}
	if a.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if a.Metadata.Slug == "" {
		errors.Metadata.Slug = errorCannotBeEmpty
	}
	if a.Metadata.Author == "" {
		errors.Metadata.Author = errorCannotBeEmpty
	}
	if len(a.Metadata.Title) > maxTitleLen {
		errors.Metadata.Title = errorTagetToBig
	}
	if len(a.Metadata.Author) > maxAuthorLen {
		errors.Metadata.Author = errorTagetToBig
	}
	if len(a.Metadata.Slug) > maxSlugLen {
		errors.Metadata.Slug = errorTagetToBig
	}
	re := regexp.MustCompile("^" + router.SlugRegexp + "$")
	if !re.Match([]byte(a.Metadata.Slug)) {
		errors.Metadata.Slug = errorNotASlug
	}
	if c, err := check_id(r.Context(), id); c == 0 {
		fmt.Println(r.PostForm.Get(articleEditLayout), id, c, err)
		errors.Metadata.LayoutID = errorInvalidId
	}
	return a, errors, nil
}

func parseArticleEditContent(r *http.Request) (ParsedArticleEdit, ParsedArticleEditError, error) {
	return ParsedArticleEdit{}, ParsedArticleEditError{}, nil
}

func parseArticleEditBlockEdit(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	return ParsedArticleEdit{}, ParsedArticleEditError{}, nil
}

func parseArticleEditBlockAdd(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	return ParsedArticleEdit{}, ParsedArticleEditError{}, nil
}
