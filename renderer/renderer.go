package renderer

import (
	. "totala_reader/geometry/matrix4x4"
	graphicadapter "totala_reader/graphic_adapter"
	. "totala_reader/model"
)

// Renders OBJECT (the stuff with matrices), not just the model
type Renderer struct {
	gAdapter                   graphicadapter.GraphicBackend
	onScreenOffX, onScreenOffY int32

	frame int

	zBuffer                                [][]float64
	zBufMinX, zBufMaxX, zBufMinY, zBufMaxY int32 // for clearing changed area of the buffer only

	debugMode bool
}

func (r *Renderer) Init(adapter graphicadapter.GraphicBackend) {
	r.gAdapter = adapter
	r.onScreenOffX, r.onScreenOffY = r.gAdapter.GetRenderResolution()
	r.onScreenOffX /= 2
	r.onScreenOffY = 55 * r.onScreenOffY / 100
	r.gAdapter.Clear()
	r.initZBuffer()
}

func (r *Renderer) RenderObject(rootObject *ModelledObject) {
	r.clearZBuffer()
	r.gAdapter.Clear()
	r.drawModelledObject(rootObject, nil)
	r.frame++
}

func (r *Renderer) drawModelledObject(currObj *ModelledObject, parentWrldMtrx *Matrix4x4) {
	mdl := currObj.ModelForObject

	// Dummy animation code, safe to delete
	const rotateByDegrees = 1.0
	if currObj.Name == "base" || currObj.Name == "ground" {
		currObj.Matrix.RotateAroundY(rotateByDegrees * 3.14159265 / 180)
	}
	if currObj.Name == "charged" || currObj.Name == "fork1" {
		currObj.Matrix.RotateAroundX(-3 * rotateByDegrees * 3.14159265 / 180)
	}
	if currObj.Name == "luparm" || currObj.Name == "fork2" {
		currObj.Matrix.RotateAroundX(3 * rotateByDegrees * 3.14159265 / 180)
	}
	// Dummy animation code ended

	currWrldMtrx := currObj.Matrix.Dup() // it will be transformed to world one later, it's a dup just for now
	if parentWrldMtrx != nil {
		currWrldMtrx = parentWrldMtrx.MultiplyByMatrix4x4(currWrldMtrx)
		ox, oy, oz := parentWrldMtrx.MultiplyByXyz0Vector(mdl.XFromParent, mdl.YFromParent, mdl.ZFromParent)
		currWrldMtrx.Translate(ox, oy, oz)
	}

	// Render primitives
	if mdl.SelectionPrimitive != nil {
		r.drawSelectionPrimitive(currWrldMtrx, mdl, mdl.SelectionPrimitive)
	}
	for _, p := range mdl.Primitives {
		// Back-face culling based on normals' directions
		rotatedNormal := currWrldMtrx.MultiplyByVectorW0(p.NormalVector)
		// This expression is a dot product of the surface normal and (0; 0.894427; -0.447214) vector.
		// These "magic numbers" are vector-to-surface coords for oblique projection.
		// The values are gotten from the equation: z - y/2 = 0;
		// Y and Z are the values which give the normalized (length of 1.0) vector.
		if rotatedNormal.Y*0.894427-rotatedNormal.Z*0.447214 > 0 {
			continue
		}
		// back-face culling ended

		if mdl.SelectionPrimitive == p {
			continue
		} else if len(p.VertexIndices) == 4 {
			r.drawQuadPrimitive(currWrldMtrx, mdl, p)
		} else {
			r.drawNonquadPrimitive(currWrldMtrx, mdl, p)
		}
	}

	if currObj.Child != nil { // && len(mdl.ChildObject.Primitives) > 0 {
		r.drawModelledObject(currObj.Child, currWrldMtrx)
	}
	if currObj.Sibling != nil { // && len(mdl.SiblingObject.Primitives) > 0 {
		r.drawModelledObject(currObj.Sibling, parentWrldMtrx)
	}
}

func (r *Renderer) drawSelectionPrimitive(objWorldMatrix *Matrix4x4, obj *Model, prim *ModelSurface) {
	for i := 0; i < len(prim.VertexIndices); i++ {
		x1, y1, z1 := objWorldMatrix.MultiplyByXyz1Vector(
			(obj.Vertices[prim.VertexIndices[i]][0]),
			(obj.Vertices[prim.VertexIndices[i]][1]),
			(obj.Vertices[prim.VertexIndices[i]][2]),
		)
		x2, y2, z2 := objWorldMatrix.MultiplyByXyz1Vector(
			(obj.Vertices[prim.VertexIndices[(i+1)%len(prim.VertexIndices)]][0]),
			(obj.Vertices[prim.VertexIndices[(i+1)%len(prim.VertexIndices)]][1]),
			(obj.Vertices[prim.VertexIndices[(i+1)%len(prim.VertexIndices)]][2]),
		)

		px1, py1 := obliqueProjectionInt32(x1, y1, z1)
		px2, py2 := obliqueProjectionInt32(x2, y2, z2)
		r.gAdapter.SetColor(getTaPaletteColor(2))
		r.DrawLine(px1+r.onScreenOffX, py1+r.onScreenOffY, px2+r.onScreenOffX, py2+r.onScreenOffY)
	}
}

// Separate routine needed because trapezoids DON'T texture properly
// So we need separate triangulation (quad is split to 4 triangles, each has quad's center as a vertex)
func (r *Renderer) drawQuadPrimitive(currWrldMtrx *Matrix4x4, mdl *Model, prim *ModelSurface) {
	zeroCrds := currWrldMtrx.MultiplyByArr3Vector(prim.CenterCoords)
	for i := 0; i < len(prim.VertexIndices); i++ {
		newTriangle := &triangle{
			coords: [3][3]float64{
				zeroCrds,
				currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i]]),
				currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[(i+1)%4]]),
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
		r.Draw3dTriangleStruct(newTriangle)
	}
}

func (r *Renderer) drawNonquadPrimitive(currWrldMtrx *Matrix4x4, mdl *Model, prim *ModelSurface) {
	if len(prim.VertexIndices) < 3 {
		return
	}

	zeroCrds := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[0]])
	for i := 2; i < len(prim.VertexIndices); i++ {
		newTriangle := &triangle{
			coords: [3][3]float64{
				zeroCrds,
				currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i-1]]),
				currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i]]),
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
		r.Draw3dTriangleStruct(newTriangle)
	}
}

func (r *Renderer) Draw3dTriangleStruct(t *triangle) {
	projX0, projY0 := obliqueProjectionInt32(t.coords[0][0], t.coords[0][1], t.coords[0][2])
	projX1, projY1 := obliqueProjectionInt32(t.coords[1][0], t.coords[1][1], t.coords[1][2])
	projX2, projY2 := obliqueProjectionInt32(t.coords[2][0], t.coords[2][1], t.coords[2][2])

	/*  REDUNDANT while normal-based back-face culling is there. May be more useful if moved before triangulation.

	// Back-face culling based on the on-screen vertex draw order
	x10, y10 := projX0-projX1, projY0-projY1
	x12, y12 := projX2-projX1, projY2-projY1
	// If clockwise (determinant > 0) or collinear (determinant == 0), skip this triangle
	if x10*y12-x12*y10 >= 0 {
		return
	}
	// Back-face culling ended

	*/

	if t.texture == nil {
		r.drawRasterizedFilledTriangle(
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
		r.drawRasterizedTexturedTriangle(
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
