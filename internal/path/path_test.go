package path_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
	"github.com/stretchr/testify/require"
)

func TestFilesFiltersNestedFilesByWalkedPath(t *testing.T) {
	dir := t.TempDir()
	selected := writeTextFile(t, filepath.Join(dir, "nested", "selected.cov"))
	writeTextFile(t, filepath.Join(dir, "ignored.cov"))

	files, err := path.Files(dir, `nested[/\\].*\.cov$`, "")
	require.NoError(t, err)
	require.Equal(t, []string{selected}, files)
}

func TestFilesExcludesRelativeOutputAcrossSymlinkedPaths(t *testing.T) {
	base := t.TempDir()
	realDir := filepath.Join(base, "real")
	linkDir := filepath.Join(base, "link")
	require.NoError(t, os.Mkdir(realDir, 0o755))
	require.NoError(t, os.Symlink(realDir, linkDir))

	input := writeTextFile(t, filepath.Join(realDir, "a.out"))
	writeTextFile(t, filepath.Join(realDir, "merged.out"))

	t.Chdir(linkDir)

	files, err := path.Files(realDir, "", "merged.out")
	require.NoError(t, err)
	require.Equal(t, []string{input}, files)
}

func TestFilesAcceptsExcludedPathWithMissingParent(t *testing.T) {
	input := writeTextFile(t, filepath.Join(t.TempDir(), "a.out"))

	files, err := path.Files(filepath.Dir(input), "", filepath.Join(filepath.Dir(input), "missing", "cover.out"))
	require.NoError(t, err)
	require.Equal(t, []string{input}, files)
}

func TestFilesReportsExcludedPathSymlinkLoop(t *testing.T) {
	dir := t.TempDir()
	loop := filepath.Join(dir, "loop")
	require.NoError(t, os.Symlink(loop, loop))

	files, err := path.Files(dir, "", loop)
	require.Error(t, err)
	require.Nil(t, files)
}

func TestFilesReportsWalkedPathSymlinkLoop(t *testing.T) {
	dir := t.TempDir()
	loop := filepath.Join(dir, "loop")
	require.NoError(t, os.Symlink(loop, loop))

	files, err := path.Files(dir, "", filepath.Join(dir, "merged.out"))
	require.Error(t, err)
	require.Nil(t, files)
}

func TestFilesReportsInvalidPattern(t *testing.T) {
	files, err := path.Files(t.TempDir(), "[", "")
	require.Error(t, err)
	require.Nil(t, files)
}

func writeTextFile(t *testing.T, path string) string {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte("mode: set\n"), 0o600))

	return path
}
