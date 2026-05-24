//go:build js && wasm

package lifecycle

// Wait keeps the WebAssembly program alive after event handlers are registered.
func Wait() {
	select {}
}
