package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AppRouter struct {
	router *mux.Router
}

func newAppRouter(router *mux.Router) AppRouter {
	return AppRouter{router: router}
}

func (r AppRouter) subRouter(prefix string) AppRouter {
	return newAppRouter(r.router.PathPrefix(prefix).Subrouter())
}

func (r AppRouter) use(mwf ...mux.MiddlewareFunc) {
	r.router.Use(mwf...)
}

func (r AppRouter) get(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.router.Methods(http.MethodGet).Path(path).HandlerFunc(handler)
}

func (r AppRouter) post(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.router.Methods(http.MethodPost).Path(path).HandlerFunc(handler)
}

// func (r AppRouter) put(path string, handler func(http.ResponseWriter, *http.Request)) {
// 	r.router.Methods(http.MethodPut).Path(path).HandlerFunc(handler)
// }

// func (r AppRouter) delete(path string, handler func(http.ResponseWriter, *http.Request)) {
// 	r.router.Methods(http.MethodDelete).Path(path).HandlerFunc(handler)
// }
