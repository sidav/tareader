package raylibrenderer

import (
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

	r.totalMessages = 0
	r.gAdapter.Clear()
	r.drawObject(rootObject, 0, 0, 0)
	r.frame++
}

func (r *RaylibRenderer) drawObject(obj *model.Model, parentOffsetX, parentOffsetY, parentOffsetZ float64) {
	currentOffsetX, currentOffsetY, currentOffsetZ := obj.XFromParent+parentOffsetX,
		obj.YFromParent+parentOffsetY, obj.ZFromParent+parentOffsetZ

	for _, p := range obj.Primitives {
		if len(p.VertexIndices) == 4 {
			r.drawQuadPrimitive(obj, p, currentOffsetX, currentOffsetY, currentOffsetZ)
		} else {
			r.drawNonquadPrimitive(obj, p, currentOffsetX, currentOffsetY, currentOffsetZ)
		}
	}

	if obj.ChildObject != nil && len(obj.ChildObject.Primitives) > 0 {
		r.drawObject(obj.ChildObject, currentOffsetX, currentOffsetY, currentOffsetZ)
	}
	if obj.SiblingObject != nil && len(obj.SiblingObject.Primitives) > 0 {
		r.drawObject(obj.SiblingObject, parentOffsetX, parentOffsetY, parentOffsetZ)
	}
}

// Separate routine needed because trapezoids DON'T texture properly
// So we need separate triangulation (quad is split to 4 triangles, each has quad's center as a vertex)
func (r *RaylibRenderer) drawQuadPrimitive(obj *model.Model, prim *model.ModelSurface, offsetX, offsetY, offsetZ float64) {
	if len(prim.VertexIndices) != 4 || obj.SelectionPrimitive == prim {
		return
	}

	zerox, zeroy, zeroz := (prim.CenterCoords[0]+offsetX)*r.scaleFactor,
		(prim.CenterCoords[1]+offsetY)*r.scaleFactor,
		(prim.CenterCoords[2]+offsetZ)*r.scaleFactor
	for i := 0; i < len(prim.VertexIndices); i++ {
		newTriangle := &triangle{
			coords: [3][3]float64{
				{zerox, zeroy, zeroz},
				{
					(obj.Vertices[prim.VertexIndices[i]][0] + offsetX) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i]][1] + offsetY) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[i]][2] + offsetZ) * r.scaleFactor,
				},
				{
					(obj.Vertices[prim.VertexIndices[(i+1)%4]][0] + offsetX) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[(i+1)%4]][1] + offsetY) * r.scaleFactor,
					(obj.Vertices[prim.VertexIndices[(i+1)%4]][2] + offsetZ) * r.scaleFactor,
				},
			},
		}
		if len(prim.UVCoordinatesPerIndex) > 0 {
			newTriangle.uvCoords = [3][2]float64{
				prim.CenterUVCoords,
				prim.UVCoordinatesPerIndex[i],
				prim.UVCoordinatesPerIndex[(i+1)%4],
			}
			newTriangle.texture = prim.Texture
		} else {
			newTriangle.colorPaletteIndex = prim.Color
		}
		newTriangle.rotate(r.frame)
		newTriangle.calcMiddle()
		r.Draw3dTriangleStruct(newTriangle)
	}
}

func (r *RaylibRenderer) drawNonquadPrimitive(obj *model.Model, prim *model.ModelSurface, offsetX, offsetY, offsetZ float64) {
	if len(prim.VertexIndices) < 3 || obj.SelectionPrimitive == prim {
		return
	}

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
		r.Draw3dTriangleStruct(newTriangle)
	}
}

func (r *RaylibRenderer) Draw3dTriangleStruct(t *triangle) {
	projX0, projY0 := obliqueProjectionInt32(t.coords[0][0], t.coords[0][1], t.coords[0][2])
	projX1, projY1 := obliqueProjectionInt32(t.coords[1][0], t.coords[1][1], t.coords[1][2])
	projX2, projY2 := obliqueProjectionInt32(t.coords[2][0], t.coords[2][1], t.coords[2][2])
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
