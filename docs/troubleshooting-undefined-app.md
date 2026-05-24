# Troubleshooting: "undefined: App" Error

## Issue

When running `gofront dev`, you see:

```
build wasm: exit status 1
# gofront-app/src
src\main.go:X:Y: undefined: App
```

## Causes & Solutions

### 1. Missing `app.gox` File

**Cause:** The `src/app.gox` file doesn't exist.

**Solution:**
- Create `src/app.gox` with a basic component:

```gox
<template>
<Div>
  <H1>My App</H1>
</Div>
</template>
```

### 2. Invalid `app.gox` Content

**Cause:** The `.gox` file is empty or has syntax errors in template.

**Solution:**
- Ensure your template has valid JSX-like elements
- Validate element names: `Div`, `Span`, `Button`, `H1`, `P`, `Main`, `Section`, `Header`, `Footer`, `Text`
- Check that all opening tags have closing tags

### 3. Transpilation Not Running

**Cause:** The `.gox` file exists but wasn't transpiled before compilation.

**Solution:**
- The transpiler runs automatically before each build
- If build fails, check for transpilation errors:
  - Look at the build output in the dist error overlay
  - Ensure the `.gox` file syntax is valid

### 4. Generated File Not Being Found

**Cause:** The generated `app_gox_gen.go` file exists but isn't included in the build.

**Solution:**
- The build system automatically discovers `_gox_gen.go` files
- Ensure generated files are in the `src/` directory (not in subdirectories like `.goxgen`)
- Check that your project structure matches:

```
project-root/
├── src/
│   ├── main.go
│   ├── app.gox
│   └── app_gox_gen.go (generated automatically)
├── public/
├── go.mod
└── gofront.config.json
```

## Verification Checklist

- [ ] `src/app.gox` exists and is not empty
- [ ] `app.gox` has valid JSX-like markup in `<template>` section
- [ ] All elements in template are valid GoFront components
- [ ] No syntax errors (mismatched tags, invalid attributes)
- [ ] After fixing `.gox` file, save and wait for auto-rebuild
- [ ] Check the browser console for any runtime errors

## Debug Steps

1. **Check the error overlay in browser** for transpilation errors
2. **Look at generated `app_gox_gen.go`** to verify the `App()` function was created
3. **Verify package name** in generated file matches main package
4. **Rebuild manually**:
   ```bash
   gofront build
   ```
5. **Check src directory** for `app_gox_gen.go` file

## New .gox Format Support

If you recently updated GoFront, note that the `.gox` format now supports:

- `<script>` blocks for Go code
- `<style>` blocks for CSS
- `<template>` blocks (or template as implicit content)

Simple example:
```gox
<template>
<Div>Hello</Div>
</template>
```

Full example:
```gox
<script>
var count = 0
</script>

<style>
.app { padding: 2rem; }
</style>

<template>
<Div class="app">
  <H1>App</H1>
</Div>
</template>
```

If your `app.gox` is in the old format (no blocks), it still works—the entire file is treated as the template.
