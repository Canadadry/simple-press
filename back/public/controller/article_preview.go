package controller

import (
	"app/page"
	"app/pkg/router"
	"fmt"
	"net/http"
)

func (c *Controller) GetArticlePreview(w http.ResponseWriter, r *http.Request) error {
	slug := router.GetField(r, 0)
	if slug == "" {
		slug = "index"
	}
	a, ok, err := c.Repository.SelectArticleBySlug(r.Context(), slug)
	if err != nil {
		return fmt.Errorf("cannot select article : %w", err)
	}
	if !ok {
		return c.GetFile(w, r)
	}
	baseTemplates, err := c.Repository.SelectAllTemplate(r.Context())
	if err != nil {
		return fmt.Errorf("cannot select base template : %w", err)
	}
	files := map[string]string{}
	for _, l := range baseTemplates {
		files[l.Name] = l.Content
	}
	pageLayout, ok, err := c.Repository.SelectLayoutByID(r.Context(), a.LayoutID)
	if err != nil {
		return fmt.Errorf("cannot select page layout %d : %w", a.LayoutID, err)
	}
	if !ok {
		return fmt.Errorf("cannot found page layout %d : %w", a.LayoutID, err)
	}
	files[pageLayout.Name] = pageLayout.Content

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
		return fmt.Errorf("cannot select block dataarticle : %w", err)
	}

	blockDataView := []page.ArticleBlock{}
	for _, p := range blockDatas {
		blockDataView = append(blockDataView, page.ArticleBlock{
			Position: int(p.Position), Data: p.Data, BlockName: p.BlockName,
		})
	}
	return page.Render(w, page.Data{
		Title:         a.Title,
		Content:       a.Content,
		Files:         files,
		BlocksContent: blockSelector,
		ArticleBlocks: blockDataView,
	})
}
