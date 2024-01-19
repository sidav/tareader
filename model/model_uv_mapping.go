package model

import (
	"fmt"
	"totala_reader/geometry"
)

func (m *Model) performUvMappingOnAllSurfaces() {
	for _, prim := range m.Primitives {
		// UV-mapping for colored (= non-textured) primitives and for non-3d objects is not requred.
		if len(prim.VertexIndices) > 2 && !prim.IsColored && prim != m.SelectionPrimitive {
			m.uvMapSurface(prim)
		}
	}
}

func (m *Model) uvMapSurface(prim *ModelSurface) {
	// We take the first index (at position 0 in indices array) as (u,v) = (0,0)
	// Step 1. Collect all the coordinates.
	var allCoords [][3]float64
	for _, index := range prim.VertexIndices {
		allCoords = append(allCoords, [3]float64{
			m.Vertices[index][0],
			m.Vertices[index][1],
			m.Vertices[index][2],
		})
	}

	fmt.Printf("----------\nInitial    : %.2v\n", allCoords)
	// Step 2. Transform the primitive coordinates space so that the first will be at [0, 0, 0]
	for i := 1; i < len(allCoords); i++ {
		for j := 0; j < 3; j++ {
			allCoords[i][j] -= allCoords[0][j]
		}
	}
	allCoords[0][0], allCoords[0][1], allCoords[0][2] = 0.0, 0.0, 0.0
	fmt.Printf("Transformed: %.2v\n", allCoords)

	// It's a WIP variant based on vector math
	vec2 := geometry.NewVector3FromArr(allCoords[2])
	localX := geometry.NewVector3FromArr(allCoords[1])
	normal := geometry.CrossProduct(&localX, &vec2)
	localY := geometry.CrossProduct(&normal, &localX)
	localX.Normalize()
	localY.Normalize()

	uv := make([][2]float64, len(allCoords))
	for i := 0; i < len(allCoords); i++ {
		currCoords := geometry.NewVector3FromArr(allCoords[i])
		uv[i] = [2]float64{
			geometry.DotProduct(&currCoords, &localX),
			geometry.DotProduct(&currCoords, &localY),
		}
	}
	fmt.Printf("UV         : %.2v\n", uv)
	prim.UVCoordinatesPerIndex = uv
}
