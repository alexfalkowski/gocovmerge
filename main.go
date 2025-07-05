// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import (
	"flag"
	"log"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cmd"
)

func main() {
	flag.Parse()

	if err := cmd.Run(os.Stdout, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
