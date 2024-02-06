package matrix4x4

import (
	"fmt"
	"strconv"
)

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

func (m *Matrix4x4) Dup() *Matrix4x4 {
	m2 := &Matrix4x4{}
	copy(m2.Vals[:], m.Vals[:])
	return m2
}

func (m *Matrix4x4) Print() {
	for row := range m.Vals {
		fmt.Print("|")
		for column := range m.Vals[row] {
			fmt.Printf("%8s ", strconv.FormatFloat(m.Vals[row][column], 'f', 2, 64))
		}
		fmt.Printf("|\n")
	}
}
