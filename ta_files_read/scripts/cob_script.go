package scripts

type CobScript struct {
	Pieces             []string
	ProcedureNames     []string
	ProcedureAddresses []int32
	RawCode            []int32
}

// Deprecated
func (cs *CobScript) FindFuncAddressByName(name string) int32 { // Use for debug only. Needs to be more performant for real use
	for i := range cs.ProcedureNames {
		if cs.ProcedureNames[i] == name {
			return cs.ProcedureAddresses[i]
		}
	}
	return -1
}
