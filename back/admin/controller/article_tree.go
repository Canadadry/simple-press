package controller

import (
	"app/admin/repository"
	"app/admin/view"
	"app/pkg/http/httpresponse"
	"app/pkg/router"
	"app/pkg/sqlutil"
	"fmt"
	"net/http"
	"strings"
)

func (c *Controller) TreeArticle(w http.ResponseWriter, r *http.Request) error {
	path := router.GetField(r, 0)
	if strings.Contains(path, "..") {
		return httpresponse.NotFound(w)
	}
	query := path
	if len(query) > 0 && query[len(query)-1] != '/' {
		query = query + "/"
	}
	articles, folders, err := c.Repository.SelectArticleTree(r.Context(), query)
	if err != nil {
		return fmt.Errorf("cannot fetch article tree : %w", err)
	}
	if len(articles) == 0 && len(folders) == 0 {
		return httpresponse.NotFound(w)
	}
	return view.ArticlesTreeOk(w, view.ArticleTreeData{
		Path: path,
		Articles: sqlutil.Map(articles, func(from repository.Article) view.ArticleListData {
			return view.ArticleListData{
				Title:   from.Title,
				Date:    from.Date,
				Author:  from.Author,
				Content: from.Content,
				Slug:    from.Slug,
				Draft:   from.Draft,
			}
		}),
		Folders: folders,
	})
}
