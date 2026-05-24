package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Scaffold(root string) error {
	moduleContent, err := moduleFile(root)
	if err != nil {
		return err
	}
	files := map[string]string{
		filepath.Join(root, "go.mod"):                          moduleContent,
		filepath.Join(root, "src", "app.go"):                   appFile(),
		filepath.Join(root, "src", "styles.css"):               stylesFile(),
		filepath.Join(root, "src", "components", "counter.go"): componentFile(),
		filepath.Join(root, "src", "pages", "home.go"):         pageFile(),
		filepath.Join(root, "public", "index.html"):            indexFile(),
		filepath.Join(root, "tailwind.config.js"):              tailwindConfigFile(),
		filepath.Join(root, "gofront.config.json"):             configFile(),
		filepath.Join(root, ".gitignore"):                      gitignoreFile(),
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func moduleFile(root string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	moduleName := sanitizeModuleName(filepath.Base(absRoot))
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
	return "package main\n\nimport \"gofront\"\n\nfunc main() {\n\tgofront.Run(func() {\n\t\tcount := gofront.State(0)\n\n\t\tgofront.Query(\"#increment\").OnClick(func() {\n\t\t\tcount.Set(count.Get() + 1)\n\t\t\tgofront.Query(\"#counter\").SetText(count.Get())\n\t\t})\n\n\t\tgofront.Query(\"#reset\").OnClick(func() {\n\t\t\tcount.Set(0)\n\t\t\tgofront.Query(\"#counter\").SetText(count.Get())\n\t\t})\n\n\t\tgofront.Query(\"body\").AddClass(\"gofront-ready\")\n\t})\n}\n"
}

func componentFile() string {
	return "package components\n\nimport \"gofront\"\n\nfunc Counter() gofront.Component {\n\treturn gofront.Div(\n\t\tgofront.Text(\"Tailwind Counter\"),\n\t).With(\n\t\tgofront.Class(\"gofront-card max-w-2xl p-8 text-center\"),\n\t)\n}\n"
}

func pageFile() string {
	return "package pages\n\nimport \"gofront\"\n\nfunc HomePage() gofront.Component {\n\treturn gofront.Div(\n\t\tgofront.Div(\n\t\t\tgofront.Text(\"GoFront Counter\"),\n\t\t).With(\n\t\t\tgofront.Class(\"mb-4 text-sm font-medium uppercase tracking-[0.35em] text-cyan-300/80\"),\n\t\t),\n\t\tgofront.Div(\n\t\t\tgofront.Text(\"0\"),\n\t\t).With(\n\t\t\tgofront.ID(\"counter\"),\n\t\t\tgofront.Class(\"text-7xl font-black tracking-tight text-cyan-300\"),\n\t\t),\n\t)\n}\n"
}
