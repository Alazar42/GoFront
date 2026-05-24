# GoFront

GoFront is a frontend framework that allows developers to build browser applications using Go instead of JavaScript.

## Features

- Go-powered frontend development
- HTML integration
- WebAssembly runtime
- Reactive state management
- Component system
- Routing
- Hot reload
- Modern browser support

## CLI Commands

```bash
gofront help
```

Available commands:

- `gofront init`
- `gofront build`
- `gofront dev`
- `gofront serve`
- `gofront version`

Command help:

```bash
gofront <command> -h
```

## Installation

### Requirements

- Go 1.24+
- Git

Verify installation:

```bash
go version
```

### Install GoFront CLI

```bash
git clone https://github.com/gofront/gofront.git
cd gofront
go build -o gofront ./cmd/gofront
```

Windows:

```bash
move gofront.exe C:\GoFront\
```

Verify:

```bash
gofront version
```

## Create a New Project

```bash
gofront init my-app
```

Project structure:

```text
my-app/
├── src/
│   ├── app.go
│   ├── components/
│   └── pages/
├── public/
│   └── index.html
├── dist/
└── gofront.config.json
```

## Project Config

GoFront reads `gofront.config.json` by default:

```json
{
    "name": "GoFront App",
    "src": "src",
    "public": "public",
    "dist": "dist"
}
```

You can override path values with CLI flags.

## HTML Setup

Create `public/index.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>GoFront App</title>
</head>
<body>

    <h1 id="counter">0</h1>
    <button id="increment">Increment</button>

    <script src="app.go"></script>

</body>
</html>
```

## Writing Go Frontend Code

Create `src/app.go`:

```go
package main

import "gofront"

func main() {
    count := gofront.State(0)

    button := gofront.Query("#increment")

    button.OnClick(func() {
        count.Set(count.Get() + 1)
        gofront.Query("#counter").SetText(count.Get())
    })

    gofront.Wait()
}
```

## Development Mode

```bash
gofront dev
```

Custom dev options:

```bash
gofront dev --addr :5173 --poll 500
```

## Production Build

```bash
gofront build
```

Custom build paths:

```bash
gofront build --src src --public public --dist dist
```

Output:

```text
dist/
├── index.html
├── app.wasm
├── gofront.js
└── assets/
```

## Serve Build Output

```bash
gofront serve
```

Serve custom root/address:

```bash
gofront serve --root dist --addr :8080
```

## Typical Workflow

```bash
gofront init my-app
cd my-app
gofront dev
```

Production:

```bash
gofront build
gofront serve
```

## Notes

- In browser apps, call `gofront.Wait()` in `main()` to keep WebAssembly runtime alive.
- If your browser logs `lockdown-install.js` SES warnings, that is typically from an extension, not GoFront.

## Routing

```go
gofront.Route("/", HomePage)
gofront.Route("/about", AboutPage)
```

## State Management

```go
count := gofront.State(0)
count.Set(10)
value := count.Get()
```

## License

MIT
