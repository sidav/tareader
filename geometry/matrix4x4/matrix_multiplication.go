package matrix4x4

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
