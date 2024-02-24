package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"time"
	graphicadapter "totala_reader/graphic_adapter"
	"totala_reader/model"
	"totala_reader/renderer"
	binaryreader "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/scripts"
	"totala_reader/ta_files_read/texture"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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
	rend := renderer.Renderer{}
	rend.Init(gAdapter)

	textures := readAllGAFsFromDirectory("game_files/files_gaf")

	if strings.Contains(openedFile, ".3do") {

		modelInTAFormat := object3d.ReadObjectFromReader(r, 0)
		fmt.Printf("{\n%s}\n", modelInTAFormat.ToString(0))

		mdl := model.NewModelFrom3doObject3d(modelInTAFormat, textures)
		object := model.CreateObjectFromModel(mdl)
		object.Print(0)

		for i := 0; i < 3; i++ {
			gAdapter.BeginFrame()
			gAdapter.Clear()
			rend.ShowPalette()
			gAdapter.EndFrame()
			gAdapter.Flush()
			time.Sleep(time.Second / 10)
		}
		var totalDuration time.Duration
		totalFrames := 0

		for !rl.IsKeyDown(rl.KeyEscape) {
			start := time.Now()
			// gAdapter.BeginFrame()
			rend.RenderObject(object)
			// gAdapter.EndFrame()
			totalFrames++
			timeSince := time.Since(start)
			totalDuration += timeSince
			pp("Total frames %d; current done in %v ~> %d FPS (mean %v ~> %d FPS)",
				totalFrames, timeSince, int(time.Second/timeSince),
				totalDuration/time.Duration(totalFrames), int(time.Second/(totalDuration/time.Duration(totalFrames))))
			gAdapter.Flush()
			timeSince = time.Since(start)
			pp("                         with flush: %v ~> %d FPS",
				timeSince, int(time.Second/timeSince))
			// time.Sleep(time.Microsecond)
			// gAdapter.Clear()
		}
	}
	if strings.Contains(strings.ToLower(openedFile), ".gaf") {
		fmt.Printf("Opening texture\n")
		gafEntries := texture.ReadTextureFromReader(r)
		for _, ge := range gafEntries {
			gAdapter.Clear()
			rend.DrawGafFrame(ge)
			gAdapter.Flush()
			fmt.Printf("%s\n", ge.Name)
			time.Sleep(1 * time.Second / 2)
		}
	}
	if strings.Contains(strings.ToLower(openedFile), ".cob") {
		fmt.Printf("Opening script %s\n", openedFile)
		scripts.ReadCobFileFromReader(r)
	}
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
