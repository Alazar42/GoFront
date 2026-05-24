package cli

import (
	"flag"
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
