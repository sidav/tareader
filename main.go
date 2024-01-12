package main

import (
	"fmt"
	binaryreader "totala_reader/binary_reader"
	"totala_reader/object3d"
)

func main() {
	r := &binaryreader.Reader{}
	r.ReadFromFile("armsy.3do")

	obj := object3d.ReadObjectFromReader(r, 0)
	fmt.Printf(obj.ToString(0))
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
