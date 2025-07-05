// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
)

func main() {
	var (
		outputFile string
		out        io.Writer
	)

	flag.StringVar(&outputFile, "o", "", "output file")
	flag.Parse()

	if len(outputFile) > 0 {
		f, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}

		out = f
	} else {
		out = os.Stdout
	}

	if err := cmd.Run(out, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
