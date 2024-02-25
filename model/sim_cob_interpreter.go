package model

import (
	"math/rand"
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
	// case opcodes.CI_RETURN:
	// 	disasmText = "RETURN"
	case opcodes.CI_ALLOC_LOCAL_VAR:
		disasmText = "?? ALLOC LOCAL VAR ??"
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
		disasmText = sprint("IF (%d < %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_CMP_LEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a <= b)
		disasmText = sprint("IF (%d <= %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_CMP_EQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a == b)
		disasmText = sprint("IF (%d == %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_CMP_NEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a != b)
		disasmText = sprint("IF (%d != %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_CMP_GREATER:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a > b)
		disasmText = sprint("IF (%d > %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_CMP_GEQ:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.PushBool(a >= b)
		disasmText = sprint("IF (%d >= %d) PUSH 1 ELSE PUSH 0", a, b)
	case opcodes.CI_BITWISE_OR:
		b := t.DataStack.PopWord()
		a := t.DataStack.PopWord()
		t.DataStack.Push(a | b)
		disasmText = sprint("BITWISE OR [%d | %d] (pushing %d)", a, b, t.DataStack.Peek())
		cobPanic("Check!")
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
		disasmText = sprint("MUL: %X / %X (pushing %X)", a, b, t.DataStack.Peek())
	case opcodes.CI_RAND:
		to := t.DataStack.PopWord()
		from := t.DataStack.PopWord()
		val := from + rand.Int31n(to-from)
		t.DataStack.Push(val)
		disasmText = sprint("RANDOM [%d...%d] (pushing %d)", from, to, val)
		cobPanic("Check!")
	// case opcodes.CI_SIGNAL:
	// 	// Destroy all the threads with passing mask.
	// 	disasmText = "SIGNAL [signal]"
	// case opcodes.CI_SLEEP:
	// remainingTicks := t.DataStack.PopWord()
	// if remainingTicks < 0 {
	// 	panic("Negative sleep duration")
	// }
	// if remainingTicks == 0 {
	// 	disasmText = "SLEEP ENDED"
	// } else {
	// 	disasmText = sprint("SLEEP (%d ticks remaining)", remainingTicks)
	// 	ipIncrement = 0
	// 	t.DataStack.Push(remainingTicks - 1)
	// }

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
	case opcodes.CI_PUSH_LOCAL_VAR:
		t.DataStack.Push(t.LVars[nextval1])
		disasmText = sprint("PUSH LOCAL VAR #%d (equals %d)", nextval1, t.DataStack.Peek())
		ipIncrement = 2
	case opcodes.CI_POP_LOCAL_VAR:
		val := t.DataStack.PopWord()
		t.LVars[nextval1] = val
		disasmText = sprint("POP TO LOCAL VAR #%d (#%d = %d)", nextval1, nextval1, val)
		ipIncrement = 2
	// case opcodes.CI_PUSH_STATIC_VAR:
	// 	disasmText = sprint("PUSH STATIC VAR #%d", nextval1)
	// 	ipIncrement = 2
	// case opcodes.CI_POP_STATIC_VAR:
	// 	disasmText = sprint("POP TO STATIC VAR #%d", nextval1)
	// 	ipIncrement = 2
	// case opcodes.CI_EMIT_SFX_FROM_PIECE:
	// 	disasmText = sprint("EMIT SFX FROM PIECE #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_EXPLODE_PIECE:
	// 	disasmText = sprint("EXPLODE PIECE #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_SHOW_OBJECT:
	// 	disasmText = sprint("SHOW OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_HIDE_OBJECT:
	// 	disasmText = sprint("HIDE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_CACHE:
	// 	disasmText = sprint("CACHE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_DONTCACHE:
	// 	disasmText = sprint("DISABLE CACHE FOR #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_SHADE:
	// 	disasmText = sprint("SHADE OBJECT #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2
	// case opcodes.CI_DONTSHADE:
	// 	disasmText = sprint("DISABLE SHADE FOR #%02d ('%s')", nextval1, so.Script.Pieces[nextval1])
	// 	ipIncrement = 2

	// 2 arguments
	// case opcodes.CI_MOVE_OBJECT:
	// 	disasmText = sprint("MOVE OBJECT #%d BY AXIS #%d [position, dir]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_ROTATE_OBJECT:
	// 	disasmText = sprint("ROTATE OBJECT #%d BY AXIS #%d [speed, dir]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_SPIN_OBJECT:
	// 	disasmText = sprint("SPIN OBJECT #%d BY AXIS #%d [speed, dir]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_STOP_SPIN_OBJECT:
	// 	disasmText = sprint("STOP SPINNING OBJECT #%d BY AXIS #%d [deceleration]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_MOVE_NOW:
	// 	disasmText = sprint("MOVE NOW OBJECT #%02d BY AXIS #%d [position]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_TURN_NOW:
	// 	disasmText = sprint("TURN NOW OBJECT #%02d BY AXIS #%d [angle]", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_WAIT_FOR_TURN:
	// 	disasmText = sprint("WAIT FOR TURN OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_WAIT_FOR_MOVE:
	// 	disasmText = sprint("WAIT FOR MOVE OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_START_SCRIPT:
	// 	sName := so.Script.ProcedureNames[nextval1]
	// 	// IMPORTANT: new threads should be created with the current (i.e. inherited) signal mask.
	// 	disasmText = sprint("NEW THREAD FOR SCRIPT #%d ('%s') WITH %d PARAMS FROM STACK", nextval1, sName, nextval2)
	// 	ipIncrement = 3
	// case opcodes.CI_CALL_SCRIPT:
	// 	sName := so.Script.ProcedureNames[nextval1]
	// 	disasmText = sprint("CALL SCRIPT #%d ('%s') WITH %d PARAMS FROM STACK", nextval1, sName, nextval2)
	// 	ipIncrement = 3

	// Unimplemented stuff:
	default:
		// disasmText = sprint("<0x%08X (%s)>", opcode, sprintInt32AsBigEndianHex(opcode))
		disasmText = sprint("Unknown opcode < 0x%08X > (next words 0x%08X and 0x%08X)", opcode, nextval1, nextval2)
	}

	print("  IP %04X:  "+disasmText+"\n", t.IP)
	t.IP += ipIncrement
}
