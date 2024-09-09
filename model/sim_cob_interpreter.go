package model

import (
	"math/rand"
	"strings"
	"totala_reader/model/cob"
	"totala_reader/ta_files_read/scripts/opcodes"
)

func (so *SimObject) CobExecAllThreads() {
	for i := range so.CobMachine.Threads {
		if so.CobMachine.Threads[i].Active {
			// Skip thread if it sleeps
			if so.CobMachine.Threads[i].SleepTicksRemaining > 0 {
				// Ticks are calculated at SLEEP command execution routine, see below.
				so.CobMachine.Threads[i].SleepTicksRemaining--
				continue
			}

			for so.cobStepThread(&so.CobMachine.Threads[i], i) {
				// just continue
			}
		}
	}
}

// Returns true if execution should continue
func (so *SimObject) cobStepThread(t *cob.CobThread, threadNum int) bool {
	var continueExec = true
	var ipIncrement int32 = 1
	var nextval1, nextval2 int32
	opcode := so.Script.RawCode[t.IP]
	if int(t.IP)+1 < len(so.Script.RawCode) {
		nextval1 = so.Script.RawCode[t.IP+1]
	}
	if int(t.IP)+2 < len(so.Script.RawCode) {
		nextval2 = so.Script.RawCode[t.IP+2]
	}
	var disasmText string

	switch opcode {
	// No arguments
	case opcodes.CI_RETURN:
		ipBefore := t.IP
		t.DoReturn()
		ipAfter := t.IP
		if t.Active {
			disasmText = sprint("RETURN (from 0x%04X to 0x%04X)", ipBefore, ipAfter)
		} else {
			disasmText = sprint("RETURN (call stack empty, deactivating the thread)")
		}
		continueExec = false
		ipIncrement = 0
	case opcodes.CI_ALLOC_LOCAL_VAR:
		disasmText = "ALLOC LOCAL VAR"
		t.DoAllocNewLocalVar()
	// case opcodes.CI_GET_VALUE:
	// 	disasmText = "GET VALUE [port]"
	// case opcodes.CI_GET_VALUE_WITH_ARGS:
	// 	disasmText = "? GET VALUE WITH ARGS [arg1 arg2 arg3 arg4 port] ?"
	// case opcodes.CI_SET_VALUE:
	// 	disasmText = "SET VALUE [val port]"
	case opcodes.CI_CMP_LESS:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a < b)
		disasmText = sprint("IF (%d < %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_CMP_LEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a <= b)
		disasmText = sprint("IF (%d <= %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_CMP_EQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a == b)
		disasmText = sprint("IF (%d == %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_CMP_NEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a != b)
		disasmText = sprint("IF (%d != %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_CMP_GREATER:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a > b)
		disasmText = sprint("IF (%d > %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_CMP_GEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a >= b)
		disasmText = sprint("IF (%d >= %d) PUSH 1 ELSE PUSH 0 (pushed %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_BITWISE_OR:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a | b)
		disasmText = sprint("BITWISE OR [%d | %d] (pushing %d)", a, b, t.DataStack.Peek())
	case opcodes.CI_LOGICAL_NOT: // TODO: check if it really logical (alternative is bitwise)
		a := t.DataStack.PopWord()
		t.DataStack.Push(1 ^ a)
		disasmText = sprint("? LOGICAL NOT %d (pushing %d) ?", a, t.DataStack.Peek())
	case opcodes.CI_SETSIGMASK:
		// Set a mask for thread-killing routine (SIGNAL opcode)
		t.SigMask = t.DataStack.PopWord()
		disasmText = sprint("SET SIGMASK 0x%08X", t.SigMask)
	case opcodes.CI_ADD:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a + b)
		disasmText = sprint("ADD: %X + %X (pushing %X)", a, b, t.DataStack.Peek())
	case opcodes.CI_SUB:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a - b)
		disasmText = sprint("SUB: %X - %X (pushing %X)", a, b, t.DataStack.Peek())
	case opcodes.CI_MUL:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a * b)
		disasmText = sprint("MUL: %X * %X (pushing %X)", a, b, t.DataStack.Peek())
	case opcodes.CI_DIV:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a / b)
		disasmText = sprint("DIV: %X / %X (pushing %X)", a, b, t.DataStack.Peek())
	case opcodes.CI_RAND:
		to := t.DataStack.PopWord()
		from := t.DataStack.PopWord()
		val := from + rand.Int31n(to-from)
		t.DataStack.Push(val)
		disasmText = sprint("RANDOM [%d...%d] (pushing %d)", from, to, val)
	case opcodes.CI_SIGNAL:
		// Destroy all the threads by signal mask.
		mask := t.DataStack.PopWord()
		sigResult := so.CobMachine.Signal(mask)
		disasmText = sprint("SIGNAL [0x%08X], stopped threads: [%s]", mask, sigResult)
		continueExec = false
	case opcodes.CI_SLEEP:
		duration := t.DataStack.PopWord()
		if duration < 0 {
			cobPanic("Negative sleep duration")
		}

		// TODO: investigate units for sleep calculation... It's either not ticks or there is a bug
		// Current implementation: consider sleep value as milliseconds

		// It's duration / (1000/simTicksPerSecond); here 1000/SimTicksPerSecond is amount of ms per COB tick
		// +500 is here for proper integer rounding
		ticksToSleep := (SimTicksPerSecond*duration + 500) / 1000
		disasmText = sprint("SLEEP FOR %dms (%d ticks)", duration, ticksToSleep)
		t.SetSleep(ticksToSleep)
		continueExec = false

	// 1 argument
	case opcodes.CI_PUSH_CONST:
		disasmText = sprint("PUSH CONSTANT 0x%08X (dec: %4d)", nextval1, nextval1)
		t.DataStack.Push(nextval1)
		ipIncrement = 2
	case opcodes.CI_JMP:
		disasmText = sprint("JMP TO 0x%04X", nextval1)
		ipIncrement = 0
		t.IP = nextval1
	case opcodes.CI_JMP_IF_FALSE:
		val := t.DataStack.PopWord()
		disasmText = sprint("IF %d == 0 JMP TO 0x%04X", val, nextval1)
		if val == 0 {
			ipIncrement = 0
			t.IP = nextval1
			break
		}
		ipIncrement = 2
	case opcodes.CI_PUSH_LOCAL_VAR: // UNSURE if the implementation is correct
		t.DataStack.Push(t.GetCurrentScopeLocalVar(nextval1))
		disasmText = sprint("PUSH LOCAL VAR #%d (pushed %d)", nextval1, t.DataStack.Peek())
		ipIncrement = 2
	case opcodes.CI_POP_LOCAL_VAR: // UNSURE if the implementation is correct
		val := t.DataStack.PopWord()
		t.SetCurrentScopeLocalVar(val, nextval1)
		disasmText = sprint("POP TO LOCAL VAR #%d (var $%d = %d)", nextval1, nextval1, val)
		ipIncrement = 2
	case opcodes.CI_PUSH_STATIC_VAR:
		t.DataStack.Push(so.CobMachine.SVars[nextval1])
		disasmText = sprint("PUSH STATIC VAR #%d (pushed %d)", nextval1, t.DataStack.Peek())
		ipIncrement = 2
	case opcodes.CI_POP_STATIC_VAR:
		so.CobMachine.SVars[nextval1] = t.DataStack.PopWord()
		disasmText = sprint("POP TO STATIC VAR #%d ($%d = %d)", nextval1, nextval1, so.CobMachine.SVars[nextval1])
		ipIncrement = 2

	case opcodes.CI_START_SCRIPT:
		if nextval2 > 0 {
			cobPanic("Please implement arguments for START (%d args requested)", nextval2)
		}
		sName := so.Script.ProcedureNames[nextval1]
		// IMPORTANT: new threads should be created with the current (i.e. inherited) signal mask.
		so.CobMachine.AllocNewThread(so.Script.ProcedureAddresses[nextval1], t.SigMask)
		disasmText = sprint("NEW THREAD: script #%d ('%s') WITH %d PARAMS FROM STACK", nextval1, sName, nextval2)
		ipIncrement = 3

	case opcodes.CI_CALL_SCRIPT:
		if nextval2 > 0 {
			cobPanic("Please implement arguments for CALL (%d args requested)", nextval2)
		}
		ipIncrement = 0 // DON'T auto-increase the IP, it will be manually increased below
		sName := so.Script.ProcedureNames[nextval1]
		disasmText = sprint("FROM 0x%04X CALL SCRIPT #%d ('%s') AT ADDR 0x%04X WITH %d PARAMS FROM STACK\n",
			t.IP, nextval1, sName, so.Script.ProcedureAddresses[nextval1], nextval2)
		t.IP += 3
		t.DoCall(so.Script.ProcedureAddresses[nextval1], nextval2)

	case opcodes.CI_SHOW_OBJECT:
		so.PiecesMapping[nextval1].Hidden = false
		disasmText = sprint("SHOW OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2

	case opcodes.CI_HIDE_OBJECT:
		so.PiecesMapping[nextval1].Hidden = true
		disasmText = sprint("HIDE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2

	case opcodes.CI_SPIN_OBJECT:
		speed := t.DataStack.PopWord()
		acceleration := t.DataStack.PopWord()
		if acceleration != 0 {
			cobPanic("Accelerated spin unimplemented")
		}
		disasmText = sprint("SPIN OBJECT #%d BY AXIS #%d [speed %d, acc %d]", nextval1, nextval2, speed, acceleration)
		so.PiecesMapping[nextval1].SetSpin(nextval2, speed, acceleration)
		ipIncrement = 3

	case opcodes.CI_ROTATE_OBJECT:
		angle := t.DataStack.PopWord()
		speed := t.DataStack.PopWord()
		disasmText = sprint("ROTATE OBJECT #%d BY AXIS #%d [to angle %d, speed %d]", nextval1, nextval2, angle, speed)
		so.PiecesMapping[nextval1].SetTurn(nextval2, speed, angle)
		ipIncrement = 3

	case opcodes.CI_MOVE_OBJECT:
		position := t.DataStack.PopWord()
		speed := t.DataStack.PopWord()
		disasmText = sprint("MOVE OBJECT #%d BY AXIS #%d [position %d, speed %d]", nextval1, nextval2, position, speed)
		so.PiecesMapping[nextval1].SetMove(nextval2, position, speed)
		ipIncrement = 3

	case opcodes.CI_MOVE_NOW:
		position := t.DataStack.PopWord()
		disasmText = sprint("MOVE NOW OBJECT #%02d BY AXIS #%d [position %d]", nextval1, nextval2, position)
		so.PiecesMapping[nextval1].moveNow(nextval2, position)
		ipIncrement = 3

	case opcodes.CI_TURN_NOW:
		angle := t.DataStack.PopWord()
		disasmText = sprint("TURN NOW OBJECT #%02d BY AXIS #%d [angle %d]", nextval1, nextval2, angle)
		so.PiecesMapping[nextval1].turnNow(nextval2, angle)
		ipIncrement = 3

	case opcodes.CI_STOP_SPIN_OBJECT:
		deceleration := t.DataStack.PopWord()
		if deceleration != 0 {
			cobPanic("Decelerated StopSpin unimplemented")
		}
		disasmText = sprint("STOP SPINNING OBJECT #%d BY AXIS #%d [deceleration %d]", nextval1, nextval2, deceleration)
		so.PiecesMapping[nextval1].IsSpinning[nextval2] = false
		so.PiecesMapping[nextval1].TurnSpeed[nextval2] = 0.0
		ipIncrement = 3

	case opcodes.CI_WAIT_FOR_TURN:
		if so.PiecesMapping[nextval1].CurrentTurn[nextval2] != so.PiecesMapping[nextval1].TargetTurn[nextval2] {
			disasmText = sprint("WAIT FOR TURN OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 0
			continueExec = false
		} else {
			disasmText = sprint("END WAIT FOR TURN OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 3
		}

	case opcodes.CI_WAIT_FOR_MOVE:
		if so.PiecesMapping[nextval1].CurrentMove[nextval2] != so.PiecesMapping[nextval1].TargetMove[nextval2] {
			disasmText = sprint("WAIT FOR MOVE OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 0
			continueExec = false
		} else {
			disasmText = sprint("END WAIT FOR MOVE OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 3
		}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Unimplemented or postponed stuff:
	case opcodes.CI_EMIT_SFX_FROM_PIECE:
		disasmText = sprint("UNIMPLEMENTED: EMIT SFX FROM PIECE #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_EXPLODE_PIECE:
		disasmText = sprint("UNIMPLEMENTED: EXPLODE PIECE #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_CACHE:
		disasmText = sprint("UNIMPLEMENTED: CACHE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_DONTCACHE:
		disasmText = sprint("UNIMPLEMENTED: DISABLE CACHE FOR #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_SHADE:
		disasmText = sprint("UNIMPLEMENTED: SHADE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_DONTSHADE:
		disasmText = sprint("UNIMPLEMENTED: DISABLE SHADE FOR #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
		ipIncrement = 2
	case opcodes.CI_ATTACH_UNIT:
		pieceNum := t.DataStack.PopWord()
		unitId := t.DataStack.PopWord()
		disasmText = sprint("UNIMPLEMENTED: ATTACH UNIT %d TO PIECE %d", unitId, pieceNum)
	case opcodes.CI_DROP_UNIT:
		unitId := t.DataStack.PopWord()
		disasmText = sprint("UNIMPLEMENTED: DROP UNIT %d", unitId)

	default:
		disasmText = sprint("IP 0x%04X -> Unknown opcode < 0x%08X > (next words 0x%08X and 0x%08X)", t.IP, opcode, nextval1, nextval2)
	}

	spaces := strings.Repeat("    ", threadNum)
	print("%sTrd %d -> IP %04X:  %s\n", spaces, threadNum, t.IP, disasmText)
	t.IP += ipIncrement
	return continueExec
}
