//go:build js && wasm

package gofront

import (
	"fmt"
	"syscall/js"
)

func storageGet(scope, key string) string {
	storage := js.Global().Get(scope + "Storage")
	if storage.Truthy() {
		return storage.Call("getItem", key).String()
	}
	return ""
}

func storageSet(scope, key, value string) {
	storage := js.Global().Get(scope + "Storage")
	if storage.Truthy() {
		storage.Call("setItem", key, value)
	}
}

func storageDelete(scope, key string) {
	storage := js.Global().Get(scope + "Storage")
	if storage.Truthy() {
		storage.Call("removeItem", key)
	}
}

func storageClear(scope string) {
	storage := js.Global().Get(scope + "Storage")
	if storage.Truthy() {
		storage.Call("clear")
	}
}

func historyBack()               { js.Global().Get("history").Call("back") }
func historyForward()            { js.Global().Get("history").Call("forward") }
func historyPush(path string)    { js.Global().Get("history").Call("pushState", nil, "", path) }
func historyReplace(path string) { js.Global().Get("history").Call("replaceState", nil, "", path) }

func locationHref() string     { return js.Global().Get("location").Get("href").String() }
func locationPathname() string { return js.Global().Get("location").Get("pathname").String() }
func locationHash() string     { return js.Global().Get("location").Get("hash").String() }
func locationReload()          { js.Global().Get("location").Call("reload") }

func fetchImpl(url string) (Response, error) {
	promise := js.Global().Call("fetch", url)
	resultCh := make(chan Response, 1)
	errCh := make(chan error, 1)

	onSuccess := js.FuncOf(func(this js.Value, args []js.Value) any {
		resp := args[0]
		status := resp.Get("status").Int()
		textPromise := resp.Call("text")
		onText := js.FuncOf(func(this js.Value, args []js.Value) any {
			resultCh <- Response{Status: status, Body: args[0].String()}
			return nil
		})
		textPromise.Call("then", onText)
		return nil
	})

	onError := js.FuncOf(func(this js.Value, args []js.Value) any {
		errCh <- fmt.Errorf("fetch failed")
		return nil
	})

	defer onSuccess.Release()
	defer onError.Release()
	promise.Call("then", onSuccess).Call("catch", onError)

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return Response{}, err
	}
}
