package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Adapt(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
