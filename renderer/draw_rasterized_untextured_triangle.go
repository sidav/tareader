package renderer

func (r *Renderer) drawRasterizedFilledTriangle(ost *onScreenTriangle) {
	ost.reorderCoordsYAscSkipUv()
	var currx1, currx2 int32
	dx1 := ost.x1 - ost.x0
	dy1 := ost.y1 - ost.y0
	dx2 := ost.x2 - ost.x0
	dy2 := ost.y2 - ost.y0

	dy1f := float64(dy1)
	dy2f := float64(dy2)
	dz1 := (ost.z1 - ost.z0) / dy1f
	dz2 := (ost.z2 - ost.z0) / dy2f

	currz1 := ost.z0
	currz2 := ost.z0

	curry := ost.y0
	for curry < ost.y1 {
		currx1 = ost.x0 + dx1*(curry-ost.y0)/dy1
		currx2 = ost.x0 + dx2*(curry-ost.y0)/dy2

		r.HLineZBufRstr(currx1, currx2, curry, ost.color, currz1, currz2)
		currz1 += dz1
		currz2 += dz2
		curry++
	}

	dx1 = ost.x2 - ost.x1
	dy1 = ost.y2 - ost.y1

	dy1f = float64(dy1)
	dz1 = (ost.z2 - ost.z1) / dy1f
	currz1 = ost.z1

	for curry < ost.y2 {
		currx1 = ost.x1 + dx1*(curry-ost.y1)/dy1
		currx2 = ost.x0 + dx2*(curry-ost.y0)/dy2

		r.HLineZBufRstr(currx1, currx2, curry, ost.color, currz1, currz2)
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
