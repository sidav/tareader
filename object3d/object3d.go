package object3d

import (
	binaryreader "totala_reader/binary_reader"
)

type Object struct {
	ObjectName         string
	XFromParent        int
	YFromParent        int
	ZFromParent        int
	Vertexes           []Vertex3d
	Primitives         []*Primitive
	SelectionPrimitive *Primitive
	ChildObject        *Object
	SiblingObject      *Object
}

func (o *Object) ParentOffsetAsFloats() (float64, float64, float64) {
	return FixedPointToFloat(o.XFromParent), FixedPointToFloat(o.YFromParent), FixedPointToFloat(o.ZFromParent)
}

func ReadObjectFromReader(r *binaryreader.Reader, modelOffset int) *Object {
	// metadata from the file descriptor
	// VersionSignature := r.ReadIntFromBytesArray(modelOffset, 0)
	NumberOfVertexes := r.ReadIntFromBytesArray(modelOffset, 4)
	NumberOfPrimitives := r.ReadIntFromBytesArray(modelOffset, 8)
	OffsetToselectionPrimitive := r.ReadIntFromBytesArray(modelOffset, 12)
	XFromParent := r.ReadIntFromBytesArray(modelOffset, 16)
	YFromParent := r.ReadIntFromBytesArray(modelOffset, 20)
	ZFromParent := r.ReadIntFromBytesArray(modelOffset, 24)
	OffsetToObjectName := r.ReadIntFromBytesArray(modelOffset, 28)
	// always0:                    r.ReadIntFromBytesArray(modelOffset, 32)
	OffsetToVertexArray := r.ReadIntFromBytesArray(modelOffset, 36)
	OffsetToPrimitiveArray := r.ReadIntFromBytesArray(modelOffset, 40)
	OffsetToSiblingObject := r.ReadIntFromBytesArray(modelOffset, 44)
	OffsetToChildObject := r.ReadIntFromBytesArray(modelOffset, 48)

	obj := &Object{
		XFromParent: XFromParent,
		YFromParent: YFromParent,
		ZFromParent: ZFromParent,
		ObjectName:  r.ReadNullTermStringFromBytesArray(0, OffsetToObjectName),
		Vertexes:    ReadVertexesFromReader(r, OffsetToVertexArray, NumberOfVertexes),
		Primitives:  ReadPrimitivesArrayFromReader(r, OffsetToPrimitiveArray, NumberOfPrimitives),
	}

	if OffsetToselectionPrimitive != -1 && len(obj.Primitives) > 0 {
		// Bugged? Is it the index or is it really byte-offset?
		obj.SelectionPrimitive = obj.Primitives[OffsetToselectionPrimitive]
	}
	if OffsetToChildObject != 0 {
		obj.ChildObject = ReadObjectFromReader(r, OffsetToChildObject)
	}
	if OffsetToSiblingObject != 0 {
		obj.SiblingObject = ReadObjectFromReader(r, OffsetToSiblingObject)
	}

	return obj
}
