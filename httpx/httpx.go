package httpx

import (
	"context"
	"net/http"

	"github.com/josestg/httprouter"
)

// contextType represents a type for any context for httpx.
type contextType struct {
	name string
}

// Params contains request path parameters.
type Params struct {
	p httprouter.Params
}

// Get gets params by var name.
func (p Params) Get(name string) string { return p.p.ByName(name) }

var paramsContextKey = &contextType{name: "params"}

// ParamsFromContext gets Params from context.
func ParamsFromContext(ctx context.Context) Params {
	p, _ := ctx.Value(paramsContextKey).(Params)
	return p
}

// contextWithParams creates anew context with params in it.
func contextWithParams(ctx context.Context, p Params) context.Context {
	return context.WithValue(ctx, paramsContextKey, p)
}

// Handler is just like http.Handler but this Handler returns an error.
// This is useful to remove redundancies when handling an error into
// http error response.
type Handler interface {
	// ServeHTTP handles incoming http request a response to w.
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// HandlerFunc is an adapter that enabled an ordinary function
// to implement the Handler interface.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP calls fn(w, r)
func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error { return fn(w, r) }

// ServeMux is a http router.
type ServeMux struct {
	internal *httprouter.Router
	chain    *Chain
}

// NewServeMux creates a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		internal: httprouter.New(),
		chain:    NewChain(),
	}

}

// NewServeMuxWithChain creates a new ServeMux with a chain of middlewares.
func NewServeMuxWithChain(chain *Chain) *ServeMux {
	return &ServeMux{
		internal: httprouter.New(),
		chain:    chain,
	}
}

// Handle registers a new Handler.
func (mux *ServeMux) Handle(method, path string, handler Handler, middlewares ...Middleware) {
	mux.internal.Handle(method, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		ctx := contextWithParams(r.Context(), Params{p: p})
		return mux.chain.Extend(middlewares...).Then(handler).ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandleFunc registers an ordinary function as a Handler.
func (mux *ServeMux) HandleFunc(method, path string, fn HandlerFunc, middlewares ...Middleware) {
	mux.Handle(method, path, fn, middlewares...)
}

// ServeHTTP implements the http.Handler to make it compatible with
// net/http Handler.
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.internal.ServeHTTP(w, r)
}
