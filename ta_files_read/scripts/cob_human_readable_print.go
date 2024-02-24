package scripts

import "strings"

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
		case CI_RETURN:
			disasmText = "RETURN"
		case CI_ALLOC_LOCAL_VAR:
			disasmText = "ALLOC LOCAL VAR"
		case CI_GET_VALUE:
			disasmText = "GET VALUE [port]"
		case CI_SET_VALUE:
			disasmText = "SET VALUE [val port]"
		case CI_CMP_LESS:
			disasmText = "COMPARE LESS"
		case CI_CMP_LEQ:
			disasmText = "COMPARE LESS OR EQUAL"
		case CI_CMP_EQ:
			disasmText = "COMPARE EQUAL"
		case CI_CMP_NEQ:
			disasmText = "COMPARE NOT EQUAL"
		case CI_OR:
			disasmText = "BITWISE OR"
		case CI_NEG:
			disasmText = "BITWISE NEGATE"
		case CI_SIGNAL:
			disasmText = "? SIGNAL [signal] ?"
		case CI_SETSIGMASK:
			disasmText = "? SET SIGNAL MASK [mask] ?"
		case CI_SUB:
			disasmText = "SUB [A B]"
		case CI_MUL:
			disasmText = "MULTIPLY [A B]"
		case CI_RAND:
			disasmText = "RANDOM"
		case CI_SLEEP:
			disasmText = "SLEEP"

		// 1 argument
		case CI_PUSH_CONST:
			disasmText = sprint("PUSH CONSTANT 0x%08X (dec: %4d)", nextval1, nextval1)
			ipIncrement = 2
		case CI_JMP:
			disasmText = sprint("JMP TO 0x%04X", nextval1)
			ipIncrement = 2
		case CI_JMP_IF_FALSE:
			disasmText = sprint("IF ZERO JMP TO 0x%04X", nextval1)
			ipIncrement = 2
		case CI_PUSH_LOCAL_VAR:
			disasmText = sprint("PUSH LOCAL VAR #%d", nextval1)
			ipIncrement = 2
		case CI_POP_LOCAL_VAR:
			disasmText = sprint("POP TO LOCAL VAR #%d", nextval1)
			ipIncrement = 2
		case CI_EMIT_SFX_FROM_PIECE:
			disasmText = sprint("EMIT SFX FROM PIECE #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_EXPLODE_PIECE:
			disasmText = sprint("EXPLODE PIECE #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_PUSH_STATIC_VAR:
			disasmText = sprint("PUSH STATIC VAR #%d", nextval1)
			ipIncrement = 2
		case CI_POP_STATIC_VAR:
			disasmText = sprint("POP TO STATIC VAR #%d", nextval1)
			ipIncrement = 2
		case CI_SHOW_OBJECT:
			disasmText = sprint("SHOW OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_HIDE_OBJECT:
			disasmText = sprint("HIDE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_CACHE:
			disasmText = sprint("CACHE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_DONTCACHE:
			disasmText = sprint("DISABLE CACHE FOR #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_SHADE:
			disasmText = sprint("SHADE OBJECT #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2
		case CI_DONTSHADE:
			disasmText = sprint("DISABLE SHADE FOR #%02d ('%s')", nextval1, cs.Pieces[nextval1])
			ipIncrement = 2

		// 2 arguments
		case CI_MOVE_OBJECT:
			disasmText = sprint("MOVE OBJECT #%d BY AXIS #%d [position, dir]", nextval1, nextval2)
			ipIncrement = 3
		case CI_ROTATE_OBJECT:
			disasmText = sprint("ROTATE OBJECT #%d BY AXIS #%d [speed, dir]", nextval1, nextval2)
			ipIncrement = 3
		case CI_SPIN_OBJECT:
			disasmText = sprint("SPIN OBJECT #%d BY AXIS #%d [speed, dir]", nextval1, nextval2)
			ipIncrement = 3
		case CI_STOP_SPIN_OBJECT:
			disasmText = sprint("STOP SPINNING OBJECT #%d BY AXIS #%d [deceleration]", nextval1, nextval2)
			ipIncrement = 3
		case CI_MOVE_NOW:
			disasmText = sprint("MOVE NOW OBJECT #%02d BY AXIS #%d [position]", nextval1, nextval2)
			ipIncrement = 3
		case CI_TURN_NOW:
			disasmText = sprint("TURN NOW OBJECT #%02d BY AXIS #%d [angle]", nextval1, nextval2)
			ipIncrement = 3
		case CI_WAIT_FOR_MOVE:
			disasmText = sprint("WAIT FOR MOVE OBJECT #%02d BY AXIS #%d", nextval1, nextval2)
			ipIncrement = 3
		case CI_START_SCRIPT:
			disasmText = sprint("? START SCRIPT #%d WITH %d PARAMS ?", nextval1, nextval2)
			ipIncrement = 3
		case CI_CALL_SCRIPT:
			disasmText = sprint("? CALL SCRIPT #%d WITH %d PARAMS ?", nextval1, nextval2)
			ipIncrement = 3

		// Unimplemented stuff:
		default:
			disasmText = sprint("<0x%08X (%s)>", opcode, sprintInt32AsBigEndianHex(opcode))
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
