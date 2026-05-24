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
		return fmt.Errorf("unknown command %q", args[0])
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
	_, err := fmt.Fprintln(os.Stdout, "GoFront CLI\n\nCommands:\n  gofront init [name]\n  gofront build [--src dir] [--public dir] [--dist dir]\n  gofront dev [--addr :3000] [--src dir] [--public dir] [--dist dir] [--poll 1000]\n  gofront serve [--addr :3000] [--root dist]\n  gofront version\n\nRun `gofront <command> -h` for command options.")
	return err
}

func requireProjectRoot() error {
	if _, err := os.Stat(filepath.Join("go.mod")); err != nil {
		return fmt.Errorf("run this command from a GoFront project root")
	}
	return nil
}
