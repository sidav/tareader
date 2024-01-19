package geometry

import "math"

func Rotate3dCoordsAroundX(x, y, z, radians float64) (rx, ry, rz float64) {
	return x,
		y*math.Cos(radians) - z*math.Sin(radians),
		y*math.Sin(radians) + z*math.Cos(radians)
}

func Rotate3dCoordsAroundY(x, y, z, radians float64) (rx, ry, rz float64) {
	return x*math.Cos(radians) - z*math.Sin(radians),
		y,
		x*math.Sin(radians) + z*math.Cos(radians)
}

func Rotate3dCoordsAroundZ(x, y, z, radians float64) (rx, ry, rz float64) {
	return x*math.Cos(radians) - y*math.Sin(radians),
		x*math.Sin(radians) + y*math.Cos(radians),
		z
}
