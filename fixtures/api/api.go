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
	type AddArticleRsp struct {
		Slug string `json:"slug"`
	}
	article := AddArticleRsp{}
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
