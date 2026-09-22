//go:build ignore

package main

import (
	"os"
	"path/filepath"

	"github.com/cashtro/cashtro/internal/orgmap"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	write(filepath.Join(root, "docs/FLEET_MAP.md"), []byte(orgmap.Markdown()))
	write(filepath.Join(root, "inventory/joint-map.json"), orgmap.JSON())
}

func write(path string, body []byte) {
	if err := os.WriteFile(path, body, 0o644); err != nil {
		panic(err)
	}
}
