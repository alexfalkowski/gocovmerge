// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import (
	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
	"github.com/alexfalkowski/gocovmerge/v2/internal/flag"
	"github.com/alexfalkowski/gocovmerge/v2/internal/io"
	"github.com/alexfalkowski/gocovmerge/v2/internal/log"
)

func main() {
	logger := log.NewLogger()

	v, err := flag.Parse()
	if err != nil {
		logger.Fatal(err.Error())
	}

	files, err := v.Files()
	if err != nil {
		logger.Fatal(err.Error())
	}

	out, err := io.Output(v.Out())
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := cmd.Run(files, out); err != nil {
		logger.Fatal(err.Error())
	}
}
