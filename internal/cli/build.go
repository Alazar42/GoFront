package cli

import (
	"flag"
	"fmt"

	"github.com/Alazar42/GoFront/internal/compiler"
	"github.com/Alazar42/GoFront/internal/config"
)

func runBuild(args []string) error {
	if err := requireProjectRoot(); err != nil {
		return err
	}
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: gofront build [options]\n\nBuilds frontend sources into WebAssembly output.")
		fmt.Fprintln(fs.Output(), "\nExamples:\n  gofront build\n  gofront build --src web --public static --dist out")
		fs.PrintDefaults()
	}
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
