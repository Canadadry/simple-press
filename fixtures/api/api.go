package api

import (
	"app/admin/serializer"
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
	article := serializer.ArticleAdded{}
	rsp := map[int]any{
		http.StatusCreated:    &article,
		http.StatusBadRequest: nil,
	}
	st, err := c.client.Post(c.ctx, "/admin/article/add", serializer.ArticleAdded{
		Title:  title,
		Author: author,
	}, rsp)
	if err != nil {
		return "", fmt.Errorf("cannot add article : %w", err)
	}
	if st != http.StatusCreated {
		return "", fmt.Errorf("cannot add article invalid status code %d", st)
	}
	return article.Slug, nil
}

func (c *Client) AddLayout(name string) (int64, error) {
	layout := serializer.LayoutAdded{}
	rsp := map[int]any{
		http.StatusCreated:    &layout,
		http.StatusBadRequest: nil,
	}
	st, err := c.client.Post(c.ctx, "/admin/layout/add", serializer.LayoutAdded{
		Name: name,
	}, rsp)
	if err != nil {
		return 0, fmt.Errorf("cannot add layout : %w", err)
	}
	if st != http.StatusCreated {
		return 0, fmt.Errorf("cannot add layout invalid status code %d", st)
	}
	return layout.ID, nil
}
