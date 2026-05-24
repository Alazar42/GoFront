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
- Tailwind CSS build stage
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

The generated starter is a Tailwind-powered counter page with a polished card layout, increment and reset actions, and ready-to-edit source files.

Project structure:

```text
my-app/
├── src/
│   ├── app.go
│   ├── styles.css
│   ├── components/
│   └── pages/
├── public/
│   └── index.html
├── dist/
├── tailwind.config.js
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

## Tailwind CSS

GoFront automatically builds Tailwind when it finds `src/styles.css`.

Default files:

- `src/styles.css`
- `tailwind.config.js`

The build pipeline will generate `dist/styles.css` and link it into the built HTML.

Example `src/styles.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
    html {
        @apply bg-slate-950 text-slate-100 antialiased;
    }

    body {
        @apply min-h-screen;
    }
}
```

Example `tailwind.config.js`:

```js
module.exports = {
    content: [
        './src/**/*.go',
        './src/**/*.html',
        './public/**/*.html'
    ],
    theme: {
        extend: {}
    },
    plugins: []
}
```

GoFront looks for `tailwindcss` first, then uses `npx tailwindcss@3.4.17` if the local binary is not installed.

## Writing Go Frontend Code

Create `src/app.go`:

`main` inside `src/app.go` is the application entrypoint and is responsible for running your app.

```go
package main

import "gofront"

func main() {
    gofront.Run(func() {
        count := gofront.State(0)

        button := gofront.Query("#increment")

        button.OnClick(func() {
            count.Set(count.Get() + 1)
            gofront.Query("#counter").SetText(count.Get())
        })
    })
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

- Use `gofront.Run(func() { ... })` as your app entrypoint. It keeps the WebAssembly runtime alive automatically.
- `gofront.Wait()` is still available for advanced/manual lifecycle control.
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
