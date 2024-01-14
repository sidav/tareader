package texture

import (
	"fmt"
	tafilesread "totala_reader/ta_files_read"
)

type GafTexture struct {
}

func ReadTextureFromReader(r *tafilesread.Reader) {
	// Reading header
	version := r.ReadIntFromBytesArray(0, 0)
	entries := r.ReadIntFromBytesArray(0, 4)
	always0 := r.ReadIntFromBytesArray(0, 8)
	fmt.Printf("Version %d, entries %d, should be zero %d\n", version, entries, always0)
	var entryPointers []int
	for i := 0; i < entries; i++ {
		entryPointers = append(entryPointers, r.ReadIntFromBytesArray(12, i*4))
	}
	fmt.Printf("Pointers acquired.\n")
	for _, off := range entryPointers {
		frames := r.ReadUint16FromBytesArray(off, 0)
		always1 := r.ReadUint16FromBytesArray(off, 2)
		always0 = r.ReadIntFromBytesArray(off, 4)
		name := r.ReadFixedLengthStringFromBytesArray(off, 8, 32)
		fmt.Printf("GAF entry at offset %d:\n  Name \"%s\", %d frames, %d should be one, %d should be zero\n",
			off, name, frames, always1, always0)
		// Need to read GAF frame entry
	}
}
