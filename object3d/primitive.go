package object3d

import (
	binaryreader "totala_reader/binary_reader"
)

type Primitive struct {
	ColorIndex    int
	IsColored     bool
	TextureName   string
	vertexIndices []int
}

func ReadPrimitiveFromReader(r *binaryreader.Reader, primOffset int) *Primitive {
	prim := Primitive{}

	ColorIndex := r.ReadIntFromBytesArray(primOffset, 0)
	NumberOfVertexIndexes := r.ReadIntFromBytesArray(primOffset, 4)
	// Always_0 := r.ReadIntFromBytesArray(primOffset, 8)
	OffsetToVertexIndexArray := r.ReadIntFromBytesArray(primOffset, 12)
	OffsetToTextureName := r.ReadIntFromBytesArray(primOffset, 16)
	// Unknown_1 := r.ReadIntFromBytesArray(primOffset, 20)
	// Unknown_2 := r.ReadIntFromBytesArray(primOffset, 24)
	IsColored := r.ReadIntFromBytesArray(primOffset, 28)

	prim.ColorIndex = ColorIndex
	prim.IsColored = IsColored != 0
	prim.TextureName = r.ReadNullTermStringFromBytesArray(OffsetToTextureName, 0)
	// reading the vertex indices array
	prim.vertexIndices = make([]int, NumberOfVertexIndexes)
	for i := 0; i < NumberOfVertexIndexes; i++ {
		prim.vertexIndices[i] = r.ReadUint16FromBytesArray(OffsetToVertexIndexArray, i*2)
	}
	return &prim
}

func ReadPrimitivesArrayFromReader(r *binaryreader.Reader, primArrOffset, total int) []*Primitive {
	var primArr []*Primitive

	for i := 0; i < total; i++ {
		primArr = append(primArr, ReadPrimitiveFromReader(r, primArrOffset+(i*32)))
	}

	return primArr
}
