package object3d

import (
	"fmt"
	"strings"
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
	ObjectName string
	Vertexes   []Vertex3d
	Primitives []*Primitive
}

func (o *Object) Print(tabAmount int) {
	spaces := strings.Repeat(" ", tabAmount)
	result := spaces + "{\n"

	result += fmt.Sprintf(spaces+"  Object name: %s,\n", o.ObjectName)
	result += fmt.Sprintf(spaces+"  XFromParent: %d,\n", o.XFromParent)
	result += fmt.Sprintf(spaces+"  YFromParent: %d,\n", o.YFromParent)
	result += fmt.Sprintf(spaces+"  ZFromParent: %d,\n", o.ZFromParent)
	result += fmt.Sprintf(spaces+"  Vertexes (%d total): [\n", len(o.Vertexes))
	for index, v := range o.Vertexes {
		result += fmt.Sprintf(spaces+"    %d: %d, %d, %d\n", index, v.x, v.y, v.z)
	}
	result += fmt.Sprintf(spaces + "  ]\n")

	result += fmt.Sprintf(spaces+"  Primitives (%d total): [\n", len(o.Primitives))
	for _, prim := range o.Primitives {
		result += prim.ToString(tabAmount + 4)
	}
	result += fmt.Sprintf(spaces + "  ]\n")
	result += spaces + o.gatherParsedPrimitiveMetadata()
	result += fmt.Sprintf(spaces + "}\n")

	fmt.Printf(result)
}

func (obj *Object) gatherParsedPrimitiveMetadata() string {
	str := "Primitives metadata: "
	// find maxIndex vertex index
	minIndex, maxIndex := 65536, 0
	minVertices, maxVertices := 65536, 0
	for _, p := range obj.Primitives {
		for _, ind := range p.vertexIndices {
			if ind < minIndex {
				minIndex = ind
			}
			if ind > maxIndex {
				maxIndex = ind
			}
		}
		numVerts := len(p.vertexIndices)
		if numVerts > maxVertices {
			maxVertices = numVerts
		}
		if numVerts < minVertices {
			minVertices = numVerts
		}
	}
	str += fmt.Sprintf("Vertex counts: %d-%d, vertex indices: %d-%d\n", minVertices, maxVertices, minIndex, maxIndex)
	return str
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
	return obj
}
