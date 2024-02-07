package renderer

import (
	"fmt"
	"time"
)

func (r *Renderer) DebugFlush() {
	if r.debugMode {
		r.gAdapter.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (r *Renderer) DebugPrint(text string, args ...interface{}) {
	if r.debugMode {
		fmt.Printf(text, args...)
	}
}
