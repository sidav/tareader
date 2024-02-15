package renderer

import (
	"totala_reader/ta_files_read/texture"
)

var xedge [1080][2]float64
var zedge [1080][2]float64

// For textures U is X, and V is Y.
var uedge [1080][2]float64
var vedge [1080][2]float64

func (r *Renderer) drawEdgebufTexturedTriangle(ost *onScreenTriangle) {
	ost.reorderCoordsYAsc()
	r.bufferEdgeTextured(ost.x0, ost.y0, ost.x1, ost.y1, ost.z0, ost.z1, ost.u0, ost.u1, ost.v0, ost.v1)
	r.bufferEdgeTextured(ost.x1, ost.y1, ost.x2, ost.y2, ost.z1, ost.z2, ost.u1, ost.u2, ost.v1, ost.v2)
	r.bufferEdgeTextured(ost.x2, ost.y2, ost.x0, ost.y0, ost.z2, ost.z0, ost.u2, ost.u0, ost.v2, ost.v0)
	if ost.y0 < 0 {
		ost.y0 = 0
	}
	if ost.y2 >= int32(len(xedge)) {
		ost.y2 = int32(len(xedge) - 1)
	}
	for y := ost.y0; y <= ost.y2; y++ {
		r.HLineTexturedZBuf(int32(xedge[y][0]), int32(xedge[y][1]), y, ost.texture)
	}
}

func (r *Renderer) bufferEdgeTextured(x1, y1, x2, y2 int32, z1, z2, u1, u2, v1, v2 float64) {
	side := 0
	if y1 >= y2 {
		side = 1
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
	}
	xslope := float64(x2-x1) / float64(y2-y1)
	zslope := (z2 - z1) / float64(y2-y1)
	uslope := (u2 - u1) / float64(y2-y1)
	vslope := (v2 - v1) / float64(y2-y1)
	if y1 == y2 {
		xslope = float64(x2 - x1)
		zslope = (z2 - z1)
		uslope = (u2 - u1)
		vslope = (v2 - v1)
	}

	currX := float64(x1)
	currZ := z1
	currU := u1
	currV := v1
	for y := y1; y <= y2; y++ {
		if y >= 0 && y < int32(len(xedge)) {
			xedge[y][side] = currX
			zedge[y][side] = currZ
			uedge[y][side] = currU
			vedge[y][side] = currV
		}
		currX += xslope
		currZ += zslope
		currU += uslope
		currV += vslope
	}
}

func (r *Renderer) HLineTexturedZBuf(x1, x2, y int32, texture *texture.GafEntry) {
	z1, z2 := zedge[y][0], zedge[y][1]
	u1, u2 := uedge[y][0], uedge[y][1]
	v1, v2 := vedge[y][0], vedge[y][1]
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
