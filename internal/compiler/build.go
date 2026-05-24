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
	// transpile any .gox files into generated .go files before discovering sources
	if err := TranspileGoxFiles(srcDir); err != nil {
		return err
	}

	files, err := DiscoverFrontendFiles(srcDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no frontend source files (.go or .gox) found under %s", srcDir)
	}
	stylesPresent, err := buildStyles(srcDir, publicDir, distDir)
	if err != nil {
		return err
	}
	if err := buildWasm(srcDir, distDir); err != nil {
		return err
	}
	if err := writeLoader(distDir); err != nil {
		return err
	}
	if err := writeIndex(publicDir, distDir, stylesPresent); err != nil {
		return err
	}
	_ = ClearBuildError(distDir)
	return copyAssets(publicDir, distDir, stylesPresent)
}

func DiscoverFrontendFiles(srcDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_gox_gen.go") {
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

func buildStyles(srcDir, publicDir, distDir string) (bool, error) {
	srcStyles := filepath.Join(srcDir, "styles.css")
	publicStyles := filepath.Join(publicDir, "styles.css")
	outputStyles := filepath.Join(distDir, "styles.css")

	if fileExists(srcStyles) {
		return true, compileTailwindCSS(srcStyles, outputStyles, srcDir, publicDir)
	}

	if fileExists(publicStyles) {
		return true, copyFile(publicStyles, outputStyles)
	}

	return false, nil
}

func compileTailwindCSS(inputPath, outputPath, srcDir, publicDir string) error {
	configPath, cleanup, err := resolveTailwindConfig(srcDir, publicDir)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	commandName, commandArgs, err := tailwindCommand(inputPath, outputPath, configPath)
	if err != nil {
		return err
	}
	cmd := exec.Command(commandName, commandArgs...)
	cmd.Env = append(os.Environ(), "NODE_ENV=production")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build styles: %w\n%s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func resolveTailwindConfig(srcDir, publicDir string) (string, func(), error) {
	for _, candidate := range []string{"tailwind.config.js", "tailwind.config.cjs", "tailwind.config.mjs"} {
		if fileExists(candidate) {
			return candidate, nil, nil
		}
	}

	file, err := os.CreateTemp(".", ".gofront-tailwind-*.cjs")
	if err != nil {
		return "", nil, err
	}
	// Tailwind content patterns: match class names in Go source and generated HTML
	// Pattern explanation:
	// - src/**/*.{go,gox,html}: files to scan
	// - extract: regex patterns to find class names in code
	//   - "class=\\\"[^\\\"]*\\\"" matches HTML class attributes
	//   - GoFront\\.Class\\(\\\"[^\\\"]*\\\" matches GoFront.Class() calls
	//   - on:[a-z]*=\\{[^}]*\\} matches on:event={handler} attributes
	content := fmt.Sprintf(`module.exports = {
  content: [
    './%s/**/*.{go,gox,html}',
    './%s/**/*.{go,gox,html}',
    './%s/**/*.{go,gox,html}'
  ],
  safelist: [
    { pattern: /bg-(gradient|slate|cyan|red).*/ },
    { pattern: /text-(slate|cyan|4xl|5xl|6xl|7xl|center|slate-950).*/ },
    { pattern: /hover:.*/ },
    { pattern: /active:.*/ },
    { pattern: /flex.*/ },
    { pattern: /gap-.*/ },
    { pattern: /px-.*/ },
    { pattern: /py-.*/ },
    { pattern: /rounded-.*/ },
    { pattern: /transition.*/ },
    { pattern: /transform.*/ },
    { pattern: /scale-.*/ },
    { pattern: /border.*/ },
    { pattern: /shadow.*/ },
    { pattern: /backdrop-.*/ },
    { pattern: /min-h-.*/ },
    { pattern: /max-w-.*/ },
    { pattern: /mx-auto/ },
    { pattern: /space-y-.*/ },
    { pattern: /mb-.*/ },
    { pattern: /mt-.*/ },
    { pattern: /p-.*/ },
    { pattern: /bg-clip-text/ },
    { pattern: /bg-white.*/ },
    { pattern: /font-.*/ },
    { pattern: /uppercase/ },
    { pattern: /tracking-.*/ },
    { pattern: /from-.*/ },
    { pattern: /to-.*/ },
    { pattern: /via-.*/ },
    { pattern: /h-screen/ },
    { pattern: /justify-.*/ },
    { pattern: /items-.*/ },
    { pattern: /divide-.*/ }
  ],
  theme: { extend: {} },
  plugins: []
};
`, filepath.ToSlash(filepath.Clean(srcDir)), filepath.ToSlash(filepath.Clean(srcDir)), filepath.ToSlash(filepath.Clean(publicDir)))
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return "", nil, err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return "", nil, err
	}
	return file.Name(), func() { _ = os.Remove(file.Name()) }, nil
}

func tailwindCommand(inputPath, outputPath, configPath string) (string, []string, error) {
	if tailwindPath, err := exec.LookPath("tailwindcss"); err == nil {
		return tailwindPath, []string{"-i", inputPath, "-o", outputPath, "--minify", "--config", configPath}, nil
	}
	if npxPath, err := exec.LookPath("npx"); err == nil {
		return npxPath, []string{"--yes", "tailwindcss@3.4.17", "-i", inputPath, "-o", outputPath, "--minify", "--config", configPath}, nil
	}
	return "", nil, fmt.Errorf("tailwindcss not found. install tailwindcss or npm/npx to build src/styles.css")
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

func writeIndex(publicDir, distDir string, stylesPresent bool) error {
	input := []byte(defaultIndex())
	if source, err := os.ReadFile(filepath.Join(publicDir, "index.html")); err == nil {
		input = source
	}
	if stylesPresent {
		input = ensureStylesheetLink(input)
	}
	output := bytes.ReplaceAll(input, []byte(`<script src="app.go"></script>`), []byte(`<script src="gofront.js"></script>`))
	if !bytes.Contains(output, []byte("gofront.js")) {
		output = append(output, []byte("\n<script src=\"gofront.js\"></script>\n")...)
	}
	if stylesPresent {
		output = ensureStylesheetLink(output)
	}
	return os.WriteFile(filepath.Join(distDir, "index.html"), output, 0o644)
}

func copyAssets(publicDir, distDir string, skipStyles bool) error {
	return filepath.Walk(publicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(publicDir, path)
		if err != nil {
			return err
		}
		if rel == "index.html" {
			return nil
		}
		if skipStyles && rel == "styles.css" {
			return nil
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
	return "<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>GoFront App</title>\n</head>\n<body>\n  <div id=\"app\"></div>\n  <script src=\"gofront.js\"></script>\n</body>\n</html>\n"
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

func ensureStylesheetLink(input []byte) []byte {
	linkTag := []byte(`<link rel="stylesheet" href="styles.css">`)
	if bytes.Contains(input, linkTag) {
		return input
	}
	closingHead := []byte("</head>")
	if idx := bytes.Index(input, closingHead); idx >= 0 {
		result := make([]byte, 0, len(input)+len(linkTag)+1)
		result = append(result, input[:idx]...)
		result = append(result, []byte("  ")...)
		result = append(result, linkTag...)
		result = append(result, byte('\n'))
		result = append(result, input[idx:]...)
		return result
	}
	return append([]byte(string(linkTag)+"\n"), input...)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func copyFile(srcPath, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func defaultDir(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return filepath.Clean(value)
}
