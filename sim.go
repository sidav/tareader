package main

import (
	"time"
	graphicadapter "totala_reader/graphic_adapter"
	"totala_reader/model"
	"totala_reader/renderer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func sim3doModelAndScripts(openedFileName string, gAdapter graphicadapter.GraphicBackend) {
	rend := renderer.Renderer{}
	rend.Init(gAdapter)

	textures := readAllGAFsFromDirectory()
	cavedogModel, cobScript := loadModelAndCobByFilename(openedFileName)
	object := model.SimObject{}
	object.InitFromCavedogData(cavedogModel, textures, cobScript)
	object.ModelState.Print(0)

	var totalDuration time.Duration
	totalFrames := 0

	for !rl.IsKeyDown(rl.KeyEscape) {
		start := time.Now()
		// gAdapter.BeginFrame()
		rend.RenderObject(object.ModelState)
		// gAdapter.EndFrame()
		totalFrames++
		timeSince := time.Since(start)
		if false {
			totalDuration += timeSince
			pp("Total frames %d; current done in %v ~> %d FPS (mean %v ~> %d FPS)",
				totalFrames, timeSince, int(time.Second/timeSince),
				totalDuration/time.Duration(totalFrames), int(time.Second/(totalDuration/time.Duration(totalFrames))))
		}
		gAdapter.Flush()
		if false {
			timeSince = time.Since(start)
			pp("                         with flush: %v ~> %d FPS",
				timeSince, int(time.Second/timeSince))
		}
		object.CobExecAllThreads()
		object.PerformScriptedMovementStep()
	}
}
