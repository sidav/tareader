package main

import (
	"flag"
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
	var doProfiling bool
	var exportGAFtoPNG bool
	flag.BoolVar(&doProfiling, "p", false, "Perform CPU pprof recording")
	flag.BoolVar(&exportGAFtoPNG, "export", false, "Export GAF images as PNGs if GAF is opened")
	flag.Parse()

	if doProfiling {
		pp("Enabling the profiler.")
		f, err := os.Create("cpu.pprof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var openedFile string
	// Workaround for go flags unable to get along with os.Args
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		openedFile = arg
	}
	if openedFile == "" {
		pp("No 3do/gaf/cob file provided, exiting.")
		return
	}

	r := &binaryreader.Reader{}
	r.ReadFromFile(openedFile)

	if strings.Contains(strings.ToLower(openedFile), ".cob") {
		fmt.Printf("Disassembling the script %s\n", openedFile)
		scripts.ReadCobFileFromReader(r)
		return
	}

	gAdapter := &graphicadapter.RaylibBackend{}
	var scale int32 = 1
	gAdapter.Init(1366, 768)
	gAdapter.SetInternalResolution(1366/scale, 768/scale)

	if strings.Contains(openedFile, ".3do") {
		onlyShow3doModel(r, gAdapter)
		return
	}
	if strings.Contains(strings.ToLower(openedFile), ".gaf") {
		onlyShowGafContents(r, gAdapter, exportGAFtoPNG)
		return
	}
}

func onlyShow3doModel(r *binaryreader.Reader, gAdapter graphicadapter.GraphicBackend) {
	textures := readAllGAFsFromDirectory("game_files/files_gaf")
	rend := renderer.Renderer{}
	rend.Init(gAdapter)
	pp("Opening model %s\n", r.FileName)
	modelInTAFormat := object3d.ReadObjectFromReader(r, 0)
	pp("{\n%s}", modelInTAFormat.ToString(0))

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
	}
}

func onlyShowGafContents(r *binaryreader.Reader, gAdapter graphicadapter.GraphicBackend, export bool) {
	pp("Opening texture %s\n", r.FileName)
	gafEntries := texture.ReadTextureFromReader(r, true)
	if export {
		for _, e := range gafEntries {
			e.Export()
		}
	} else {
		rend := renderer.Renderer{}
		rend.Init(gAdapter)
		for _, ge := range gafEntries {
			gAdapter.Clear()
			rend.DrawGafFrame(ge)
			gAdapter.Flush()
			fmt.Printf("%s\n", ge.Name)
			time.Sleep(1 * time.Second / 2)
		}
	}
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
