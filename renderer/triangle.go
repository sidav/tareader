package renderer

import (
	"totala_reader/ta_files_read/texture"
)

type triangle struct {
	coords            [3][3]float64
	uvCoords          [3][2]float64
	texture           *texture.GafEntry
	colorPaletteIndex byte
}
