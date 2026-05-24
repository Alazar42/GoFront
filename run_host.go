//go:build !js || !wasm

package gofront

// Run executes app setup on non-wasm targets without blocking.
func Run(setup func()) {
	if setup != nil {
		setup()
	}
}
