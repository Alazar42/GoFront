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
	// Style block contains CSS (scoped to this component via .module.css)
	// Template block contains JSX-like structure
	return `<script>
import (
	"strconv"
	"sync"
)

var count = 0
var bindOnce sync.Once

func renderCount() {
	_ = GoFront.SetText("#counter-value", strconv.Itoa(count))
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
		GoFront.RegisterEvent("#increment", "click", func() {
			incrementCount()
		})
		GoFront.RegisterEvent("#reset", "click", func() {
			resetCount()
		})
	})
	renderCount()
}
</script>

<style>
.app-root {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}
</style>

<template>
<Div class="app-root">
  <Div id="counter-value" class="text-7xl font-black text-cyan-300">0</Div>
  <Div class="mt-4">
    <Button id="increment">Increment</Button>
    <Button id="reset">Reset</Button>
  </Div>
</Div>
</template>
`
}

func componentFile() string {
	// .gox component template with optional script, style, and template blocks
	return `<style>
.gofront-card {
  border-radius: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  padding: 1rem;
  max-width: 50rem;
  text-align: center;
}
</style>

<template>
<Div class="gofront-card">
  <Text>Counter Component</Text>
</Div>
</template>
`
}

func pageFile() string {
	// .gox page template (pages are automatically suffixed with "Page" in function name)
	return `<style>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>

<template>
<Div class="home-page">
  <Div class="mb-4 text-sm font-medium uppercase tracking-widest text-cyan-300/80">GoFront Starter</Div>
  <Div id="counter" class="text-7xl font-black tracking-tight text-cyan-300">0</Div>
  <Div class="mt-4">
    <Button id="increment">Increment</Button>
    <Button id="reset">Reset</Button>
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
