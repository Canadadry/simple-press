package public

import (
	"app/pkg/http/httpresponse"
	"app/pkg/stacktrace"
	"app/pkg/urlutil"
	"net/http"
)

func BadCredentialHandler(w http.ResponseWriter, r *http.Request) error {
	return httpresponse.Unauthorized(w)
}

func HandleError(devMode bool) func(int, error, http.ResponseWriter, *http.Request) {
	return func(statusCode int, err error, w http.ResponseWriter, r *http.Request) {
		if statusCode == http.StatusNotFound {
			httpresponse.RouteNotFound(w, urlutil.GetFullPath(r.URL))
			return
		}
		if statusCode == http.StatusMethodNotAllowed {
			httpresponse.MethodNotAllowed(w)
			return
		}
		var trace any
		stError, ok := err.(*stacktrace.Error)
		if ok {
			trace = stError.Frames
		}
		msg := map[string]any{
			"method":      r.Method,
			"url":         urlutil.GetFullPath(r.URL),
			"status_code": statusCode,
			"error":       err.Error(),
			"trace":       trace,
		}

		if devMode {
			httpresponse.Json(w, http.StatusInternalServerError, msg)
			return
		}
		httpresponse.ServerError(w)
	}
}
