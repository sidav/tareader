package renderer

import (
	. "totala_reader/geometry/matrix4x4"
	graphicadapter "totala_reader/graphic_adapter"
	. "totala_reader/model"
)

// Renders OBJECT (the stuff with matrices), not just the model
type Renderer struct {
	gAdapter graphicadapter.GraphicBackend

	zBuffer [][]float64

	frame int

	onScreenOffX, onScreenOffY int32

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
		} else if len(p.VertexIndices) > 2 {
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
	zeroOsx, zeroOsy := obliqueProjectionInt32Arr(zeroCrds)
	ost := &onScreenTriangle{}
	for i := 0; i < len(prim.VertexIndices); i++ {
		bIndex := (i + 1) % 4
		rotatedCoordsA := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i]])
		rotatedCoordsB := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[bIndex]])

		osxa, osya := obliqueProjectionInt32Arr(rotatedCoordsA)
		osxb, osyb := obliqueProjectionInt32Arr(rotatedCoordsB)

		ost.x0 = zeroOsx + r.onScreenOffX
		ost.y0 = zeroOsy + r.onScreenOffY
		ost.z0 = zeroCrds[1]
		ost.x1 = osxa + r.onScreenOffX
		ost.y1 = osya + r.onScreenOffY
		ost.z1 = rotatedCoordsA[1]
		ost.x2 = osxb + r.onScreenOffX
		ost.y2 = osyb + r.onScreenOffY
		ost.z2 = rotatedCoordsB[1]

		if len(prim.UVCoordinatesPerIndex) > 0 {
			ost.u0 = prim.CenterUVCoords[0]
			ost.u1 = prim.UVCoordinatesPerIndex[i][0]
			ost.u2 = prim.UVCoordinatesPerIndex[bIndex][0]
			ost.v0 = prim.CenterUVCoords[1]
			ost.v1 = prim.UVCoordinatesPerIndex[i][1]
			ost.v2 = prim.UVCoordinatesPerIndex[bIndex][1]
			ost.texture = prim.Texture
			r.drawRasterizedTexturedTriangle(ost)
		} else {
			ost.color = prim.Color
			r.drawRasterizedFilledTriangle(ost)
		}
	}
}

func (r *Renderer) drawNonquadPrimitive(currWrldMtrx *Matrix4x4, mdl *Model, prim *ModelSurface) {
	zeroCrds := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[0]])
	zeroOsx, zeroOsy := obliqueProjectionInt32Arr(zeroCrds)
	ost := &onScreenTriangle{}
	for i := 2; i < len(prim.VertexIndices); i++ {
		rotatedCoordsA := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i-1]])
		rotatedCoordsB := currWrldMtrx.MultiplyByArr3Vector(mdl.Vertices[prim.VertexIndices[i]])
		osxa, osya := obliqueProjectionInt32Arr(rotatedCoordsA)
		osxb, osyb := obliqueProjectionInt32Arr(rotatedCoordsB)

		ost.x0 = zeroOsx + r.onScreenOffX
		ost.y0 = zeroOsy + r.onScreenOffY
		ost.z0 = zeroCrds[1]
		ost.x1 = osxa + r.onScreenOffX
		ost.y1 = osya + r.onScreenOffY
		ost.z1 = rotatedCoordsA[1]
		ost.x2 = osxb + r.onScreenOffX
		ost.y2 = osyb + r.onScreenOffY
		ost.z2 = rotatedCoordsB[1]

		if len(prim.UVCoordinatesPerIndex) > 0 {
			ost.u0 = prim.CenterUVCoords[0]
			ost.u1 = prim.UVCoordinatesPerIndex[i-1][0]
			ost.u2 = prim.UVCoordinatesPerIndex[i][0]
			ost.v0 = prim.CenterUVCoords[1]
			ost.v1 = prim.UVCoordinatesPerIndex[i-1][1]
			ost.v2 = prim.UVCoordinatesPerIndex[i][1]
			ost.texture = prim.Texture
			r.drawRasterizedTexturedTriangle(ost)
		} else {
			ost.color = prim.Color
			r.drawRasterizedFilledTriangle(ost)
		}
	}
}
