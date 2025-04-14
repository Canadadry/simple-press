package router

import (
	"app/pkg/http/httpresponse"
	"app/pkg/urlutil"
	"fmt"
	"net/http"
)

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

		msg := fmt.Sprintf("cannot serve %s:%s : error %d\n%s", r.Method, urlutil.GetFullPath(r.URL), statusCode, err)
		if devMode {
			fmt.Fprintf(w, "%s", msg)
			return
		}
		httpresponse.ServerError(w)
	}
}
