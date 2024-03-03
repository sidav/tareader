package model

import "math"

func (piece *ModelledObject) SetSpin(axis, cobSpeed, cobAcceleration int32) {
	piece.TurnSpeed[axis] = CavedogAngleSpeedToFloatRadians(cobSpeed)
	piece.IsSpinning[axis] = true
}

func (piece *ModelledObject) SetTurn(axis, cobAngularSpeed, cobAngle int32) {
	targAngle := CavedogAngleToFloatRadians(cobAngle)
	turnSpeed := CavedogAngleSpeedToFloatRadians(cobAngularSpeed)
	piece.TargetTurn[axis] = targAngle
	piece.TurnSpeed[axis] = turnSpeed
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
	// All rotations will be local. TODO: think how to optimize it
	// (a variant: maybe apply translations by piece.CurrentMove only in the renderer, together with XYZFromParent?)
	piece.Matrix.DiscardTranslations()
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

	////////////////////////////////////////////////////
	// perform turning
	piece.doScriptedTurn()

	///////////////////////////////////////////////////
	// perform moving
	for axis := range piece.CurrentMove {
		if math.Abs(piece.TargetMove[axis]-piece.CurrentMove[axis]) < math.Abs(piece.MoveSpeed[axis]) {
			piece.CurrentMove[axis] += piece.TargetMove[axis] - piece.CurrentMove[axis]
		} else {
			piece.CurrentMove[axis] += piece.MoveSpeed[axis]
		}
	}
	// "Un-local" all the operations.
	piece.Matrix.Translate(piece.CurrentMove[0], piece.CurrentMove[1], piece.CurrentMove[2])
}

func (piece *ModelledObject) doScriptedTurn() {
	for axis := range piece.CurrentTurn {

		var turnVal float64
		delta := piece.TargetTurn[axis] - piece.CurrentTurn[axis]
		if delta > pi {
			delta -= pi2
		} else if delta <= -pi {
			delta += pi2
		}
		if math.Abs(delta) < piece.TurnSpeed[axis] {
			turnVal = delta
			// piece.TurnSpeed[axis] = 0
		} else if delta < 0 {
			turnVal = -piece.TurnSpeed[axis]
		} else if delta > 0 {
			turnVal = piece.TurnSpeed[axis]
		}
		piece.CurrentTurn[axis] += turnVal
		switch axis {
		case 0:
			piece.Matrix.RotateAroundX(turnVal)
		case 1:
			piece.Matrix.RotateAroundY(turnVal)
		case 2:
			// TODO: change matrix operations so that the minus won't be needed?
			piece.Matrix.RotateAroundZ(-turnVal)
		}
		// for piece.CurrentTurn[axis] < 0 {
		// 	piece.CurrentTurn[axis] += pi2
		// }
		// for piece.CurrentTurn[axis] > pi2 {
		// 	piece.CurrentTurn[axis] -= pi2
		// }
	}
}

// Instant movement
func (piece *ModelledObject) moveNow(axis, cobPos int32) {
	var movement [3]float64
	moveFloat := CavedogPositionToFloatPosition(cobPos)
	movement[axis] = moveFloat - piece.CurrentMove[axis]
	piece.CurrentMove[axis] = moveFloat
	// The following is maybe redundant, as the transition is applied in performScriptedMovement?
	piece.Matrix.Translate(movement[0], movement[1], movement[2])
}

// Instant turn
func (piece *ModelledObject) turnNow(axis, cobAngularPos int32) {
	angle := CavedogAngleToFloatRadians(cobAngularPos)
	delta := angle - piece.CurrentTurn[axis]
	if delta > pi {
		delta -= pi2
	} else if delta <= -pi {
		delta += pi2
	}
	piece.CurrentTurn[axis] += delta
	switch axis {
	case 0:
		piece.Matrix.RotateAroundX(delta)
	case 1:
		piece.Matrix.RotateAroundY(delta)
	case 2:
		// TODO: change matrix operations so that the minus won't be needed?
		piece.Matrix.RotateAroundZ(-delta)
	}
}
