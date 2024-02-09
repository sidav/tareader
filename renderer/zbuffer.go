package renderer

import "math"

func (r *Renderer) canDrawOverZBufferAt(x, y int32, depth float64) bool {
	if x < 0 || x >= int32(len(r.zBuffer)) || y < 0 || y >= int32(len(r.zBuffer[0])) {
		return false
	}

	// There was <= instead of <, but this variant rendered obsolete since back-face culling was implemented.
	// (Back-face culling solves see-through texturing of two-side primitives better than this zbuffer stuff did)
	return r.zBuffer[x][y] < depth
}

func (r *Renderer) setZBufferValueAt(val float64, x, y int32) {
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

func (r *Renderer) initZBuffer() {
	w, h := r.gAdapter.GetRenderResolution()
	r.zBuffer = make([][]float64, w)
	for i := range r.zBuffer {
		r.zBuffer[i] = make([]float64, h)
	}
	r.zBufMinX, r.zBufMaxX = int32(len(r.zBuffer)), 0
	r.zBufMinY, r.zBufMaxY = int32(len(r.zBuffer[0])), 0
	for i := range r.zBuffer {
		for j := range r.zBuffer[i] {
			r.zBuffer[i][j] = -math.MaxFloat64
		}
	}
}

func (r *Renderer) clearZBuffer() {
	for i := r.zBufMinX; i <= r.zBufMaxX; i++ {
		for j := r.zBufMinY; j <= r.zBufMaxY; j++ {
			r.zBuffer[i][j] = -math.MaxFloat64
		}
	}
	r.zBufMinX, r.zBufMaxX = int32(len(r.zBuffer)), 0
	r.zBufMinY, r.zBufMaxY = int32(len(r.zBuffer[0])), 0
}
