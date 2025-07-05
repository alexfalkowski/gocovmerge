// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

func main() {
	var (
		out     string
		dir     string
		pattern string

		files  []string
		output io.Writer
	)

	flag.StringVar(&out, "o", "", "output file (if missing stdout)")
	flag.StringVar(&dir, "d", "", "directory of files (if missing paths passed in)")
	flag.StringVar(&pattern, "p", "", "pattern to filter directory (if missing all files)")
	flag.Parse()

	if len(out) > 0 {
		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}

		output = f
	} else {
		output = os.Stdout
	}

	if len(dir) > 0 {
		f, err := path.Files(dir, pattern)
		if err != nil {
			log.Fatal(err)
		}

		files = f
	} else {
		files = flag.Args()
	}

	if err := cmd.Run(output, files); err != nil {
		log.Fatal(err)
	}
}
