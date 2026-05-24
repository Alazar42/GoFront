# GoFront v0.0.4 Sample Scaffold

This is an example project structure showing the new features in v0.0.4.

## Directory Structure

```
project/
├── src/
│   ├── main.go                    # Entry point
│   ├── app.gox                    # Main app component (new format)
│   ├── app.module.css             # Generated from app.gox <style> block
│   ├── styles.css                 # Global styles
│   ├── components/
│   │   ├── counter.gox            # Counter component (new format)
│   │   └── counter.module.css     # Generated scoped styles
│   └── pages/
│       └── home.gox               # Home page component
├── public/
│   ├── index.html
│   └── favicon.ico
├── dist/                          # Build output (generated)
├── go.mod
├── gofront.config.json
├── tailwind.config.js
└── .gitignore
```

## Example Files

### src/main.go
```go
package main

import GoFront "github.com/Alazar42/GoFront"

func main() {
	GoFront.Run(func() {
		// Mount the generated app component into body
		GoFront.Mount("body", App)
	})
}
```

### src/app.gox (New Format)
```gox
<script>
// App state and handlers
var page = "home"

func navigateTo(p string) {
  page = p
}
</script>

<style>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
  color: white;
}

.app-main {
  flex: 1;
  padding: 2rem;
}
</style>

<template>
<Div class="app-container">
  <Header class="app-header">
    <H1>GoFront v0.0.4 App</H1>
  </Header>
  
  <Main class="app-main">
    <!-- Page content -->
    <Section>
      <P>Welcome to GoFront with enhanced component format!</P>
    </Section>
  </Main>
  
  <Footer class="app-footer">
    <P>Made with GoFront</P>
  </Footer>
</Div>
</template>
```

### src/components/counter.gox (New Format)
```gox
<script>
// Counter component state
var count = 0

func increment() {
  count++
  // Update DOM
}

func decrement() {
  if count > 0 {
    count--
  }
}

func reset() {
  count = 0
}
</script>

<style>
.counter-card {
  background: white;
  border-radius: 0.5rem;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  max-width: 400px;
  margin: 0 auto;
}

.counter-display {
  font-size: 3rem;
  font-weight: bold;
  text-align: center;
  color: #667eea;
  margin: 1rem 0;
}

.button-group {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  flex-wrap: wrap;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5568d3;
  transform: translateY(-2px);
}

.btn-secondary {
  background: #e0e7ff;
  color: #667eea;
}

.btn-secondary:hover {
  background: #c7d2fe;
}
</style>

<template>
<Div class="counter-card">
  <H1>Counter</H1>
  
  <Div class="counter-display" id="count">
    <Text>0</Text>
  </Div>
  
  <Div class="button-group">
    <Button class="btn btn-secondary" id="decrement">−</Button>
    <Button class="btn btn-primary" id="increment">+</Button>
    <Button class="btn btn-secondary" id="reset">Reset</Button>
  </Div>
</Div>
</template>
```

### src/pages/home.gox (New Format)
```gox
<style>
.home-page {
  max-width: 800px;
  margin: 0 auto;
}

.hero {
  text-align: center;
  padding: 2rem 0;
}

.hero h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 2rem 0;
}

.feature-list li {
  padding: 0.75rem 0;
  border-bottom: 1px solid #e5e7eb;
}

.feature-list li:before {
  content: "✓ ";
  color: #10b981;
  font-weight: bold;
  margin-right: 0.5rem;
}
</style>

<template>
<Div class="home-page">
  <Section class="hero">
    <H1>Welcome to GoFront</H1>
    <P>Build interactive web apps with Go and JSX-like syntax</P>
  </Section>
  
  <Section>
    <H2>v0.0.4 Features</H2>
    <Div class="feature-list">
      <Div>✓ Enhanced .gox component format</Div>
      <Div>✓ Script blocks for Go code</Div>
      <Div>✓ Scoped CSS modules</Div>
      <Div>✓ JSX-like template syntax</Div>
      <Div>✓ Live hot reload</Div>
      <Div>✓ Tailwind CSS support</Div>
    </Div>
  </Section>
</Div>
</template>
```

### src/styles.css (Global Styles)
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  html {
    @apply bg-slate-50 text-slate-900 antialiased;
  }

  body {
    @apply min-h-screen overflow-x-hidden;
  }

  ::selection {
    @apply bg-purple-200 text-slate-900;
  }
}

@layer components {
  .gofront-container {
    @apply max-w-6xl mx-auto px-4;
  }
}
```

### public/index.html
```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="styles.css">
  <title>GoFront v0.0.4</title>
</head>
<body>
  <main id="app"></main>
  <script src="gofront.js"></script>
</body>
</html>
```

### go.mod
```
module gofront-app

go 1.24

require github.com/Alazar42/GoFront v0.0.4
```

### gofront.config.json
```json
{
  "name": "GoFront v0.0.4 Demo",
  "src": "src",
  "public": "public",
  "dist": "dist"
}
```

### tailwind.config.js
```javascript
module.exports = {
  content: [
    './src/**/*.go',
    './src/**/*.gox',
    './public/**/*.html'
  ],
  theme: {
    extend: {}
  },
  plugins: []
};
```

### .gitignore
```
dist/
node_modules/
*.wasm
*.exe
.idea/
.vscode/
*_gox_gen.go
*.module.css
```

## Getting Started

```bash
# Install GoFront v0.0.4
go install github.com/Alazar42/GoFront/cmd/gofront@v0.0.4

# Or use the binary directly
./releases/v0.0.4/gofront.exe init my-app

# Start development
cd my-app
gofront dev

# Open browser to http://localhost:3000
```

## Key Changes from v0.0.3

1. **Component Format**: New Svelte-like syntax with `<script>`, `<style>`, `<template>` blocks
2. **Style Generation**: Automatic `.module.css` file generation from style blocks
3. **Build Pipeline**: Fixed package detection for generated components
4. **Scaffolding**: Updated templates show best practices with new format
5. **Backwards Compatibility**: Old format still fully supported

## Documentation

Run `gofront --help` for command reference. See the main README for more information.
