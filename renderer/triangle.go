package renderer

import (
	"totala_reader/ta_files_read/texture"
)

type triangle struct {
	texture           *texture.GafEntry
	coords            [3][3]float64
	uvCoords          [3][2]float64
	colorPaletteIndex byte
}
