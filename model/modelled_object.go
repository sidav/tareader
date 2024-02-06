package model

import (
	"fmt"
	"strings"
	"totala_reader/geometry/matrix4x4"
)

// It's NOT a model, it's something which has a model. Something with its own rotation matrices.
// One model can have multiple ModelledObjects, but a single ModelledObject can belong only to one Model.
// TODO: rename! There is too much stuff called "object" already.
type ModelledObject struct {
	Name           string
	Matrix         *matrix4x4.Matrix4x4
	Child          *ModelledObject
	Sibling        *ModelledObject
	ModelForObject *Model
}

func (mo *ModelledObject) Print(offset int) {
	fmt.Printf("%s'%s'\n", strings.Repeat(" ", offset), mo.Name)
	if mo.Child != nil {
		mo.Child.Print(offset + 2)
	}
	if mo.Sibling != nil {
		mo.Sibling.Print(offset)
	}
}

// TODO: translate the matrices by offsets in model objects.
func CreateObjectFromModel(m *Model) *ModelledObject {
	return initObjectFromModel(m, nil)
}

func initObjectFromModel(m *Model, parentObj *ModelledObject) *ModelledObject {
	obj := &ModelledObject{
		Name:           m.ObjectName,
		Matrix:         matrix4x4.NewUnitMatrix(),
		ModelForObject: m,
	}
	const scale = 6
	if parentObj == nil {
		obj.Matrix.Scale(scale)
	} else {
		// obj.Matrix.Translate(m.XFromParent, m.YFromParent, m.ZFromParent)
		// obj.Matrix.Translate(m.XFromParent*scale, m.YFromParent*scale, m.ZFromParent*scale)
	}
	if m.ChildObject != nil {
		obj.Child = initObjectFromModel(m.ChildObject, obj)
	}
	if m.SiblingObject != nil {
		obj.Sibling = initObjectFromModel(m.SiblingObject, parentObj)
	}
	return obj
}
