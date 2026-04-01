package main

import (
	"errors"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
	"github.com/alexfalkowski/gocovmerge/v2/internal/flag"
	"github.com/alexfalkowski/gocovmerge/v2/internal/io"
	"github.com/alexfalkowski/gocovmerge/v2/internal/log"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	logger := log.NewLogger(stderr)

	v, err := flag.Parse(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		return 1
	}

	files, err := v.Files()
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	out := io.Output(v.Out(), stdout)

	if err := cmd.Run(files, out); err != nil {
		logger.Error(err.Error())
		return 1
	}

	if err := out.Close(); err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
