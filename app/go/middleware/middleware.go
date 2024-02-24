package middleware

import "net/http"

type Middleware struct {
	middlewares []func(http.ResponseWriter, *http.Request)
}

func NewMiddleWare(fns ...func(w http.ResponseWriter, r *http.Request)) Middleware {
	return Middleware{
		middlewares: fns,
	}
}

func (m *Middleware) Chain(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mFn := range m.middlewares {
			mFn(w, r)
		}
		fn(w, r)
	}
}
