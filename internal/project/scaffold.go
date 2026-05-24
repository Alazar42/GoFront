package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Scaffold(root string, moduleName string) error {
	moduleContent, err := moduleFile(root, moduleName)
	if err != nil {
		return err
	}
	files := map[string]string{
		filepath.Join(root, "go.mod"):                           moduleContent,
		filepath.Join(root, "src", "main.go"):                   mainGoFile(),
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

	// Run go mod download to generate go.sum
	if err := runGoModDownload(root); err != nil {
		return fmt.Errorf("failed to download modules: %w", err)
	}
	return nil
}

func moduleFile(root string, moduleName string) (string, error) {
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
		// released dependency. Update this tag when releasing newer versions.
		builder.WriteString("\nrequire github.com/Alazar42/GoFront v0.0.3\n")
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
	return "dist/\n*.wasm\n*.exe\n.idea/\n.vscode/\n"
}

func indexFile() string {
	return "<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <link rel=\"stylesheet\" href=\"styles.css\">\n  <title>GoFront Counter</title>\n</head>\n<body class=\"bg-slate-950 text-slate-100\">\n  <main class=\"relative mx-auto flex min-h-screen max-w-4xl items-center justify-center px-6 py-12\">\n    <div class=\"pointer-events-none absolute inset-0 overflow-hidden\">\n      <div class=\"absolute left-1/2 top-0 h-72 w-72 -translate-x-1/2 rounded-full bg-cyan-400/20 blur-3xl\"></div>\n      <div class=\"absolute bottom-0 right-0 h-72 w-72 rounded-full bg-indigo-500/20 blur-3xl\"></div>\n    </div>\n\n    <section class=\"relative w-full max-w-2xl rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl shadow-cyan-950/30 backdrop-blur-xl sm:p-12\">\n      <div class=\"mb-8 flex flex-col gap-3\">\n        <p class=\"text-sm font-medium uppercase tracking-[0.35em] text-cyan-300/80\">GoFront Starter</p>\n        <h1 class=\"text-4xl font-semibold tracking-tight sm:text-5xl\">Tailwind Counter</h1>\n        <p class=\"max-w-xl text-sm leading-6 text-slate-300 sm:text-base\">A clean starter page for new GoFront apps with Tailwind-powered layout, counter actions, and a polished first impression.</p>\n      </div>\n\n      <div class=\"flex flex-col items-center gap-8 rounded-2xl border border-white/10 bg-slate-950/60 px-6 py-10 text-center sm:px-10\">\n        <div class=\"space-y-2\">\n          <p class=\"text-sm uppercase tracking-[0.3em] text-slate-400\">Current Count</p>\n          <h2 id=\"counter\" class=\"text-7xl font-black tracking-tight text-cyan-300 sm:text-8xl\">0</h2>\n        </div>\n\n        <div class=\"flex flex-wrap items-center justify-center gap-4\">\n          <button id=\"increment\" class=\"inline-flex items-center justify-center rounded-full bg-cyan-400 px-6 py-3 text-sm font-semibold text-slate-950 transition duration-200 hover:-translate-y-0.5 hover:bg-cyan-300 hover:shadow-lg hover:shadow-cyan-400/25 focus:outline-none focus:ring-2 focus:ring-cyan-300 focus:ring-offset-2 focus:ring-offset-slate-950\">Increment</button>\n          <button id=\"reset\" class=\"inline-flex items-center justify-center rounded-full border border-white/15 bg-white/5 px-6 py-3 text-sm font-semibold text-slate-100 transition duration-200 hover:-translate-y-0.5 hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-white/20 focus:ring-offset-2 focus:ring-offset-slate-950\">Reset</button>\n        </div>\n\n        <p class=\"max-w-md text-sm leading-6 text-slate-400\">This starter shows how styles, DOM updates, and event wiring work together in a GoFront app.</p>\n      </div>\n    </section>\n  </main>\n  <script src=\"app.go\"></script>\n</body>\n</html>\n"
}

func stylesFile() string {
	return "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n\n@layer base {\n  html {\n    @apply bg-slate-950 text-slate-100 antialiased;\n  }\n\n  body {\n    @apply min-h-screen overflow-x-hidden;\n  }\n\n  ::selection {\n    @apply bg-cyan-300 text-slate-950;\n  }\n}\n\n@layer components {\n  .gofront-card {\n    @apply rounded-3xl border border-white/10 bg-white/5 shadow-2xl shadow-cyan-950/30 backdrop-blur-xl;\n  }\n}\n"
}

func tailwindConfigFile() string {
	return "module.exports = {\n  content: [\n    './src/**/*.go',\n    './src/**/*.html',\n    './public/**/*.html'\n  ],\n  theme: {\n    extend: {}\n  },\n  plugins: []\n};\n"
}

func appFile() string {
	// .gox component template (JSX-like) for the main app
	return "<Div class=\"app-root\">\n  <Div id=\"counter\" class=\"text-7xl font-black text-cyan-300\">0</Div>\n  <Div class=\"mt-4\">\n    <Button id=\"increment\">Increment</Button>\n    <Button id=\"reset\">Reset</Button>\n  </Div>\n</Div>\n"
}

func componentFile() string {
	// .gox component template
	return "<Div class=\"gofront-card max-w-2xl p-8 text-center\">\n  <Text>Tailwind Counter</Text>\n</Div>\n"
}

func pageFile() string {
	// .gox page template
	return "<Div>\n  <Div class=\"mb-4 text-sm font-medium uppercase tracking-[0.35em] text-cyan-300/80\">GoFront Starter</Div>\n  <Div id=\"counter\" class=\"text-7xl font-black tracking-tight text-cyan-300\">0</Div>\n  <Div class=\"mt-4\">\n    <Button id=\"increment\">Increment</Button>\n    <Button id=\"reset\">Reset</Button>\n  </Div>\n</Div>\n"
}

func mainGoFile() string {
	return "package main\n\nimport GoFront \"github.com/Alazar42/GoFront\"\n\nfunc main() {\n\tGoFront.Run(func() {\n\t\t// Mount the generated app component into body\n\t\tGoFront.Mount(\"body\", App)\n\t})\n}\n"

}

func runGoModDownload(root string) error {
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = root
	return cmd.Run()
}
