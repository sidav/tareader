package graphicadapter

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibBackend struct {
	currColor        color.RGBA
	TargetTexture    rl.RenderTexture2D
	realW, realH     int32
	renderW, renderH int32
}

func (rb *RaylibBackend) Init(w, h int32) {
	rl.InitWindow(w, h, "RENDERER")
	rb.realW, rb.realH = w, h
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyEscape)
}

func (rb *RaylibBackend) SetInternalResolution(w, h int32) {
	rb.TargetTexture = rl.LoadRenderTexture(w, h)
	rb.renderW, rb.renderH = w, h
	rl.SetTextureFilter(rb.TargetTexture.Texture, rl.FilterAnisotropic16x)
}

func (rb *RaylibBackend) GetRealResolution() (int32, int32) {
	return rb.realW, rb.realH
}

func (rb *RaylibBackend) GetRenderResolution() (int32, int32) {
	return rb.renderW, rb.renderH
}

func (rb *RaylibBackend) Clear() {
	rl.ClearBackground(rl.Black)
}

func (rb *RaylibBackend) BeginFrame() {
	rl.BeginTextureMode(rb.TargetTexture)
}

func (rb *RaylibBackend) EndFrame() {
	rl.EndTextureMode()
}

func (rb *RaylibBackend) Flush() {
	rl.BeginDrawing()
	rl.DrawTexturePro(rb.TargetTexture.Texture, rl.Rectangle{
		X:      0,
		Y:      float32(rb.TargetTexture.Texture.Height),
		Width:  float32(rb.TargetTexture.Texture.Width),
		Height: float32(-rb.TargetTexture.Texture.Height),
	},
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(rl.GetScreenWidth()),
			Height: float32(rl.GetScreenHeight()),
		},
		rl.Vector2{
			X: 0,
			Y: 0,
		},
		0,
		color.RGBA{255, 255, 255, 255})
	rl.EndDrawing()
}

func (rb *RaylibBackend) SetColor(r, g, b uint8) {
	rb.currColor.R = r
	rb.currColor.G = g
	rb.currColor.B = b
	rb.currColor.A = 255
	//currColor = color.RGBA{
	//	R: r,
	//	G: g,
	//	B: b,
	//	A: 255,
	//}
}

func (rb *RaylibBackend) DrawPoint(x, y int32) {
	rl.DrawPixel(x, y, rb.currColor)
}

func (rb *RaylibBackend) FillRect(x, y, w, h int) {
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rb.currColor)
}

func (rb *RaylibBackend) VerticalLine(x, y0, y1 int) {
	rl.DrawLine(int32(x), int32(y0), int32(x), int32(y1), rb.currColor)
}

func (rb *RaylibBackend) LoadImageAsRlTexture(filename string) rl.Texture2D {
	// file, _ := os.Open(filename)
	// img, _ := png.Decode(file)
	// file.Close()
	// rlImg := rl.NewImageFromImage(img)
	rlImg := rl.LoadImage(filename)
	return rl.LoadTextureFromImage(rlImg)
}

func (rb *RaylibBackend) DrawRlTextureAt(tex rl.Texture2D, leftx, topy int32) {
	rl.DrawTexture(tex, leftx, topy, rl.RayWhite)
}
