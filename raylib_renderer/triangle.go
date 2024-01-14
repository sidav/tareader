package raylibrenderer

type triangle struct {
	coords            [3][3]float64
	middleY, middleZ  float64
	colorPaletteIndex int
}

func (t *triangle) calcMiddle() {
	for i := range t.coords {
		t.middleY += t.coords[i][1]
		t.middleZ += t.coords[i][2]
	}
	t.middleY /= 3
	t.middleZ /= 3
}
