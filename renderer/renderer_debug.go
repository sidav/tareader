package renderer

import (
	"fmt"
	"time"
)

func (r *ModelRenderer) DebugFlush() {
	if r.debugMode {
		r.gAdapter.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (r *ModelRenderer) DebugPrint(text string, args ...interface{}) {
	if r.debugMode {
		fmt.Printf(text, args...)
	}
}
