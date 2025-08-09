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
	var (
		out     string
		dir     string
		pattern string
	)
	logger := log.NewLogger()

	flag.StringVar(&out, "o", "", "output file (if missing stdout)")
	flag.StringVar(&dir, "d", "", "directory of files (if missing paths passed in)")
	flag.StringVar(&pattern, "p", "", "pattern to filter directory (if missing all files)")
	flag.Parse()

	files, err := flag.Files(dir, pattern)
	if err != nil {
		logger.Fatal(err.Error())
	}

	output, err := io.Output(out)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := cmd.Run(files, output); err != nil {
		logger.Fatal(err.Error())
	}
}
