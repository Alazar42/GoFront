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

## Production Build

```bash
gofront build
```

Output:

```text
dist/
├── index.html
├── app.wasm
├── gofront.js
└── assets/
```

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
