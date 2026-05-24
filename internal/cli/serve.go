package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"gofront/internal/config"
	"gofront/internal/server"
)

func runServe(args []string) error {
	if err := requireProjectRoot(); err != nil {
		return err
	}
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: gofront serve [options]\n\nServes static output files from a dist directory.")
		fmt.Fprintln(fs.Output(), "\nExamples:\n  gofront serve\n  gofront serve --addr :8080 --root ./dist")
		fs.PrintDefaults()
	}
	configPath := fs.String("config", config.DefaultPath, "path to gofront config file")
	addr := fs.String("addr", ":3000", "server address")
	root := fs.String("root", "", "directory to serve")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}
	if *root != "" {
		cfg.Dist = *root
	}

	return server.Serve(server.Options{Addr: *addr, Root: cfg.Dist, Logger: os.Stdout, Handler: http.DefaultServeMux})
}
