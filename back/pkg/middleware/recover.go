package middleware

import (
	"app/pkg/router"
	"app/pkg/trycatch"
	"net/http"
)

func Recoverer() func(next router.HandlerFunc) router.HandlerFunc {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			return trycatch.Catch(func() error {
				return next(w, r)
			})
		}
	}
}
