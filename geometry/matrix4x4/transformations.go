package matrix4x4

import "math"

func (m1 *Matrix4x4) RotateAroundX(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	tm1, tm2 := m1.Vals[0][1], m1.Vals[0][2]
	m1.Vals[0][1] = tm1*cos + tm2*(-sin)
	m1.Vals[0][2] = tm1*sin + tm2*cos
	tm1, tm2 = m1.Vals[1][1], m1.Vals[1][2]
	m1.Vals[1][1] = tm1*cos + tm2*(-sin)
	m1.Vals[1][2] = tm1*sin + tm2*cos
	tm1, tm2 = m1.Vals[2][1], m1.Vals[2][2]
	m1.Vals[2][1] = tm1*cos + tm2*(-sin)
	m1.Vals[2][2] = tm1*sin + tm2*cos
	tm1, tm2 = m1.Vals[3][1], m1.Vals[3][2]
	m1.Vals[3][1] = tm1*cos + tm2*(-sin)
	m1.Vals[3][2] = tm1*sin + tm2*cos
}

func (m1 *Matrix4x4) RotateAroundY(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	tm1, tm2 := m1.Vals[0][0], m1.Vals[0][2]
	m1.Vals[0][0] = tm1*cos + tm2*(sin)
	m1.Vals[0][2] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[1][0], m1.Vals[1][2]
	m1.Vals[1][0] = tm1*cos + tm2*(sin)
	m1.Vals[1][2] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[2][0], m1.Vals[2][2]
	m1.Vals[2][0] = tm1*cos + tm2*(sin)
	m1.Vals[2][2] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[3][0], m1.Vals[3][2]
	m1.Vals[3][0] = tm1*cos + tm2*(sin)
	m1.Vals[3][2] = tm1*(-sin) + tm2*cos
}

func (m1 *Matrix4x4) RotateAroundZ(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	tm1, tm2 := m1.Vals[0][0], m1.Vals[0][1]
	m1.Vals[0][0] = tm1*cos + tm2*(sin)
	m1.Vals[0][1] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[1][0], m1.Vals[1][1]
	m1.Vals[1][0] = tm1*cos + tm2*(sin)
	m1.Vals[1][1] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[2][0], m1.Vals[2][1]
	m1.Vals[2][0] = tm1*cos + tm2*(sin)
	m1.Vals[2][1] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m1.Vals[3][0], m1.Vals[3][1]
	m1.Vals[3][0] = tm1*cos + tm2*(sin)
	m1.Vals[3][1] = tm1*(-sin) + tm2*cos
}
