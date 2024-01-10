package object3d

import binaryreader "totala_reader/binary_reader"

type Vertex3d struct {
	x, y, z int
}

func ReadVertexesFromReader(r *binaryreader.Reader, vertexArrayOffset, vertexCount int) []Vertex3d {
	var vertexArray []Vertex3d
	for vInd := 0; vInd < vertexCount*3; vInd += 4 * 3 {
		vertexArray = append(vertexArray, Vertex3d{
			r.ReadIntFromBytesArray(vertexArrayOffset, vInd),
			r.ReadIntFromBytesArray(vertexArrayOffset, vInd+4),
			r.ReadIntFromBytesArray(vertexArrayOffset, vInd+8),
		})
	}
	return vertexArray
}
