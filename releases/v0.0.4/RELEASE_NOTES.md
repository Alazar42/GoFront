# GoFront v0.0.4 Release Notes

## New Features

### Enhanced .gox Component Format
- Added support for **Svelte-like component structure** with `<script>`, `<style>`, and `<template>` blocks
- Script blocks for Go code (variables, functions, event handlers)
- Style blocks that automatically generate scoped `.module.css` files
- Template blocks for JSX-like markup (backwards compatible with old format)

### Improved Project Scaffolding
- New scaffold templates demonstrate the enhanced `.gox` format
- Better organized component examples with proper structure
- Updated `.gitignore` to exclude generated files (`*_gox_gen.go`, `*.module.css`)

### Generator Improvements
- Fixed package detection for generated components
- Generated files now placed directly in source directory for proper compilation
- Fixes "undefined App" errors in development builds

## Migration from v0.0.3

The new format is **fully backwards compatible**. Existing `.gox` files continue to work:

### Old Format (Still Works)
```gox
<Div class="app">
  <H1>Hello</H1>
</Div>
```

### New Format (Enhanced)
```gox
<script>
var count = 0

func increment() {
  count++
}
</script>

<style>
.app { padding: 2rem; }
</style>

<template>
<Div class="app">
  <H1>Counter: <Text>0</Text></H1>
  <Button id="increment">+</Button>
</Div>
</template>
```

## Usage

```bash
# Initialize a new project
gofront init my-app
cd my-app

# Start development server with hot reload
gofront dev

# Build for production
gofront build
```

## Key Improvements

✓ Script and style blocks enable true component encapsulation  
✓ Automatic CSS module generation from style blocks  
✓ Fixed generation pipeline for proper package detection  
✓ Better error messages for build failures  
✓ Tailwind CSS support maintained and improved  

## Documentation

- **Component Format**: See `docs/gox-format.md` for detailed format documentation
- **Troubleshooting**: See `docs/troubleshooting-undefined-app.md` for common issues
- **Getting Started**: See `README.md` for setup instructions

## Supported Elements

- `<Div>`, `<Span>`, `<Button>`, `<H1>` (through heading levels)
- `<P>`, `<Main>`, `<Section>`, `<Header>`, `<Footer>`
- `<Text>` for text content

## Attributes

- `id="..."` - Element ID
- `class="..."` - CSS classes
- Custom attributes supported

## Breaking Changes

None. v0.0.4 is fully backwards compatible with v0.0.3.
