package api

import (
	"app/admin/view"
	"app/pkg/http/httpcaller"
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	client httpcaller.Caller
	ctx    context.Context
}

func New(client httpcaller.Caller) *Client {
	return &Client{client: client, ctx: context.Background()}
}

func (c *Client) AddArticle(title, author string) (string, error) {
	article := view.ArticleAddData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &article,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, "/admin/articles/add", view.ArticleAddData{
		Title:  title,
		Author: author,
	}, rsp)
	if err != nil {
		return "", fmt.Errorf("cannot add article : %w", err)
	}
	if st != http.StatusCreated {
		return "", fmt.Errorf("cannot add article invalid status code  %d\n%v", st, errs)
	}
	return article.Slug, nil
}

func (c *Client) AddLayout(name string) (int64, error) {
	layout := view.LayoutAddData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &layout,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, "/admin/layouts/add", view.LayoutAddData{
		Name: name,
	}, rsp)
	if err != nil {
		return 0, fmt.Errorf("cannot add layout : %w", err)
	}
	if st != http.StatusCreated {
		return 0, fmt.Errorf("cannot add layout invalid status code %d\n%v", st, errs)
	}
	return layout.ID, nil
}

func (c *Client) AddTemplate(name string) (int64, error) {
	template := view.TemplateAddData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &template,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, "/admin/templates/add", view.TemplateAddData{
		Name: name,
	}, rsp)
	if err != nil {
		return 0, fmt.Errorf("cannot add template : %w", err)
	}
	if st != http.StatusCreated {
		return 0, fmt.Errorf("cannot add template invalid status code %d\n%v", st, errs)
	}
	return template.ID, nil
}

func (c *Client) AddBlock(name string) (string, error) {
	block := view.BlockAddData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &block,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, "/admin/blocks/add", view.BlockAddData{
		Name: name,
	}, rsp)
	if err != nil {
		return "", fmt.Errorf("cannot add block : %w", err)
	}
	if st != http.StatusCreated {
		return "", fmt.Errorf("cannot add block invalid status code %d\n%v", st, errs)
	}
	return block.Name, nil
}
