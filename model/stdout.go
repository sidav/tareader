package model

import "fmt"

func print(text string, args ...interface{}) {
	fmt.Printf(text, args...)
}

func sprint(text string, args ...interface{}) string {
	return fmt.Sprintf(text, args...)
}

func cobPanic(text string, args ...interface{}) {
	panic(sprint("COB MACHINE ERROR: "+text, args...))
}
