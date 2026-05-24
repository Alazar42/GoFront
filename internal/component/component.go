package component

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"gofront/runtime"
)

type Attribute struct {
	Name  string
	Value string
}

type Component struct {
	tag      string
	text     any
	attrs    []Attribute
	children []Component
}

func Text(value any) Component { return Component{tag: "#text", text: value} }
func ElementNode(tag string, children ...Component) Component {
	return Component{tag: tag, children: children}
}
func Div(children ...Component) Component     { return ElementNode("div", children...) }
func Span(children ...Component) Component    { return ElementNode("span", children...) }
func Button(children ...Component) Component  { return ElementNode("button", children...) }
func H1(children ...Component) Component      { return ElementNode("h1", children...) }
func P(children ...Component) Component       { return ElementNode("p", children...) }
func Main(children ...Component) Component    { return ElementNode("main", children...) }
func Section(children ...Component) Component { return ElementNode("section", children...) }
func Header(children ...Component) Component  { return ElementNode("header", children...) }
func Footer(children ...Component) Component  { return ElementNode("footer", children...) }

func ID(value string) Attribute    { return Attribute{Name: "id", Value: value} }
func Class(value string) Attribute { return Attribute{Name: "class", Value: value} }
func Attr(name, value string) Attribute {
	return Attribute{Name: name, Value: value}
}

func (c Component) With(attrs ...Attribute) Component {
	result := c
	result.attrs = append(append([]Attribute(nil), c.attrs...), attrs...)
	return result
}

func (c Component) HTML() string { return c.render() }

func Mount(selector string, view func() Component) {
	render := func() {
		_ = runtime.MountHTML(selector, view().HTML())
	}
	runtime.SetRenderCallback(render)
	render()
}

func Render(selector string, view func() Component) { Mount(selector, view) }

func (c Component) render() string {
	if c.tag == "#text" {
		return html.EscapeString(stringify(c.text))
	}
	var builder strings.Builder
	builder.WriteString("<")
	builder.WriteString(c.tag)
	for _, attr := range sortAttributes(c.attrs) {
		builder.WriteByte(' ')
		builder.WriteString(attr.Name)
		builder.WriteString("=\"")
		builder.WriteString(html.EscapeString(attr.Value))
		builder.WriteString("\"")
	}
	builder.WriteString(">")
	for _, child := range c.children {
		builder.WriteString(child.render())
	}
	builder.WriteString("</")
	builder.WriteString(c.tag)
	builder.WriteString(">")
	return builder.String()
}

func sortAttributes(attrs []Attribute) []Attribute {
	result := append([]Attribute(nil), attrs...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func stringify(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return html.EscapeString(fmt.Sprint(v))
	}
}
