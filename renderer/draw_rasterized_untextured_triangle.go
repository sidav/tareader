package renderer

func (r *Renderer) drawRasterizedFilledTriangle(x0, y0, x1, y1, x2, y2 int32, z0, z1, z2 float64, color byte) {
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
		z0, z1 = z1, z0
	}
	if y0 > y2 {
		x0, x2 = x2, x0
		y0, y2 = y2, y0
		z0, z2 = z2, z0
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
	}

	var currx1, currx2 int32
	dx1 := x1 - x0
	dy1 := y1 - y0
	dx2 := x2 - x0
	dy2 := y2 - y0

	dy1f := float64(dy1)
	dy2f := float64(dy2)
	dz1 := (z1 - z0) / dy1f
	dz2 := (z2 - z0) / dy2f

	currz1 := z0
	currz2 := z0

	curry := y0
	for curry < y1 {
		currx1 = x0 + dx1*(curry-y0)/dy1
		currx2 = x0 + dx2*(curry-y0)/dy2

		r.HLineZBufRstr(currx1, currx2, curry, color, currz1, currz2)
		currz1 += dz1
		currz2 += dz2
		curry++
	}

	dx1 = x2 - x1
	dy1 = y2 - y1

	dy1f = float64(dy1)
	dz1 = (z2 - z1) / dy1f
	currz1 = z1

	for curry < y2 {
		currx1 = x1 + dx1*(curry-y1)/dy1
		currx2 = x0 + dx2*(curry-y0)/dy2

		r.HLineZBufRstr(currx1, currx2, curry, color, currz1, currz2)
		currz1 += dz1
		currz2 += dz2
		curry++
	}
}

func (r *Renderer) HLineZBufRstr(x1, x2, y int32, color byte, z1, z2 float64) {
	if x1 > x2 {
		x1, x2 = x2, x1
		z1, z2 = z2, z1
	}
	zinc := (z2 - z1) / float64(x2-x1)
	for x := x1; x <= x2; x++ {
		if r.canDrawOverZBufferAt(x, y, z1) {
			r.gAdapter.SetColor(getTaPaletteColor(uint8(color)))
			r.setZBufferValueAt(z1, x, y)
			r.gAdapter.DrawPoint(x, y)
		}
		z1 += zinc
	}
}
