package raylibrenderer

import (
	"fmt"
	"sort"
	"time"
	"totala_reader/model"
	"totala_reader/raylib_renderer/middleware"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRenderer struct {
	onScreenOffX, onScreenOffY int32
	fontSize                   int32
	scaleFactor                float64
	totalMessages              int32

	zBuffer [1920][1080]float64

	trianglesBatch []*triangle
}

func (r *RaylibRenderer) Init() {
	r.onScreenOffX, r.onScreenOffY = middleware.GetWindowSize()
	r.onScreenOffX /= 2
	r.onScreenOffY /= 2
	r.fontSize = 32
	r.scaleFactor = 4.0
	middleware.Clear()
	middleware.Flush()
}

func (r *RaylibRenderer) DrawModel(rootObject *model.Model) {

	for i := range r.zBuffer {
		for j := range r.zBuffer[i] {
			r.zBuffer[i][j] = -100000.0
		}
	}

	r.trianglesBatch = nil
	r.totalMessages = 0
	middleware.Clear()
	// middleware.Flush()
	middleware.SetColor(getTaPaletteColor(4))
	r.drawObject(rootObject, 0, 0, 0)
	r.DrawTrianglesBatch()
	// middleware.Flush()
}

func (r *RaylibRenderer) drawObject(obj *model.Model, parentOffsetX, parentOffsetY, parentOffsetZ float64) {
	currentOffsetX, currentOffsetY, currentOffsetZ := obj.XFromParent+parentOffsetX,
		obj.YFromParent+parentOffsetY, obj.ZFromParent+parentOffsetZ

	for _, p := range obj.Primitives {
		r.drawPrimitive(obj, p, currentOffsetX, currentOffsetY, currentOffsetZ)
	}

	if false {
		rl.DrawText(fmt.Sprintf("OBJECT: %s\n", obj.ObjectName), 0, (r.fontSize+2)*r.totalMessages, r.fontSize, rl.White)
		rl.DrawLine(0, r.totalMessages*(r.fontSize+2)+r.fontSize,
			340, r.totalMessages*(r.fontSize+2)+r.fontSize, rl.Red)
		rl.DrawLine(340, r.totalMessages*(r.fontSize+2)+r.fontSize,
			int32(currentOffsetX*r.scaleFactor)+r.onScreenOffX, int32(currentOffsetZ*r.scaleFactor)+r.onScreenOffY, rl.Red)
		r.totalMessages++
		// middleware.Flush()
		time.Sleep(time.Second)
	}

	if obj.ChildObject != nil && len(obj.ChildObject.Primitives) > 0 {
		middleware.SetColor(getTaPaletteColor(5))
		r.drawObject(obj.ChildObject, currentOffsetX, currentOffsetY, currentOffsetZ)
	}
	if obj.SiblingObject != nil && len(obj.SiblingObject.Primitives) > 0 {
		middleware.SetColor(getTaPaletteColor(4))
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
		newTriangle.calcMiddle()
		r.trianglesBatch = append(r.trianglesBatch, newTriangle)
	}
}

func (r *RaylibRenderer) DrawTrianglesBatch() {
	// first of all, sort the triangles
	sort.Slice(r.trianglesBatch, func(x, y int) bool {
		my1, mz1 := r.trianglesBatch[x].middleY, r.trianglesBatch[x].middleZ
		my2, mz2 := r.trianglesBatch[y].middleY, r.trianglesBatch[y].middleZ
		return mz2 < mz1
		return (-mz1 - my1/2) < (-mz2 - my2/2)
		// return r.trianglesBatch[x].coords[2][1] < r.trianglesBatch[y].coords[2][1]
	})

	// draw the sorted triangles
	var projX0, projY0, projX1, projY1, projX2, projY2 int32
	for i, t := range r.trianglesBatch {
		middleware.SetColor(getTaPaletteColor(uint8(i%3 + 3)))
		// middleware.SetColor(getTaPaletteColor(4))
		projX0, projY0 = obliqueProjectionInt32(t.coords[0][0], t.coords[0][1], t.coords[0][2])
		projX1, projY1 = obliqueProjectionInt32(t.coords[1][0], t.coords[1][1], t.coords[1][2])
		projX2, projY2 = obliqueProjectionInt32(t.coords[2][0], t.coords[2][1], t.coords[2][2])
		r.polygon(
			projX0+r.onScreenOffX,
			projY0+r.onScreenOffY,
			projX1+r.onScreenOffX,
			projY1+r.onScreenOffY,
			projX2+r.onScreenOffX,
			projY2+r.onScreenOffY,
			t.coords[0][1],
			t.coords[1][1],
			t.coords[2][1],
		)
		// rl.DrawLine(projX0+r.onScreenOffX,
		// 	projY0+r.onScreenOffY,
		// 	projX1+r.onScreenOffX,
		// 	projY1+r.onScreenOffY, rl.White)
		// rl.DrawLine(projX0+r.onScreenOffX,
		// 	projY0+r.onScreenOffY,
		// 	projX2+r.onScreenOffX,
		// 	projY2+r.onScreenOffY, rl.White)
		// rl.DrawLine(projX2+r.onScreenOffX,
		// 	projY2+r.onScreenOffY,
		// 	projX1+r.onScreenOffX,
		// 	projY1+r.onScreenOffY, rl.White)

		// middleware.Flush()
		// time.Sleep(10 * time.Millisecond / 10000)
	}
}
