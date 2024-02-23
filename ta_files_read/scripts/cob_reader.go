package scripts

import (
	tafilesread "totala_reader/ta_files_read"
)

func ReadCobFileFromReader(r *tafilesread.Reader) {
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

	print("Descriptor:\n  Version: %d, SC %d, PC %d,\n", version, scriptsCount, piecesCount)
	print("  CI offset %08X, SI offset %08X,\nPiece offset %08X, code offset %08X\n",
		offsetToScriptCodeIndices,
		offsetToScriptNamesIndices,
		offsetToPieceNamesIndices,
		offsetToRawCode,
	)

	print("Pieces: \n")
	for pNum := 0; pNum < piecesCount; pNum++ {
		pieceNameOffset := r.ReadIntFromBytesArray(offsetToPieceNamesIndices, pNum*4)
		pieceName := r.ReadNullTermStringFromBytesArray(pieceNameOffset, 0)
		print("  %s\n", pieceName)
	}

	print("Scripts: \n")
	for sNum := 0; sNum < scriptsCount; sNum++ {
		script := &CobScript{}
		scriptNameOffset := r.ReadIntFromBytesArray(offsetToScriptNamesIndices, sNum*4)
		script.Name = r.ReadNullTermStringFromBytesArray(scriptNameOffset, 0)
		// Offset to a script is calculated by: OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber] * 4)
		// ScriptNumber is given in int32s, so the final formula is OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber*4] * 4)
		currentScriptCodeOffset := r.ReadIntFromBytesArray(offsetToScriptCodeIndices, sNum*4) * 4

		print("  %-12s at 0x%08X (local index 0x%08X); ", script.Name,
			offsetToRawCode+currentScriptCodeOffset, currentScriptCodeOffset)

		script.RawCode = readScriptFromOpenedCOB(r, offsetToRawCode, currentScriptCodeOffset)
		script.PrintHumanReadableDisassembly()
	}
}

func readScriptFromOpenedCOB(r *tafilesread.Reader, offsetToRawCode, currentScriptCodeOffset int) []int32 {
	var rawCode []int32
	// Reading the script itself:
	currWord := -1
	currInstrOffset := 0
	for currWord != CI_RET {
		currWord = r.ReadIntFromBytesArray(offsetToRawCode+currentScriptCodeOffset, currInstrOffset)
		rawCode = append(rawCode, int32(currWord))
		currInstrOffset += 4
	}
	print("RET found at instr offset %d (total %d words).\n", currInstrOffset, currInstrOffset/4)
	return rawCode
}
