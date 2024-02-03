package raylibrenderer

func (r *RaylibRenderer) drawedge(x1, y1, x2, y2 int32, z1, z2 float64) {
	side := 0
	xslope := float64(x2-x1) / float64(y2-y1)
	zslope := (z2 - z1) / float64(y2-y1)
	if y1 >= y2 {
		side = 1
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
	}
	currX := float64(x1)
	currZ := float64(z1)
	for y := y1; y <= y2; y++ {
		if y >= 0 && y < int32(len(xedge)) {
			xedge[y][side] = currX
			zpos[y][side] = currZ
		}
		currX += xslope
		currZ += zslope
	}
}

func (r *RaylibRenderer) drawFilledTriangle(x1, y1, x2, y2, x3, y3 int32, z1, z2, z3 float64, color byte) {

	var minx, maxx int32
	r.drawedge(x1, y1, x2, y2, z1, z2)
	r.drawedge(x2, y2, x3, y3, z2, z3)
	r.drawedge(x3, y3, x1, y1, z3, z1)
	miny := y1
	if miny > y2 {
		miny = y2
	}
	if miny > y3 {
		miny = y3
	}
	if miny < 0 {
		miny = 0
	}

	maxy := y1
	if maxy < y2 {
		maxy = y2
	}
	if maxy < y3 {
		maxy = y3
	}
	minx = x1
	if minx > x2 {
		minx = x2
	}
	if minx > x3 {
		minx = x3
	}
	maxx = x1
	if maxx < x2 {
		maxx = x2
	}
	if maxx < x3 {
		maxx = x3
	}
	for y := miny; y <= maxy; y++ {
		r.HLineZBuf(int32(xedge[y][0]), int32(xedge[y][1]), y, color)
	}
}

func (r *RaylibRenderer) HLineZBuf(x1, x2, y int32, color byte) {
	z1 := zpos[y][0]
	z2 := zpos[y][1]
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
