package model

import (
	"totala_reader/geometry"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/texture"
)

// Non-TA format (floats everywhere) for ease of rendering
type Model struct {
	SelectionPrimitive                    *ModelSurface
	ChildObject                           *Model
	SiblingObject                         *Model
	ObjectName                            string
	Vertices                              [][3]float64
	Primitives                            []*ModelSurface
	XFromParent, YFromParent, ZFromParent float64
}

func NewModelFrom3doObject3d(obj *object3d.Object, allTextures []*texture.GafEntry) *Model {
	model := &Model{
		ObjectName: obj.ObjectName,
	}
	model.XFromParent, model.YFromParent, model.ZFromParent = obj.ParentOffsetAsFloats()
	for _, v := range obj.Vertexes {
		x, y, z := v.ToFloats()
		model.Vertices = append(model.Vertices, [3]float64{x, y, z})
	}
	for _, p := range obj.Primitives {
		newSurf := newModelSurfaceFrom3doPrimitive(p, allTextures)
		model.Primitives = append(model.Primitives, newSurf)
		// TODO: have selection primitive saved separately from other ones
		if p == obj.SelectionPrimitive {
			model.SelectionPrimitive = newSurf
		}
	}
	// Calculate and store centers for all surfaces
	model.calcCenterOfAllSurfaces()
	// Calculate UV-mapping
	model.performUvMappingOnAllSurfaces()
	// Calculate normals
	model.calcNormalsForAllSurfaces()
	if obj.ChildObject != nil {
		model.ChildObject = NewModelFrom3doObject3d(obj.ChildObject, allTextures)
	}
	if obj.SiblingObject != nil {
		model.SiblingObject = NewModelFrom3doObject3d(obj.SiblingObject, allTextures)
	}
	return model
}

func (m *Model) calcCenterOfAllSurfaces() {
	for _, surf := range m.Primitives {
		for _, index := range surf.VertexIndices {
			surf.CenterCoords[0] += m.Vertices[index][0]
			surf.CenterCoords[1] += m.Vertices[index][1]
			surf.CenterCoords[2] += m.Vertices[index][2]
		}
		for i := range surf.CenterCoords {
			surf.CenterCoords[i] /= float64(len(surf.VertexIndices))
		}
	}
}

func (m *Model) calcNormalsForAllSurfaces() {
	for _, surf := range m.Primitives {
		if len(surf.VertexIndices) < 3 {
			continue
		}
		// Take coords 0, 1 and 2.
		coords0 := geometry.NewVector3FromArr(m.Vertices[surf.VertexIndices[0]])
		coords1 := geometry.NewVector3FromArr(m.Vertices[surf.VertexIndices[1]])
		coords2 := geometry.NewVector3FromArr(m.Vertices[surf.VertexIndices[2]])
		// Move them relative to vertex 1.
		coords0.Sub(coords1)
		coords2.Sub(coords1)
		// Find the cross product.
		surf.NormalVector = geometry.CrossProduct(&coords0, &coords2)
		surf.NormalVector.Normalize()
	}
}
