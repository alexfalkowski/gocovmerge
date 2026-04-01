package io

import (
	"io"
	"os"
)

// Writer is the output sink used by the CLI.
type Writer = io.Writer

// Output returns a writer for out.
//
// If out is non-empty, Output creates (or truncates) the file at that path and
// returns the resulting writer. Otherwise it returns stdout.
func Output(out string, stdout Writer) (Writer, error) {
	if len(out) > 0 {
		return os.Create(out)
	}
	return stdout, nil
}
