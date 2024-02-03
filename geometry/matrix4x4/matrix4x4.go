package matrix4x4

type Matrix4x4 struct {
	Vals [4][4]float64
}

func NewUnitMatrix() *Matrix4x4 {
	return &Matrix4x4{
		Vals: [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func (m *Matrix4x4) ResetToUnitMatrix() {
	m.Vals = [4][4]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}
