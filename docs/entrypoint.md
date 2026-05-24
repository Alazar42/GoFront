# App Entrypoint

In a GoFront app, the `main` function is inside `src/app.go`.

That `main` function is the entrypoint responsible for running the app.

When the WebAssembly module starts, Go executes `main`, and your app setup starts from there.

## Keep the WASM runtime alive

Use `gofront.Run(func(){ ... })` inside `main` so the WebAssembly runtime remains active after your handlers register. If `main` returns the Go program exits and the wasm module will no longer handle events.

Example:

```go
package main

import "gofront"

func main() {
	gofront.Run(func() {
		// app setup: register handlers, mount components, etc.
	})
}
```

If you prefer manual control, `gofront.Wait()` is still available to block until an explicit shutdown.
