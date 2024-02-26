package cob

const (
	maxLocalVars = 8
)

// Each scripted entity can have up to 8 those "threads"
type CobThread struct {
	IP        int32               // Instruction pointer
	LVars     [maxLocalVars]int32 // Local variables
	DataStack Stack               // Data stack
	RetStack  Stack               // Return stack

	SigMask             int32 // Signal mask (needed for thread stop)
	Active              bool
	SleepTicksRemaining int32
}

func (ct *CobThread) reset() {
	ct.DataStack.reset()
	ct.RetStack.reset()
	ct.SleepTicksRemaining = 0
}

func (ct *CobThread) SetSleep(duration int32) {
	ct.SleepTicksRemaining = duration
}
