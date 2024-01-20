package raylibrenderer

import (
	"totala_reader/raylib_renderer/middleware"
	"totala_reader/ta_files_read/texture"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draw the GAF entry, just to look at it.
func (r *RaylibRenderer) DrawGafFrame(gafEntry *texture.GafEntry) {
	const pixelSize = 3
	middleware.Clear()

	rl.DrawText(gafEntry.Name, 0, 0, 32, rl.White)
	offsetFromTop := int32(36)
	for _, f := range gafEntry.Frames {
		for x := range f.Pixels {
			for y := range f.Pixels[x] {
				middleware.SetColor(getTaPaletteColor(f.Pixels[x][y]))
				r.FillRect(int32(x)*pixelSize, int32(y)*pixelSize+offsetFromTop, pixelSize, pixelSize)
			}
		}
		offsetFromTop += pixelSize*int32(len(f.Pixels[0])) + pixelSize
	}
}
