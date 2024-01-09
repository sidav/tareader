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
	return int(binary.LittleEndian.Uint32(mr.fileBytes[baseOffset+offset : baseOffset+offset+4]))
}

func (mr *Reader) ReadNullTermStringFromBytesArray(baseOffset, offset int) string {
	buff := bytes.NewBufferString("")
	index := 0
	for index < 10000 {
		byteHere := mr.fileBytes[baseOffset+offset+index]
		if byteHere == 0x00 {
			return buff.String()
		}
		buff.WriteByte(byteHere)
		index++
	}
	panic("No end on N-T string in sight!")
}
