package routes

import (
	"BTC-price-restful/auth"
	"BTC-price-restful/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []route

func init() {
	register("GET", "/", controllers.GetDefault, nil)

	register("GET", "/getServiceMap", controllers.GetServiceMap, auth.TokenMiddleware)
	register("GET", "/getLatestPrice/{service}", controllers.GetLatestPrice, auth.TokenMiddleware)
	register("GET", "/getLatestAllPrice", controllers.GetLatestAllPrice, auth.TokenMiddleware)

	register("POST", "/user/register", controllers.Register, nil)
	register("POST", "/user/login", controllers.Login, nil)
}

// NewRouter create one 'mux.NewRouter()' and register handle funcs.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		if route.Middleware == nil {
			r.HandleFunc(route.Pattern, route.Handler).Methods(route.Method)
		} else {
			r.Handle(route.Pattern, route.Middleware(route.Handler)).Methods(route.Method)
		}
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, route{method, pattern, handler, middleware})
}
