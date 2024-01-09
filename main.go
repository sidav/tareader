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
	showAllChildsAndSiblings(r, obj, 0)
	pp("Texture names: %s", r.ReadNullTermStringFromBytesArray(0, 52))
}

func showAllChildsAndSiblings(r *binaryreader.Reader, obj *object3d.Object, recursionDepth int) {
	obj.Print(recursionDepth * 2)
	if obj.OffsetToChildObject != 0 {
		pp("Child (depth %d) at %d: ", recursionDepth, obj.OffsetToChildObject)
		showAllChildsAndSiblings(r, object3d.ReadObjectFromReader(r, obj.OffsetToChildObject), recursionDepth+1)
	}
	if obj.OffsetToSiblingObject != 0 {
		pp("Sibling (depth %d) at %d: ", recursionDepth, obj.OffsetToSiblingObject)
		showAllChildsAndSiblings(r, object3d.ReadObjectFromReader(r, obj.OffsetToSiblingObject), recursionDepth+1)
	}
}

func pp(str string, args ...interface{}) {
	fmt.Printf(str+"\n", args...)
}
