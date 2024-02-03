package model

import (
	"fmt"
	"math"
	"totala_reader/geometry"
)

func (m *Model) performUvMappingOnAllSurfaces() {
	for _, prim := range m.Primitives {
		// UV-mapping for colored (= non-textured) primitives and for non-3d objects is not requred.
		if prim.Texture != nil {
			m.uvMapSurface(prim)
			prim.calculateUvCenterCoords()
		}
	}
}

func (m *Model) uvMapSurface(prim *ModelSurface) {
	// It's very simple if the polygon is quad:
	if len(prim.VertexIndices) == 4 {
		prim.UVCoordinatesPerIndex = [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
		return
	} else {
		// This panic left as I'm yet to see a model with non-quad primitive; need explicit warning for that.
		panic(fmt.Sprintf("REMOVE THIS PANIC: Primitive with %d vertices detected!", len(prim.VertexIndices)))
	}

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
	fmt.Printf("Lcplne crds: %.2v\n", uv)

	// Normalize UV coords (transform them to range [0.0 - 1.0])
	minU, minV := math.MaxFloat64, math.MaxFloat64
	for i := range uv {
		if uv[i][0] < minU {
			minU = uv[i][0]
		}
		if uv[i][1] < minV {
			minV = uv[i][1]
		}
	}
	for i := range uv {
		uv[i][0] -= minU
		uv[i][1] -= minV
	}

	maxU, maxV := 0.0, 0.0
	for i := range uv {
		if uv[i][0] > maxU {
			maxU = uv[i][0]
		}
		if uv[i][1] > maxV {
			maxV = uv[i][1]
		}
	}
	for i := range uv {
		uv[i][0] /= maxU
		uv[i][1] /= maxV
	}
	fmt.Printf("UV normalzd: %.2v\n", uv)

	prim.UVCoordinatesPerIndex = uv
}

func (prim *ModelSurface) calculateUvCenterCoords() {
	for i := range prim.UVCoordinatesPerIndex {
		prim.CenterUVCoords[0] += prim.UVCoordinatesPerIndex[i][0]
		prim.CenterUVCoords[1] += prim.UVCoordinatesPerIndex[i][1]
	}
	prim.CenterUVCoords[0] /= float64(len(prim.UVCoordinatesPerIndex))
	prim.CenterUVCoords[1] /= float64(len(prim.UVCoordinatesPerIndex))
}
