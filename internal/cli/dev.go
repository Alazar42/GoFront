package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"gofront/internal/compiler"
	"gofront/internal/config"
	"gofront/internal/server"
)

func runDev(args []string) error {
	if err := requireProjectRoot(); err != nil {
		return err
	}

	fs := flag.NewFlagSet("dev", flag.ContinueOnError)
	fs.SetOutput(flag.CommandLine.Output())
	configPath := fs.String("config", config.DefaultPath, "path to gofront config file")
	addr := fs.String("addr", ":3000", "server address")
	src := fs.String("src", "", "source directory")
	public := fs.String("public", "", "public directory")
	dist := fs.String("dist", "", "output directory")
	pollMS := fs.Int("poll", 1000, "watch polling interval in milliseconds")
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
	if *pollMS < 100 {
		*pollMS = 100
	}

	reload := server.NewReloadHub()
	if err := compiler.Build(compiler.Options{Dev: true, ReloadHub: reload, SrcDir: cfg.Src, PublicDir: cfg.Public, DistDir: cfg.Dist}); err != nil {
		return err
	}

	go watchProject(reload, cfg, time.Duration(*pollMS)*time.Millisecond)
	return server.Serve(server.Options{Addr: *addr, Root: cfg.Dist, Dev: true, Logger: os.Stdout, Handler: http.DefaultServeMux, ReloadHub: reload})
}

func watchProject(reload *server.ReloadHub, cfg config.ProjectConfig, interval time.Duration) {
	last := compiler.SnapshotProject(cfg.Src, cfg.Public)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		current := compiler.SnapshotProject(cfg.Src, cfg.Public)
		if current == last {
			continue
		}
		if err := compiler.Build(compiler.Options{Dev: true, ReloadHub: reload, SrcDir: cfg.Src, PublicDir: cfg.Public, DistDir: cfg.Dist}); err != nil {
			_ = compiler.WriteBuildError(cfg.Dist, err.Error())
			fmt.Fprintln(os.Stderr, err)
			reload.Broadcast()
			continue
		}
		_ = compiler.ClearBuildError(cfg.Dist)
		last = current
		reload.Broadcast()
	}
}
