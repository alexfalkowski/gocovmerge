package cover_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
)

func TestAddProfileRejectsDifferentModesAcrossFiles(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "set", block(1, 1, 1, 2, 1)),
	}

	_, err := cover.AddProfile(profiles, profile("b.go", "atomic", block(1, 1, 1, 2, 1)))
	if !errors.Is(err, cover.ErrInvalidMode) {
		t.Fatalf("expected ErrInvalidMode, got %v", err)
	}
}

func TestWriteProfilesRejectsDifferentModes(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "set", block(1, 1, 1, 2, 1)),
		profile("b.go", "atomic", block(1, 1, 1, 2, 1)),
	}

	err := cover.WriteProfiles(profiles, &bytes.Buffer{})
	if !errors.Is(err, cover.ErrInvalidMode) {
		t.Fatalf("expected ErrInvalidMode, got %v", err)
	}
}

func TestAddProfileAppendsTrailingBlocks(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "set", block(1, 1, 1, 2, 1)),
	}

	merged, err := cover.AddProfile(profiles, profile("a.go", "set",
		block(1, 1, 1, 2, 0),
		block(2, 1, 2, 2, 1),
	))
	if err != nil {
		t.Fatalf("expected merge to succeed, got %v", err)
	}

	if len(merged[0].Blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(merged[0].Blocks))
	}

	if got := merged[0].Blocks[1]; got != block(2, 1, 2, 2, 1) {
		t.Fatalf("expected trailing block to be appended, got %+v", got)
	}
}

func TestAddProfileRejectsOverlapBefore(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "set", block(1, 1, 1, 5, 1)),
	}

	_, err := cover.AddProfile(profiles, profile("a.go", "set", block(1, 4, 1, 10, 1)))
	if err == nil || !strings.Contains(err.Error(), "overlap before") {
		t.Fatalf("expected overlap before error, got %v", err)
	}
}

func TestAddProfileRejectsOverlapAfter(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "set", block(2, 1, 2, 5, 1)),
	}

	_, err := cover.AddProfile(profiles, profile("a.go", "set", block(1, 4, 2, 2, 1)))
	if err == nil || !strings.Contains(err.Error(), "overlap after") {
		t.Fatalf("expected overlap after error, got %v", err)
	}
}

func TestAddProfileRejectsDifferentNumStmtForSameBlock(t *testing.T) {
	profiles := []*cover.Profile{
		profile("a.go", "count", block(1, 1, 1, 2, 1)),
	}

	_, err := cover.AddProfile(profiles, profile("a.go", "count", blockWithNumStmt(1, 1, 1, 2, 2, 1)))
	if err == nil || !strings.Contains(err.Error(), "inconsistent NumStmt") {
		t.Fatalf("expected inconsistent NumStmt error, got %v", err)
	}
}

func profile(fileName, mode string, blocks ...cover.ProfileBlock) *cover.Profile {
	return &cover.Profile{FileName: fileName, Mode: mode, Blocks: blocks}
}

func block(startLine, startCol, endLine, endCol, count int) cover.ProfileBlock {
	return blockWithNumStmt(startLine, startCol, endLine, endCol, 1, count)
}

func blockWithNumStmt(startLine, startCol, endLine, endCol, numStmt, count int) cover.ProfileBlock {
	return cover.ProfileBlock{
		StartLine: startLine,
		StartCol:  startCol,
		EndLine:   endLine,
		EndCol:    endCol,
		NumStmt:   numStmt,
		Count:     count,
	}
}
