package public

import (
	"app/pkg/clock"
	"app/pkg/middleware"
	"app/pkg/router"
	"app/pkg/sqlutil"
	"app/public/controller"
	"app/public/repository"
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
	r.Use(middleware.Logger(services.Out, services.Clock, func(path string) bool {
		return false
	}))
	r.Use(middleware.AutoCloseRequestBody)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer())
	r.Error(HandleError(true))

	c, err := controller.New(repository.Repository{
		Db: services.Db,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create admin controller : %w", err)
	}

	r.Get("/:slug", c.GetArticlePreview)

	return r.ServeHTTP, nil
}
