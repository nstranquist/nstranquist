package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	check := flag.Bool("check", false, "exit 1 if generated files are stale")
	root := flag.String("root", ".", "repository root")
	flag.Parse()

	if err := run(*root, *check); err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}
}

func run(root string, check bool) error {
	cat, err := loadCatalog(filepath.Join(root, "catalog.yaml"))
	if err != nil {
		return err
	}

	files := map[string]string{
		filepath.Join(root, "README.md"):                   renderREADME(cat),
		filepath.Join(root, "assets", "banner-dark.svg"):   renderBanner(cat, themeDark()),
		filepath.Join(root, "assets", "banner-light.svg"):  renderBanner(cat, themeLight()),
		filepath.Join(root, "assets", "catalog-dark.svg"):  renderCatalogBoard(cat, themeDark()),
		filepath.Join(root, "assets", "catalog-light.svg"): renderCatalogBoard(cat, themeLight()),
		filepath.Join(root, "preview.html"):                renderPreview(cat),
	}
	for _, p := range cat.Featured {
		slug := slugFor(p)
		files[filepath.Join(root, "assets", "cards", slug+"-dark.svg")] = renderCard(p, themeDark())
		files[filepath.Join(root, "assets", "cards", slug+"-light.svg")] = renderCard(p, themeLight())
	}

	var stale []string
	for path, body := range files {
		if check {
			existing, err := os.ReadFile(path)
			if err != nil {
				stale = append(stale, path+" (missing)")
				continue
			}
			if !bytes.Equal(existing, []byte(body)) {
				stale = append(stale, path)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	if check && len(stale) > 0 {
		return fmt.Errorf("generated files are stale:\n  %s", strings.Join(stale, "\n  "))
	}
	return nil
}

func slugFor(p Product) string {
	name := strings.TrimPrefix(p.ID, "product.")
	return strings.ReplaceAll(name, "_", "-")
}
