package raylibrenderer

import (
	"fmt"
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
}

func (r *RaylibRenderer) Init() {
	r.onScreenOffX, r.onScreenOffY = middleware.GetWindowSize()
	r.onScreenOffX /= 2
	r.onScreenOffY /= 2
	r.fontSize = 32
	r.scaleFactor = 5.0
	middleware.Clear()
	middleware.Flush()
}

func (r *RaylibRenderer) DrawModel(rootObject *model.Model) {
	r.totalMessages = 0
	middleware.Clear()
	// middleware.Flush()
	middleware.SetColor(getTaPaletteColor(4))
	r.drawObject(rootObject, 0, 0, 0)
}

func (r *RaylibRenderer) drawObject(obj *model.Model, parentOffsetX, parentOffsetY, parentOffsetZ float64) {
	currentOffsetX, currentOffsetY, currentOffsetZ := obj.XFromParent+parentOffsetX,
		obj.YFromParent+parentOffsetY, obj.ZFromParent+parentOffsetZ

	for _, p := range obj.Primitives {
		r.drawPrimitive(obj, p, currentOffsetX, currentOffsetY, currentOffsetZ)
	}

	rl.DrawText(fmt.Sprintf("OBJECT: %s\n", obj.ObjectName), 0, (r.fontSize+2)*r.totalMessages, r.fontSize, rl.White)
	rl.DrawLine(0, r.totalMessages*(r.fontSize+2)+r.fontSize,
		340, r.totalMessages*(r.fontSize+2)+r.fontSize, rl.Red)
	rl.DrawLine(340, r.totalMessages*(r.fontSize+2)+r.fontSize,
		int32(currentOffsetX*r.scaleFactor)+r.onScreenOffX, int32(currentOffsetZ*r.scaleFactor)+r.onScreenOffY, rl.Red)
	r.totalMessages++
	middleware.Flush()
	time.Sleep(time.Second / 2)

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
	projectedCoords := make([][2]int32, len(prim.VertexIndices))
	for i, vInd := range prim.VertexIndices {
		vx, vy, vz := obj.Vertices[vInd][0], obj.Vertices[vInd][1], obj.Vertices[vInd][2]
		vx *= r.scaleFactor
		vy *= r.scaleFactor
		vz *= r.scaleFactor
		vx += offsetX * r.scaleFactor
		vy += offsetY * r.scaleFactor
		vz += offsetZ * r.scaleFactor
		projectedCoords[i][0], projectedCoords[i][1] = obliqueProjectionInt32(vx, vy, vz)

		projectedCoords[i][0] += r.onScreenOffX
		projectedCoords[i][1] += r.onScreenOffY
	}

	if obj.SelectionPrimitive != prim {
		middleware.FillPolygon(projectedCoords)
	}
	for i := range projectedCoords {
		color := rl.White
		if obj.SelectionPrimitive == prim {
			color = rl.Green
		}
		rl.DrawLine(
			projectedCoords[i][0],
			projectedCoords[i][1],
			projectedCoords[(i+1)%len(projectedCoords)][0],
			projectedCoords[(i+1)%len(projectedCoords)][1],
			color,
		)
	}
}
