//go:build !js || !wasm

package lifecycle

// Wait is a no-op on non-wasm targets.
func Wait() {}
