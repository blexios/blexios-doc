package renderer

import "blexios/internal/doc/model"

type Renderer interface {
	Render(document model.Document, outputPath string) error
}
