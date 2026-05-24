# GoFront .gox Component Format

The `.gox` format is a JSX-like template syntax for building GoFront components. It supports script blocks, style blocks, and template markup similar to Svelte components.

## Basic Structure

A `.gox` file can have up to three sections:

1. **`<script>`** - Go code (variables, functions, event handlers)
2. **`<style>`** - CSS styles scoped to the component
3. **`<template>`** - JSX-like markup (required)

### Minimal Example

```gox
<template>
<Div>
  <H1>Hello World</H1>
</Div>
</template>
```

### Full Example with All Sections

```gox
<script>
// Component state and handlers
var count = 0

func increment() {
  count++
  // Re-render logic here
}

func reset() {
  count = 0
}
</script>

<style>
.counter-container {
  padding: 2rem;
  border-radius: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.counter-display {
  font-size: 3rem;
  font-weight: bold;
  color: white;
  text-align: center;
}

.button-group {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 1rem;
}

.button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-weight: 600;
}

.button-primary {
  background: #4CAF50;
  color: white;
}

.button-secondary {
  background: #f44336;
  color: white;
}
</style>

<template>
<Div class="counter-container">
  <Div class="counter-display" id="counter">
    <Text>0</Text>
  </Div>
  <Div class="button-group">
    <Button class="button button-primary" id="increment">Increment</Button>
    <Button class="button button-secondary" id="reset">Reset</Button>
  </Div>
</Div>
</template>
```

## Supported Elements

The template section supports these JSX-like elements:

- `<Div>` - Creates a div element
- `<Span>` - Creates a span element
- `<Button>` - Creates a button element
- `<H1>` through heading tags
- `<P>` - Creates a paragraph
- `<Main>` - Creates a main element
- `<Section>` - Creates a section element
- `<Header>` - Creates a header element
- `<Footer>` - Creates a footer element
- `<Text>` - Creates a text node

## Attributes

Supported attributes in template elements:

- `id="..."` - Sets the element ID
- `class="..."` - Sets CSS classes
- Other attributes can be added using `any_attr="value"` syntax

### Example with Attributes

```gox
<template>
<Div id="main-container" class="flex gap-4">
  <H1 id="title" class="text-2xl font-bold">Welcome</H1>
  <Button class="btn btn-primary">Click Me</Button>
</Div>
</template>
```

## Component Naming

Generated function names are derived from the filename:

- `app.gox` → `App()` function
- `header.gox` → `Header()` function
- `button.gox` → `Button()` function
- `pages/home.gox` → `HomePage()` function
- `pages/about.gox` → `AboutPage()` function

## Style Files

When you include a `<style>` block, GoFront automatically creates a `.module.css` file:

- `app.gox` + `<style>` block → `app.module.css`
- `components/card.gox` + `<style>` block → `components/card.module.css`

These CSS files are automatically created in the same directory and can be imported or referenced in your styles.

## Script Block

The `<script>` block contains Go code that will be included in the generated `.go` file. This is where you define:

- Component state variables
- Event handler functions
- Helper functions
- Any other Go code needed by the component

**Note:** Don't include package declarations or imports—GoFront handles these automatically.

```gox
<script>
// State variables
var isOpen = false
var items []string

// Event handlers
func toggleOpen() {
  isOpen = !isOpen
}

func addItem(item string) {
  items = append(items, item)
}
</script>

<template>
<Div>
  <!-- Template uses the state and handlers -->
</Div>
</template>
```

## Template Nesting

Elements can be nested to create component hierarchies:

```gox
<template>
<Div class="card">
  <Header>
    <H1>Title</H1>
  </Header>
  <Main class="card-content">
    <P>Content goes here</P>
    <Button>Action</Button>
  </Main>
  <Footer>
    <Span>Footer text</Span>
  </Footer>
</Div>
</template>
```

## Build Process

When you run `gofront dev` or `gofront build`:

1. GoFront scans for `.gox` files
2. Extracts `<script>`, `<style>`, and `<template>` blocks
3. Generates `*_gox_gen.go` files with components
4. Creates `.module.css` files for styles
5. Compiles the WebAssembly module
6. Includes styles in the final bundle

Generated files (`*_gox_gen.go` and `*.module.css`) are automatically gitignored.

## Best Practices

1. **Keep components focused** - One component per file
2. **Use meaningful names** - File names should describe the component
3. **Organize in folders** - Use `components/`, `pages/`, `layouts/` folders
4. **Style scoping** - Use class names to scope styles to components
5. **Reusable logic** - Extract common patterns into helper functions
