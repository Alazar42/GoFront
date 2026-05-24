//go:build js && wasm

package gofront

// Run executes app setup and keeps the wasm runtime alive.
func Run(setup func()) {
	if setup != nil {
		setup()
	}
	Wait()
}
