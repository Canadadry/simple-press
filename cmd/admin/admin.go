package admin

import (
	"app/admin"
	"app/config"
	"app/pkg/clock"
	"app/pkg/dbconn"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	Action = "admin"
)

func Run(c config.Parameters) error {
	db, err := dbconn.Open(c.DatabaseUrl)
	if err != nil {
		return err
	}

	var out io.Writer
	out = os.Stdout
	if c.Out != "" {
		f, err := os.OpenFile(c.Out, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("cant open out : %w", err)
		}
		defer f.Close()
		out = f
	}

	rt, err := admin.GetRouter(admin.Services{
		Db:    db,
		Clock: clock.Real{},
		Out:   out,
	})

	if err != nil {
		return err
	}
	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", c.Port),
		Handler: rt,
	}

	fmt.Println("starting admin server", "endpoint", fmt.Sprintf(":%d", c.Port))

	return server.ListenAndServe()
}
