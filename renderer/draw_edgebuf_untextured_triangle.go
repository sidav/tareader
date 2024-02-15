package renderer

func (r *Renderer) drawEdgebufFilledTriangle(ost *onScreenTriangle) {
	ost.reorderCoordsYAscSkipUv()

	r.bufferEdge(ost.x0, ost.y0, ost.x1, ost.y1, ost.z0, ost.z1)
	r.bufferEdge(ost.x1, ost.y1, ost.x2, ost.y2, ost.z1, ost.z2)
	r.bufferEdge(ost.x2, ost.y2, ost.x0, ost.y0, ost.z2, ost.z0)
	if ost.y0 < 0 {
		ost.y0 = 0
	}
	if ost.y2 >= int32(len(xedge)) {
		ost.y2 = int32(len(xedge) - 1)
	}
	for y := ost.y0; y <= ost.y2; y++ {
		r.HLineZBuf(int32(xedge[y][0]), int32(xedge[y][1]), y, ost.color)
	}
}

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
