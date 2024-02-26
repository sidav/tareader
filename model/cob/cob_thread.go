package cob

const (
	maxLocalVars = 8
)

// Each COB machine can have up to 8 those "threads"
type CobThread struct {
	IP        int32               // Instruction pointer
	LVars     [maxLocalVars]int32 // Local variables
	DataStack Stack               // Data stack
	CallStack Stack               // Call stack

	SigMask             int32 // Signal mask (needed for thread stop)
	Active              bool
	SleepTicksRemaining int32
}

func (ct *CobThread) reset() {
	ct.DataStack.reset()
	ct.CallStack.reset()
	ct.SleepTicksRemaining = 0
}

func (ct *CobThread) SetSleep(duration int32) {
	ct.SleepTicksRemaining = duration
}

func (ct *CobThread) DoCall(callAddress int32, params int32) {
	if params > 0 {
		panic("COB VPU: unimplemented params count")
	}
	ct.CallStack.Push(ct.IP)
	ct.IP = callAddress
}

func (ct *CobThread) DoReturn() {
	if ct.CallStack.stackSize > 0 {
		ct.IP = ct.CallStack.PopWord()
		return
	}
	ct.Active = false
}
