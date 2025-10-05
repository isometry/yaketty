package main

import (
	"embed"

	"github.com/isometry/yaketty/internal/library"
)

//go:embed personas/*.yaml scenarios/*.yaml
var embeddedFS embed.FS

func init() {
	library.SetEmbeddedFS(embeddedFS)
}
