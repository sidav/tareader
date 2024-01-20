package raylibrenderer

import "math"

func (r *RaylibRenderer) canDrawOverZBufferAt(x, y int32, depth float64) bool {
	return r.zBuffer[x][y] < depth
}

func (r *RaylibRenderer) setZBufferValueAt(val float64, x, y int32) {
	r.zBuffer[x][y] = val
	if x < r.zBufMinX {
		r.zBufMinX = x
	}
	if x > r.zBufMaxX {
		r.zBufMaxX = x
	}
	if y < r.zBufMinY {
		r.zBufMinY = y
	}
	if y > r.zBufMaxY {
		r.zBufMaxY = y
	}
}

func (r *RaylibRenderer) initZBuffer() {
	r.zBufMinX, r.zBufMaxX = int32(len(r.zBuffer)), 0
	r.zBufMinY, r.zBufMaxY = int32(len(r.zBuffer[0])), 0
	for i := range r.zBuffer {
		for j := range r.zBuffer[i] {
			r.zBuffer[i][j] = -math.MaxFloat64
		}
	}
}

func (r *RaylibRenderer) clearZBuffer() {
	for i := r.zBufMinX; i <= r.zBufMaxX; i++ {
		for j := r.zBufMinY; j <= r.zBufMaxY; j++ {
			r.zBuffer[i][j] = -math.MaxFloat64
		}
	}
	r.zBufMinX, r.zBufMaxX = int32(len(r.zBuffer)), 0
	r.zBufMinY, r.zBufMaxY = int32(len(r.zBuffer[0])), 0
}
