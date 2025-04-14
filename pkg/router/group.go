package router

import (
	"fmt"
	"net/http"
)

type middleware func(HandlerFunc) HandlerFunc

type Group struct {
	routes       []Route
	middlewares  []middleware
	errorHandler RoutingErrorHandler
	Cors         CorsOption
}

func (g *Group) Error(errorHandler RoutingErrorHandler) {
	g.errorHandler = errorHandler
}

func (g *Group) Use(fn middleware) {
	g.middlewares = append(g.middlewares, fn)
}

func (g *Group) add(r Route) {
	h := r.handler
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		fn := g.middlewares[i]
		h = fn(h)
	}
	r.handler = h
	g.routes = append(g.routes, r)
}

func (g *Group) Get(pattern string, handler HandlerFunc) {
	g.add(Get(pattern, handler))
}
func (g *Group) Post(pattern string, handler HandlerFunc) {
	g.add(Post(pattern, handler))
}
func (g *Group) Patch(pattern string, handler HandlerFunc) {
	g.add(Patch(pattern, handler))
}
func (g *Group) Delete(pattern string, handler HandlerFunc) {
	g.add(Delete(pattern, handler))
}

func (g *Group) Group(init func(g *Group)) {
	sub := &Group{
		middlewares: append([]middleware(nil), g.middlewares...),
	}
	init(sub)
	g.routes = Merge(g.routes, sub.routes)
}

func (g *Group) Mount(pattern string, h http.Handler) {
	g.Get(pattern+".*", NopHandler(h.ServeHTTP))
}

func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if g.errorHandler == nil {
		g.errorHandler = func(statusCode int, err error, w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Cannot serve %s:%s : error %d\n%s", r.Method, r.URL.Path, statusCode, err)
		}
	}
	h := ServeRoutesAndHandleErrorWith(g.routes, g.errorHandler, g.Cors)
	h(w, r)
}
