package cover_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
	"github.com/alexfalkowski/gocovmerge/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestWriteProfilesRejectsDifferentModes(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(1, 1, 1, 2, 1)),
		test.Profile("b.go", "atomic", test.Block(1, 1, 1, 2, 1)),
	}

	err := cover.WriteProfiles(profiles, &bytes.Buffer{})
	require.ErrorIs(t, err, cover.ErrInvalidMode)
}

func TestAddProfileAppendsTrailingBlocks(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(1, 1, 1, 2, 1)),
	}

	merged, err := cover.AddProfile(profiles, test.Profile("a.go", "set",
		test.Block(1, 1, 1, 2, 0),
		test.Block(2, 1, 2, 2, 1),
	))
	require.NoError(t, err)
	require.Len(t, merged[0].Blocks, 2)
	require.Equal(t, test.Block(2, 1, 2, 2, 1), merged[0].Blocks[1])
}

func TestAddProfileRejectsOverlapAfter(t *testing.T) {
	profiles := []*cover.Profile{
		test.Profile("a.go", "set", test.Block(2, 1, 2, 5, 1)),
	}

	_, err := cover.AddProfile(profiles, test.Profile("a.go", "set", test.Block(1, 4, 2, 2, 1)))
	require.ErrorContains(t, err, "overlap after")
}
