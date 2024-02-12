package renderer

import "math"

// The camera in TA works by these coordinates (looked at upside down):
// +------------> X
// |
// |
// |     Y points up at the camera
// |
// |
// V Z

// oblique/orthographic projection
func obliqueProjection(x, y, z float64) (osx, osy float64) {
	return x, z - y/2
}

func obliqueProjectionInt32(x, y, z float64) (osx, osy int32) {
	return int32(math.Round(x)), int32(math.Round(-z - y/2))
}
