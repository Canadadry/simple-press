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
)

type Services struct {
	Db    sqlutil.DBTX
	Clock clock.Clock
	Out   io.Writer
}

func GetRouter(services Services) (http.HandlerFunc, error) {
	r := router.Group{}
	r.Use(middleware.Logger(services.Out, services.Clock))
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

	r.Get("/admin/articles", c.GetArticleList)
	r.Get("/admin/article/add", c.GeArticleAdd)
	r.Post("/admin/article/add", c.PostArticleAdd)
	r.Get("/admin/articles/:slug/edit", c.GeArticleEdit)
	r.Post("/admin/articles/:slug/edit", c.PostArticleEdit)

	return r.ServeHTTP, nil
}
