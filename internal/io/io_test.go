package io_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexfalkowski/gocovmerge/v2/internal/io"
	"github.com/stretchr/testify/require"
)

func TestFileOutputReplacesDestinationOnlyOnClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cover.out")
	require.NoError(t, os.WriteFile(path, []byte("previous\n"), 0o600))

	out := io.Output(path, &bytes.Buffer{})
	_, err := out.Write([]byte("mode: set\n"))
	require.NoError(t, err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "previous\n", string(got))

	require.NoError(t, out.Close())

	got, err = os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "mode: set\n", string(got))
}

func TestFileOutputRemovesTempFileWhenCommitFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, strings.Repeat("a", 300))

	out := io.Output(path, &bytes.Buffer{})
	_, err := out.Write([]byte("mode: set\n"))
	require.NoError(t, err)

	require.Error(t, out.Close())

	matches, err := filepath.Glob(filepath.Join(dir, ".gocovmerge-*"))
	require.NoError(t, err)
	require.Empty(t, matches)
}
