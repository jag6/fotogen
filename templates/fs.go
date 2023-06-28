package templates

import "embed"

//go:embed pages/*.html users/*.html *.html
var FS embed.FS
