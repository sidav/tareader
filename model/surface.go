package model

import "totala_reader/ta_files_read/object3d"

// TODO: split into triangles?..
type ModelSurface struct {
	VertexIndices []int
}

func newModelSurfaceFrom3doPrimitive(p *object3d.Primitive) *ModelSurface {
	ms := &ModelSurface{}
	ms.VertexIndices = make([]int, len(p.VertexIndices))
	copy(ms.VertexIndices, p.VertexIndices)
	return ms
}
