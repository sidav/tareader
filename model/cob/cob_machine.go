package cob

import "strconv"

// import "totala_reader/ta_files_read/scripts"

const (
	maxCobThreads = 8 // Value from the TA itself, according to the TA community
	maxStaticVars = 8
)

// Each scripted entity has one
type CobMachine struct {
	// ExecutedScript *scripts.CobScript
	SVars   [maxStaticVars]int32 // Static variables - TODO: check if they belong to the thread, not to the whole machine
	Threads [maxCobThreads]CobThread
}

// Search for any inactive thread, reset it and run from instruction in ip.
func (cm *CobMachine) AllocNewThread(ip, mask int32) {
	for i := range cm.Threads {
		if !cm.Threads[i].Active {
			cm.Threads[i].reset()
			cm.Threads[i].IP = ip
			cm.Threads[i].SigMask = mask
			cm.Threads[i].Active = true
			return
		}
	}
	panic("COB VPU error: can't alloc a new thread - thread space full.")
}

// SIGNAL [val] opcode. Destroys all threads with mask & val != 0.
func (cm *CobMachine) Signal(mask int32) string {
	retString := ""
	for i := range cm.Threads {
		if cm.Threads[i].SigMask&mask != 0 {
			cm.Threads[i].Active = false
			retString += strconv.Itoa(i) + " "
		}
	}
	// TODO: don't return the value
	if retString == "" {
		return "none"
	}
	return " " + retString
}
