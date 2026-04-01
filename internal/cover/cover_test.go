package cover_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
	"github.com/alexfalkowski/gocovmerge/v2/internal/test"
)

func TestWriteProfilesRejectsDifferentModes(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(1, 1, 1, 2, 1)),
		test.Profile("b.go", "atomic", test.Block(1, 1, 1, 2, 1)),
	}

	err := cover.WriteProfiles(profiles, &bytes.Buffer{})
	if !errors.Is(err, cover.ErrInvalidMode) {
		t.Fatalf("expected ErrInvalidMode, got %v", err)
	}
}

func TestAddProfileAppendsTrailingBlocks(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(1, 1, 1, 2, 1)),
	}

	merged, err := cover.AddProfile(profiles, test.Profile("a.go", "set",
		test.Block(1, 1, 1, 2, 0),
		test.Block(2, 1, 2, 2, 1),
	))
	if err != nil {
		t.Fatalf("expected merge to succeed, got %v", err)
	}

	if len(merged[0].Blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(merged[0].Blocks))
	}

	if got := merged[0].Blocks[1]; got != test.Block(2, 1, 2, 2, 1) {
		t.Fatalf("expected trailing block to be appended, got %+v", got)
	}
}

func TestAddProfileRejectsOverlapAfter(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(2, 1, 2, 5, 1)),
	}

	_, err := cover.AddProfile(profiles, test.Profile("a.go", "set", test.Block(1, 4, 2, 2, 1)))
	if err == nil || !strings.Contains(err.Error(), "overlap after") {
		t.Fatalf("expected overlap after error, got %v", err)
	}
}
