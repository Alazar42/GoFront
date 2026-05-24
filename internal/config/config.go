package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const DefaultPath = "gofront.config.json"

type ProjectConfig struct {
	Name   string `json:"name"`
	Src    string `json:"src"`
	Public string `json:"public"`
	Dist   string `json:"dist"`
}

func Defaults() ProjectConfig {
	return ProjectConfig{
		Name:   "GoFront App",
		Src:    "src",
		Public: "public",
		Dist:   "dist",
	}
}

func Load(path string) (ProjectConfig, error) {
	if path == "" {
		path = DefaultPath
	}
	cfg := Defaults()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return ProjectConfig{}, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ProjectConfig{}, err
	}
	normalize(&cfg)
	return cfg, nil
}

func normalize(cfg *ProjectConfig) {
	if cfg.Name == "" {
		cfg.Name = "GoFront App"
	}
	if cfg.Src == "" {
		cfg.Src = "src"
	}
	if cfg.Public == "" {
		cfg.Public = "public"
	}
	if cfg.Dist == "" {
		cfg.Dist = "dist"
	}
	cfg.Src = filepath.Clean(cfg.Src)
	cfg.Public = filepath.Clean(cfg.Public)
	cfg.Dist = filepath.Clean(cfg.Dist)
}
