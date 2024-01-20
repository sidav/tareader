package raylibrenderer

func (r *RaylibRenderer) canDrawOverZBufferAt(x, y int32, depth float64) bool {
	if r.zBufferReverse {
		return -r.zBuffer[x][y] < depth
	}
	return r.zBuffer[x][y] < depth
}

func (r *RaylibRenderer) setZBufferValueAt(val float64, x, y int32) {
	if r.zBufferReverse {
		val = -val
	}
	r.zBuffer[x][y] = val
}

func (r *RaylibRenderer) flipZBuffer() {
	r.zBufferReverse = !r.zBufferReverse
}
