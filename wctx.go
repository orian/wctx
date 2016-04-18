package wctx

import (
	"github.com/julienschmidt/httprouter"
	paramsHelper "github.com/orian/params"
	"golang.org/x/net/context"

	"log"
	"net/http"
)

type HandleFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandleFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type ContextProvider func(*http.Request) context.Context

type Router struct {
	R              *httprouter.Router
	ContextFactory ContextProvider
	m              []Middleware
}

var paramsKey = "params key"

func (r *Router) wrap(h HandleFunc) httprouter.Handle {
	// builds wrapper
	for i := len(r.m) - 1; i >= 0; i-- {
		h = r.m[i](h)
	}
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx := r.ContextFactory(req)
		if p != nil {
			ctx = context.WithValue(ctx, &paramsKey, p)
		}
		h(ctx, w, req)
	})
}

func FromContext(c context.Context) (paramsHelper.Params, bool) {
	p, ok := c.Value(&paramsKey).(httprouter.Params)
	return paramsHelper.NewFromHttpRouter(p), ok
}

func DefaultContextFactory(req *http.Request) context.Context {
	return context.TODO()
}

func New() *Router {
	return &Router{httprouter.New(), DefaultContextFactory, nil}
}

func (r *Router) DELETE(path string, handle HandleFunc) {
	r.R.DELETE(path, r.wrap(handle))
}

func (r *Router) GET(path string, handle HandleFunc) {
	r.R.GET(path, r.wrap(handle))
}

func (r *Router) HEAD(path string, handle HandleFunc) {
	r.R.HEAD(path, r.wrap(handle))
}

func (r *Router) Handle(method, path string, handle HandleFunc) {
	r.R.Handle(method, path, r.wrap(handle))
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	if len(r.m) > 0 {
		log.Printf("warning: when registering standard http.Handler the middlewares are skipped")
	}
	r.R.Handler(method, path, handler)
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	if len(r.m) > 0 {
		log.Printf("warning: when registering standard http.HandleFunc the middlewares are skipped")
	}
	r.R.HandlerFunc(method, path, handler)
}

func (r *Router) Lookup(method, path string) (httprouter.Handle, httprouter.Params, bool) {
	return r.R.Lookup(method, path)
}

func (r *Router) OPTIONS(path string, handle HandleFunc) {
	r.R.OPTIONS(path, r.wrap(handle))
}

func (r *Router) PATCH(path string, handle HandleFunc) {
	r.R.PATCH(path, r.wrap(handle))
}

func (r *Router) POST(path string, handle HandleFunc) {
	r.R.POST(path, r.wrap(handle))
}

func (r *Router) PUT(path string, handle HandleFunc) {
	r.R.PUT(path, r.wrap(handle))
}

func (r *Router) ServeFiles(path string, root http.FileSystem) {
	r.R.ServeFiles(path, root)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.R.ServeHTTP(w, req)
}

type Middleware func(HandleFunc) HandleFunc

func (r *Router) WithMiddleware(m Middleware) *Router {
	newR := *r
	newR.m = make([]Middleware, 0, len(r.m)+1)
	newR.m = append(newR.m, r.m...)
	newR.m = append(newR.m, m)
	return &newR
}
