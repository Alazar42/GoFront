//go:build !js || !wasm

package gofront

// Wait is a no-op on non-wasm targets.
func Wait() {}
