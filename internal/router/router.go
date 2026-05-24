package router

import (
	"github.com/Alazar42/GoFront/internal/component"
	"github.com/Alazar42/GoFront/runtime"
)

type RouteHandler func() component.Component

var routes = map[string]RouteHandler{}

func Route(path string, handler RouteHandler) { routes[path] = handler }

func StartRouter(selector string) {
	component.Mount(selector, func() component.Component {
		if handler, ok := routes[runtime.CurrentPath()]; ok {
			return handler()
		}
		if handler, ok := routes["*"]; ok {
			return handler()
		}
		return component.Div(component.Text("Not Found"))
	})
	runtime.SetPathCallback(func() { runtime.RequestRender() })
}

func Navigate(path string) { _ = runtime.Navigate(path) }
