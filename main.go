package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	graphicadapter "totala_reader/graphic_adapter"
	"totala_reader/model"
	raylibrenderer "totala_reader/raylib_renderer"
	binaryreader "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/texture"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	openedFile := "game_files/files_3do/armlab.3do"
	if len(os.Args) > 1 {
		openedFile = os.Args[1]
	}
	r := &binaryreader.Reader{}
	r.ReadFromFile(openedFile)

	gAdapter := &graphicadapter.RaylibBackend{}
	var scale int32 = 1
	gAdapter.Init(1366, 768)
	gAdapter.SetInternalResolution(1366/scale, 768/scale)
	rend := raylibrenderer.RaylibRenderer{}
	rend.Init(gAdapter)

	textures := readAllGAFsFromDirectory("game_files/files_gaf")

	if strings.Contains(openedFile, ".3do") {

		obj := object3d.ReadObjectFromReader(r, 0)
		fmt.Printf("{\n%s}\n", obj.ToString(0))

		model := model.NewModelFrom3doObject3d(obj, textures)

		for i := 0; i < 10; i++ {
			gAdapter.BeginFrame()
			rend.ShowPalette()
			gAdapter.EndFrame()
			gAdapter.Flush()
			time.Sleep(time.Second / 10)
		}
		var totalDuration time.Duration
		totalFrames := 0

		for !rl.IsKeyDown(rl.KeyEscape) {
			start := time.Now()
			gAdapter.BeginFrame()
			rend.DrawModel(model)
			gAdapter.EndFrame()
			totalDuration += time.Since(start)
			totalFrames++
			pp("Total frames %d; current done in %v (mean %v ~> %d FPS)",
				totalFrames, time.Since(start),
				totalDuration/time.Duration(totalFrames), int(time.Second/(totalDuration/time.Duration(totalFrames))))
			gAdapter.Flush()
			// time.Sleep(time.Microsecond)
			// gAdapter.Clear()
		}
	}
	if strings.Contains(strings.ToLower(openedFile), ".gaf") {
		fmt.Printf("Opening texture\n")
		gafEntries := texture.ReadTextureFromReader(r)
		for _, ge := range gafEntries {
			rend.DrawGafFrame(ge)
			gAdapter.Flush()
			time.Sleep(1 * time.Second / 10)
		}
	}
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
