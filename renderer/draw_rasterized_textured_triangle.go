package renderer

import "totala_reader/ta_files_read/texture"

func (r *Renderer) drawRasterizedTexturedTriangle(x0, y0, x1, y1, x2, y2 int32, z0, z1, z2, u0, u1, u2, v0, v1, v2 float64, texture *texture.GafEntry) {
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
		z0, z1 = z1, z0
		u0, u1 = u1, u0
		v0, v1 = v1, v0
	}
	if y0 > y2 {
		x0, x2 = x2, x0
		y0, y2 = y2, y0
		z0, z2 = z2, z0
		u0, u2 = u2, u0
		v0, v2 = v2, v0
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
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
	du1 := (u1 - u0) / dy1f
	du2 := (u2 - u0) / dy2f
	dv1 := (v1 - v0) / dy1f
	dv2 := (v2 - v0) / dy2f

	currz1 := z0
	currz2 := z0
	curru1 := u0
	curru2 := u0
	currv1 := v0
	currv2 := v0

	curry := y0
	for curry < y1 {
		currx1 = x0 + dx1*(curry-y0)/dy1
		currx2 = x0 + dx2*(curry-y0)/dy2
		r.HLineTexturedZBufRstr(currx1, currx2, curry, currz1, currz2, curru1, curru2, currv1, currv2, texture)
		currz1 += dz1
		currz2 += dz2
		curru1 += du1
		curru2 += du2
		currv1 += dv1
		currv2 += dv2
		curry++
	}

	dx1 = x2 - x1
	dy1 = y2 - y1

	dy1f = float64(dy1)
	dz1 = (z2 - z1) / dy1f
	du1 = (u2 - u1) / dy1f
	dv1 = (v2 - v1) / dy1f
	currz1 = z1
	curru1 = u1
	currv1 = v1

	for curry < y2 {
		currx1 = x1 + dx1*(curry-y1)/dy1
		currx2 = x0 + dx2*(curry-y0)/dy2

		r.HLineTexturedZBufRstr(currx1, currx2, curry, currz1, currz2, curru1, curru2, currv1, currv2, texture)
		currz1 += dz1
		currz2 += dz2
		curru1 += du1
		curru2 += du2
		currv1 += dv1
		currv2 += dv2
		curry++
	}
}

func (r *Renderer) HLineTexturedZBufRstr(x1, x2, y int32, z1, z2, u1, u2, v1, v2 float64, texture *texture.GafEntry) {
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
