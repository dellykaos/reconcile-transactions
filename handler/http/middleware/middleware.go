package middleware

import "github.com/julienschmidt/httprouter"

type Middleware func(httprouter.Handle) httprouter.Handle

// PrependMiddleware prepends the middleware to the handler
func PrependMiddleware(handler httprouter.Handle, middleware ...Middleware) httprouter.Handle {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
