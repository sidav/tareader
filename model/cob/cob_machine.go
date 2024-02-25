package cob

// import "totala_reader/ta_files_read/scripts"

const maxCobThreads = 8

// Each scripted entity has one
type CobMachine struct {
	// ExecutedScript *scripts.CobScript
	Threads [maxCobThreads]CobThread
}
