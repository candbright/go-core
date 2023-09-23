package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	Routes() []Route
}

type Router struct {
	RootPath string
	Routes   []Route
}

type Route struct {
	Method  string
	Path    string
	Handler func(context *gin.Context)
}

func (router Router) RegisterHandler(engine gin.IRouter) {
	group := engine.Group(router.RootPath)
	for _, route := range router.Routes {
		switch route.Method {
		case http.MethodGet:
			group.GET(route.Path, route.Handler)
		case http.MethodPost:
			group.POST(route.Path, route.Handler)
		case http.MethodPut:
			group.PUT(route.Path, route.Handler)
		case http.MethodPatch:
			group.PATCH(route.Path, route.Handler)
		case http.MethodDelete:
			group.DELETE(route.Path, route.Handler)
		case http.MethodOptions:
			group.OPTIONS(route.Path, route.Handler)
		case http.MethodHead:
			group.HEAD(route.Path, route.Handler)
		default:
			group.Any(route.Path, route.Handler)
		}
	}
}

func RegisterHandler(engine gin.IRouter, rootPath string, handler Handler) {
	Router{
		RootPath: rootPath,
		Routes:   handler.Routes(),
	}.RegisterHandler(engine)
}
