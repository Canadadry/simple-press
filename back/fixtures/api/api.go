package api

import (
	"app/admin/view"
	"app/pkg/http/httpcaller"
	"context"
	"fmt"
	"io"
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

func (c *Client) EditArticleContent(slug, content string) error {
	article := view.ArticleEditData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusOK:         &article,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, fmt.Sprintf("/admin/articles/%s/edit/content", slug), view.ArticleEditData{
		Content: content,
	}, rsp)
	if err != nil {
		return fmt.Errorf("cannot edit article content : %w", err)
	}
	if st != http.StatusOK {
		return fmt.Errorf("cannot edit article content invalid status code %d\n%v", st, errs)
	}
	return nil
}

func (c *Client) EditArticleBlockAdd(slug string, id, position int) (int64, error) {
	data := view.ArticleAddBlockData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &data,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, fmt.Sprintf("/admin/articles/%s/edit/block_add", slug), map[string]any{
		"new_block": id,
		"position":  position,
	}, rsp)
	if err != nil {
		return 0, fmt.Errorf("cannot edit article block_add : %w", err)
	}
	if st != http.StatusCreated {
		return 0, fmt.Errorf("cannot edit article block_add invalid status code %d\n%v", st, errs)
	}
	return data.ID, nil
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

func (c *Client) EditLayout(name, content string) error {
	layout := view.LayoutEditData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusOK:         &layout,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, fmt.Sprintf("/admin/layouts/%s/edit", name), view.LayoutEditData{
		Name:    name,
		Content: content,
	}, rsp)
	if err != nil {
		return fmt.Errorf("cannot edit layout : %w", err)
	}
	if st != http.StatusOK {
		return fmt.Errorf("cannot edit layout invalid status code %d\n%v", st, errs)
	}
	return nil
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

func (c *Client) EditTemplate(name, content string) error {
	template := view.TemplateEditData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusOK:         &template,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, fmt.Sprintf("/admin/templates/%s/edit", name), view.TemplateEditData{
		Name:    name,
		Content: content,
	}, rsp)
	if err != nil {
		return fmt.Errorf("cannot edit template : %w", err)
	}
	if st != http.StatusOK {
		return fmt.Errorf("cannot edit template invalid status code %d\n%v", st, errs)
	}
	return nil
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

func (c *Client) EditBlock(data view.BlockEditData) error {
	block := view.BlockEditData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusOK:         &block,
		http.StatusBadRequest: &errs,
	}
	st, err := c.client.Post(c.ctx, fmt.Sprintf("/admin/blocks/%s/edit", data.Name), data, rsp)
	if err != nil {
		return fmt.Errorf("cannot edit block : %w", err)
	}
	if st != http.StatusOK {
		return fmt.Errorf("cannot edit block invalid status code %d\n%v", st, errs)
	}
	return nil
}

func (c *Client) AddFile(filename string, file io.ReadCloser) (int64, error) {
	defer file.Close()
	fileData := view.FileAddData{}
	errs := map[string]any{}
	rsp := map[int]any{
		http.StatusCreated:    &fileData,
		http.StatusBadRequest: &errs,
	}
	r, ct, err := httpcaller.CreateMultiPartForm(map[string][]httpcaller.FormValue{
		"content": []httpcaller.FormValue{
			{Filename: filename, File: file},
		},
		"name": []httpcaller.FormValue{
			{String: filename},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("cannot create multipart request  : %w", err)
	}
	client := c.client.WithHeader("Content-Type", ct).WithHeader("Accept", "application/json")
	st, err := client.Post(c.ctx, "/admin/files/add", r, rsp)
	if err != nil {
		return 0, fmt.Errorf("cannot add files : %w", err)
	}
	if st != http.StatusCreated {
		return 0, fmt.Errorf("cannot add files invalid status code  %d\n%v", st, errs)
	}
	return fileData.ID, nil
}
