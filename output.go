package main

import (
	"io"
	"os"
)

func output(out string) (io.Writer, error) {
	if len(out) > 0 {
		f, err := os.Create(out)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return os.Stdout, nil
}
