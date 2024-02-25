package cob

const (
	maxSVars = 16
	maxLVars = 16
)

// Each scripted entity can have up to 8 those "threads"
type CobThread struct {
	IP        int32           // Instruction pointer
	SVars     [maxSVars]int32 // Static variables - TODO: check if they belong not to the thread, but to the whole entity
	LVars     [maxLVars]int32 // Local variables
	DataStack Stack           // Data stack
	RetStack  Stack           // Return stack

	SigMask int32 // Signal mask (needed for thread stop)
	Active  bool
}
