package model

const SimTicksPerSecond = 60
const pi = 3.1415926535897932
const pi2 = 2 * pi

func CavedogAngleSpeedToFloatRadians(cavedogAngle int32) float64 {
	const factor = pi2 / 65536.0 / SimTicksPerSecond
	return float64(cavedogAngle) * factor
}

func CavedogAngleToFloatRadians(cavedogAngle int32) float64 {
	const factor = pi2 / 65536.0
	return float64(cavedogAngle) * factor
}

func CavedogPositionToFloatPosition(cavedogDistance int32) float64 {
	const factor = 1 / 65536.0
	return float64(cavedogDistance) * factor
}

func CavedogPositionSpeedToFloat(cavedogSpeed int32) float64 {
	return CavedogPositionToFloatPosition(cavedogSpeed) / SimTicksPerSecond
}
