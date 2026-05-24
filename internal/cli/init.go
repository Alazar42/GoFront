package cli

import (
	"flag"
	"fmt"

	"github.com/Alazar42/GoFront/internal/project"
)

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: gofront init [name] [--name my-app]\n\nCreate a new GoFront project scaffold.")
		fs.PrintDefaults()
	}
	name := fs.String("name", "", "target directory name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	remaining := fs.Args()
	if len(remaining) > 1 {
		return fmt.Errorf("init accepts at most one path argument")
	}
	target := "."
	if *name != "" {
		target = *name
	}
	if len(remaining) == 1 && remaining[0] != "" {
		target = remaining[0]
	}
	moduleName := target
	return project.Scaffold(target, moduleName)
}
