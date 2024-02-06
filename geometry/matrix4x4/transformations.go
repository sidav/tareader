package matrix4x4

import "math"

/*
Important stuff which made me spend a lot of time debugging:
The matrix being transformed is on the RIGHT here.
if transform matrix is R, and the matrix is A, transformation apply is R*A and !!NOT A*R!!
We need to premultiply transformation matrix by A (as in R*A) because we work with column-vectors v -
as we multiply A by v, not v by A. So, the vector with all the R1, R2... transformations will be:
v_transformed = Rn * ... * R2 * R1 * A * v      <- notice the order consistency
Any other order will break transformations.
This code results in A_transformed = R3 * R2 * R1 * A.


RENDERING PIPELINE ORDER (IMPORTANT! This knowledge costed me a whole weekend of trial and error!)
The rendering with parent-child model structure is recursive (pseudo-code):
	func RenderObject(object, parentWorldMatrix)
		// If there is no parent, parentWorldMatrix is simply an unit matrix.
		currentObjectLocalMatrix = object.localMatrix
		currentObjectWorldMatrix = describedSteps(parentWorldMatrix) // see below
		DrawEverythingForObject(object, currentObjectWorldMatrix)
		RenderObject(object.child, currentObjectWorldMatrix) // notice that we pass WORLD matrix here

At each step, we have current object local matrix COM and parent matrix (after all transformations): PTM.
Local object has an offset from the parent: vector3 (XFromParent, YFromParent, ZFromParent).
describedSteps above are used to transform current object matrix to world space matrix,
while taking all parent(s) transforms (rotations etc) into account.

	1. Transform/rotate the current offset by multiplying the parent matrix by the offset.
	IMPORTANT: offset shouldn't be translated, only scaled and rotated, so it is multiplied as (x, y, z, 0) vector
	(the last 0 is important to skip translations here).
		ox, oy, oz := parentMtrx.MultiplyByXyz0Vector(XFromParent, YFromParent, ZFromParent)
	2. Transforming the current local matrix by multiplying the parent matrix by the current local matrix
	(this step applies PARENT rotations/translations to the current object)
		transformed = parentMtrx.MultiplyByMatrix4x4(currentMatrix)
	3. Translating the transforming matrix by "parent-transformed" current object offset.
	We do it on step 3 and not before because this offset is LOCAL and related to the parent,
	so it should not be affected by CURRENT object local transformations (it's affected by parent ones only).
		currentObjectWorldMatrix = currentMatrix.Translate(ox, oy, oz)

this currentObjectWorldMatrix should be later used to transform all local "current" object vertices for rendering.
The transformation is performed via multiplying the currentObjectWorldMatrix by vertex coords.
Coords SHOULD be affected by translation, so matrix should be multiplied by (vx, vy, vz, 1), notice the 1 at the end.
*/

// All funcs in this file modify the matrix!
// All multiplications here are left-to-right with the matrix being transformed on the right! (TRANSFORM * CURRMATRIX)
func (m *Matrix4x4) Scale(factor float64) {
	m.Vals[0][0] *= factor
	m.Vals[0][1] *= factor
	m.Vals[0][2] *= factor
	m.Vals[0][3] *= factor
	// row 2
	m.Vals[1][0] *= factor
	m.Vals[1][1] *= factor
	m.Vals[1][2] *= factor
	m.Vals[1][3] *= factor
	// row 3
	m.Vals[2][0] *= factor
	m.Vals[2][1] *= factor
	m.Vals[2][2] *= factor
	m.Vals[2][3] *= factor

	//// Unoptimized but readable variant:
	// transformMtrx := Matrix4x4{
	// 	Vals: [4][4]float64{
	// 		{scale, 0, 0, 0},
	// 		{0, scale, 0, 0},
	// 		{0, 0, scale, 0},
	// 		{0, 0, 0, 1},
	// 	},
	// }
	// m.Vals = transformMtrx.MultiplyByMatrix4x4(m).Vals
}

func (m *Matrix4x4) Translate(x, y, z float64) {
	// row 1
	m.Vals[0][0] += x * m.Vals[3][0]
	m.Vals[0][1] += x * m.Vals[3][1]
	m.Vals[0][2] += x * m.Vals[3][2]
	m.Vals[0][3] += x * m.Vals[3][3]
	// row 2
	m.Vals[1][0] += y * m.Vals[3][0]
	m.Vals[1][1] += y * m.Vals[3][1]
	m.Vals[1][2] += y * m.Vals[3][2]
	m.Vals[1][3] += y * m.Vals[3][3]
	// row 3
	m.Vals[2][0] += z * m.Vals[3][0]
	m.Vals[2][1] += z * m.Vals[3][1]
	m.Vals[2][2] += z * m.Vals[3][2]
	m.Vals[2][3] += z * m.Vals[3][3]
	// row 4 is unchanged

	//// Unoptimized but readable variant:
	// transformMtrx := Matrix4x4{
	// 	Vals: [4][4]float64{
	// 		{1, 0, 0, x},
	// 		{0, 1, 0, y},
	// 		{0, 0, 1, z},
	// 		{0, 0, 0, 1},
	// 	},
	// }
	// m.Vals = transformMtrx.MultiplyByMatrix4x4(m).Vals
}

func (m *Matrix4x4) RotateAroundX(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	tm1, tm2 := m.Vals[1][0], m.Vals[2][0]
	m.Vals[1][0] = tm1*cos + tm2*sin
	m.Vals[2][0] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m.Vals[1][1], m.Vals[2][1]
	m.Vals[1][1] = tm1*cos + tm2*sin
	m.Vals[2][1] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m.Vals[1][2], m.Vals[2][2]
	m.Vals[1][2] = tm1*cos + tm2*sin
	m.Vals[2][2] = tm1*(-sin) + tm2*cos
	tm1, tm2 = m.Vals[1][3], m.Vals[2][3]
	m.Vals[1][3] = tm1*cos + tm2*sin
	m.Vals[2][3] = tm1*(-sin) + tm2*cos

	//// Unoptimized but readable variant:
	// transformMtrx := Matrix4x4{
	// 	Vals: [4][4]float64{
	// 		{1, 0, 0, 0},
	// 		{0, cos, sin, 0},
	// 		{0, -sin, cos, 0},
	// 		{0, 0, 0, 1},
	// 	},
	// }
	// m.Vals = transformMtrx.MultiplyByMatrix4x4(m).Vals
}

func (m *Matrix4x4) RotateAroundY(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)

	tm1, tm2 := m.Vals[0][0], m.Vals[2][0]
	m.Vals[0][0] = tm1*cos + tm2*(-sin)
	m.Vals[2][0] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][1], m.Vals[2][1]
	m.Vals[0][1] = tm1*cos + tm2*(-sin)
	m.Vals[2][1] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][2], m.Vals[2][2]
	m.Vals[0][2] = tm1*cos + tm2*(-sin)
	m.Vals[2][2] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][3], m.Vals[2][3]
	m.Vals[0][3] = tm1*cos + tm2*(-sin)
	m.Vals[2][3] = tm1*(sin) + tm2*cos

	//// Unoptimized but readable variant:
	// transformMtrx := Matrix4x4{
	// 	Vals: [4][4]float64{
	// 		{cos, 0, -sin, 0},
	// 		{0, 1, 0, 0},
	// 		{sin, 0, cos, 0},
	// 		{0, 0, 0, 1},
	// 	},
	// }
	// m.Vals = transformMtrx.MultiplyByMatrix4x4(m).Vals
}

func (m *Matrix4x4) RotateAroundZ(radians float64) {
	sin := math.Sin(radians)
	cos := math.Cos(radians)

	tm1, tm2 := m.Vals[0][0], m.Vals[1][0]
	m.Vals[0][0] = tm1*cos + tm2*(-sin)
	m.Vals[1][0] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][1], m.Vals[1][1]
	m.Vals[0][1] = tm1*cos + tm2*(-sin)
	m.Vals[1][1] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][2], m.Vals[1][2]
	m.Vals[0][2] = tm1*cos + tm2*(-sin)
	m.Vals[1][2] = tm1*(sin) + tm2*cos
	tm1, tm2 = m.Vals[0][3], m.Vals[1][3]
	m.Vals[0][3] = tm1*cos + tm2*(-sin)
	m.Vals[1][3] = tm1*(sin) + tm2*cos

	//// Unoptimized but readable variant:
	// transformMtrx := Matrix4x4{
	// 	Vals: [4][4]float64{
	// 		{cos, -sin, 0, 0},
	// 		{sin, cos, 0, 0},
	// 		{0, 0, 1, 0},
	// 		{0, 0, 0, 1},
	// 	},
	// }
	// m.Vals = transformMtrx.MultiplyByMatrix4x4(m).Vals
}
