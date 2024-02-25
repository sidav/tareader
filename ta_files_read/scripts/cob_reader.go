package scripts

import (
	tafilesread "totala_reader/ta_files_read"
)

func ReadCobFileFromReader(r *tafilesread.Reader) *CobScript {
	// Reading the header
	version := r.ReadIntFromBytesArray(0, 0)
	scriptsCount := r.ReadIntFromBytesArray(0, 4)
	piecesCount := r.ReadIntFromBytesArray(0, 8)
	// unknown0 := r.ReadIntFromBytesArray(0, 12)
	// unknown1 := r.ReadIntFromBytesArray(0, 16)
	// always0 := r.ReadIntFromBytesArray(0, 20)
	offsetToScriptCodeIndices := r.ReadIntFromBytesArray(0, 24)
	offsetToScriptNamesIndices := r.ReadIntFromBytesArray(0, 28)
	offsetToPieceNamesIndices := r.ReadIntFromBytesArray(0, 32)
	offsetToRawCode := r.ReadIntFromBytesArray(0, 36)
	// unknown3 := r.ReadIntFromBytesArray(0, 40)

	print(
		"Descriptor:\n"+ // <-- just to ease code alignment
			"  Version: %d\n"+
			"  Total scripts  %-8d   Total pieces     %-8d\n"+
			"  SCI offset     %08X   SNI offset       %08X\n"+
			"  Pieces offset  %08X   Raw code offset  %08X\n",
		version, scriptsCount, piecesCount,
		offsetToScriptCodeIndices,
		offsetToScriptNamesIndices,
		offsetToPieceNamesIndices,
		offsetToRawCode,
	)

	script := &CobScript{}
	print("Pieces: \n")
	for pNum := 0; pNum < piecesCount; pNum++ {
		pieceNameOffset := r.ReadIntFromBytesArray(offsetToPieceNamesIndices, pNum*4)
		pieceName := r.ReadNullTermStringFromBytesArray(pieceNameOffset, 0)
		script.Pieces = append(script.Pieces, pieceName)
		print("%3d:  %s\n", pNum, pieceName)
	}

	print("Reading scripts descriptors: \n")
	for sNum := 0; sNum < scriptsCount; sNum++ {
		scriptNameOffset := r.ReadIntFromBytesArray(offsetToScriptNamesIndices, sNum*4)
		scriptName := r.ReadNullTermStringFromBytesArray(scriptNameOffset, 0)
		// Offset to a script is calculated by: OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber] * 4)
		// ScriptNumber is given in int32s, so the final formula is OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber*4] * 4)
		currentScriptCodeOffset := r.ReadIntFromBytesArray(offsetToScriptCodeIndices, sNum*4)

		print("%3d:  %-20s at 0x%08X (local index 0x%08X); \n", sNum, scriptName,
			offsetToRawCode+currentScriptCodeOffset*4, currentScriptCodeOffset*4)

		script.ProcedureNames = append(script.ProcedureNames, scriptName)
		script.ProcedureAddresses = append(script.ProcedureAddresses, int32(currentScriptCodeOffset))
	}

	print("Reading raw code...\n")
	script.RawCode = readRawCodeFromCOB(r, offsetToRawCode, offsetToScriptCodeIndices)
	print("COB parsed.\n")
	return script
}

func readRawCodeFromCOB(r *tafilesread.Reader, offsetToRawCode, readUntilAddress int) []int32 {
	var rawCode []int32
	// Reading the script itself:
	var currWord int
	currInstrOffset := 0
	for offsetToRawCode+currInstrOffset < readUntilAddress {
		currWord = r.ReadIntFromBytesArray(offsetToRawCode, currInstrOffset)
		rawCode = append(rawCode, int32(currWord))
		currInstrOffset += 4
	}
	return rawCode
}
