package model

import (
	"fmt"
	"time"
)

func print(text string, args ...interface{}) {
	fmt.Printf(text, args...)
}

func sprint(text string, args ...interface{}) string {
	return fmt.Sprintf(text, args...)
}

func printwait(delayms int, text string, args ...interface{}) {
	print(text, args...)
	time.Sleep(time.Duration(delayms) * time.Millisecond)
}

func cobPanic(text string, args ...interface{}) {
	panic(sprint("COB MACHINE ERROR: "+text, args...))
}
