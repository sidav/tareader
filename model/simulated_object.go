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
	CobState   cob.CobMachine
	Script     *scripts.CobScript
}

func (so *SimObject) InitFromCavedogData(cavedogModel *object3d.Object,
	textures []*texture.GafEntry, cobScript *scripts.CobScript) {

	modelgeometry := NewModelFrom3doObject3d(cavedogModel, textures)
	so.ModelState = CreateObjectFromModel(modelgeometry)
	if cobScript != nil { // TODO: make this obligatory
		so.Script = cobScript
		so.CobState.Threads[0].Active = true
		// TODO: call 'Create' COB subprogram here?
	}
}
