package cob

const stackMaxDepth = 16

type Stack struct {
	data      [stackMaxDepth]int32
	stackSize int
}

func (cs *Stack) Push(word int32) {
	cs.data[cs.stackSize] = word
	cs.stackSize++
}

func (cs *Stack) PopWord() int32 {
	cs.stackSize--
	return cs.data[cs.stackSize]
}
