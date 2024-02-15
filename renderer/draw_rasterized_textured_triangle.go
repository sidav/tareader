package renderer

import "totala_reader/ta_files_read/texture"

func (r *Renderer) drawRasterizedTexturedTriangle(ost *onScreenTriangle) {
	ost.reorderCoordsYAsc()

	var currx1, currx2 int32
	dx1 := ost.x1 - ost.x0
	dy1 := ost.y1 - ost.y0
	dx2 := ost.x2 - ost.x0
	dy2 := ost.y2 - ost.y0

	dy1f := float64(dy1)
	dy2f := float64(dy2)
	dz1 := (ost.z1 - ost.z0) / dy1f
	dz2 := (ost.z2 - ost.z0) / dy2f
	du1 := (ost.u1 - ost.u0) / dy1f
	du2 := (ost.u2 - ost.u0) / dy2f
	dv1 := (ost.v1 - ost.v0) / dy1f
	dv2 := (ost.v2 - ost.v0) / dy2f

	currz1 := ost.z0
	currz2 := ost.z0
	curru1 := ost.u0
	curru2 := ost.u0
	currv1 := ost.v0
	currv2 := ost.v0

	curry := ost.y0
	for curry < ost.y1 {
		currx1 = ost.x0 + dx1*(curry-ost.y0)/dy1
		currx2 = ost.x0 + dx2*(curry-ost.y0)/dy2
		r.HLineTexturedZBufRstr(currx1, currx2, curry, currz1, currz2, curru1, curru2, currv1, currv2, ost.texture)
		currz1 += dz1
		currz2 += dz2
		curru1 += du1
		curru2 += du2
		currv1 += dv1
		currv2 += dv2
		curry++
	}

	dx1 = ost.x2 - ost.x1
	dy1 = ost.y2 - ost.y1

	dy1f = float64(dy1)
	dz1 = (ost.z2 - ost.z1) / dy1f
	du1 = (ost.u2 - ost.u1) / dy1f
	dv1 = (ost.v2 - ost.v1) / dy1f
	currz1 = ost.z1
	curru1 = ost.u1
	currv1 = ost.v1

	for curry < ost.y2 {
		currx1 = ost.x1 + dx1*(curry-ost.y1)/dy1
		currx2 = ost.x0 + dx2*(curry-ost.y0)/dy2

		r.HLineTexturedZBufRstr(currx1, currx2, curry, currz1, currz2, curru1, curru2, currv1, currv2, ost.texture)
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
