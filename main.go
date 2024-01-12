package main

import (
	"fmt"
	"os"
	binaryreader "totala_reader/binary_reader"
	"totala_reader/object3d"
)

func main() {
	readModel := "armsy.3do"
	if len(os.Args) > 1 {
		readModel = os.Args[1]
	}
	r := &binaryreader.Reader{}
	r.ReadFromFile(readModel)

	obj := object3d.ReadObjectFromReader(r, 0)
	fmt.Printf("{\n%s}\n", obj.ToString(0))
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
