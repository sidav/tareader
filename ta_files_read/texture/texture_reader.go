package texture

import (
	"fmt"
	tafilesread "totala_reader/ta_files_read"
)

type GafEntry struct {
	Name   string
	Frames []*GafFrame
}

type GafFrame struct {
	Pixels [][]uint8 // each value is an index from palette
}

func ReadTextureFromReader(r *tafilesread.Reader) []*GafEntry {
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

	entriesArray := make([]*GafEntry, entries)

	for index, off := range entryPointers {
		frames := r.ReadUint16FromBytesArray(off, 0)
		always1 := r.ReadUint16FromBytesArray(off, 2)
		always0 = r.ReadIntFromBytesArray(off, 4)
		name := r.ReadFixedLengthStringFromBytesArray(off, 8, 32)
		fmt.Printf("GAF entry #%d at offset %d:\n", index, off)
		fmt.Printf("  Name \"%s\", %d frames, %d should be one, %d should be zero\n",
			name, frames, always1, always0)

		entry := &GafEntry{
			Name: name,
		}
		// Read each GAF frame entry for GAF entry
		for gfe := 0; gfe < frames; gfe++ {
			ptrFrameEntry := r.ReadIntFromBytesArray(off, 40)
			unknown := r.ReadIntFromBytesArray(off, 44)
			fmt.Printf("    GAF frame entry %d:\n", gfe)
			fmt.Printf("      Pointer to the data: %d, unknown value: %d\n", ptrFrameEntry, unknown)
			entry.Frames = append(entry.Frames, readGafFrameData(r, ptrFrameEntry))
		}
		entriesArray[index] = entry
	}
	return entriesArray
}

func readGafFrameData(r *tafilesread.Reader, offset int) *GafFrame {
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

	frame := &GafFrame{}

	// read the raw data itself
	if compressed {
		panic("Compressed frame data reading not implemented yet!")
	} else {
		frame.Pixels = readUncompressedPixels(r, ptrFrameData, width, height)
	}
	return frame
}

func readUncompressedPixels(r *tafilesread.Reader, offset, width, height int) [][]uint8 {
	// Row and column indices (width and height) must be swapped, as the data is written row-by-row and read column-by-column
	var pixels = make([][]uint8, width)
	for i := 0; i < width; i++ {
		pixels[i] = make([]uint8, height)
		for j := 0; j < height; j++ {
			index := i + j*width
			pixels[i][j] = r.ReadByteFromBytesArray(offset, index)
		}
	}
	return pixels
}
