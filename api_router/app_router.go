package api_router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AppRouter struct {
	Router *mux.Router
}

func NewAppRouter(router *mux.Router) AppRouter {
	return AppRouter{Router: router}
}

// func (r AppRouter) Router() *mux.Router {
// 	return r.router
// }

func (r AppRouter) SubRouter(prefix string) AppRouter {
	return NewAppRouter(r.Router.PathPrefix(prefix).Subrouter())
}

func (r AppRouter) Use(mwf ...mux.MiddlewareFunc) {
	r.Router.Use(mwf...)
}

func (r AppRouter) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.Router.Methods(http.MethodGet).Path(path).HandlerFunc(handler)
}

func (r AppRouter) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.Router.Methods(http.MethodPost).Path(path).HandlerFunc(handler)
}

// func (r AppRouter) put(path string, handler func(http.ResponseWriter, *http.Request)) {
// 	r.Router.Methods(http.MethodPut).Path(path).HandlerFunc(handler)
// }

func (r AppRouter) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.Router.Methods(http.MethodDelete).Path(path).HandlerFunc(handler)
}
