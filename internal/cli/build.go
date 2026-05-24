package cli

import (
	"flag"

	"gofront/internal/compiler"
	"gofront/internal/config"
)

func runBuild(args []string) error {
	if err := requireProjectRoot(); err != nil {
		return err
	}
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
	configPath := fs.String("config", config.DefaultPath, "path to gofront config file")
	src := fs.String("src", "", "source directory")
	public := fs.String("public", "", "public directory")
	dist := fs.String("dist", "", "output directory")
	release := fs.Bool("release", true, "enable release build mode")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}
	if *src != "" {
		cfg.Src = *src
	}
	if *public != "" {
		cfg.Public = *public
	}
	if *dist != "" {
		cfg.Dist = *dist
	}

	return compiler.Build(compiler.Options{
		Release:   *release,
		SrcDir:    cfg.Src,
		PublicDir: cfg.Public,
		DistDir:   cfg.Dist,
	})
}
