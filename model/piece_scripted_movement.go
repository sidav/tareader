package model

func (piece *ModelledObject) SetSpin(axis, cobSpeed, cobAcceleration int32) {
	piece.TurnSpeed[axis] = CavedogAngleSpeedToFloatRadians(cobSpeed)
	piece.IsSpinning[axis] = true
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
}
