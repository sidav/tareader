package scripts

import (
	"encoding/binary"
	"fmt"
)

func print(text string, args ...interface{}) {
	fmt.Printf(text, args...)
}

func sprint(text string, args ...interface{}) string {
	return fmt.Sprintf(text, args...)
}

func sprintInt32AsBigEndianHex(val int32) string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(val))
	returned := ""
	for i := range bs {
		returned = returned + fmt.Sprintf("%02X", bs[i])
		if i < len(bs)-1 {
			returned += " "
		}
	}
	return returned
}
