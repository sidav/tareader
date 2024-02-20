package scripts

import (
	"fmt"
	tafilesread "totala_reader/ta_files_read"
)

type CobScript struct {
}

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

	fmt.Printf("Descriptor:\n  Version: %d, SC %d, PC %d,\n", version, scriptsCount, piecesCount)
	fmt.Printf("  CI offset %08X, SI offset %08X,\nPiece offset %08X, code offset %08X\n",
		offsetToScriptCodeIndices,
		offsetToScriptNamesIndices,
		offsetToPieceNamesIndices,
		offsetToRawCode,
	)

	fmt.Printf("Pieces: \n")
	for pNum := 0; pNum < piecesCount; pNum++ {
		pieceNameOffset := r.ReadIntFromBytesArray(offsetToPieceNamesIndices, pNum*4)
		pieceName := r.ReadNullTermStringFromBytesArray(pieceNameOffset, 0)
		fmt.Printf("  %s\n", pieceName)
	}

	fmt.Printf("Scripts: \n")
	for sNum := 0; sNum < scriptsCount; sNum++ {
		scriptNameOffset := r.ReadIntFromBytesArray(offsetToScriptNamesIndices, sNum*4)
		scriptName := r.ReadNullTermStringFromBytesArray(scriptNameOffset, 0)
		// Offset to a script is calculated by: OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber] * 4)
		// ScriptNumber is given in int32s, so the final formula is OffsetToScriptCode + (ScriptCodeIndexArray[ScriptNumber*4] * 4)
		currentScriptCodeOffset := r.ReadIntFromBytesArray(offsetToScriptCodeIndices, sNum*4) * 4

		fmt.Printf("  %-12s at 0x%08X (local index 0x%08X)\n", scriptName,
			offsetToRawCode+currentScriptCodeOffset, currentScriptCodeOffset)

		readScriptFromOpenedCOB(r, offsetToRawCode, currentScriptCodeOffset)
	}
}

func readScriptFromOpenedCOB(r *tafilesread.Reader, offsetToRawCode, currentScriptCodeOffset int) {
	// Reading the script itself:
	currOpcode := -1
	currInstrOffset := 0
	for currOpcode != CI_RET {
		currOpcode = r.ReadIntFromBytesArray(offsetToRawCode+currentScriptCodeOffset, currInstrOffset)
		if currInstrOffset/4 < 10 {
			// fmt.Printf("    0x%08X;\n", currOpcode)
		}
		currInstrOffset += 4
	}
	fmt.Printf("    RET found; script end reached with instr offset %d (total %d instructions).\n", currInstrOffset, currInstrOffset/4)
}
