package renderer

import (
	"math"
	"totala_reader/geometry"
	"totala_reader/ta_files_read/texture"
)

type triangle struct {
	coords            [3][3]float64
	texture           *texture.GafEntry
	uvCoords          [3][2]float64
	middleY, middleZ  float64
	colorPaletteIndex byte
}

func (t *triangle) rotate(degrees int) {
	// only around Y for now
	radians := float64(degrees) * math.Pi / 180
	for i := range t.coords {
		t.coords[i][0], t.coords[i][1], t.coords[i][2] = geometry.Rotate3dCoordsAroundY(t.coords[i][0],
			t.coords[i][1], t.coords[i][2], radians)
	}
}

func (t *triangle) calcMiddle() {
	for i := range t.coords {
		t.middleY += t.coords[i][1]
		t.middleZ += t.coords[i][2]
	}
	t.middleY /= 3
	t.middleZ /= 3
}
