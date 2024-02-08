package renderer

// ZBuffered horizontal scanline fill triangle algorithm.
func (r *Renderer) drawScanlineFilledTriangle(x0, y0, x1, y1, x2, y2 int32, z0, z1, z2 float64, color byte) {
	var a, b, y, last int32
	var az, bz float64 // for z-buffering
	// Sort coordinates by Y order (y2 >= y1 >= y0)
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
		z0, z1 = z1, z0
	}
	if y1 > y2 {
		x2, y2, x1, y1 = x1, y1, x2, y2
		z2, z1 = z1, z2
	}
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
		z0, z1 = z1, z0
	}

	if y0 == y2 { // All on same line case
		a = x0
		b = x0
		az = z0
		bz = z0
		if x1 < a {
			a = x1
			az = z1
		} else if x1 > b {
			b = x1
			bz = z1
		}
		if x2 < a {
			a = x2
			az = z2
		} else if x2 > b {
			b = x2
			bz = z2
		}
		r.HLineZBufNoArr(a, b, y0, color, az, bz)
		return
	}

	dx01 := x1 - x0
	dy01 := y1 - y0
	dz01 := z1 - z0
	dx02 := x2 - x0
	dy02 := y2 - y0
	dz02 := z2 - z0
	dx12 := x2 - x1
	dy12 := y2 - y1
	dz12 := z2 - z1
	var sa, sb int32
	var saz, sbz float64

	// For upper part of triangle, find scanline crossings for segment
	// 0-1 and 0-2.  If y1=y2 (flat-bottomed triangle), the scanline y
	// is included here (and second loop will be skipped, avoiding a /
	// error there), otherwise scanline y1 is skipped here and handle
	// in the second loop...which also avoids a /0 error here if y0=y
	// (flat-topped triangle)
	if y1 == y2 {
		last = y1 // Include y1 scanline
	} else {
		last = y1 - 1
	} // Skip it

	for y = y0; y <= last; y++ {
		a = x0 + sa/dy01
		b = x0 + sb/dy02
		az = z0 + saz/float64(dy01)
		bz = z0 + sbz/float64(dy02)
		sa += dx01
		sb += dx02
		saz += dz01
		sbz += dz02
		// longhand a = x0 + (x1 - x0) * (y - y0) / (y1 - y0)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		r.HLineZBufNoArr(a, b, y, color, az, bz)
	}

	// For lower part of triangle, find scanline crossings for segment
	// 0-2 and 1-2.  This loop is skipped if y1=y2
	sa = dx12 * (y - y1)
	sb = dx02 * (y - y0)
	saz = dz12 * float64(y-y1)
	sbz = dz02 * float64(y-y0)

	for ; y <= y2; y++ {
		a = x1 + sa/dy12
		b = x0 + sb/dy02
		az = z1 + saz/float64(dy12)
		bz = z0 + sbz/float64(dy02)
		sa += dx12
		sb += dx02
		saz += dz12
		sbz += dz02
		// longhand a = x1 + (x2 - x1) * (y - y1) / (y2 - y1)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		r.HLineZBufNoArr(a, b, y, color, az, bz)
	}
}

func (r *Renderer) HLineZBufNoArr(x1, x2, y int32, color byte, z1, z2 float64) {
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
