package model

import "totala_reader/geometry/matrix4x4"

// It's NOT a model, it's something which has a model. Something with its own rotation matrices.
// One model can have multiple ModelledObjects, but a single ModelledObject can belong only to one Model.
// TODO: rename! There is too much stuff called "object" already.
type ModelledObject struct {
	Name           string
	Matrix         matrix4x4.Matrix4x4
	Child          *ModelledObject
	Sibling        *ModelledObject
	ModelForObject *Model
}

// TODO: translate the matrices by offsets in model objects.
func NewObjectFromModel(m *Model) *ModelledObject {
	obj := &ModelledObject{
		Name:           m.ObjectName,
		Matrix:         *matrix4x4.NewUnitMatrix(),
		ModelForObject: m,
	}
	if m.ChildObject != nil {
		obj.Child = NewObjectFromModel(m.ChildObject)
	}
	if m.SiblingObject != nil {
		obj.Sibling = NewObjectFromModel(m.SiblingObject)
	}
	return obj
}
