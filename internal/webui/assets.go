package webui

import (
	"embed"
)

//go:embed index.html
//go:embed assets
var Assets embed.FS
