package raylibrenderer

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
