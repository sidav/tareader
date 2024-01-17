package model

import (
	"fmt"
	"strings"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/texture"
)

// TODO: split into triangles?..
type ModelSurface struct {
	VertexIndices []int
	Texture       *texture.GafEntry
}

func newModelSurfaceFrom3doPrimitive(p *object3d.Primitive, allTextures []*texture.GafEntry) *ModelSurface {
	ms := &ModelSurface{}
	ms.VertexIndices = make([]int, len(p.VertexIndices))
	copy(ms.VertexIndices, p.VertexIndices)
	// Assign GAF texture to this surface
	if p.TextureName != "" {
		for _, tex := range allTextures {
			if strings.ToLower(tex.Name) == strings.ToLower(p.TextureName) {
				ms.Texture = tex
				break
			}
		}
		if ms.Texture == nil {
			print(fmt.Sprintf("Texture '%s' (bytes %v) not found.\n", p.TextureName, []byte(p.TextureName)))
			// panic()
		}
	}
	return ms
}
