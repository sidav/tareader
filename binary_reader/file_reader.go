package binaryreader

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
)

type Reader struct {
	fileBytes []byte
}

func (mr *Reader) ReadFromFile(fileName string) {
	file, _ := os.Open(fileName)
	defer file.Close()
	// Get the file size
	stat, _ := file.Stat()
	// Read the file into a byte slice
	mr.fileBytes = make([]byte, stat.Size())
	bufio.NewReader(file).Read(mr.fileBytes)

	// fill the reader
}

func (mr *Reader) ReadIntFromBytesArray(baseOffset, offset int) int {
	// fmt.Printf("Reading INT32 at 0x%X (%d+%d): ", baseOffset+offset, baseOffset, offset)
	// fmt.Printf("Got %x\n", mr.fileBytes[baseOffset+offset:baseOffset+offset+4])
	uint32Value := binary.LittleEndian.Uint32(mr.fileBytes[baseOffset+offset : baseOffset+offset+4])
	int32Value := int32(uint32Value)
	return int(int32Value)
}

func (mr *Reader) ReadUint16FromBytesArray(baseOffset, offset int) int {
	// fmt.Printf("Reading UINT16 at 0x%X (%d+%d): ", baseOffset+offset, baseOffset, offset)
	// fmt.Printf("Got %x\n", mr.fileBytes[baseOffset+offset:baseOffset+offset+2])
	uint16Value := binary.LittleEndian.Uint16(mr.fileBytes[baseOffset+offset : baseOffset+offset+2])
	return int(uint16Value)
}

func (mr *Reader) ReadNullTermStringFromBytesArray(baseOffset, offset int) string {
	var buff bytes.Buffer
	index := 0
	for index < 256 {
		byteHere := mr.fileBytes[baseOffset+offset+index]
		if byteHere == 0x00 {
			return buff.String()
		}
		buff.WriteByte(byteHere)
		index++
	}
	if index == 0 {
		return ""
	}
	panic("Null-terminated string longer than 256 bytes!")
}
