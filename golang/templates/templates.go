package templates

import (
	"embed"
	"io/fs"
)

//go:embed template.*
var EmbeddedTemplates embed.FS

// GetTemplatesFS returns a filesystem for use with http.FS or for reading files.
func GetTemplatesFS() (fs.FS, error) {
	subFS, err := fs.Sub(EmbeddedTemplates, ".")
	if err != nil {
		return nil, err
	}
	return subFS, nil
}
