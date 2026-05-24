package cli

import (
	"flag"
	"fmt"

	"gofront/internal/project"
)

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
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
	return project.Scaffold(target)
}
