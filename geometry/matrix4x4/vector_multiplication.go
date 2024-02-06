package matrix4x4

import (
	"totala_reader/geometry"
)

// TODO: maybe use pointers?..
func (m *Matrix4x4) MultiplyByVector(v geometry.Vector3) geometry.Vector3 {
	return geometry.Vector3{
		X: m.Vals[0][0]*v.X + m.Vals[0][1]*v.Y + m.Vals[0][2]*v.Z + m.Vals[0][3],
		Y: m.Vals[1][0]*v.X + m.Vals[1][1]*v.Y + m.Vals[1][2]*v.Z + m.Vals[1][3],
		Z: m.Vals[2][0]*v.X + m.Vals[2][1]*v.Y + m.Vals[2][2]*v.Z + m.Vals[2][3],
		// the next value is simply unused
	}
}

// Multiplies the matrix by (arr[0], arr[1], arr[2], 1) column-vector
// Useful for multiplying by VERTICES (as positions in space; so the translations will be applied)
func (m *Matrix4x4) MultiplyByArr3Vector(arr [3]float64) [3]float64 {
	return [3]float64{
		m.Vals[0][0]*arr[0] + m.Vals[0][1]*arr[1] + m.Vals[0][2]*arr[2] + m.Vals[0][3],
		m.Vals[1][0]*arr[0] + m.Vals[1][1]*arr[1] + m.Vals[1][2]*arr[2] + m.Vals[1][3],
		m.Vals[2][0]*arr[0] + m.Vals[2][1]*arr[1] + m.Vals[2][2]*arr[2] + m.Vals[2][3],
	}
}

// Multiplies the matrix by (x, y, z, 1) column-vector
// Useful for multiplying by VERTICES (as positions in space; so the translations will be applied)
func (m *Matrix4x4) MultiplyByXyz1Vector(x, y, z float64) (float64, float64, float64) {
	return m.Vals[0][0]*x + m.Vals[0][1]*y + m.Vals[0][2]*z + m.Vals[0][3],
		m.Vals[1][0]*x + m.Vals[1][1]*y + m.Vals[1][2]*z + m.Vals[1][3],
		m.Vals[2][0]*x + m.Vals[2][1]*y + m.Vals[2][2]*z + m.Vals[2][3]
}

// Multiplies the matrix by (x, y, z, 0) column-vector
// Useful for multiplying by VECTORS (as directions from origin; so the translations won't be applied)
func (m *Matrix4x4) MultiplyByXyz0Vector(x, y, z float64) (float64, float64, float64) {
	return m.Vals[0][0]*x + m.Vals[0][1]*y + m.Vals[0][2]*z,
		m.Vals[1][0]*x + m.Vals[1][1]*y + m.Vals[1][2]*z,
		m.Vals[2][0]*x + m.Vals[2][1]*y + m.Vals[2][2]*z
}
