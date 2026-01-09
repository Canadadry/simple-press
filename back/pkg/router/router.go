package router

import (
	"app/pkg/stacktrace"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	DigitRegExp      = "([0-9]+)"
	UuidV4Regexp     = "([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})"
	EmailRegexp      = "(\\S+@\\S+\\.\\S+)"
	StringRegExp     = "([^/]+)"
	TokenRegexp      = "([a-zA-Z0-9]+)"
	MultiTokenRegexp = "((?:[a-zA-Z0-9_]+)(?:;(?:[a-zA-Z0-9_]+))*)"
	SlugRegexp       = "([-_a-zA-Z0-9\\.]+)"
	JwtRegexp        = "([a-zA-Z0-9_=]+\\.[a-zA-Z0-9_=]+\\.[a-zA-Z0-9_\\-\\+\\/=]*)"
	Base64Regexp     = "((?:[A-Za-z\\d+/]{4})*(?:[A-Za-z\\d+/]{3}=|[A-Za-z\\d+/]{2}==)?)"
	PathRegexp       = "(([A-Za-z0-9_-][A-Za-z0-9_.-]*)(\\/([A-Za-z0-9_-][A-Za-z0-9_.-]*))*)"
	AnyRegexp        = "(.*)"
)

var tags = map[string]string{
	":digit":    DigitRegExp,
	":string":   StringRegExp,
	":uuid":     UuidV4Regexp,
	":email":    EmailRegexp,
	":token":    TokenRegexp,
	":multitok": MultiTokenRegexp,
	":slug":     SlugRegexp,
	":jwt":      JwtRegexp,
	":b64":      Base64Regexp,
	":path":     PathRegexp,
	":any":      AnyRegexp,
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func NopHandler(fn http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		fn(w, r)
		return nil
	}
}

type Route struct {
	name    string
	method  string
	regex   *regexp.Regexp
	handler HandlerFunc
}

func newRoute(method, pattern string, handler HandlerFunc) Route {
	name := method + ":" + pattern
	for k, v := range tags {
		pattern = strings.ReplaceAll(pattern, k, v)
	}
	if strings.Contains(pattern, ":") {
		panic("cannot replace all pattern in " + name)
	}
	return Route{name, method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Get(pattern string, handler HandlerFunc) Route {
	return newRoute(http.MethodGet, pattern, handler)
}
func Post(pattern string, handler HandlerFunc) Route {
	return newRoute(http.MethodPost, pattern, handler)
}
func Patch(pattern string, handler HandlerFunc) Route {
	return newRoute(http.MethodPatch, pattern, handler)
}
func Delete(pattern string, handler HandlerFunc) Route {
	return newRoute(http.MethodDelete, pattern, handler)
}

func Merge(sliceOfRoutes ...[]Route) []Route {
	totalNumberOfRoute := 0
	for _, routes := range sliceOfRoutes {
		totalNumberOfRoute += len(routes)
	}
	merged := make([]Route, 0, totalNumberOfRoute)
	for _, routes := range sliceOfRoutes {
		merged = append(merged, routes...)
	}
	return merged
}

type RoutingErrorHandler func(statusCode int, err error, w http.ResponseWriter, r *http.Request)

type CorsOption struct {
	Origin  string
	Headers []string
	Methods []string
}

func ServeRoutesAndHandleErrorWith(routes []Route, errorHandler RoutingErrorHandler, cors CorsOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allow []string
		matched := []string{}
		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				// fmt.Println("router has", route.name, "got", r.Method, r.URL.Path)
				matched = append(matched, route.name)
				if r.Method != route.method {
					allow = append(allow, route.method)
					continue
				}
				w.Header().Add("Access-Control-Allow-Origin", cors.Origin)
				w.Header().Add("Access-Control-Allow-Credentials", "true")
				ctx := context.WithValue(r.Context(), ctxKey{}, ctxValue{matches, route.name})
				err := route.handler(w, r.WithContext(ctx))
				if err != nil {
					errorHandler(
						http.StatusInternalServerError,
						err, w, r,
					)
				}
				return
			}
		}
		if len(allow) == 0 {
			err := fmt.Errorf("route not found")
			errorHandler(http.StatusNotFound, err, w, r)
			return
		}

		if r.Method == http.MethodOptions && cors.Origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", cors.Origin)
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Headers", strings.Join(cors.Headers, ", "))
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(cors.Methods, ", "))
			w.WriteHeader(http.StatusOK)
			return
		}

		allowed := strings.Join(allow, ", ")
		w.Header().Set("Allow", allowed)
		fmt.Println("##### not alloew on ######", matched)
		err := stacktrace.Errorf(
			"method not allowed : allow only %s referer %s useragent %s",
			allowed,
			r.Header.Get("Referer"),
			r.Header.Get("User-Agent"),
		)
		errorHandler(http.StatusMethodNotAllowed, err, w, r)
	}
}

type ctxKey struct{}
type ctxValue struct {
	Value   []string
	Pattern string
}

func GetField(r *http.Request, index int) string {
	value, _ := r.Context().Value(ctxKey{}).(ctxValue)
	if len(value.Value) <= index+1 {
		return ""
	}
	return value.Value[index+1]
}

func GetFieldAsInt(r *http.Request, index int) (int, error) {
	return strconv.Atoi(GetField(r, index))
}

func GetPattern(r *http.Request) string {
	value, _ := r.Context().Value(ctxKey{}).(ctxValue)
	return value.Pattern
}
