package renderer

import (
	"totala_reader/ta_files_read/texture"
)

type onScreenTriangle struct {
	texture    *texture.GafEntry
	x0, x1, x2 int32
	y0, y1, y2 int32
	z0, z1, z2 float64
	u0, u1, u2 float64
	v0, v1, v2 float64
	color      byte
}

func (ost *onScreenTriangle) areCoordsClockwise() bool {
	x10, y10 := ost.x0-ost.x1, ost.y0-ost.y1
	x12, y12 := ost.x2-ost.x1, ost.y2-ost.y1
	// Return true if clockwise (determinant > 0) or collinear (determinant == 0)
	return x10*y12-x12*y10 >= 0
}

func (ost *onScreenTriangle) reorderCoordsYAsc() {
	if ost.y0 > ost.y1 {
		ost.x0, ost.x1 = ost.x1, ost.x0
		ost.y0, ost.y1 = ost.y1, ost.y0
		ost.z0, ost.z1 = ost.z1, ost.z0
		ost.u0, ost.u1 = ost.u1, ost.u0
		ost.v0, ost.v1 = ost.v1, ost.v0
	}
	if ost.y0 > ost.y2 {
		ost.x0, ost.x2 = ost.x2, ost.x0
		ost.y0, ost.y2 = ost.y2, ost.y0
		ost.z0, ost.z2 = ost.z2, ost.z0
		ost.u0, ost.u2 = ost.u2, ost.u0
		ost.v0, ost.v2 = ost.v2, ost.v0
	}
	if ost.y1 > ost.y2 {
		ost.x1, ost.x2 = ost.x2, ost.x1
		ost.y1, ost.y2 = ost.y2, ost.y1
		ost.z1, ost.z2 = ost.z2, ost.z1
		ost.u1, ost.u2 = ost.u2, ost.u1
		ost.v1, ost.v2 = ost.v2, ost.v1
	}
}

func (ost *onScreenTriangle) reorderCoordsYAscSkipUv() {
	if ost.y0 > ost.y1 {
		ost.x0, ost.x1 = ost.x1, ost.x0
		ost.y0, ost.y1 = ost.y1, ost.y0
		ost.z0, ost.z1 = ost.z1, ost.z0
	}
	if ost.y0 > ost.y2 {
		ost.x0, ost.x2 = ost.x2, ost.x0
		ost.y0, ost.y2 = ost.y2, ost.y0
		ost.z0, ost.z2 = ost.z2, ost.z0
	}
	if ost.y1 > ost.y2 {
		ost.x1, ost.x2 = ost.x2, ost.x1
		ost.y1, ost.y2 = ost.y2, ost.y1
		ost.z1, ost.z2 = ost.z2, ost.z1
	}
}
