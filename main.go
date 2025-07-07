// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import (
	"log"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
	"github.com/alexfalkowski/gocovmerge/v2/internal/flag"
	"github.com/alexfalkowski/gocovmerge/v2/internal/io"
)

func main() {
	var (
		out     string
		dir     string
		pattern string
	)

	flag.StringVar(&out, "o", "", "output file (if missing stdout)")
	flag.StringVar(&dir, "d", "", "directory of files (if missing paths passed in)")
	flag.StringVar(&pattern, "p", "", "pattern to filter directory (if missing all files)")
	flag.Parse()

	files, err := flag.Files(dir, pattern)
	if err != nil {
		log.Fatal(err)
	}

	output, err := io.Output(out)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Run(files, output); err != nil {
		log.Fatal(err)
	}
}
