package test

import "github.com/alexfalkowski/gocovmerge/v2/internal/cover"

// Profile constructs a cover.Profile for tests.
func Profile(fileName, mode string, blocks ...cover.ProfileBlock) *cover.Profile {
	return &cover.Profile{FileName: fileName, Mode: mode, Blocks: blocks}
}

// Block constructs a cover.ProfileBlock with NumStmt set to 1.
func Block(startLine, startCol, endLine, endCol, count int) cover.ProfileBlock {
	return BlockWithNumStmt(startLine, startCol, endLine, endCol, 1, count)
}

// BlockWithNumStmt constructs a cover.ProfileBlock with an explicit NumStmt.
func BlockWithNumStmt(startLine, startCol, endLine, endCol, numStmt, count int) cover.ProfileBlock {
	return cover.ProfileBlock{
		StartLine: startLine,
		StartCol:  startCol,
		EndLine:   endLine,
		EndCol:    endCol,
		NumStmt:   numStmt,
		Count:     count,
	}
}
