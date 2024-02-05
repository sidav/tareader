package matrix4x4

import (
	"math"
	"totala_reader/geometry"
)

// TODO: maybe use pointers?..
func (m *Matrix4x4) MultiplyByVector(v geometry.Vector3) geometry.Vector3 {
	// Matrix is 4x4, and the vector is 1x3. So, we consider the (1,4)th value as 1. THIS MAY BE A MISTAKE.
	return geometry.Vector3{
		X: m.Vals[0][0]*v.X + m.Vals[0][1]*v.Y + m.Vals[0][2]*v.Z + m.Vals[0][3],
		Y: m.Vals[1][0]*v.X + m.Vals[1][1]*v.Y + m.Vals[1][2]*v.Z + m.Vals[1][3],
		Z: m.Vals[2][0]*v.X + m.Vals[2][1]*v.Y + m.Vals[2][2]*v.Z + m.Vals[2][3],
		// the next value is simply unused
	}
}

func (m *Matrix4x4) MultiplyByArr3Vector(arr [3]float64) [3]float64 {
	// Matrix is 4x4, and the vector is 1x3. So, we consider the (1,4)th value as 1. THIS MAY BE A MISTAKE.
	return [3]float64{
		m.Vals[0][0]*arr[0] + m.Vals[0][1]*arr[1] + m.Vals[0][2]*arr[2] + m.Vals[0][3],
		m.Vals[1][0]*arr[0] + m.Vals[1][1]*arr[1] + m.Vals[1][2]*arr[2] + m.Vals[1][3],
		m.Vals[2][0]*arr[0] + m.Vals[2][1]*arr[1] + m.Vals[2][2]*arr[2] + m.Vals[2][3],
	}
}

func (m *Matrix4x4) MultiplyByXyzVector(x, y, z float64) (float64, float64, float64) {
	// Matrix is 4x4, and the vector is 1x3. So, we consider the (1,4)th value as 1. THIS MAY BE A MISTAKE.
	return m.Vals[0][0]*x + m.Vals[0][1]*y + m.Vals[0][2]*z + m.Vals[0][3],
		m.Vals[1][0]*x + m.Vals[1][1]*y + m.Vals[1][2]*z + m.Vals[1][3],
		m.Vals[2][0]*x + m.Vals[2][1]*y + m.Vals[2][2]*z + m.Vals[2][3]
}

func (m1 *Matrix4x4) MultiplyByMatrix4x4(m2 *Matrix4x4) *Matrix4x4 {
	// matrix[line_number][column_number]
	return &Matrix4x4{
		Vals: [4][4]float64{
			// row 1
			{
				m1.Vals[0][0]*m2.Vals[0][0] + m1.Vals[0][1]*m2.Vals[1][0] + m1.Vals[0][2]*m2.Vals[2][0] + m1.Vals[0][3]*m2.Vals[3][0],
				m1.Vals[0][0]*m2.Vals[0][1] + m1.Vals[0][1]*m2.Vals[1][1] + m1.Vals[0][2]*m2.Vals[2][1] + m1.Vals[0][3]*m2.Vals[3][1],
				m1.Vals[0][0]*m2.Vals[0][2] + m1.Vals[0][1]*m2.Vals[1][2] + m1.Vals[0][2]*m2.Vals[2][2] + m1.Vals[0][3]*m2.Vals[3][2],
				m1.Vals[0][0]*m2.Vals[0][3] + m1.Vals[0][1]*m2.Vals[1][3] + m1.Vals[0][2]*m2.Vals[2][3] + m1.Vals[0][3]*m2.Vals[3][3],
			},
			// row 2
			{
				m1.Vals[1][0]*m2.Vals[0][0] + m1.Vals[1][1]*m2.Vals[1][0] + m1.Vals[1][2]*m2.Vals[2][0] + m1.Vals[1][3]*m2.Vals[3][0],
				m1.Vals[1][0]*m2.Vals[0][1] + m1.Vals[1][1]*m2.Vals[1][1] + m1.Vals[1][2]*m2.Vals[2][1] + m1.Vals[1][3]*m2.Vals[3][1],
				m1.Vals[1][0]*m2.Vals[0][2] + m1.Vals[1][1]*m2.Vals[1][2] + m1.Vals[1][2]*m2.Vals[2][2] + m1.Vals[1][3]*m2.Vals[3][2],
				m1.Vals[1][0]*m2.Vals[0][3] + m1.Vals[1][1]*m2.Vals[1][3] + m1.Vals[1][2]*m2.Vals[2][3] + m1.Vals[1][3]*m2.Vals[3][3],
			},
			// row 3
			{
				m1.Vals[2][0]*m2.Vals[0][0] + m1.Vals[2][1]*m2.Vals[1][0] + m1.Vals[2][2]*m2.Vals[2][0] + m1.Vals[2][3]*m2.Vals[3][0],
				m1.Vals[2][0]*m2.Vals[0][1] + m1.Vals[2][1]*m2.Vals[1][1] + m1.Vals[2][2]*m2.Vals[2][1] + m1.Vals[2][3]*m2.Vals[3][1],
				m1.Vals[2][0]*m2.Vals[0][2] + m1.Vals[2][1]*m2.Vals[1][2] + m1.Vals[2][2]*m2.Vals[2][2] + m1.Vals[2][3]*m2.Vals[3][2],
				m1.Vals[2][0]*m2.Vals[0][3] + m1.Vals[2][1]*m2.Vals[1][3] + m1.Vals[2][2]*m2.Vals[2][3] + m1.Vals[2][3]*m2.Vals[3][3],
			},
			// row 4
			{
				m1.Vals[3][0]*m2.Vals[0][0] + m1.Vals[3][1]*m2.Vals[1][0] + m1.Vals[3][2]*m2.Vals[2][0] + m1.Vals[3][3]*m2.Vals[3][0],
				m1.Vals[3][0]*m2.Vals[0][1] + m1.Vals[3][1]*m2.Vals[1][1] + m1.Vals[3][2]*m2.Vals[2][1] + m1.Vals[3][3]*m2.Vals[3][1],
				m1.Vals[3][0]*m2.Vals[0][2] + m1.Vals[3][1]*m2.Vals[1][2] + m1.Vals[3][2]*m2.Vals[2][2] + m1.Vals[3][3]*m2.Vals[3][2],
				m1.Vals[3][0]*m2.Vals[0][3] + m1.Vals[3][1]*m2.Vals[1][3] + m1.Vals[3][2]*m2.Vals[2][3] + m1.Vals[3][3]*m2.Vals[3][3],
			},
		},
	}
}

func (m1 *Matrix4x4) GetRotatedAroundX(radians float64) *Matrix4x4 {
	// TODO: rewrite without returning value; change the m1 in place
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	return &Matrix4x4{
		Vals: [4][4]float64{
			// row 1
			{
				m1.Vals[0][0],
				m1.Vals[0][1]*cos + m1.Vals[0][2]*(-sin),
				m1.Vals[0][1]*sin + m1.Vals[0][2]*cos,
				m1.Vals[0][3],
			},
			// row 2
			{
				m1.Vals[1][0],
				m1.Vals[1][1]*cos + m1.Vals[1][2]*(-sin),
				m1.Vals[1][1]*sin + m1.Vals[1][2]*cos,
				m1.Vals[1][3],
			},
			// row 3
			{
				m1.Vals[2][0],
				m1.Vals[2][1]*cos + m1.Vals[2][2]*(-sin),
				m1.Vals[2][1]*sin + m1.Vals[2][2]*cos,
				m1.Vals[2][3],
			},
			// row 4
			{
				m1.Vals[3][0],
				m1.Vals[3][1]*cos + m1.Vals[3][2]*(-sin),
				m1.Vals[3][1]*sin + m1.Vals[3][2]*cos,
				m1.Vals[3][3],
			},
		},
	}
}

func (m1 *Matrix4x4) GetRotatedAroundY(radians float64) *Matrix4x4 {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	// matrix[line_number][column_number]
	return &Matrix4x4{
		Vals: [4][4]float64{
			// row 1
			{
				m1.Vals[0][0]*cos + m1.Vals[0][2]*sin,
				m1.Vals[0][1],
				m1.Vals[0][0]*(-sin) + m1.Vals[0][2]*cos,
				m1.Vals[0][3],
			},
			// row 2
			{
				m1.Vals[1][0]*cos + m1.Vals[1][2]*sin,
				m1.Vals[1][1],
				m1.Vals[1][0]*(-sin) + m1.Vals[1][2]*cos,
				m1.Vals[1][3],
			},
			// row 3
			{
				m1.Vals[2][0]*cos + m1.Vals[2][2]*sin,
				m1.Vals[2][1],
				m1.Vals[2][0]*(-sin) + m1.Vals[2][2]*cos,
				m1.Vals[2][3],
			},
			// row 4
			{
				m1.Vals[3][0]*cos + m1.Vals[3][2]*sin,
				m1.Vals[3][1],
				m1.Vals[3][0]*(-sin) + m1.Vals[3][2]*cos,
				m1.Vals[3][3],
			},
		},
	}
}

func (m1 *Matrix4x4) GetRotatedAroundZ(radians float64) *Matrix4x4 {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	// matrix[line_number][column_number]
	return &Matrix4x4{
		Vals: [4][4]float64{
			// row 1
			{
				m1.Vals[0][0]*cos + m1.Vals[0][1]*sin,
				m1.Vals[0][0]*(-sin) + m1.Vals[0][1]*cos,
				m1.Vals[0][2],
				m1.Vals[0][3],
			},
			// row 2
			{
				m1.Vals[1][0]*cos + m1.Vals[1][1]*sin,
				m1.Vals[1][0]*(-sin) + m1.Vals[1][1]*cos,
				m1.Vals[1][2],
				m1.Vals[1][3],
			},
			// row 3
			{
				m1.Vals[2][0]*cos + m1.Vals[2][1]*sin,
				m1.Vals[2][0]*(-sin) + m1.Vals[2][1]*cos,
				m1.Vals[2][2],
				m1.Vals[2][3],
			},
			// row 4
			{
				m1.Vals[3][0]*cos + m1.Vals[3][1]*sin,
				m1.Vals[3][0]*(-sin) + m1.Vals[3][1]*cos,
				m1.Vals[3][2],
				m1.Vals[3][3],
			},
		},
	}
}
