package public

import (
	"app/pkg/clock"
	"app/pkg/middleware"
	"app/pkg/router"
	"app/pkg/sqlutil"
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
	r.Error(HandleError(true))

	return r.ServeHTTP, nil
}
