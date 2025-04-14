package middleware

import (
	"app/pkg/router"
	"net/http"
)

func AutoCloseRequestBody(next router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		return next(w, r)
	}
}
