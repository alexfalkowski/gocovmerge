package io

import (
	"io"
	"os"
)

// Output returns a writer for out.
//
// If out is non-empty, Output creates (or truncates) the file at that path and
// returns the resulting writer. Otherwise it returns os.Stdout.
func Output(out string) (io.Writer, error) {
	if len(out) > 0 {
		return os.Create(out)
	}
	return os.Stdout, nil
}
