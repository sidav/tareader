package cob

// Each COB machine can have up to 8 those "threads"
type CobThread struct {
	IP        int32 // Instruction pointer
	DataStack Stack // Data stack
	CallStack Stack // Call stack

	// Local vars related stuff:
	LVars                     [maxLocalVars]int32 // Local variables container
	currentScopeLVarZeroIndex int32               // Which index is accounted as zero for current scope
	currentScopeLVarsCount    int32

	SigMask             int32 // Signal mask (needed for thread stop)
	Active              bool
	SleepTicksRemaining int32
}

func (ct *CobThread) reset() {
	ct.DataStack.reset()
	ct.CallStack.reset()
	ct.currentScopeLVarsCount = 0
	ct.currentScopeLVarZeroIndex = 0
	ct.SleepTicksRemaining = 0
}

func (ct *CobThread) SetSleep(duration int32) {
	ct.SleepTicksRemaining = duration
}

func (ct *CobThread) DoCall(callAddress int32, params int32) {
	if params > 0 {
		ct.setParamsFromStack(params)
	}
	ct.CallStack.Push(ct.IP)
	ct.CallStack.Push(ct.currentScopeLVarsCount)

	// setting the new local vars scope
	ct.currentScopeLVarZeroIndex += ct.currentScopeLVarsCount
	ct.currentScopeLVarsCount = 0
	// setting return address
	ct.IP = callAddress
}

func (ct *CobThread) DoReturn() {
	if ct.CallStack.stackSize > 0 {
		// restoring the pre-call local vars scope
		ct.currentScopeLVarsCount = ct.CallStack.PopWord()
		ct.currentScopeLVarZeroIndex = ct.currentScopeLVarZeroIndex - ct.currentScopeLVarsCount
		// setting return address
		ct.IP = ct.CallStack.PopWord()
		return
	}
	ct.Active = false
}
