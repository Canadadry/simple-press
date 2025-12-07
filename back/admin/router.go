package admin

import (
	"app/admin/assets"
	"app/admin/controller"
	"app/admin/repository"
	"app/pkg/clock"
	"app/pkg/middleware"
	"app/pkg/router"
	"app/pkg/sqlutil"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Services struct {
	Db    sqlutil.DBTX
	Clock clock.Clock
	Out   io.Writer
}

func GetRouter(services Services) (http.HandlerFunc, error) {
	r := router.Group{
		Cors: router.CorsOption{
			Origin:  "http://localhost:5173",
			Methods: []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete},
			Headers: []string{"Content-Type"},
		},
	}
	r.Use(middleware.Logger(services.Out, services.Clock, func(path string) bool {
		return strings.HasPrefix(path, "/public/")
	}))
	r.Use(middleware.AutoCloseRequestBody)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer())
	r.Mount("/public/", http.FileServer(http.FS(assets.GetPublicFiles())))
	r.Error(HandleError(true))

	c, err := controller.New(repository.Repository{
		Db:    services.Db,
		Clock: services.Clock,
	}, services.Clock)
	if err != nil {
		return nil, fmt.Errorf("cannot create admin controller : %w", err)
	}

	r.Get("/admin", c.GetDashboard)
	r.Get("/admin/articles", c.GetArticleList)
	r.Get("/admin/articles/add", c.GetArticleAdd)
	r.Post("/admin/articles/add", c.PostArticleAdd)
	r.Get("/admin/articles/:slug/edit", c.GetArticleEdit)
	r.Post("/admin/articles/:slug/edit/metadata", c.PostArticleEditMetadata)
	r.Post("/admin/articles/:slug/edit/content", c.PostArticleEditContent)
	r.Post("/admin/articles/:slug/edit/block_edit", c.PostArticleEditBlockEdit)
	r.Post("/admin/articles/:slug/edit/block_add", c.PostArticleEditBlockAdd)
	r.Get("/admin/articles/:slug/preview", c.GetArticlePreview)
	r.Get("/admin/templates", c.GetTemplateList)
	r.Get("/admin/templates/add", c.GetTemplateAdd)
	r.Post("/admin/templates/add", c.PostTemplateAdd)
	r.Get("/admin/templates/:path/edit", c.GetTemplateEdit)
	r.Post("/admin/templates/:path/edit", c.PostTemplateEdit)
	r.Get("/admin/blocks", c.GetBlockList)
	r.Get("/admin/blocks/add", c.GetBlockAdd)
	r.Post("/admin/blocks/add", c.PostBlockAdd)
	r.Get("/admin/blocks/:path/edit", c.GetBlockEdit)
	r.Post("/admin/blocks/:path/edit", c.PostBlockEdit)
	r.Get("/admin/layouts", c.GetLayoutList)
	r.Get("/admin/layouts/add", c.GetLayoutAdd)
	r.Post("/admin/layouts/add", c.PostLayoutAdd)
	r.Get("/admin/layouts/:path/edit", c.GetLayoutEdit)
	r.Post("/admin/layouts/:path/edit", c.PostLayoutEdit)
	r.Get("/admin/files", c.GetFileList)
	r.Get("/admin/file/add", c.GetFileAdd)
	r.Post("/admin/file/add", c.PostFileAdd)
	r.Get("/admin/files/:slug", c.GetFile)
	r.Get("/:slug", c.GetFile)

	return r.ServeHTTP, nil
}
