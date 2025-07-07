package io

import (
	"io"
	"os"
)

// Output to a file or stdout.
func Output(out string) (io.Writer, error) {
	if len(out) > 0 {
		return os.Create(out)
	}

	return os.Stdout, nil
}
