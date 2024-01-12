package object3d

import (
	"fmt"
	"strings"
)

func (o *Object) ToString(tabAmount int) string {
	spaces := strings.Repeat(" ", tabAmount)
	result := spaces + "{\n"

	result += fmt.Sprintf(spaces+"  Object name: %s,\n", o.ObjectName)
	result += fmt.Sprintf(spaces+"  XFromParent: %d,\n", o.XFromParent)
	result += fmt.Sprintf(spaces+"  YFromParent: %d,\n", o.YFromParent)
	result += fmt.Sprintf(spaces+"  ZFromParent: %d,\n", o.ZFromParent)
	result += fmt.Sprintf(spaces+"  Vertexes (%d total): [\n", len(o.Vertexes))
	for index, v := range o.Vertexes {
		// result += fmt.Sprintf(spaces+"    %3d: %s\n", index, v.ToString(0))
		result += fmt.Sprintf(spaces+"    %3d: %s\n", index, v.ToFloatsString(0))
	}
	result += fmt.Sprintf(spaces + "  ],\n")

	result += fmt.Sprintf(spaces+"  Primitives (%d total): [\n", len(o.Primitives))
	for _, prim := range o.Primitives {
		result += prim.ToString(tabAmount + 4)
	}
	result += fmt.Sprintf(spaces + "  ]\n")
	result += spaces + "  " + o.gatherParsedPrimitiveMetadata()

	if o.ChildObject != nil {
		result += o.ChildObject.ToString(tabAmount + 2)
	}
	if o.SiblingObject != nil {
		result += o.SiblingObject.ToString(tabAmount + 2)
	}

	result += fmt.Sprintf(spaces + "}\n")

	return result
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

func (p *Primitive) ToString(tabAmount int) string {
	spaces := strings.Repeat(" ", tabAmount)
	str := spaces + "{\n"
	str += fmt.Sprintf(spaces+"  ColorPaletteIndex: %d,\n", p.ColorIndex)
	str += fmt.Sprintf(spaces+"  IsColored: %v,\n", p.IsColored)
	str += fmt.Sprintf(spaces+"  TextureName: \"%s\",\n", p.TextureName)
	str += fmt.Sprintf(spaces+"  VertexIndices: %v,\n", p.vertexIndices)
	str += spaces + "}\n"
	return str
}

func (v *Vertex3d) ToString(tabAmount int) string {
	spaces := strings.Repeat(" ", tabAmount)
	return spaces + fmt.Sprintf("[%8d %8d %8d],", v.x, v.y, v.z)
}

func (v *Vertex3d) ToFloatsString(tabAmount int) string {
	spaces := strings.Repeat(" ", tabAmount)
	return spaces + fmt.Sprintf("[%.2f, %.2f, %.2f],", float64(v.x)/65536, float64(v.y)/65536, float64(v.z)/65536)
}

// func (v *Vertex3d) ToFixedPointsString(tabAmount int) string {
// 	spaces := strings.Repeat(" ", tabAmount)
// 	xi, xr := intToFixedPoint(v.x)
// 	yi, yr := intToFixedPoint(v.y)
// 	zi, zr := intToFixedPoint(v.z)
// 	return spaces + fmt.Sprintf("[%d.%d, %d.%d, %d.%d],", xi, xr, yi, yr, zi, zr)
// }
