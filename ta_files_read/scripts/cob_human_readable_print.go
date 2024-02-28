package scripts

import (
	"strings"
	"totala_reader/ta_files_read/scripts/opcodes"
)

func (cs *CobScript) PrintHumanReadableDisassembly() {
	ip := 0
	var ipIncrement int
	var opcode, nextval1, nextval2 int32
	var disasmText string

	// Debug only:
	var unknownOpcodeStrings []string

	for ip < len(cs.RawCode) {
		for i := range cs.ProcedureAddresses {
			if cs.ProcedureAddresses[i] == int32(ip) {
				print("  --- begin proc %s\n", cs.ProcedureNames[i])
			}
		}
		ipIncrement = 1
		opcode = cs.RawCode[ip]
		if ip+1 < len(cs.RawCode) {
			nextval1 = cs.RawCode[ip+1]
		}
		if ip+2 < len(cs.RawCode) {
			nextval2 = cs.RawCode[ip+2]
		}

		switch opcode {
		// No arguments
		case opcodes.CI_RETURN:
			disasmText = "RETURN"
		case opcodes.CI_ALLOC_LOCAL_VAR:
			disasmText = "ALLOC LOCAL VAR"
		case opcodes.CI_GET_VALUE:
			disasmText = "GET VALUE [port]"
		case opcodes.CI_GET_VALUE_WITH_ARGS:
			disasmText = "? GET VALUE WITH ARGS [arg1 arg2 arg3 arg4 port] ?"
		case opcodes.CI_SET_VALUE:
			disasmText = "SET VALUE [val port]"
		case opcodes.CI_CMP_LESS:
			disasmText = "COMPARE LESS"
		case opcodes.CI_CMP_LEQ:
			disasmText = "COMPARE LESS OR EQUAL"
		case opcodes.CI_CMP_EQ:
			disasmText = "COMPARE EQUAL"
		case opcodes.CI_CMP_NEQ:
			disasmText = "COMPARE NOT EQUAL"
		case opcodes.CI_CMP_GREATER:
			disasmText = "COMPARE GREATER"
		case opcodes.CI_CMP_GEQ:
			disasmText = "COMPARE GREATER OR EQUAL"
		case opcodes.CI_BITWISE_OR:
			disasmText = "BITWISE OR"
		case opcodes.CI_NEG:
			disasmText = "BITWISE NEGATE"
		case opcodes.CI_LOGICAL_OR:
			disasmText = "LOGICAL OR"
		case opcodes.CI_LOGICAL_XOR:
			disasmText = "LOGICAL XOR"
		case opcodes.CI_LOGICAL_AND:
			disasmText = "LOGICAL AND"
		case opcodes.CI_SIGNAL:
			// Destroy all the threads with passing mask.
			disasmText = "SIGNAL [signal]"
		case opcodes.CI_SETSIGMASK:
			// Set a mask for thread-killing routine (SIGNAL opcode)
			disasmText = "SET SIGNAL MASK [mask]"
		case opcodes.CI_ADD:
			disasmText = "ADD [A B]"
		case opcodes.CI_SUB:
			disasmText = "SUB [A B]"
		case opcodes.CI_MUL:
			disasmText = "MULTIPLY [A B]"
		case opcodes.CI_DIV:
			disasmText = "DIVIDE [A B]"
		case opcodes.CI_RAND:
			disasmText = "RANDOM [FROM TO]"
		case opcodes.CI_SLEEP:
			disasmText = "SLEEP"

		// 1 argument
		case opcodes.CI_PUSH_CONST:
			disasmText = sprint("PUSH CONSTANT 0x%08X (dec: %4d)", nextval1, nextval1)
			ipIncrement = 2
		case opcodes.CI_JMP:
			disasmText = sprint("JMP TO 0x%04X", nextval1)
			ipIncrement = 2
		case opcodes.CI_JMP_IF_FALSE:
			disasmText = sprint("IF ZERO JMP TO 0x%04X", nextval1)
			ipIncrement = 2
		case opcodes.CI_PUSH_LOCAL_VAR:
			disasmText = sprint("PUSH LOCAL VAR #%d", nextval1)
			ipIncrement = 2
		case opcodes.CI_POP_LOCAL_VAR:
			disasmText = sprint("POP TO LOCAL VAR #%d", nextval1)
			ipIncrement = 2
		case opcodes.CI_EMIT_SFX_FROM_PIECE:
			disasmText = sprint("EMIT SFX FROM PIECE #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_EXPLODE_PIECE:
			disasmText = sprint("EXPLODE PIECE #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_PUSH_STATIC_VAR:
			disasmText = sprint("PUSH STATIC VAR #%d", nextval1)
			ipIncrement = 2
		case opcodes.CI_POP_STATIC_VAR:
			disasmText = sprint("POP TO STATIC VAR #%d", nextval1)
			ipIncrement = 2
		case opcodes.CI_SHOW_OBJECT:
			disasmText = sprint("SHOW OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_HIDE_OBJECT:
			disasmText = sprint("HIDE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_CACHE:
			disasmText = sprint("CACHE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_DONTCACHE:
			disasmText = sprint("DISABLE CACHE FOR #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_SHADE:
			disasmText = sprint("SHADE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case opcodes.CI_DONTSHADE:
			disasmText = sprint("DISABLE SHADE FOR #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2

		// 2 arguments
		case opcodes.CI_MOVE_OBJECT:
			disasmText = sprint("MOVE OBJECT #%d ('%s') BY AXIS #%d [position, speed]", nextval1, cs.Pieces[nextval1], nextval2)
			ipIncrement = 3
		case opcodes.CI_ROTATE_OBJECT:
			disasmText = sprint("ROTATE OBJECT #%d ('%s') BY AXIS #%d [angle, speed]", nextval1, cs.Pieces[nextval1], nextval2)
			ipIncrement = 3
		case opcodes.CI_SPIN_OBJECT:
			disasmText = sprint("SPIN OBJECT #%d BY AXIS #%d [speed, acceleration]", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_STOP_SPIN_OBJECT:
			disasmText = sprint("STOP SPINNING OBJECT #%d BY AXIS #%d [deceleration]", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_MOVE_NOW:
			disasmText = sprint("MOVE NOW OBJECT #%02d BY AXIS #%d [position]", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_TURN_NOW:
			disasmText = sprint("TURN NOW OBJECT #%02d BY AXIS #%d [angle]", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_WAIT_FOR_TURN:
			disasmText = sprint("WAIT FOR TURN OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_WAIT_FOR_MOVE:
			disasmText = sprint("WAIT FOR MOVE OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 3
		case opcodes.CI_START_SCRIPT:
			sName := cs.ProcedureNames[nextval1]
			// IMPORTANT: new threads should be created with the current (i.e. inherited) signal mask.
			disasmText = sprint("NEW THREAD FOR SCRIPT #%d ('%s') WITH %d PARAMS FROM STACK", nextval1, sName, nextval2)
			ipIncrement = 3
		case opcodes.CI_CALL_SCRIPT:
			sName := cs.ProcedureNames[nextval1]
			disasmText = sprint("CALL SCRIPT #%d ('%s') WITH %d PARAMS FROM STACK", nextval1, sName, nextval2)
			ipIncrement = 3

		// Unimplemented stuff:
		default:
			// disasmText = sprint("<0x%08X (%s)>", opcode, sprintInt32AsBigEndianHex(opcode))
			disasmText = sprint("< 0x%08X >", opcode)
			unknownOpcodeStrings = addStringToArrUnlessPresent(disasmText, unknownOpcodeStrings)
		}

		print("  %04X:  "+disasmText+"\n", ip)
		ip += ipIncrement
	}
	if len(unknownOpcodeStrings) > 0 {
		print("  Unknown words:\n  %s\n", strings.Join(unknownOpcodeStrings, "\n  "))
		panic("Implement please")
	}
}

func addStringToArrUnlessPresent(str string, arr []string) []string {
	for i := range arr {
		if arr[i] == str {
			return arr
		}
	}
	arr = append(arr, str)
	return arr
}
