package tui

import (
	"fmt"
	"image-processing/imageio"
	"image-processing/internal/ascii_preview"
)

func (m *Model) loadImagePreview(path string) {
	file, err := imageio.OpenBmpImage(path)
	if err != nil {
		m.UIState.err = fmt.Errorf("failed to open image: %v", err)
		return
	}

	availableHeight := m.terminalSize.height

	convertOptions := ascii_preview.DefaultOptions
	convertOptions.FixedWidth = availableHeight * 2
	convertOptions.FixedHeight = availableHeight

	converter := ascii_preview.NewImageConverter()
	converted := converter.Image2ASCIIString(file, &convertOptions)

	m.imagePreview = converted
}

type terminalSize struct {
	width  int
	height int
}
