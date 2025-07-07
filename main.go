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

func files(dir, pattern string) ([]string, error) {
	if len(dir) > 0 {
		f, err := path.Files(dir, pattern)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return flag.Args(), nil
}

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

	files, err := files(dir, pattern)
	if err != nil {
		log.Fatal(err)
	}

	output, err := output(out)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Run(files, output); err != nil {
		log.Fatal(err)
	}
}
