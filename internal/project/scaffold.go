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
		filepath.Join(root, "src", "components", "counter.go"): componentFile(),
		filepath.Join(root, "src", "pages", "home.go"):         pageFile(),
		filepath.Join(root, "public", "index.html"):            indexFile(),
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
	return "<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>GoFront App</title>\n</head>\n<body>\n  <h1 id=\"counter\">0</h1>\n  <button id=\"increment\">Increment</button>\n  <script src=\"app.go\"></script>\n</body>\n</html>\n"
}

func appFile() string {
	return "package main\n\nimport \"gofront\"\n\nfunc main() {\n\tcount := gofront.State(0)\n\tbutton := gofront.Query(\"#increment\")\n\tbutton.OnClick(func() {\n\t\tcount.Set(count.Get() + 1)\n\t\tgofront.Query(\"#counter\").SetText(count.Get())\n\t})\n\n\tgofront.Wait()\n}\n"
}

func componentFile() string {
	return "package components\n\nimport \"gofront\"\n\nfunc Counter() gofront.Component {\n\treturn gofront.Div(\n\t\tgofront.Text(\"Hello\"),\n\t)\n}\n"
}

func pageFile() string {
	return "package pages\n\nimport \"gofront\"\n\nfunc HomePage() gofront.Component {\n\treturn gofront.Div(\n\t\tgofront.Text(\"Home\"),\n\t)\n}\n"
}
