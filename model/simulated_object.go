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
	CobMachine cob.CobMachine
	Script     *scripts.CobScript
}

func (so *SimObject) InitFromCavedogData(cavedogModel *object3d.Object,
	textures []*texture.GafEntry, cobScript *scripts.CobScript) {

	modelgeometry := NewModelFrom3doObject3d(cavedogModel, textures)
	so.ModelState = CreateObjectFromModel(modelgeometry)
	if cobScript != nil { // TODO: make this obligatory
		so.Script = cobScript
		// We need to call 'Create' COB subprogram here.
		// First, determine which virtual address we need
		var createFuncAddr int32 = -1
		for i := range cobScript.ProcedureNames {
			if cobScript.ProcedureNames[i] == "Create" {
				createFuncAddr = cobScript.ProcedureAddresses[i]
			}
		}
		if createFuncAddr == -1 {
			panic("COB INIT ERROR: No 'Create' func found!")
		}
		so.CobMachine.AllocNewThread(createFuncAddr, 0)
	}
}
