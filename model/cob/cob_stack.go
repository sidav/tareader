package cob

import "fmt"

// const stackMaxDepth = 16
const stackMaxDepth = 64 // TEMPORARY VALUE!!!

type Stack struct {
	data      [stackMaxDepth]int32
	stackSize int
}

func (cs *Stack) reset() {
	cs.stackSize = 0
}

func (cs *Stack) Push(word int32) {
	if cs.stackSize == stackMaxDepth {
		fmt.Print("ERROR: Stack overflow at --> ")
		return
	}
	cs.data[cs.stackSize] = word
	cs.stackSize++
}

func (cs *Stack) PushBool(val bool) {
	var setVar int32
	if val {
		setVar = 1
	}
	cs.Push(setVar)
}

func (cs *Stack) PopWord() int32 {
	if cs.stackSize == 0 {
		panic("COB VPU ERROR: attempting to pop from empty stack.")
	}
	cs.stackSize--
	return cs.data[cs.stackSize]
}

// Deprecated except for debug output.
func (cs *Stack) Peek() int32 {
	return cs.data[cs.stackSize-1]
}
