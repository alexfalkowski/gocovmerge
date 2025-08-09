// gocovmerge takes the results from multiple `go test -coverprofile` runs and
// merges them into one profile
package main

import "github.com/alexfalkowski/gocovmerge/v2/internal/log"

func main() {
	logger := log.NewLogger()
	out, dir, pattern := flags()

	if err := merge(out, dir, pattern); err != nil {
		logger.Fatal(err.Error())
	}
}
