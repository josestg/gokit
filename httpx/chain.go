package httpx

import (
	"net/http"
)

// Middleware is a function that wraps a Handler.
type Middleware func(h Handler) HandlerFunc

// Chain is a chain of middleware.
type Chain struct {
	middleware []Middleware
}

// NewChain creates a new chain of middleware.
func NewChain(middleware ...Middleware) *Chain {
	return &Chain{middleware: middleware}
}

// Then chains the middleware with h and returns a new Handler.
func (c *Chain) Then(h Handler) Handler {
	return c.chain(h)
}

// ThenFunc chains the middleware with f and returns a new Handler.
func (c *Chain) ThenFunc(f HandlerFunc) Handler {
	return c.Then(f)
}

// Extend extends existing chain with new middlewares, and returns a new copy of
// chain.
func (c *Chain) Extend(middleware ...Middleware) *Chain {
	return &Chain{middleware: append(c.middleware, middleware...)}
}

// ToHandler converts httpx.Handler to http.Handler
func (c *Chain) ToHandler(h http.Handler) http.Handler {
	chainedHandler := c.ThenFunc(func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = chainedHandler.ServeHTTP(w, r)
	})
}

func (c *Chain) chain(h Handler) Handler {
	for i := len(c.middleware) - 1; i >= 0; i-- {
		if c.middleware[i] != nil {
			h = c.middleware[i](h)
		}
	}

	return h
}
