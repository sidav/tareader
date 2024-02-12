package renderer

func (r *Renderer) bufferEdge(x1, y1, x2, y2 int32, z1, z2 float64) {
	side := 0
	if y1 >= y2 {
		side = 1
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
	}
	xslope := float64(x2-x1) / float64(y2-y1)
	zslope := (z2 - z1) / float64(y2-y1)
	if y1 == y2 {
		xslope = float64(x2 - x1)
		zslope = (z2 - z1)
	}
	currX := float64(x1)
	currZ := float64(z1)
	for y := y1; y <= y2; y++ {
		if y >= 0 && y < int32(len(xedge)) {
			xedge[y][side] = currX
			zedge[y][side] = currZ
		}
		currX += xslope
		currZ += zslope
	}
}

func (r *Renderer) drawEdgebufFilledTriangle(x0, y0, x1, y1, x2, y2 int32, z0, z1, z2 float64, color byte) {
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

	r.bufferEdge(x0, y0, x1, y1, z0, z1)
	r.bufferEdge(x1, y1, x2, y2, z1, z2)
	r.bufferEdge(x2, y2, x0, y0, z2, z0)
	if y0 < 0 {
		y0 = 0
	}
	if y2 >= int32(len(xedge)) {
		y2 = int32(len(xedge) - 1)
	}
	for y := y0; y <= y2; y++ {
		r.HLineZBuf(int32(xedge[y][0]), int32(xedge[y][1]), y, color)
	}
}

func (r *Renderer) HLineZBuf(x1, x2, y int32, color byte) {
	z1 := zedge[y][0]
	z2 := zedge[y][1]
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
