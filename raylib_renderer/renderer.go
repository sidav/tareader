package raylibrenderer

import (
	"sort"
	graphicadapter "totala_reader/graphic_adapter"
	"totala_reader/model"
)

type RaylibRenderer struct {
	gAdapter                   graphicadapter.GraphicBackend
	onScreenOffX, onScreenOffY int32
	fontSize                   int32
	scaleFactor                float64
	totalMessages              int32

	frame int

	zBuffer                                [][]float64
	zBufMinX, zBufMaxX, zBufMinY, zBufMaxY int32 // for clearing changed area of the buffer only

	trianglesBatch []*triangle

	debugMode bool
}

func (r *RaylibRenderer) Init(adapter graphicadapter.GraphicBackend) {
	r.gAdapter = adapter
	r.onScreenOffX, r.onScreenOffY = r.gAdapter.GetRenderResolution()
	r.onScreenOffX /= 2
	r.onScreenOffY = 55 * r.onScreenOffY / 100
	r.scaleFactor = 4
	r.gAdapter.Clear()
	r.initZBuffer()
}

func (r *RaylibRenderer) DrawModel(rootObject *model.Model) {

	r.clearZBuffer()

	r.trianglesBatch = nil
	r.totalMessages = 0
	r.gAdapter.Clear()
	r.drawObject(rootObject, 0, 0, 0)
	r.DrawTrianglesBatch()
	r.frame++
}

func (r *RaylibRenderer) drawObject(obj *model.Model, parentOffsetX, parentOffsetY, parentOffsetZ float64) {
	currentOffsetX, currentOffsetY, currentOffsetZ := obj.XFromParent+parentOffsetX,
		obj.YFromParent+parentOffsetY, obj.ZFromParent+parentOffsetZ

	for _, p := range obj.Primitives {
		r.drawPrimitive(obj, p, currentOffsetX, currentOffsetY, currentOffsetZ)
	}

	if obj.ChildObject != nil && len(obj.ChildObject.Primitives) > 0 {
		r.drawObject(obj.ChildObject, currentOffsetX, currentOffsetY, currentOffsetZ)
	}
	if obj.SiblingObject != nil && len(obj.SiblingObject.Primitives) > 0 {
		r.drawObject(obj.SiblingObject, parentOffsetX, parentOffsetY, parentOffsetZ)
	}
}

func (r *RaylibRenderer) drawPrimitive(obj *model.Model, prim *model.ModelSurface, offsetX, offsetY, offsetZ float64) {
	if len(prim.VertexIndices) < 3 || obj.SelectionPrimitive == prim {
		return
	}
	// fill the triangles batch
	zerox, zeroy, zeroz := (obj.Vertices[prim.VertexIndices[0]][0]+offsetX)*r.scaleFactor,
		(obj.Vertices[prim.VertexIndices[0]][1]+offsetY)*r.scaleFactor,
		(obj.Vertices[prim.VertexIndices[0]][2]+offsetZ)*r.scaleFactor
	for i := 2; i < len(prim.VertexIndices); i++ {
		newTriangle := &triangle{
			coords: [3][3]float64{
				{zerox, zeroy, zeroz},
				{
					(obj.Vertices[prim.VertexIndices[i-1]][0] + offsetX) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i-1]][1] + offsetY) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i-1]][2] + offsetZ) * r.scaleFactor,
				},
				{
					(obj.Vertices[prim.VertexIndices[i]][0] + offsetX) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i]][1] + offsetY) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i]][2] + offsetZ) * r.scaleFactor,
				},
			},
		}
		if len(prim.UVCoordinatesPerIndex) > 0 {
			newTriangle.uvCoords = [3][2]float64{
				prim.UVCoordinatesPerIndex[0],
				prim.UVCoordinatesPerIndex[i-1],
				prim.UVCoordinatesPerIndex[i],
			}
			newTriangle.texture = prim.Texture
		} else {
			newTriangle.colorPaletteIndex = prim.Color
		}
		newTriangle.rotate(r.frame)
		newTriangle.calcMiddle()
		r.trianglesBatch = append(r.trianglesBatch, newTriangle)
	}
}

func (r *RaylibRenderer) DrawTrianglesBatch() {
	// first of all, sort the triangles
	// Disabled (unneeded?)
	if false {
		sort.Slice(r.trianglesBatch, func(x, y int) bool {
			mz1 := r.trianglesBatch[x].middleZ
			mz2 := r.trianglesBatch[y].middleZ
			return mz2 < mz1
		})
	}

	// draw the sorted triangles
	var projX0, projY0, projX1, projY1, projX2, projY2 int32
	for _, t := range r.trianglesBatch {
		projX0, projY0 = obliqueProjectionInt32(t.coords[0][0], t.coords[0][1], t.coords[0][2])
		projX1, projY1 = obliqueProjectionInt32(t.coords[1][0], t.coords[1][1], t.coords[1][2])
		projX2, projY2 = obliqueProjectionInt32(t.coords[2][0], t.coords[2][1], t.coords[2][2])
		if t.texture == nil {
			r.drawFilledTriangle(
				projX0+r.onScreenOffX,
				projY0+r.onScreenOffY,
				projX1+r.onScreenOffX,
				projY1+r.onScreenOffY,
				projX2+r.onScreenOffX,
				projY2+r.onScreenOffY,
				t.coords[0][1],
				t.coords[1][1],
				t.coords[2][1],
				t.colorPaletteIndex,
			)
		} else {
			r.drawTexturedTriangle(
				projX0+r.onScreenOffX,
				projY0+r.onScreenOffY,
				projX1+r.onScreenOffX,
				projY1+r.onScreenOffY,
				projX2+r.onScreenOffX,
				projY2+r.onScreenOffY,
				t.coords[0][1],
				t.coords[1][1],
				t.coords[2][1],
				t.uvCoords[0][0],
				t.uvCoords[1][0],
				t.uvCoords[2][0],
				t.uvCoords[0][1],
				t.uvCoords[1][1],
				t.uvCoords[2][1],
				t.texture,
			)
		}
	}
}
