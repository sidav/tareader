package object3d

// converts 4-byte signed int32 (written to int) to Float
// 2-byte integer part, 2-byte fractional part
func FixedPointToFloat(x int) float64 {
	return float64(x) / 65536.0
}
