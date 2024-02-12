package geometry

import (
	"fmt"
	"math"
	"strconv"
)

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

func (a *Vector3) Print() {
	fmt.Printf("(%s  %5s  %5s)\n",
		strconv.FormatFloat(a.X, 'f', 2, 64),
		strconv.FormatFloat(a.Y, 'f', 2, 64),
		strconv.FormatFloat(a.Z, 'f', 2, 64),
	)
}

func (v *Vector3) Sub(v2 Vector3) {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
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
