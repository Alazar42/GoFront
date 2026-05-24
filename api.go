package gofront

import (
	"time"

	"github.com/Alazar42/GoFront/internal/browser"
	"github.com/Alazar42/GoFront/internal/component"
	"github.com/Alazar42/GoFront/internal/element"
	"github.com/Alazar42/GoFront/internal/lifecycle"
	"github.com/Alazar42/GoFront/internal/router"
	"github.com/Alazar42/GoFront/internal/state"
)

type Response = browser.Response
type Storage = browser.Storage
type ConsoleAPI = browser.ConsoleAPI
type HistoryAPI = browser.HistoryAPI
type LocationAPI = browser.LocationAPI
type TimerHandle = browser.TimerHandle

type Attribute = component.Attribute
type Component = component.Component
type Element = element.Element
type RouteHandler = router.RouteHandler
type StateValue[T any] = state.StateValue[T]

func Fetch(url string) (Response, error) { return browser.Fetch(url) }
func LocalStorage() Storage              { return browser.LocalStorage() }
func SessionStorage() Storage            { return browser.SessionStorage() }
func Console() ConsoleAPI                { return browser.Console() }
func History() HistoryAPI                { return browser.History() }
func Location() LocationAPI              { return browser.Location() }

func Timer(duration time.Duration, callback func()) TimerHandle {
	return browser.Timer(duration, callback)
}
func Interval(period time.Duration, callback func()) TimerHandle {
	return browser.Interval(period, callback)
}

func Query(selector string) *Element { return element.Query(selector) }

func Text(value any) Component { return component.Text(value) }
func ElementNode(tag string, children ...Component) Component {
	return component.ElementNode(tag, children...)
}
func Div(children ...Component) Component     { return component.Div(children...) }
func Span(children ...Component) Component    { return component.Span(children...) }
func Button(children ...Component) Component  { return component.Button(children...) }
func H1(children ...Component) Component      { return component.H1(children...) }
func P(children ...Component) Component       { return component.P(children...) }
func Main(children ...Component) Component    { return component.Main(children...) }
func Section(children ...Component) Component { return component.Section(children...) }
func Header(children ...Component) Component  { return component.Header(children...) }
func Footer(children ...Component) Component  { return component.Footer(children...) }

func ID(value string) Attribute         { return component.ID(value) }
func Class(value string) Attribute      { return component.Class(value) }
func Attr(name, value string) Attribute { return component.Attr(name, value) }

func Mount(selector string, view func() Component)  { component.Mount(selector, view) }
func Render(selector string, view func() Component) { component.Render(selector, view) }

func Route(path string, handler RouteHandler) { router.Route(path, handler) }
func StartRouter(selector string)             { router.StartRouter(selector) }
func Navigate(path string)                    { router.Navigate(path) }

func State[T any](initial T) *StateValue[T] { return state.State(initial) }

func Run(setup func()) { lifecycle.Run(setup) }
func Wait()            { lifecycle.Wait() }
