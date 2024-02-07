package renderer

import (
	"totala_reader/ta_files_read/texture"
)

// ZBuffered horizontal scanline fill triangle algorithm.
func (r *Renderer) scanlineTexturedTriangle(x0, y0, x1, y1, x2, y2 int32, z0, z1, z2, u0, u1, u2, v0, v1, v2 float64, texture *texture.GafEntry) {
	var a, b, y, last int32
	var az, bz float64         // for z-buffering
	var au, bu, av, bv float64 // for texture mapping
	// Sort coordinates by Y order (y2 >= y1 >= y0)
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
		z0, z1 = z1, z0
		u0, u1 = u1, u0
		v0, v1 = v1, v0
	}
	if y1 > y2 {
		x2, y2, x1, y1 = x1, y1, x2, y2
		z2, z1 = z1, z2
		u2, u1 = u1, u2
		v2, v1 = v1, v2
	}
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
		z0, z1 = z1, z0
		u0, u1 = u1, u0
		v0, v1 = v1, v0
	}

	if y0 == y2 { // All on same line case
		a = x0
		b = x0
		az = z0
		bz = z0
		au = u0
		bu = u0
		av = v0
		bv = v0
		if x1 < a {
			a = x1
			az = z1
			au = u1
			av = v1
		} else if x1 > b {
			b = x1
			bz = z1
			bu = u1
			bv = v1
		}
		if x2 < a {
			a = x2
			az = z2
			au = u2
			av = v2
		} else if x2 > b {
			b = x2
			bz = z2
			bu = u2
			bv = v2
		}
		r.HLineTexturedZBufNoArr(a, b, y, az, bz, au, bu, av, bv, texture)
		return
	}

	// Diff between points 0 and 1
	dx01 := x1 - x0
	dy01 := y1 - y0
	dz01 := z1 - z0
	du01 := u1 - u0
	dv01 := v1 - v0
	// Diff between points 0 and 2
	dx02 := x2 - x0
	dy02 := y2 - y0
	dz02 := z2 - z0
	du02 := u2 - u0
	dv02 := v2 - v0
	// Diff between points 1 and 2
	dx12 := x2 - x1
	dy12 := y2 - y1
	dz12 := z2 - z1
	du12 := u2 - u1
	dv12 := v2 - v1
	// floats to reduce type conversions:
	dy01float := float64(dy01)
	dy02float := float64(dy02)
	dy12float := float64(dy12)

	var sa, sb int32
	var saz, sbz float64
	var sau, sbu, sav, sbv float64

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
		az = z0 + saz/dy01float
		bz = z0 + sbz/dy02float
		au = u0 + sau/dy01float
		bu = u0 + sbu/dy02float
		av = v0 + sav/dy01float
		bv = v0 + sbv/dy02float
		sa += dx01
		sb += dx02
		saz += dz01
		sbz += dz02
		sau += du01
		sbu += du02
		sav += dv01
		sbv += dv02
		// longhand a = x0 + (x1 - x0) * (y - y0) / (y1 - y0)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		r.HLineTexturedZBufNoArr(a, b, y, az, bz, au, bu, av, bv, texture)
	}

	// For lower part of triangle, find scanline crossings for segment
	// 0-2 and 1-2.  This loop is skipped if y1=y2
	sa = dx12 * (y - y1)
	sb = dx02 * (y - y0)
	saz = dz12 * float64(y-y1)
	sbz = dz02 * float64(y-y0)
	sau = du12 * float64(y-y1)
	sbu = du02 * float64(y-y0)
	sav = dv12 * float64(y-y1)
	sbv = dv02 * float64(y-y0)

	for ; y <= y2; y++ {
		a = x1 + sa/dy12
		b = x0 + sb/dy02
		az = z1 + saz/dy12float
		bz = z0 + sbz/dy02float
		au = u1 + sau/dy12float
		bu = u0 + sbu/dy02float
		av = v1 + sav/dy12float
		bv = v0 + sbv/dy02float
		sa += dx12
		sb += dx02
		saz += dz12
		sbz += dz02
		sau += du12
		sbu += du02
		sav += dv12
		sbv += dv02
		// longhand a = x1 + (x2 - x1) * (y - y1) / (y2 - y1)
		//          b = x0 + (x2 - x0) * (y - y0) / (y2 - y0)
		r.HLineTexturedZBufNoArr(a, b, y, az, bz, au, bu, av, bv, texture)
	}
}

func (r *Renderer) HLineTexturedZBufNoArr(x1, x2, y int32, z1, z2, u1, u2, v1, v2 float64, texture *texture.GafEntry) {
	if x1 > x2 {
		x1, x2 = x2, x1
		z1, z2 = z2, z1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
	}
	zinc := (z2 - z1) / float64(x2-x1)
	uinc := (u2 - u1) / float64(x2-x1)
	vinc := (v2 - v1) / float64(x2-x1)

	// Real texture coord for max U and V.
	// (-0.5) here because it's -1 (as max coord can't be equal to size) added with +0.5 (for texture subpixel alignment)
	maxUReal := float64(len(texture.Frames[0].Pixels)) - 0.5
	maxVReal := float64(len(texture.Frames[0].Pixels[0])) - 0.5
	for x := x1; x <= x2; x++ {
		if r.canDrawOverZBufferAt(x, y, z1) {
			uCoord := int(maxUReal * u1)
			vCoord := int(maxVReal * v1)
			r.gAdapter.SetColor(getTaPaletteColor(texture.Frames[0].Pixels[uCoord][vCoord]))
			r.setZBufferValueAt(z1, x, y)
			r.gAdapter.DrawPoint(x, y)
		}
		z1 += zinc
		u1 += uinc
		v1 += vinc
	}
}
