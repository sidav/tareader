package model

import (
	"totala_reader/model/cob"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/scripts"
	"totala_reader/ta_files_read/texture"
)

// Top-level unit struct: something scriptable and with a model.
type SimObject struct {
	ModelState *ModelledObject

	CobMachine    cob.CobMachine
	Script        *scripts.CobScript
	PiecesMapping []*ModelledObject // maps numbered piece name to a piece, so that scripts could address them directly
}

func (so *SimObject) InitFromCavedogData(cavedogModel *object3d.Object,
	textures []*texture.GafEntry, cobScript *scripts.CobScript) {

	modelgeometry := NewModelFrom3doObject3d(cavedogModel, textures)
	so.ModelState = CreateObjectFromModel(modelgeometry)

	if cobScript != nil { // TODO: make this obligatory (crash if nil)
		so.Script = cobScript
		so.mapPieces()
		// We need to call 'Create' COB subprogram here.
		// First, determine which virtual address we need
		var createFuncAddr int32 = cobScript.FindFuncAddressByName("Create")
		if createFuncAddr == -1 {
			panic("COB INIT ERROR: No 'Create' func found!")
		}
		so.CobMachine.AllocNewThread(createFuncAddr, 0)

		// Debug purposes: run GO script to show some animation.
		goFuncAddr := cobScript.FindFuncAddressByName("Go")
		if goFuncAddr != -1 {
			so.CobMachine.AllocNewThread(goFuncAddr, 0)
		}
	}
}

func (so *SimObject) mapPieces() {
	print("Mapping the pieces for COB script:\n")
	so.PiecesMapping = make([]*ModelledObject, len(so.Script.Pieces), len(so.Script.Pieces))
	for i := range so.Script.Pieces {
		print("    Searching %s...", so.Script.Pieces[i])
		so.PiecesMapping[i] = so.ModelState.findPieceByName(so.Script.Pieces[i])
		if so.PiecesMapping[i] == nil {
			panic("Piece mapping failure. Can't find piece " + so.Script.Pieces[i])
		}
		print(" ok.\n")
	}
}

func (so *SimObject) PerformScriptedMovementStep() {
	for _, piece := range so.PiecesMapping {
		piece.performScriptedMovement()
	}
}
