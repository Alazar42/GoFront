package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Alazar42/GoFront/internal/version"
)

func Scaffold(root string, moduleName string) error {
	resolvedModuleName, err := resolveModuleName(root, moduleName)
	if err != nil {
		return err
	}
	moduleContent, err := moduleFile(root, resolvedModuleName)
	if err != nil {
		return err
	}
	files := map[string]string{
		filepath.Join(root, "go.mod"):                           moduleContent,
		filepath.Join(root, "src", "main.go"):                   mainGoFile(resolvedModuleName),
		filepath.Join(root, "src", "app.gox"):                   appFile(),
		filepath.Join(root, "src", "styles.css"):                stylesFile(),
		filepath.Join(root, "src", "components", "counter.gox"): componentFile(),
		filepath.Join(root, "src", "pages", "home.gox"):         pageFile(),
		filepath.Join(root, "public", "index.html"):             indexFile(),
		filepath.Join(root, "tailwind.config.js"):               tailwindConfigFile(),
		filepath.Join(root, "gofront.config.json"):              configFile(),
		filepath.Join(root, ".gitignore"):                       gitignoreFile(),
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	// Run go mod tidy to generate a complete go.sum for the scaffolded project.
	if err := runGoModTidy(root); err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}
	return nil
}

func resolveModuleName(root string, moduleName string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	moduleName = sanitizeModuleName(moduleName)
	if moduleName == "" {
		moduleName = sanitizeModuleName(filepath.Base(absRoot))
	}
	if moduleName == "" || moduleName == "." {
		moduleName = "gofront-app"
	}
	return moduleName, nil
}

func moduleFile(root string, moduleName string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	moduleName, err = resolveModuleName(absRoot, moduleName)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("module %s\n\ngo 1.24\n", moduleName))

	parent := filepath.Dir(absRoot)
	parentModPath := filepath.Join(parent, "go.mod")
	data, err := os.ReadFile(parentModPath)
	if err == nil && strings.Contains(string(data), "module gofront") {
		rel, relErr := filepath.Rel(absRoot, parent)
		if relErr == nil {
			builder.WriteString("\nrequire gofront v0.0.0\n")
			builder.WriteString(fmt.Sprintf("replace gofront => %s\n", filepath.ToSlash(rel)))
		}
	} else {
		// When no local gofront module is present, scaffold against the
		// published module path and tag so generated apps can `go get` the
		// released dependency. Version is automatically kept in sync with the CLI version.
		builder.WriteString(fmt.Sprintf("\nrequire github.com/Alazar42/GoFront v%s\n", version.Version))
	}

	return builder.String(), nil
}

func sanitizeModuleName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, " ", "-")
	var out strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func configFile() string {
	return "{\n  \"name\": \"GoFront App\",\n  \"src\": \"src\",\n  \"public\": \"public\",\n  \"dist\": \"dist\"\n}\n"
}

func gitignoreFile() string {
	return "dist/\n*.wasm\n*.exe\n.idea/\n.vscode/\nsrc/.goxgen/\n*.module.css\n"
}

func indexFile() string {
	return "<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <link rel=\"stylesheet\" href=\"styles.css\">\n  <title>GoFront App</title>\n</head>\n<body>\n  <div id=\"app\"></div>\n  <script src=\"gofront.js\"></script>\n</body>\n</html>\n"
}

func stylesFile() string {
	return "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n\n@layer base {\n  html {\n    @apply bg-slate-950 text-slate-100 antialiased;\n  }\n\n  body {\n    @apply min-h-screen overflow-x-hidden;\n  }\n\n  ::selection {\n    @apply bg-cyan-300 text-slate-950;\n  }\n}\n\n@layer components {\n  .gofront-card {\n    @apply rounded-3xl border border-white/10 bg-white/5 shadow-2xl shadow-cyan-950/30 backdrop-blur-xl;\n  }\n}\n"
}

func tailwindConfigFile() string {
	return "module.exports = {\n  content: [\n    './src/**/*.go',\n    './src/**/*.html',\n    './public/**/*.html'\n  ],\n  theme: {\n    extend: {}\n  },\n  plugins: []\n};\n"
}

func appFile() string {
	// .gox component with script, styles, and template blocks (Svelte-like format)
	// Script block contains Go code (variables, event handlers)
	// Template block contains JSX-like structure with on:click={handler} syntax
	return `<script>
import (
	"strconv"
	"sync"
)

var count = 0
var bindOnce sync.Once

func renderCount() {
	_ = GoFront.SetText("#counter-display", strconv.Itoa(count))
}

func incrementCount() {
	count++
	renderCount()
}

func resetCount() {
	count = 0
	renderCount()
}

func init() {
	bindOnce.Do(func() {
		renderCount()
	})
}
</script>

<template>
<Div class="flex h-screen items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
  <Div class="gofront-card p-12">
    <Div class="mb-2 text-center text-sm font-semibold uppercase tracking-widest text-cyan-300/70">Counter App</Div>
    <Div class="text-center text-7xl font-black tracking-tight text-cyan-300" id="counter-display">0</Div>
    <Div class="mt-8 flex gap-4 justify-center">
      <Button on:click={incrementCount} class="px-8 py-3 bg-cyan-500 hover:bg-cyan-600 text-slate-950 font-bold rounded-lg transition transform hover:scale-105 active:scale-95">+</Button>
      <Button on:click={resetCount} class="px-8 py-3 bg-red-500 hover:bg-red-600 text-slate-950 font-bold rounded-lg transition transform hover:scale-105 active:scale-95">Reset</Button>
    </Div>
  </Div>
</Div>
</template>
`
}

func componentFile() string {
	// .gox component: Counter display component with on:click event handlers
	return `<script>
import "strconv"

var localCount = 0

func handleIncrement() {
	localCount++
	_ = GoFront.SetText("#local-counter", strconv.Itoa(localCount))
}

func handleDecrement() {
	if localCount > 0 {
		localCount--
	}
	_ = GoFront.SetText("#local-counter", strconv.Itoa(localCount))
}
</script>

<template>
<Div class="gofront-card p-8">
  <Div class="text-sm font-semibold uppercase tracking-widest text-cyan-300/70 mb-4">Reusable Counter</Div>
  <Div class="text-5xl font-bold text-cyan-300 mb-6 text-center" id="local-counter">0</Div>
  <Div class="flex gap-3 justify-center">
    <Button on:click={handleDecrement} class="px-6 py-2 bg-slate-700 hover:bg-slate-600 text-slate-100 font-semibold rounded-lg transition hover:bg-slate-500">−</Button>
    <Button on:click={handleIncrement} class="px-6 py-2 bg-slate-700 hover:bg-slate-600 text-slate-100 font-semibold rounded-lg transition hover:bg-slate-500">+</Button>
  </Div>
</Div>
</template>
`
}

func pageFile() string {
	// .gox page template showcasing layout
	return `<template>
<Div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-8">
  <Div class="max-w-2xl mx-auto">
    <Div class="mb-12 text-center">
      <Div class="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-cyan-300 to-cyan-400 mb-2">GoFront + Tailwind</Div>
      <Div class="text-slate-400">Beautiful React-like components with Go WebAssembly</Div>
    </Div>
    <Div class="space-y-6">
      <Div class="gofront-card p-8">
        <Div class="text-sm font-semibold uppercase tracking-widest text-cyan-300/70 mb-4">Welcome</Div>
        <Div class="text-slate-200 leading-relaxed">
          <Text>This is a GoFront app. Build your UI with .gox components and Go!</Text>
        </Div>
      </Div>
    </Div>
  </Div>
</Div>
</template>
`
}

func mainGoFile(moduleName string) string {
	moduleName = sanitizeModuleName(moduleName)
	if moduleName == "" {
		moduleName = "gofront-app"
	}
	return fmt.Sprintf("package main\n\nimport (\n\tGoFront \"github.com/Alazar42/GoFront\"\n\tappgen \"%s/src/.goxgen\"\n)\n\nfunc main() {\n\tGoFront.Run(func() {\n\t\t// Mount the generated app component into the app root\n\t\tGoFront.Mount(\"#app\", appgen.App)\n\t})\n}\n", moduleName)

}

func runGoModTidy(root string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = root
	return cmd.Run()
}
