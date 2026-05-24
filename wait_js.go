//go:build js && wasm

package gofront

// Wait keeps the WebAssembly program alive after event handlers are registered.
func Wait() {
	select {}
}
