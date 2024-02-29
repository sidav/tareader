package model

import "math"

func (piece *ModelledObject) SetSpin(axis, cobSpeed, cobAcceleration int32) {
	piece.TurnSpeed[axis] = CavedogAngleSpeedToFloatRadians(cobSpeed)
	piece.IsSpinning[axis] = true
}

func (piece *ModelledObject) SetMove(axis, cobPos, cobSpeed int32) {
	// Notice that it seems the speed is always positive
	speed := CavedogPositionSpeedToFloat(cobSpeed)
	target := CavedogPositionToFloatPosition(cobPos)

	if axis != 1 { // It's important but I dunno why is it like that for X and Z axii. TODO: investigate.
		target = -target
	}

	if piece.CurrentMove[axis] > target {
		speed = -speed
	}
	piece.TargetMove[axis] = target
	piece.MoveSpeed[axis] = speed
}

func (piece *ModelledObject) performScriptedMovement() {
	// perform spinning
	if piece.IsSpinning[0] {
		piece.Matrix.RotateAroundX(piece.TurnSpeed[0])
	}
	if piece.IsSpinning[1] {
		piece.Matrix.RotateAroundY(piece.TurnSpeed[1])
	}
	if piece.IsSpinning[2] {
		piece.Matrix.RotateAroundZ(piece.TurnSpeed[2])
	}

	// perform moving
	var move [3]float64
	for axis := range piece.CurrentMove {
		if math.Abs(piece.TargetMove[axis]-piece.CurrentMove[axis]) < math.Abs(piece.MoveSpeed[axis]) {
			move[axis] = piece.TargetMove[axis] - piece.CurrentMove[axis]
		} else {
			move[axis] = piece.MoveSpeed[axis]
		}
		piece.CurrentMove[axis] += move[axis]
	}
	piece.Matrix.Translate(move[0], move[1], move[2])
}
