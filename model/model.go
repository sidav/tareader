package model

import "totala_reader/ta_files_read/object3d"

// Non-TA format (floats everywhere) for ease of rendering
type Model struct {
	Vertices                              [][3]float64
	ObjectName                            string
	XFromParent, YFromParent, ZFromParent float64
	Primitives                            []*ModelSurface
	SelectionPrimitive                    *ModelSurface
	ChildObject                           *Model
	SiblingObject                         *Model
}

func NewModelFrom3doObject3d(obj *object3d.Object) *Model {
	model := &Model{
		ObjectName: obj.ObjectName,
	}
	model.XFromParent, model.YFromParent, model.ZFromParent = obj.ParentOffsetAsFloats()
	for _, v := range obj.Vertexes {
		x, y, z := v.ToFloats()
		model.Vertices = append(model.Vertices, [3]float64{x, y, z})
	}
	for _, p := range obj.Primitives {
		newSurf := newModelSurfaceFrom3doPrimitive(p)
		model.Primitives = append(model.Primitives, newSurf)
		// TODO: have selection primitive saved separately from other ones
		if p == obj.SelectionPrimitive {
			model.SelectionPrimitive = newSurf
		}
	}
	if obj.ChildObject != nil {
		model.ChildObject = NewModelFrom3doObject3d(obj.ChildObject)
	}
	if obj.SiblingObject != nil {
		model.SiblingObject = NewModelFrom3doObject3d(obj.SiblingObject)
	}
	return model
}
