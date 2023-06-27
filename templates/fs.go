package templates

import "embed"

//go:embed components/*.html pages/*.html users/*.html *.html
var FS embed.FS
