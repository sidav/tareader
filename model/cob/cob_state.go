package cob

const maxCobThreads = 8

// Each scripted entity has one
type CobState struct {
	Threads [maxCobThreads]CobThread
}
