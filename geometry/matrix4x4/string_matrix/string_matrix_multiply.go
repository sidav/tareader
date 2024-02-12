package stringmatrix

import "fmt"

func usageExample() {
	MatrixStringsMultiply(
		[][]string{
			{"cos", "0", "-sin", "0"},
			{"0", "1", "0", "0"},
			{"sin", "0", "cos", "0"},
			{"0", "0", "0", "1"},
		},
		[][]string{
			{"m1[0][0]", "m1[0][1]", "m1[0][2]", "m1[0][3]"},
			{"m1[1][0]", "m1[1][1]", "m1[1][2]", "m1[1][3]"},
			{"m1[2][0]", "m1[2][1]", "m1[2][2]", "m1[2][3]"},
			{"m1[3][0]", "m1[3][1]", "m1[3][2]", "m1[3][3]"},
		},
	)
}

// I'm SO TIRED of writing formulas myself, so...
// This func shows and returns matrix multiplication result in strings (formula), simplifying the result for 0s and 1s
// It should ease code optimizations for matrices
func MatrixStringsMultiply(m1, m2 [][]string) [][]string {
	result := make([][]string, len(m1))
	for r := range result {
		result[r] = make([]string, len(m2[0]))
	}

	for rowInResult := range m1 {
		for columnInM2 := range m2[0] {
			currString := ""
			for rowInM2 := range m2 {
				m1val := m1[rowInResult][rowInM2]
				m2val := m2[rowInM2][columnInM2]
				// optimization: remove multiply by zero
				if m1val == "0" || m2val == "0" {
					continue
				}
				if len(currString) > 0 {
					currString += " + "
				}
				// optimization: remove "* 1" parts
				if m1val == "1" {
					currString += m2val
				} else if m2val == "1" {
					currString += m1val
				} else {
					currString += fmt.Sprintf("%s*%s", m1val, m2val)
				}
			}
			result[rowInResult][columnInM2] = currString
		}
	}
	printStringMatrix(result)
	return result
}

func printStringMatrix(m [][]string) {
	maxLength := 0
	for i := range m {
		for j := range m[i] {
			if len(m[i][j]) > maxLength {
				maxLength = len(m[i][j])
			}
		}
	}
	formatStr := "%" + fmt.Sprintf("%ds", maxLength+2)

	for i := range m {
		fmt.Print("[")
		for j := range m[i] {
			fmt.Printf(formatStr, m[i][j])
			if j < len(m[i])-1 {
				fmt.Print(",  ")
			}
		}
		fmt.Println("]")
	}
}
