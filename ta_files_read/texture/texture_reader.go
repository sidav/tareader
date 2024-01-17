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
	for index, off := range entryPointers {
		frames := r.ReadUint16FromBytesArray(off, 0)
		always1 := r.ReadUint16FromBytesArray(off, 2)
		always0 = r.ReadIntFromBytesArray(off, 4)
		name := r.ReadFixedLengthStringFromBytesArray(off, 8, 32)
		fmt.Printf("GAF entry #%d at offset %d:\n", index, off)
		fmt.Printf("  Name \"%s\", %d frames, %d should be one, %d should be zero\n",
			name, frames, always1, always0)

		// Read each GAF frame entry for GAF entry
		for gfe := 0; gfe < frames; gfe++ {
			ptrFrameEntry := r.ReadIntFromBytesArray(off, 40)
			unknown := r.ReadIntFromBytesArray(off, 44)
			fmt.Printf("    GAF frame entry %d:\n", gfe)
			fmt.Printf("      Pointer to the data: %d, unknown value: %d\n", ptrFrameEntry, unknown)
			readGafFrameData(r, ptrFrameEntry)
		}
	}
}

func readGafFrameData(r *tafilesread.Reader, offset int) {
	width := r.ReadUint16FromBytesArray(offset, 0)
	height := r.ReadUint16FromBytesArray(offset, 2)
	xPos := r.ReadUint16FromBytesArray(offset, 4)
	yPos := r.ReadUint16FromBytesArray(offset, 6)
	unknownByte := r.ReadByteFromBytesArray(offset, 8)
	compressed := r.ReadByteFromBytesArray(offset, 9) != 0
	framePointers := r.ReadUint16FromBytesArray(offset, 10)
	unknown2 := r.ReadIntFromBytesArray(offset, 12)
	ptrFrameData := r.ReadIntFromBytesArray(offset, 16)
	unknown3 := r.ReadIntFromBytesArray(offset, 20)
	fmt.Printf("      GAF Frame Data: \n")
	fmt.Printf("        %dx%dpx, xPos %d, yPos %d\n", width, height, xPos, yPos)
	fmt.Printf("        Unknown1: %d, compressed %v\n", unknownByte, compressed)
	fmt.Printf("        Frame pointers: %d, Unknown2 %d\n", framePointers, unknown2)
	fmt.Printf("        PtrFrameData: %d, Unknown3 %d\n", ptrFrameData, unknown3)
	// read the raw data itself
}
