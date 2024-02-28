package model

const SimTicksPerSecond = 60
const pi = 3.1415926535897932
const pi2 = 2 * pi

func CavedogAngleSpeedToFloatRadians(cavedogAngle int32) float64 {
	const factor = pi / 65536.0 / SimTicksPerSecond
	return float64(cavedogAngle) * factor
}

func CavedogAngleToFloatRadians(cavedogAngle int32) float64 {
	const factor = pi / 65536.0
	return float64(cavedogAngle) * factor
}

func CobWordToFloatPosition(cavedogDistance int32) float64 {
	// return float64(cavedogDistance) / 65536.0
	return float64(cavedogDistance) / 163840 // I dunno where does this magic number come from
	// TODO: investigate signed/unsigned case?
}
