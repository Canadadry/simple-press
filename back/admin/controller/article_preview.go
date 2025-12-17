package controller

import (
	"app/page"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"context"
	"fmt"
	"net/http"
)

func (c *Controller) getPages(ctx context.Context, query string, offset, limit int) []page.Page {
	list, err := c.Repository.GetArticleList(ctx, query, limit, offset)
	if err != nil {
		return nil
	}

	pages := []page.Page{}
	for _, a := range list {
		pages = append(pages, page.Page{
			Title:       a.Title,
			Author:      a.Author,
			Description: a.Content,
			Slug:        a.Slug,
		})
	}
	return pages
}

func (c *Controller) GetArticlePreview(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}

	baseTemplates, err := c.Repository.SelectAllTemplate(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base template : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseTemplates {
		files[l.Name] = l.Content
	}
	layout, ok, err := c.Repository.SelectLayoutByID(r.Context(), a.LayoutID)
	if err != nil {
		return fmt.Errorf("cannot select layout %d : %w", a.LayoutID, err)
	}
	if !ok {
		return httpresponse.NotFound(w)
	}
	files[layout.Name] = layout.Content
	keys := []string{}
	for name := range files {
		keys = append(keys, name)
	}

	blocks, err := c.Repository.SelectAllBlock(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select all layouts : %w", err)
	}
	blockSelector := map[string]string{}
	for _, b := range blocks {
		blockSelector[b.Name] = b.Content
	}

	blockDatas, err := c.Repository.SelectBlockDataByArticle(r.Context(), a.ID)
	if err != nil {
		return fmt.Errorf("cannot select block data article : %w", err)
	}

	blockDataView := []page.ArticleBlock{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, page.ArticleBlock{
			Position: int(p.Position), Data: p.Data, BlockName: p.BlockName,
		})
	}

	err = page.Render(w, page.Data{
		Title:         a.Title,
		Content:       a.Content,
		Files:         files,
		BlocksContent: blockSelector,
		ArticleBlocks: blockDataView,
		PageFtecher: func(query string, offset int, limit int) []page.Page {
			return c.getPages(r.Context(), query, offset, limit)
		},
	})
	if err != nil {
		return fmt.Errorf("cannot render article : %w", err)
	}
	return nil
}
