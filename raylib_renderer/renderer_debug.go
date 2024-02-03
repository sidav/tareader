package raylibrenderer

import (
	"fmt"
	"time"
)

func (r *RaylibRenderer) DebugFlush() {
	if r.debugMode {
		r.gAdapter.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (r *RaylibRenderer) DebugPrint(text string, args ...interface{}) {
	if r.debugMode {
		fmt.Printf(text, args...)
	}
}
