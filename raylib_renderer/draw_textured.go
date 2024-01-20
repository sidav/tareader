package raylibrenderer

import (
	"math"
	"totala_reader/raylib_renderer/middleware"
	"totala_reader/ta_files_read/texture"
)

// For textures U is X, and V is Y.
var uedge [1080][2]float64
var vedge [1080][2]float64

func (r *RaylibRenderer) drawEdgeTextured(x1, y1, x2, y2 int32, z1, z2, u1, u2, v1, v2 float64) {
	side := 0
	xslope := float64(x2-x1) / float64(y2-y1)
	zslope := (z2 - z1) / float64(y2-y1)
	uslope := (u2 - u1) / float64(y2-y1)
	vslope := (v2 - v1) / float64(y2-y1)
	if y1 >= y2 {
		side = 1
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		z1, z2 = z2, z1
		u1, u2 = u2, u1
		v1, v2 = v2, v1
	}
	currX := float64(x1)
	currZ := z1
	currU := u1
	currV := v1
	for y := y1; y <= y2; y++ {
		if y >= 0 && y < int32(len(xedge)) {
			xedge[y][side] = currX
			zpos[y][side] = currZ
			uedge[y][side] = currU
			vedge[y][side] = currV
		}
		currX += xslope
		currZ += zslope
		currU += uslope
		currV += vslope
	}
}

func (r *RaylibRenderer) drawTexturedTriangle(x1, y1, x2, y2, x3, y3 int32, z1, z2, z3, u1, u2, u3, v1, v2, v3 float64, texture *texture.GafEntry) {
	var minx, maxx int32
	r.drawEdgeTextured(x1, y1, x2, y2, z1, z2, u1, u2, v1, v2)
	r.drawEdgeTextured(x2, y2, x3, y3, z2, z3, u2, u3, v2, v3)
	r.drawEdgeTextured(x3, y3, x1, y1, z3, z1, u3, u1, v3, v1)
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
		r.HLineTexturedZBuf(int32(xedge[y][0]), int32(xedge[y][1]), y, texture)
	}
}

func (r *RaylibRenderer) HLineTexturedZBuf(x1, x2, y int32, texture *texture.GafEntry) {
	z1, z2 := zpos[y][0], zpos[y][1]
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

	for x := x1; x <= x2; x++ {
		if r.canDrawOverZBufferAt(x, y, z1) {
			texW, texH := len(texture.Frames[0].Pixels), len(texture.Frames[0].Pixels[0])
			uCoord := int(math.Abs(math.Round(float64(texW)*u1))) % texW
			vCoord := int(math.Abs(math.Round(float64(texH)*v1))) % texH
			middleware.SetColor(getTaPaletteColor(texture.Frames[0].Pixels[uCoord][vCoord]))
			r.setZBufferValueAt(z1, x, y)
			middleware.DrawPoint(x, y)
		}
		z1 += zinc
		u1 += uinc
		v1 += vinc
	}
}
