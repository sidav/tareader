package model

import (
	"totala_reader/model/cob"
	"totala_reader/ta_files_read/scripts/opcodes"
)

func (so *SimObject) CobStepAllThreads() {
	for i := range so.CobState.Threads {
		if so.CobState.Threads[i].Active {
			so.cobStepThread(&so.CobState.Threads[i])
		}
	}
}

func (so *SimObject) cobStepThread(t *cob.CobThread) {
	instruction := so.Script.RawCode[t.IP]
	switch instruction {
	case opcodes.CI_RETURN:
		// It's tougher than it seems. so TODO
		return
	// Opcodes in TODO backlog:
	case opcodes.CI_CACHE, opcodes.CI_DONTCACHE, opcodes.CI_SHADE, opcodes.CI_DONTSHADE:
		return
	}
}
