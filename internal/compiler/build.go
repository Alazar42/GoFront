package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	Release   bool
	Dev       bool
	SrcDir    string
	PublicDir string
	DistDir   string
	ReloadHub interface{ Broadcast() }
}

func Build(opts Options) error {
	srcDir := defaultDir(opts.SrcDir, "src")
	publicDir := defaultDir(opts.PublicDir, "public")
	distDir := defaultDir(opts.DistDir, "dist")

	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return err
	}
	files, err := DiscoverFrontendFiles(srcDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no frontend .go files found under %s", srcDir)
	}
	if err := buildWasm(srcDir, distDir); err != nil {
		return err
	}
	if err := writeLoader(distDir); err != nil {
		return err
	}
	if err := writeIndex(publicDir, distDir); err != nil {
		return err
	}
	_ = ClearBuildError(distDir)
	return copyAssets(publicDir, distDir)
}

func DiscoverFrontendFiles(srcDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func buildWasm(srcDir, distDir string) error {
	buildPath := "./" + filepath.ToSlash(filepath.Clean(srcDir))
	cmd := exec.Command("go", "build", "-o", filepath.Join(distDir, "app.wasm"), buildPath)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build wasm: %w\n%s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func writeLoader(distDir string) error {
	goRoot, err := goRoot()
	if err != nil {
		return err
	}
	loaderPaths := []string{
		filepath.Join(goRoot, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(goRoot, "misc", "wasm", "wasm_exec.js"),
	}
	var content []byte
	for _, loaderPath := range loaderPaths {
		content, err = os.ReadFile(loaderPath)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	bootstrap := "\n\nconst gofront = async () => {\n  const go = new Go();\n  const response = await fetch('app.wasm');\n  const bytes = await response.arrayBuffer();\n  const { instance } = await WebAssembly.instantiate(bytes, go.importObject);\n  go.run(instance);\n};\n\nconst __gofront_showOverlay = (message) => {\n  let overlay = document.getElementById('__gofront_error_overlay');\n  if (!overlay) {\n    overlay = document.createElement('pre');\n    overlay.id = '__gofront_error_overlay';\n    overlay.style.position = 'fixed';\n    overlay.style.inset = '0';\n    overlay.style.margin = '0';\n    overlay.style.padding = '24px';\n    overlay.style.zIndex = '2147483647';\n    overlay.style.background = 'rgba(16, 18, 27, 0.96)';\n    overlay.style.color = '#ffb4b4';\n    overlay.style.fontFamily = 'ui-monospace, SFMono-Regular, Menlo, Consolas, monospace';\n    overlay.style.fontSize = '14px';\n    overlay.style.whiteSpace = 'pre-wrap';\n    document.body.appendChild(overlay);\n  }\n  overlay.textContent = message;\n  overlay.style.display = 'block';\n};\n\nconst __gofront_hideOverlay = () => {\n  const overlay = document.getElementById('__gofront_error_overlay');\n  if (overlay) {\n    overlay.remove();\n  }\n};\n\nconst __gofront_checkOverlay = async () => {\n  try {\n    const response = await fetch('/__gofront_error', { cache: 'no-store' });\n    const message = (await response.text()).trim();\n    if (message) {\n      __gofront_showOverlay(message);\n    } else {\n      __gofront_hideOverlay();\n    }\n  } catch (error) {\n    __gofront_hideOverlay();\n  }\n};\n\nif (typeof window !== 'undefined') {\n  if (window.__GOFRONT_DEV__ && typeof EventSource !== 'undefined') {\n    const source = new EventSource('/__gofront_reload');\n    source.onmessage = () => window.location.reload();\n    setInterval(() => { __gofront_checkOverlay().catch(() => {}); }, 500);\n  }\n  window.addEventListener('load', () => { gofront().catch(console.error); });\n}\n"
	return os.WriteFile(filepath.Join(distDir, "gofront.js"), append(content, []byte(bootstrap)...), 0o644)
}

func writeIndex(publicDir, distDir string) error {
	input := []byte(defaultIndex())
	if source, err := os.ReadFile(filepath.Join(publicDir, "index.html")); err == nil {
		input = source
	}
	output := bytes.ReplaceAll(input, []byte(`<script src="app.go"></script>`), []byte(`<script src="gofront.js"></script>`))
	if !bytes.Contains(output, []byte("gofront.js")) {
		output = append(output, []byte("\n<script src=\"gofront.js\"></script>\n")...)
	}
	return os.WriteFile(filepath.Join(distDir, "index.html"), output, 0o644)
}

func copyAssets(publicDir, distDir string) error {
	src := filepath.Join(publicDir, "assets")
	if _, err := os.Stat(src); err != nil {
		return nil
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(publicDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(distDir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

func SnapshotProject(srcDir, publicDir string) string {
	var builder strings.Builder
	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		builder.WriteString(path)
		builder.WriteString(":")
		builder.WriteString(info.ModTime().UTC().String())
		builder.WriteByte('\n')
		return nil
	})
	_ = filepath.Walk(publicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		builder.WriteString(path)
		builder.WriteString(":")
		builder.WriteString(info.ModTime().UTC().String())
		builder.WriteByte('\n')
		return nil
	})
	return builder.String()
}

func goRoot() (string, error) {
	cmd := exec.Command("go", "env", "GOROOT")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func defaultIndex() string {
	return "<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>GoFront App</title>\n</head>\n<body>\n  <h1 id=\"counter\">0</h1>\n  <button id=\"increment\">Increment</button>\n  <script src=\"app.go\"></script>\n</body>\n</html>\n"
}

func WriteBuildError(distDir, message string) error {
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(distDir, ".gofront-error"), []byte(message), 0o644)
}

func ClearBuildError(distDir string) error {
	if err := os.Remove(filepath.Join(distDir, ".gofront-error")); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func defaultDir(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return filepath.Clean(value)
}
