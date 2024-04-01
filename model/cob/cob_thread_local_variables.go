package cob

const (
	maxLocalVars = 8 // TODO: find out exact needed number
)

func (ct *CobThread) DoAllocNewLocalVar() {
	ct.currentScopeLVarsCount++

	// TODO: remove the following after debug
	if ct.currentScopeLVarZeroIndex+ct.currentScopeLVarsCount > int32(len(ct.LVars)) {
		panic("COB VPU: nowhere to alloc local variable")
	}
}

// Is used in conjuction with CALL and/or START-SCRIPT, should be used in OUTER scope
func (ct *CobThread) setParamsFromStack(paramsCount int32) {
	nextScopeRelativeZero := ct.currentScopeLVarZeroIndex + ct.currentScopeLVarsCount
	// Parameters seem to be passed in order in which they were pushed onto the stack
	for i := paramsCount - 1; i >= 0; i-- {
		ct.LVars[nextScopeRelativeZero+i] = ct.DataStack.PopWord()
	}
}

func (ct *CobThread) GetCurrentScopeLocalVar(varNum int32) int32 {
	return ct.LVars[ct.currentScopeLVarZeroIndex+varNum]
}

func (ct *CobThread) SetCurrentScopeLocalVar(value, varNum int32) {
	ct.LVars[ct.currentScopeLVarZeroIndex+varNum] = value
}

// func (ct *CobThread) DoPushFromLocal(varNum int32) {
// 	ct.DataStack.Push(ct.LVars[ct.currentScopeLVarZeroIndex+varNum])
// }

// func (ct *CobThread) DoPopToLocal(varNum int32) {
// 	ct.LVars[ct.currentScopeLVarZeroIndex+varNum] = ct.DataStack.PopWord()
// }
