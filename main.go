package main

import (
	"fmt"
	"os"
	"time"
	"totala_reader/model"
	raylibrenderer "totala_reader/raylib_renderer"
	"totala_reader/raylib_renderer/middleware"
	binaryreader "totala_reader/ta_files_read"
	"totala_reader/ta_files_read/object3d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	openedFile := "armsy.3do"
	if len(os.Args) > 1 {
		openedFile = os.Args[1]
	}
	r := &binaryreader.Reader{}
	r.ReadFromFile(openedFile)

	obj := object3d.ReadObjectFromReader(r, 0)
	fmt.Printf("{\n%s}\n", obj.ToString(0))

	model := model.NewModelFrom3doObject3d(obj)
	middleware.InitMiddleware(1366, 768)
	defer rl.CloseWindow()
	rend := raylibrenderer.RaylibRenderer{}
	rend.Init()

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

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
