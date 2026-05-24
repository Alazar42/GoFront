//go:build js && wasm

package runtime

import "syscall/js"

func mountHTMLImpl(selector, html string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Set("innerHTML", html)
	}
	return nil
}

func setTextImpl(selector, value string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Set("textContent", value)
	}
	return nil
}

func setHTMLImpl(selector, value string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Set("innerHTML", value)
	}
	return nil
}

func addClassImpl(selector, class string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Get("classList").Call("add", class)
	}
	return nil
}

func removeClassImpl(selector, class string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Get("classList").Call("remove", class)
	}
	return nil
}

func toggleClassImpl(selector, class string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Get("classList").Call("toggle", class)
	}
	return nil
}

func getValueImpl(selector string) (string, error) {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		return el.Get("value").String(), nil
	}
	return "", nil
}

func setValueImpl(selector, value string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Set("value", value)
	}
	return nil
}

func hideImpl(selector string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Get("style").Set("display", "none")
	}
	return nil
}

func showImpl(selector string) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Get("style").Set("display", "")
	}
	return nil
}

func addEventListenerImpl(selector, event string, handler func()) error {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		callback := js.FuncOf(func(this js.Value, args []js.Value) any {
			handler()
			return nil
		})
		el.Call("addEventListener", event, callback)
	}
	return nil
}

func currentPathImpl() string { return js.Global().Get("location").Get("pathname").String() }
func navigateImpl(path string) error {
	js.Global().Get("history").Call("pushState", nil, "", path)
	return nil
}
