package raylibrenderer

import "math"

type triangle struct {
	coords            [3][3]float64
	middleY, middleZ  float64
	colorPaletteIndex int
}

func (t *triangle) rotate(degrees int) {
	// only 2d for now
	radians := float64(degrees) * math.Pi / 180
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	for i := range t.coords {
		temp := t.coords[i][0]
		t.coords[i][0] = t.coords[i][0]*cos - t.coords[i][2]*sin
		t.coords[i][2] = temp*sin + t.coords[i][2]*cos
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
