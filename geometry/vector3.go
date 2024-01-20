package geometry

import "math"

type Vector3 struct {
	X, Y, Z float64
}

func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func NewVector3FromArr(arr [3]float64) Vector3 {
	return Vector3{X: arr[0], Y: arr[1], Z: arr[2]}
}

func (a *Vector3) Normalize() {
	length := math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	a.X /= length
	a.Y /= length
	a.Z /= length
}

func CrossProduct(a, b *Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func DotProduct(a, b *Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}
