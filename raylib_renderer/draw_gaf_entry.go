package raylibrenderer

import (
	"totala_reader/ta_files_read/texture"
)

// Draw the GAF entry, just to look at it.
func (r *RaylibRenderer) DrawGafFrame(gafEntry *texture.GafEntry) {
	const pixelSize = 3
	offsetFromTop := int32(36)
	for _, f := range gafEntry.Frames {
		for x := range f.Pixels {
			for y := range f.Pixels[x] {
				r.gAdapter.SetColor(getTaPaletteColor(f.Pixels[x][y]))
				r.FillRect(int32(x)*pixelSize, int32(y)*pixelSize+offsetFromTop, pixelSize, pixelSize)
			}
		}
		offsetFromTop += pixelSize*int32(len(f.Pixels[0])) + pixelSize
	}
}
