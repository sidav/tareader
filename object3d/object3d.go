package object3d

import (
	binaryreader "totala_reader/binary_reader"
)

type Object struct {
	// metadata from the file "header"
	VersionSignature           int
	NumberOfVertexes           int
	NumberOfPrimitives         int
	OffsetToselectionPrimitive int
	XFromParent                int
	YFromParent                int
	ZFromParent                int
	OffsetToObjectName         int
	always0                    int
	OffsetToVertexArray        int
	OffsetToPrimitiveArray     int
	OffsetToSiblingObject      int
	OffsetToChildObject        int

	// the object data itself
	ObjectName         string
	Vertexes           []Vertex3d
	Primitives         []*Primitive
	SelectionPrimitive *Primitive
	ChildObject        *Object
	SiblingObject      *Object
}

func ReadObjectFromReader(r *binaryreader.Reader, modelOffset int) *Object {
	obj := &Object{
		VersionSignature:           r.ReadIntFromBytesArray(modelOffset, 0),
		NumberOfVertexes:           r.ReadIntFromBytesArray(modelOffset, 4),
		NumberOfPrimitives:         r.ReadIntFromBytesArray(modelOffset, 8),
		OffsetToselectionPrimitive: r.ReadIntFromBytesArray(modelOffset, 12),
		XFromParent:                r.ReadIntFromBytesArray(modelOffset, 16),
		YFromParent:                r.ReadIntFromBytesArray(modelOffset, 20),
		ZFromParent:                r.ReadIntFromBytesArray(modelOffset, 24),
		OffsetToObjectName:         r.ReadIntFromBytesArray(modelOffset, 28),
		always0:                    r.ReadIntFromBytesArray(modelOffset, 32),
		OffsetToVertexArray:        r.ReadIntFromBytesArray(modelOffset, 36),
		OffsetToPrimitiveArray:     r.ReadIntFromBytesArray(modelOffset, 40),
		OffsetToSiblingObject:      r.ReadIntFromBytesArray(modelOffset, 44),
		OffsetToChildObject:        r.ReadIntFromBytesArray(modelOffset, 48),
	}

	obj.ObjectName = r.ReadNullTermStringFromBytesArray(0, obj.OffsetToObjectName)
	obj.Vertexes = ReadVertexesFromReader(r, obj.OffsetToVertexArray, obj.NumberOfVertexes)
	obj.Primitives = ReadPrimitivesArrayFromReader(r, obj.OffsetToPrimitiveArray, obj.NumberOfPrimitives)
	if obj.OffsetToselectionPrimitive != -1 && len(obj.Primitives) > 0 {
		// Bugged? Is it the index or is it really byte-offset?
		obj.SelectionPrimitive = obj.Primitives[obj.OffsetToselectionPrimitive]
	}

	if obj.OffsetToChildObject != 0 {
		obj.ChildObject = ReadObjectFromReader(r, obj.OffsetToChildObject)
	}
	if obj.OffsetToSiblingObject != 0 {
		obj.SiblingObject = ReadObjectFromReader(r, obj.OffsetToSiblingObject)
	}

	return obj
}
