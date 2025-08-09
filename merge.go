package main

import (
	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
	"github.com/alexfalkowski/gocovmerge/v2/internal/flag"
	"github.com/alexfalkowski/gocovmerge/v2/internal/io"
)

func merge(out, dir, pattern string) error {
	files, err := flag.Files(dir, pattern)
	if err != nil {
		return err
	}

	output, err := io.Output(out)
	if err != nil {
		return err
	}

	return cmd.Run(files, output)
}
