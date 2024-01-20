package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"totala_reader/model"
	raylibrenderer "totala_reader/raylib_renderer"
	"totala_reader/raylib_renderer/middleware"
	binaryreader "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/object3d"
	"totala_reader/ta_files_read/texture"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	openedFile := "armsy.3do"
	if len(os.Args) > 1 {
		openedFile = os.Args[1]
	}
	r := &binaryreader.Reader{}
	r.ReadFromFile(openedFile)

	middleware.InitMiddleware(1366, 768)
	defer rl.CloseWindow()
	rend := raylibrenderer.RaylibRenderer{}
	rend.Init()

	textures := readAllGAFsFromDirectory("game_files/files_gaf")

	if strings.Contains(openedFile, ".3do") {

		obj := object3d.ReadObjectFromReader(r, 0)
		fmt.Printf("{\n%s}\n", obj.ToString(0))

		model := model.NewModelFrom3doObject3d(obj, textures)

		// rend.ShowPalette()
		// rend.ShowPalette()
		// middleware.Flush()
		// time.Sleep(3 * time.Second)

		for !rl.IsKeyDown(rl.KeyEscape) {
			start := time.Now()
			rend.DrawModel(model)
			pp("Done in %v!", time.Since(start))
			middleware.Flush()
			time.Sleep(3 * time.Second / 100)
			middleware.Clear()
		}
	}
	if strings.Contains(strings.ToLower(openedFile), ".gaf") {
		fmt.Printf("Opening texture\n")
		gafEntries := texture.ReadTextureFromReader(r)
		for _, ge := range gafEntries {
			rend.DrawGafFrame(ge)
			middleware.Flush()
			time.Sleep(1 * time.Second / 10)
		}
	}
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
