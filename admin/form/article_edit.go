package form

import (
	"app/pkg/data"
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
	articleEditEditedBlockID   = "edited_block_id"
	articleEditNewBlock        = "new_block"
	articleEditAction          = "action"
	ArticleEditActionMetadata  = "metadata"
	ArticleEditActionContent   = "content"
	ArticleEditActionBlockEdit = "block_edit"
	ArticleEditActionBlockAdd  = "block_add"
)

type ParsedArticleEdit struct {
	Title               string
	Author              string
	Draft               bool
	Slug                string
	LayoutID            int64
	Content             string
	EditedBlockID       int64
	EditedBlockData     map[string]any
	EditedBlockPosition int
	BlockID             int64
	Action              string
}

type ParsedArticleEditError struct {
	Title               string
	Author              string
	Slug                string
	Content             string
	LayoutID            string
	EditedBlockID       string
	EditedBlockData     string
	EditedBlockPosition string
	AddedBlockID        string
	Action              string
	Form                string
}

func (pe ParsedArticleEditError) HasError() bool {
	if pe.Form != "" {
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
	if be.EditedBlockPosition != "" {
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

type ParseArticleEditParam struct {
	Request         *http.Request
	CheckLayoutID   func(context.Context, int64) (int, error)
	CheckBlockID    func(context.Context, int64) (int, error)
	GetPreviousData func(int64) (map[string]any, bool)
}

func ParseArticleEdit(param ParseArticleEditParam) (ParsedArticleEdit, ParsedArticleEditError, error) {
	err := param.Request.ParseForm()
	if err != nil {
		return ParsedArticleEdit{}, ParsedArticleEditError{}, fmt.Errorf("cannot parse form : %w", err)
	}
	switch param.Request.PostForm.Get(articleEditAction) {
	case ArticleEditActionMetadata:
		return parseArticleEditMetadata(param.Request, param.CheckLayoutID)
	case ArticleEditActionContent:
		return parseArticleEditContent(param.Request)
	case ArticleEditActionBlockEdit:
		return parseArticleEditBlockEdit(param.Request, param.GetPreviousData)
	case ArticleEditActionBlockAdd:
		return parseArticleEditBlockAdd(param.Request, param.CheckBlockID)
	}
	return ParsedArticleEdit{}, ParsedArticleEditError{
		Form: errorInvalidAction,
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
	c, err := check_id(r.Context(), id)
	if err != nil {
		return a, errors, err
	}
	if c == 0 {
		errors.LayoutID = errorInvalidId
	}
	return a, errors, nil
}

func parseArticleEditContent(r *http.Request) (ParsedArticleEdit, ParsedArticleEditError, error) {
	pae := ParsedArticleEdit{
		Content: r.PostForm.Get(articleEditContent),
		Action:  ArticleEditActionContent,
	}
	errors := ParsedArticleEditError{Action: ArticleEditActionContent}
	if pae.Content == "" {
		errors.Content = errorCannotBeEmpty
	}
	if len(pae.Content) > maxContentLen {
		errors.Content = errorTagetToBig
	}
	return pae, errors, nil
}

func parseArticleEditBlockEdit(r *http.Request, get_previous_data func(int64) (map[string]any, bool)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditEditedBlockID), 10, 64)
	a := ParsedArticleEdit{
		Action:  ArticleEditActionBlockEdit,
		BlockID: id,
	}
	errors := ParsedArticleEditError{Action: ArticleEditActionBlockEdit}
	form_data, ok := get_previous_data(id)
	if !ok {
		errors.EditedBlockID = errorInvalidId
		return a, errors, nil
	}
	form_data, err := data.ParseFormData(r, form_data)
	if err != nil {
		errors.EditedBlockData = errorInvalidJson
		return a, errors, nil
	}
	a.EditedBlockID = id
	a.EditedBlockData = form_data
	return a, errors, nil
}

func parseArticleEditBlockAdd(r *http.Request, check_id func(context.Context, int64) (int, error)) (ParsedArticleEdit, ParsedArticleEditError, error) {
	id, _ := strconv.ParseInt(r.PostForm.Get(articleEditNewBlock), 10, 64)
	a := ParsedArticleEdit{
		Action:  ArticleEditActionBlockAdd,
		BlockID: id,
	}
	errors := ParsedArticleEditError{Action: ArticleEditActionBlockAdd}
	c, err := check_id(r.Context(), id)
	if err != nil {
		return a, errors, err
	}
	if c == 0 {
		errors.AddedBlockID = errorInvalidId
	}
	return a, errors, nil
}
