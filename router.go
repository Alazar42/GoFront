package gofront

import "gofront/runtime"

type RouteHandler func() Component

var routes = map[string]RouteHandler{}

func Route(path string, handler RouteHandler) { routes[path] = handler }

func StartRouter(selector string) {
	Mount(selector, func() Component {
		if handler, ok := routes[runtime.CurrentPath()]; ok {
			return handler()
		}
		if handler, ok := routes["*"]; ok {
			return handler()
		}
		return Div(Text("Not Found"))
	})
	runtime.SetPathCallback(func() { runtime.RequestRender() })
}

func Navigate(path string) { _ = runtime.Navigate(path) }
