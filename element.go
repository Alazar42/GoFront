package gofront

import (
	"fmt"
	"strconv"

	"gofront/runtime"
)

type Element struct {
	selector string
}

func Query(selector string) *Element { return &Element{selector: selector} }

func (e *Element) SetText(value any)        { _ = runtime.SetText(e.selector, stringify(value)) }
func (e *Element) SetHTML(value string)     { _ = runtime.SetHTML(e.selector, value) }
func (e *Element) AddClass(value string)    { _ = runtime.AddClass(e.selector, value) }
func (e *Element) RemoveClass(value string) { _ = runtime.RemoveClass(e.selector, value) }
func (e *Element) ToggleClass(value string) { _ = runtime.ToggleClass(e.selector, value) }
func (e *Element) OnClick(handler func())   { runtime.RegisterEvent(e.selector, "click", handler) }
func (e *Element) OnInput(handler func())   { runtime.RegisterEvent(e.selector, "input", handler) }
func (e *Element) Value() string {
	value, _ := runtime.GetValue(e.selector)
	return value
}
func (e *Element) SetValue(value any) { _ = runtime.SetValue(e.selector, stringify(value)) }
func (e *Element) Hide()              { _ = runtime.Hide(e.selector) }
func (e *Element) Show()              { _ = runtime.Show(e.selector) }

func stringify(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(v)
	}
}
