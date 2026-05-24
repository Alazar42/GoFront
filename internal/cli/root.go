package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

const version = "0.1.0"

func Execute(args []string) error {
	if len(args) == 0 {
		return runHelp()
	}

	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "build":
		return runBuild(args[1:])
	case "dev":
		return runDev(args[1:])
	case "serve":
		return runServe(args[1:])
	case "version":
		fmt.Println(version)
		return nil
	case "help", "-h", "--help":
		return runHelp()
	default:
		return fmt.Errorf("unknown command %q\n\nRun `gofront help` to see available commands", args[0])
	}
}

func ExitWithError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func runHelp() error {
	output := "GoFront CLI\n" +
		"Build frontend apps in Go + WebAssembly.\n\n" +
		"Usage:\n" +
		"  gofront <command> [options]\n\n" +
		"Commands:\n" +
		"  init     Create a new GoFront project scaffold\n" +
		"  build    Compile frontend source into dist artifacts\n" +
		"  dev      Start dev server with auto rebuild + live reload\n" +
		"  serve    Serve an existing dist directory\n" +
		"  version  Print installed CLI version\n\n" +
		"Quick Start:\n" +
		"  gofront init my-app\n" +
		"  cd my-app\n" +
		"  gofront dev\n\n" +
		"Run `gofront <command> -h` for command-specific options."
	_, err := fmt.Fprintln(os.Stdout, output)
	return err
}

func requireProjectRoot() error {
	if _, err := os.Stat(filepath.Join("go.mod")); err != nil {
		return fmt.Errorf("run this command from a GoFront project root")
	}
	return nil
}
