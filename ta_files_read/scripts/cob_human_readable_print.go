package scripts

import "strings"

func (cs *CobScript) PrintHumanReadableDisassembly() {
	ip := 0
	var ipIncrement int
	var instruction int32 = -1
	var currentText string

	// Debug only:
	var unknownOpcodeStrings []string

	for instruction != CI_RET {
		ipIncrement = 1
		instruction = cs.RawCode[ip]

		switch instruction {
		// No arguments
		case CI_RET:
			currentText = "RETURN"
		case CI_STACK_ALLOC:
			currentText = "STACK ALLOCATE"
		case CI_GET_UNIT_VAL:
			currentText = "GET UNIT VALUE"
		case CI_CMP_LESS:
			currentText = "COMPARE LESS"
		case CI_CMP_LEQ:
			currentText = "COMPARE LESS OR EQUAL"
		case CI_CMP_EQ:
			currentText = "COMPARE EQUAL"
		case CI_CMP_NEQ:
			currentText = "COMPARE NOT EQUAL"
		case CI_OR:
			currentText = "BITWISE OR"
		case CI_NOT:
			currentText = "BITWISE NOT"
		case CI_SIGNAL:
			currentText = "? SIGNAL [signal] ?"
		case CI_SETSIGMASK:
			currentText = "? SET SIGNAL MASK [mask] ?"
		case CI_SUB:
			currentText = "SUB [A B]"
		case CI_MUL:
			currentText = "MULTIPLY [A B]"
		case CI_RAND:
			currentText = "RANDOM"
		case CI_SLEEP:
			currentText = "SLEEP"
		case CI_SET_VALUE:
			currentText = "SET VALUE [val id]"

		// 1 argument
		case CI_PUSH_CONST:
			currentText = sprint("PUSH CONSTANT 0x%08X (%d)", cs.RawCode[ip+1], cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_JMP:
			currentText = sprint("JMP TO 0x%08X (%d)", cs.RawCode[ip+1], cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_JMP_IF_FALSE:
			currentText = sprint("JMP IF ZERO TO 0x%08X (%d)", cs.RawCode[ip+1], cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_PUSH_VAR:
			currentText = sprint("PUSH VAR #%d TO STACK", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_POP_VAR:
			currentText = sprint("SET VAR #%d EQ TO STACK TOP", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_EMIT_SFX_FROM_PIECE:
			currentText = sprint("EMIT SFX FROM PIECE #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_EXPLODE_PIECE:
			currentText = sprint("EXPLODE PIECE #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_PUSH_STATIC_VAR:
			currentText = sprint("PUSH STATIC VAR #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_POP_STATIC_VAR:
			currentText = sprint("SET STATIC VAR #%d EQ TO STACK TOP", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_SHOW_OBJECT:
			currentText = sprint("SHOW OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_HIDE_OBJECT:
			currentText = sprint("HIDE OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_CACHE:
			currentText = sprint("CACHE OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_DONTCACHE:
			currentText = sprint("DON'T CACHE OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_SHADE:
			currentText = sprint("SHADE OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2
		case CI_DONTSHADE:
			currentText = sprint("DON'T SHADE OBJECT #%d", cs.RawCode[ip+1])
			ipIncrement = 2

		// 2 arguments
		case CI_MOVE_OBJECT:
			currentText = sprint("MOVE OBJECT #%d BY AXIS #%d [speed, dir]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_ROTATE_OBJECT:
			currentText = sprint("ROTATE OBJECT #%d BY AXIS #%d [speed, dir]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_SPIN_OBJECT:
			currentText = sprint("SPIN OBJECT #%d BY AXIS #%d [speed, dir]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_STOP_SPIN_OBJECT:
			currentText = sprint("STOP SPINNING OBJECT #%d BY AXIS #%d [deceleration]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_MOVE_NOW:
			currentText = sprint("MOVE NOW OBJECT #%d BY AXIS #%d [position]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_TURN_NOW:
			currentText = sprint("TURN NOW OBJECT #%d BY AXIS #%d [angle]", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_WAIT_FOR_MOVE:
			currentText = sprint("WAIT FOR MOVE OBJECT #%d BY AXIS #%d", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_START_SCRIPT:
			currentText = sprint("START SCRIPT #%d WITH %d PARAMS", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3
		case CI_CALL_SCRIPT:
			currentText = sprint("CALL SCRIPT #%d WITH %d PARAMS", cs.RawCode[ip+1], cs.RawCode[ip+2])
			ipIncrement = 3

		// Unimplemented stuff:
		default:
			currentText = sprint("<0x%08X (%s)>", instruction, sprintInt32AsBigEndianHex(instruction))
			unknownOpcodeStrings = addStringToArrUnlessPresent(currentText, unknownOpcodeStrings)
		}

		print("    %03d: "+currentText+"\n", ip)
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
